//go:build linux || freebsd || openbsd || netbsd

package mpris

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
	"github.com/godbus/dbus/v5/prop"
)

type PlayerService interface {
	Pause()
	Resume()
	Stop()
	SetVolume(float64)
	IsPlaying() bool
}

type root struct{}

func (*root) Raise() *dbus.Error {
	slog.Info("Raise requested")
	return nil
}

func (*root) Quit() *dbus.Error {
	slog.Info("Quit requested")
	os.Exit(0)
	return nil
}

type Player struct {
	service PlayerService
}

func (p *Player) Next() *dbus.Error {
	// no-op
	return nil
}

func (p *Player) Previous() *dbus.Error {
	// no-op
	return nil
}

func (p *Player) Pause() *dbus.Error {
	slog.Info("Pause requested")
	p.service.Pause()
	return nil
}

func (p *Player) PlayPause() *dbus.Error {
	slog.Info("PlayPause requested")
	if p.service.IsPlaying() {
		p.service.Pause()
	} else {
		p.service.Resume()
	}
	return nil
}

func (p *Player) Stop() *dbus.Error {
	slog.Info("Stop requested")
	p.service.Stop()
	return nil
}

func (p *Player) Play() *dbus.Error {
	slog.Info("Play requested")
	p.service.Resume()
	return nil
}

type Server struct {
	props *prop.Properties
}

func (s *Server) UpdateMetadata(metadata map[string]string) {
	values := make(map[string]dbus.Variant)
	if title := metadata["StreamTitle"]; title != "" {
		parts := strings.SplitN(title, " - ", 2)
		if len(parts) == 2 {
			values["xesam:artist"] = dbus.MakeVariant([]string{parts[0]})
			values["xesam:title"] = dbus.MakeVariant(parts[1])
		} else {
			values["xesam:title"] = dbus.MakeVariant(title)
		}
	}
	values["mpris:trackid"] = dbus.MakeVariant(dbus.ObjectPath("/org/mpris/MediaPlayer2/TrackList/NoTrack"))
	s.props.SetMust("org.mpris.MediaPlayer2.Player", "Metadata", values)
}

func (s *Server) UpdatePlaybackStatus(status string) {
	s.props.SetMust("org.mpris.MediaPlayer2.Player", "PlaybackStatus", status)
}

func (s *Server) UpdateVolume(volume float64) {
	s.props.SetMust("org.mpris.MediaPlayer2.Player", "Volume", volume)
}

func Init(service PlayerService) (*Server, error) {
	conn, err := dbus.SessionBus()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to session bus: %w", err)
	}

	busName := "org.mpris.MediaPlayer2.vlna"
	reply, err := conn.RequestName(busName, dbus.NameFlagDoNotQueue)
	if err != nil || reply != dbus.RequestNameReplyPrimaryOwner {
		return nil, fmt.Errorf("failed to request bus name %s: %w", busName, err)
	}

	propsSpec := prop.Map{
		"org.mpris.MediaPlayer2": {
			"CanQuit":      {Value: true, Writable: false, Emit: prop.EmitTrue},
			"CanRaise":     {Value: true, Writable: false, Emit: prop.EmitTrue},
			"HasTrackList": {Value: false, Writable: false, Emit: prop.EmitTrue},
			"Identity":     {Value: "Vlna", Writable: false, Emit: prop.EmitTrue},
			"DesktopEntry": {Value: "vlna", Writable: false, Emit: prop.EmitTrue},
		},
		"org.mpris.MediaPlayer2.Player": {
			"PlaybackStatus": {Value: "Stopped", Writable: false, Emit: prop.EmitTrue},
			"LoopStatus":     {Value: "None", Writable: false, Emit: prop.EmitTrue},
			"Rate":           {Value: 1.0, Writable: false, Emit: prop.EmitTrue},
			"Shuffle":        {Value: false, Writable: false, Emit: prop.EmitTrue},
			"Metadata":       {Value: map[string]dbus.Variant{}, Writable: false, Emit: prop.EmitTrue},
			"Volume": {Value: 0.75, Writable: true, Emit: prop.EmitTrue, Callback: func(change *prop.Change) *dbus.Error {
				volume, ok := change.Value.(float64)
				if !ok {
					return prop.ErrInvalidArg
				}
				service.SetVolume(volume)
				return nil
			}},
			"Position":      {Value: int64(0), Writable: false, Emit: prop.EmitConst},
			"MinimumRate":   {Value: 1.0, Writable: false, Emit: prop.EmitTrue},
			"MaximumRate":   {Value: 1.0, Writable: false, Emit: prop.EmitTrue},
			"CanControl":    {Value: true, Writable: false, Emit: prop.EmitTrue},
			"CanPlay":       {Value: true, Writable: false, Emit: prop.EmitTrue},
			"CanPause":      {Value: true, Writable: false, Emit: prop.EmitTrue},
			"CanSeek":       {Value: false, Writable: false, Emit: prop.EmitTrue},
			"CanGoNext":     {Value: false, Writable: false, Emit: prop.EmitTrue},
			"CanGoPrevious": {Value: false, Writable: false, Emit: prop.EmitTrue},
		},
	}

	props, err := prop.Export(conn, "/org/mpris/MediaPlayer2", propsSpec)
	if err != nil {
		return nil, fmt.Errorf("failed to export properties: %w", err)
	}

	player := &Player{service: service}

	if err := conn.Export(&root{}, "/org/mpris/MediaPlayer2", "org.mpris.MediaPlayer2"); err != nil {
		return nil, fmt.Errorf("failed to export root interface: %w", err)
	}
	if err := conn.Export(player, "/org/mpris/MediaPlayer2", "org.mpris.MediaPlayer2.Player"); err != nil {
		return nil, fmt.Errorf("failed to export player interface: %w", err)
	}
	node := &introspect.Node{
		Name: "/org/mpris/MediaPlayer2",
		Interfaces: []introspect.Interface{
			introspect.IntrospectData,
			prop.IntrospectData,
			{
				Name:    "org.mpris.MediaPlayer2",
				Methods: introspect.Methods(&root{}),
				Properties: props.Introspection(
					"org.mpris.MediaPlayer2",
				),
			},
			{
				Name:    "org.mpris.MediaPlayer2.Player",
				Methods: introspect.Methods(player),
				Properties: props.Introspection(
					"org.mpris.MediaPlayer2.Player",
				),
			},
		},
	}
	if err := conn.Export(
		introspect.NewIntrospectable(node),
		"/org/mpris/MediaPlayer2",
		"org.freedesktop.DBus.Introspectable",
	); err != nil {
		return nil, fmt.Errorf("failed to export introspection: %w", err)
	}

	return &Server{props: props}, nil
}

func Deinit() error {
	// SessionBus returns a shared connection to the session bus so this is okay
	// to call multiple times
	conn, err := dbus.SessionBus()
	if err != nil {
		return err
	}
	return conn.Close()
}
