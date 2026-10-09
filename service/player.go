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
	"github.com/SladkyCitron/vlna/icy"
	"github.com/SladkyCitron/vlna/minimp3"
	"github.com/SladkyCitron/vlna/mpris"
	"github.com/smallnest/ringbuffer"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type PlayerService struct {
	playbackCtx *playback.Context
	player      *playback.Player
	gain        *effect.Gain
	pausable    *aio.PausableReader
	curBody     io.ReadCloser
	cancel      context.CancelFunc
	ctx         context.Context
	rb          *ringbuffer.RingBuffer
	mpris       *mpris.Server
	isPlaying   bool
	mu          sync.Mutex
}

func NewPlayerService() *PlayerService {
	return &PlayerService{
		gain: &effect.Gain{Gain: -0.25},
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

	return nil
}

func (s *PlayerService) ServiceShutdown() error {
	return mpris.Deinit()
}

func (s *PlayerService) Play(url string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.stopInternal()

	application.Get().Event.Emit("player:status", "Loading")

	s.rb = ringbuffer.New(1024 * 1024) // 1 MB ring buffer

	ctx, cancel := context.WithCancel(s.ctx)
	s.cancel = cancel

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
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

	slog.Info("Connected to stream", "url", url, "status", resp.Status)

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
					slog.Error("Stream read error", "err", err)
				}
				break
			}
		}
	}()

	// wait for buffering
	slog.Info("Buffering...")
	for s.rb.Length() < 32*1024 {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}

	deco, err := minimp3.NewDecoder(s.rb)
	if err != nil {
		resp.Body.Close()
		return fmt.Errorf("failed to decode MP3: %w", err)
	}

	s.pausable = aio.NewPausableReader(effect.Reader(deco, s.gain))
	s.player = s.playbackCtx.NewPlayer(s.pausable)
	s.player.Play()

	s.mpris.UpdatePlaybackStatus("Playing")
	application.Get().Event.Emit("player:status", "Playing")
	s.isPlaying = true

	return nil
}

func (s *PlayerService) Pause() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.pausable != nil {
		s.pausable.Pause()
		s.mpris.UpdatePlaybackStatus("Paused")
		application.Get().Event.Emit("player:status", "Paused")
		s.isPlaying = false
	}
}

func (s *PlayerService) Resume() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.pausable != nil {
		s.pausable.Resume()
		s.mpris.UpdatePlaybackStatus("Playing")
		application.Get().Event.Emit("player:status", "Playing")
		s.isPlaying = true
	}
}

func (s *PlayerService) SetVolume(gain float64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.gain.Gain = gain - 1
	s.mpris.UpdateVolume(gain)
	slog.Info("Volume set", "gain", gain, "actualGain", s.gain.Gain)
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
	if s.curBody != nil {
		s.curBody.Close()
		s.curBody = nil
	}
	if s.player != nil {
		s.player.Stop()
		s.player = nil
	}
	s.pausable = nil
	s.mpris.UpdatePlaybackStatus("Stopped")
	application.Get().Event.Emit("player:status", "Stopped")
	s.isPlaying = false
}

//wails:ignore
func (s *PlayerService) IsPlaying() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.isPlaying
}
