package main

import (
	"log"

	"fyne.io/fyne/v2"
	app "fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("recovered: ", r)
		}
	}()

	a := app.New()

	win := a.NewWindow("Crypto App")
	win.Resize(fyne.NewSize(720, 500))
	win.CenterOnScreen()

	ui := MakeUI(win)

	cont := container.NewWithoutLayout(ui.BtnGenPublicPrivateKeys,
		ui.BtnEncrypt,
		ui.BtnDecrypt,
		ui.BtnFile,
		ui.SelectMode,
		ui.LabelFile,
		ui.BtnKeyFile,
		ui.LabelKeyFile,
		ui.BtnPublicKeyFile,
		ui.LabelPublicKeyFile,
		ui.BtnPrivateKeyFile,
		ui.LabelPrivateKeyFile,
		ui.BtnOutputDir,
		ui.LabelOutputDir,
		ui.ProgressBar,
	)

	win.SetContent(cont)

	win.ShowAndRun()
}
