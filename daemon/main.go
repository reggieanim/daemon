package main

import (
	_ "embed"

	"context"
	"embed"
	"os"

	"github.com/energye/systray"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed frontend/src/assets/images/logo-universal.png
var wailsIcon []byte

func main() {
	app := NewApp()

	err := wails.Run(&options.App{
		Title:             "daemon",
		Width:             640,
		HideWindowOnClose: true,
		Height:            480,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}

func createSystemTray(ctx context.Context) func() {

	return func() {
		systray.SetIcon(wailsIcon)
		show := systray.AddMenuItem("Show", "Show The Window")
		systray.AddSeparator()
		exit := systray.AddMenuItem("Exit", "Quit The Program")
		show.Click(func() {
			runtime.WindowShow(ctx)
		})

		exit.Click(func() {
			os.Exit(0)
		})
		systray.SetOnClick(func(menu systray.IMenu) {
			runtime.WindowShow(ctx)
			menu.ShowMenu()
		})
	}
}
