package main

import (
	"embed"
	"log/slog"
	"os"

	"log"

	"github.com/SladkyCitron/slogcolor"
	"github.com/SladkyCitron/vlna/service"
	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	logger := slog.New(slogcolor.NewHandler(os.Stderr, slogcolor.DefaultOptions))
	slog.SetDefault(logger)

	app := application.New(application.Options{
		Name:        "Vlna",
		Description: "Internet radio player",
		Services: []application.Service{
			application.NewService(service.NewIPInfoService()),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
		Logger: logger,
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: "Vlna",
		// Window sized to the golden ratio (1000 / 618 ≈ 1.618).
		Width:  1000,
		Height: 618,
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundColour: application.NewRGB(6, 7, 15),
		URL:              "/",
		Frameless:        true,
	})

	// Run the application. This blocks until the application has been exited.
	err := app.Run()

	// If an error occurred while running the application, log it and exit.
	if err != nil {
		log.Fatal(err)
	}
}
