package main

import (
	"context"
	"os"

	"github.com/energye/systray"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

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
