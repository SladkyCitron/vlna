package service

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/SladkyCitron/resona/afmt"
	"github.com/SladkyCitron/resona/aio"
	"github.com/SladkyCitron/resona/effect"
	"github.com/SladkyCitron/resona/freq"
	"github.com/SladkyCitron/resona/playback"
	_ "github.com/SladkyCitron/resona/playback/driver/oto"
	"github.com/SladkyCitron/vlna/discordrichpresence"
	"github.com/SladkyCitron/vlna/icy"
	"github.com/SladkyCitron/vlna/minimp3"
	"github.com/SladkyCitron/vlna/mpris"
	"github.com/smallnest/ringbuffer"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type PlayerService struct {
	configService *ConfigService
	playbackCtx   *playback.Context
	player        *playback.Player
	volume        *effect.Volume
	pausable      *aio.PausableReader
	curBody       io.ReadCloser
	cancel        context.CancelFunc
	ctx           context.Context
	streamCtx     context.Context
	rb            *ringbuffer.RingBuffer
	mpris         *mpris.Server
	discord       *discordrichpresence.Client
	startTime     time.Time
	startTimeOnce sync.Once
	curStation    Station
	isPlaying     bool
	streaming     bool
	mu            sync.Mutex
}

const discordAppID = "1559881627709218936"

func NewPlayerService(configService *ConfigService) *PlayerService {
	return &PlayerService{
		configService: configService,
		volume:        &effect.Volume{Volume: -15}, // default 75% = -15 dB
	}
}

func (s *PlayerService) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	var err error
	s.playbackCtx, err = playback.NewContext(afmt.Format{
		SampleRate:  44100 * freq.Hertz,
		NumChannels: 2,
	})
	if err != nil {
		return fmt.Errorf("failed to create playback context: %w", err)
	}
	s.ctx = ctx

	s.mpris, err = mpris.Init(s)
	if err != nil {
		return fmt.Errorf("failed to initialize MPRIS: %w", err)
	}

	s.discord = discordrichpresence.NewClient(discordAppID)
	err = s.discord.Connect()
	if err != nil {
		slog.Warn("Failed to connect to Discord Rich Presence", "err", err)
	} else {
		slog.Info("Connected to Discord Rich Presence")
	}

	return nil
}

func (s *PlayerService) ServiceShutdown() error {
	if err := s.discord.ClearActivity(); err != nil {
		return err
	}
	if err := s.discord.Close(); err != nil {
		return err
	}

	return mpris.Deinit()
}

