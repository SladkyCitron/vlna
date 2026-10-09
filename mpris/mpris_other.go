//go:build !linux && !freebsd && !openbsd && !netbsd

package mpris

type PlayerService interface {
	Pause()
	Resume()
	Stop()
	SetVolume(float64)
	IsPlaying() bool
}

type Server struct{}

func Init(PlayerService) (*Server, error) {
	return &Server{}, nil
}

func (*Server) UpdateMetadata(map[string]string) {
}

func (*Server) UpdatePlaybackStatus(string) {
}

func (*Server) UpdateVolume(float64) {
}

func Deinit() error {
	return nil
}
