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
	"github.com/SladkyCitron/resona/audio"
	"github.com/SladkyCitron/resona/freq"
	"github.com/SladkyCitron/resona/playback"
	_ "github.com/SladkyCitron/resona/playback/driver/oto"
	"github.com/SladkyCitron/vlna/icy"
	"github.com/SladkyCitron/vlna/minimp3"
	"github.com/smallnest/ringbuffer"
	"github.com/wailsapp/wails/v3/pkg/application"
)

type PlayerService struct {
	playbackCtx *playback.Context
	player      *playback.Player
	source      *audio.Source
	curBody     io.ReadCloser
	cancel      context.CancelFunc
	ctx         context.Context
	rb          *ringbuffer.RingBuffer
	mu          sync.Mutex
}

func NewPlayerService() *PlayerService {
	return &PlayerService{}
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

	return nil
}

func (s *PlayerService) Play(url string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.stopInternal()

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

	s.source = audio.NewSource(deco)
	s.player = s.playbackCtx.NewPlayer(s.source)
	s.player.Play()

	return nil
}

func (s *PlayerService) Pause() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.source != nil {
		s.source.Pause()
	}
}

func (s *PlayerService) Resume() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.source != nil {
		s.source.Resume()
	}
}

func (s *PlayerService) SetVolume(gain float64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.source != nil {
		s.source.SetVolume(gain)
	}
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
	s.source = nil
}