func (s *PlayerService) Play(station Station) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.stopInternal()

	s.startTimeOnce.Do(func() {
		s.startTime = time.Now()
	})

	application.Get().Event.Emit("player:status", "Loading")

	s.rb = ringbuffer.New(1024 * 512) // 512 KB ring buffer

	ctx, cancel := context.WithCancel(s.ctx)
	s.cancel = cancel
	s.streamCtx = ctx

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, station.URLResolved, nil)
	if err != nil {
		return fmt.Errorf("failed to create network request: %w", err)
	}
	req.Header.Set("Icy-MetaData", "1")
	req.Header.Set("User-Agent", getUserAgent())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch stream: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	slog.Info("Connected to stream", "url", station.URLResolved, "status", resp.Status)

	s.curBody = resp.Body

	_metaint := resp.Header.Get("icy-metaint")
	metaint := 0
	if _metaint != "" {
		metaint, err = strconv.Atoi(_metaint)
		if err != nil {
			return fmt.Errorf("failed to parse icy-metaint: %w", err)
		}
	}

	slog.Info("ICY metadata interval", "metaint", metaint)

	reader := icy.NewReader(resp.Body, metaint, func(m map[string]string) {
		slog.Info("Received ICY metadata", "metadata", m)
		s.mpris.UpdateMetadata(m)
		application.Get().Event.Emit("player:icy-metadata", m)
	})

	// This read code was partly assisted by Copilot
	readErr := make(chan error, 1)
	go func() {
		tmp := make([]byte, 4096)
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			if s.rb.Free() < 8192 {
				time.Sleep(20 * time.Millisecond)
				continue
			}

			n, err := reader.Read(tmp)
			if n > 0 {
				_, _ = s.rb.Write(tmp[:n])
			}
			if err != nil {
				if ctx.Err() == nil {
					application.Get().Event.Emit("player:error", fmt.Sprintf("Stream read error: %v", err))
					slog.Error("Stream read error", "err", err)
					readErr <- err

					s.mu.Lock()
					if s.streaming && s.streamCtx == ctx {
						s.stopInternal()
					}
					s.mu.Unlock()
				}
				break
			}
		}
	}()

	// wait for buffering
	s.mu.Unlock()
	slog.Info("Buffering...")
	for s.rb.Length() < 32*1024 {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-readErr:
			// clean up
			cancel()
			s.cancel = nil
			_ = resp.Body.Close()
			s.curBody = nil
			s.streamCtx = nil
			s.streaming = false
			s.mpris.UpdatePlaybackStatus("Stopped")
			application.Get().Event.Emit("player:status", "Stopped")
			return fmt.Errorf("stream disconnected while buffering: %w", err)
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
	s.mu.Lock()

	deco, err := minimp3.NewDecoder(s.rb)
	if err != nil {
		resp.Body.Close()
		return fmt.Errorf("failed to decode MP3: %w", err)
	}

	s.pausable = aio.NewPausableReader(effect.Reader(deco, s.volume))
	s.player = s.playbackCtx.NewPlayer(s.pausable)
	s.player.Play()

	s.isPlaying = true
	s.streaming = true
	s.curStation = station
	application.Get().Event.Emit("player:status", "Playing")
	s.mpris.UpdatePlaybackStatus("Playing")
	s.UpdateDiscordRPC()

	return nil
}

func (s *PlayerService) UpdateDiscordRPC() {
	if !s.configService.GetConfig().DiscordRPC {
		_ = s.discord.ClearActivity()
		return
	}

	if s.isPlaying {
		activity := discordrichpresence.NewActivity().
			State(s.curStation.Name).
			//Details(s.curStation.Name).
			Type(2).
			StartTime(s.startTime).
			Build()
		if err := s.discord.SetActivity(activity); err != nil {
			slog.Error("Failed to set Discord Rich Presence activity", "err", err)
		}
	} else {
		if err := s.discord.ClearActivity(); err != nil {
			slog.Error("Failed to clear Discord Rich Presence activity", "err", err)
		}
	}
}

func (s *PlayerService) Pause() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.pausable != nil {
		s.isPlaying = false
		s.pausable.Pause()
		s.mpris.UpdatePlaybackStatus("Paused")
		application.Get().Event.Emit("player:status", "Paused")
		s.UpdateDiscordRPC()
	}
}

func (s *PlayerService) Resume() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.pausable != nil {
		s.isPlaying = true
		s.pausable.Resume()
		s.mpris.UpdatePlaybackStatus("Playing")
		application.Get().Event.Emit("player:status", "Playing")
		s.UpdateDiscordRPC()
	}
}

func (s *PlayerService) SetVolume(vol float64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if vol <= 0 {
		s.volume.Volume = -100 // mute
	} else {
		s.volume.Volume = -60 * (1 - vol)
	}
	s.mpris.UpdateVolume(vol)
	slog.Info("Volume set", "vol", vol, "actualVol", s.volume.Volume)
}

func (s *PlayerService) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.stopInternal()
}

func (s *PlayerService) stopInternal() {
	slog.Info("Stopping streaming")
	if s.cancel != nil {
		s.cancel()
		s.cancel = nil
	}
	s.streamCtx = nil
	s.streaming = false
	if s.curBody != nil {
		s.curBody.Close()
		s.curBody = nil
	}
	if s.player != nil {
		s.player.Stop()
		s.player = nil
	}
	s.pausable = nil
	s.isPlaying = false
	s.mpris.UpdatePlaybackStatus("Stopped")
	application.Get().Event.Emit("player:status", "Stopped")
}

//wails:ignore
func (s *PlayerService) IsPlaying() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.isPlaying
}
