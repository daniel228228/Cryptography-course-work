package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	internalUI "crypto-app/ui"
)

type UI struct {
	BtnGenPublicPrivateKeys *widget.Button
	BtnEncrypt              *widget.Button
	BtnDecrypt              *widget.Button

	SelectMode *widget.Select
	Mode       string

	BtnFile   *widget.Button
	LabelFile *widget.Label
	File      string

	BtnKeyFile   *widget.Button
	LabelKeyFile *widget.Label
	KeyFile      string

	BtnPublicKeyFile   *widget.Button
	LabelPublicKeyFile *widget.Label
	PublicKeyFile      string

	BtnPrivateKeyFile   *widget.Button
	LabelPrivateKeyFile *widget.Label
	PrivateKeyFile      string

	BtnOutputDir   *widget.Button
	LabelOutputDir *widget.Label
	OutputDir      string

	ProgressBar *widget.ProgressBar
}

func MakeUI(win fyne.Window) *UI {
	ui := &UI{}

	ui.BtnGenPublicPrivateKeys = widget.NewButton("Сгенерировать ключи", func() {
		dialog.ShowInformation("Success", "Генерация ключей...", win)

		if err := internalUI.Gen(); err != nil {
			dialog.ShowError(err, win)
		} else {
			dialog.ShowInformation("Success", "Ключи сохранены в папке keys", win)
		}
	})
	ui.BtnGenPublicPrivateKeys.Resize(fyne.NewSize(200, 40))
	ui.BtnGenPublicPrivateKeys.Move(fyne.NewPos(30, 370))

	ui.BtnEncrypt = widget.NewButton("Зашифровать файл", func() {
		if err := internalUI.Encrypt(true, ui.File, ui.PublicKeyFile, ui.OutputDir, ui.Mode, ui.progress); err != nil {
			dialog.ShowError(err, win)
		} else {
			dialog.ShowInformation("Success", "Успешно зашифровано!", win)
		}
	})
	ui.BtnEncrypt.Resize(fyne.NewSize(210, 40))
	ui.BtnEncrypt.Move(fyne.NewPos(260, 370))

	ui.BtnDecrypt = widget.NewButton("Расшифровать файл", func() {
		if err := internalUI.Decrypt(true, ui.File, ui.KeyFile, ui.PublicKeyFile, ui.PrivateKeyFile, ui.OutputDir, ui.progress); err != nil {
			dialog.ShowError(err, win)
		} else {
			dialog.ShowInformation("Success", "Успешно расшифровано!", win)
		}
	})
	ui.BtnDecrypt.Resize(fyne.NewSize(210, 40))
	ui.BtnDecrypt.Move(fyne.NewPos(490, 370))

	ui.BtnFile = widget.NewButton("Выбрать файл", func() {
		showFileSelectionDialog(win, &ui.File, ui.LabelFile)
	})
	ui.BtnFile.Resize(fyne.NewSize(210, 40))
	ui.BtnFile.Move(fyne.NewPos(30, 40))

	ui.LabelFile = widget.NewLabel("")
	ui.LabelFile.Resize(fyne.NewSize(400, 40))
	ui.LabelFile.Move(fyne.NewPos(270, 40))

	ui.SelectMode = widget.NewSelect(
		[]string{"ECB", "CBC", "CFB", "OFB"},
		func(s string) {
			ui.Mode = s
		},
	)
	ui.SelectMode.Resize(fyne.NewSize(210, 40))
	ui.SelectMode.Move(fyne.NewPos(30, 290))

	ui.BtnKeyFile = widget.NewButton("Выбрать ключ", func() {
		showFileSelectionDialog(win, &ui.KeyFile, ui.LabelKeyFile)
	})
	ui.BtnKeyFile.Resize(fyne.NewSize(210, 40))
	ui.BtnKeyFile.Move(fyne.NewPos(30, 90))

	ui.LabelKeyFile = widget.NewLabel("")
	ui.LabelKeyFile.Resize(fyne.NewSize(400, 40))
	ui.LabelKeyFile.Move(fyne.NewPos(270, 90))

	ui.BtnPublicKeyFile = widget.NewButton("Выбрать публичный ключ", func() {
		showFileSelectionDialog(win, &ui.PublicKeyFile, ui.LabelPublicKeyFile)
	})
	ui.BtnPublicKeyFile.Resize(fyne.NewSize(210, 40))
	ui.BtnPublicKeyFile.Move(fyne.NewPos(30, 140))

	ui.LabelPublicKeyFile = widget.NewLabel("")
	ui.LabelPublicKeyFile.Resize(fyne.NewSize(400, 40))
	ui.LabelPublicKeyFile.Move(fyne.NewPos(270, 140))

	ui.BtnPrivateKeyFile = widget.NewButton("Выбрать приватный ключ", func() {
		showFileSelectionDialog(win, &ui.PrivateKeyFile, ui.LabelPrivateKeyFile)
	})
	ui.BtnPrivateKeyFile.Resize(fyne.NewSize(210, 40))
	ui.BtnPrivateKeyFile.Move(fyne.NewPos(30, 190))

	ui.LabelPrivateKeyFile = widget.NewLabel("")
	ui.LabelPrivateKeyFile.Resize(fyne.NewSize(400, 40))
	ui.LabelPrivateKeyFile.Move(fyne.NewPos(270, 190))

	ui.BtnOutputDir = widget.NewButton("Папка для результатов", func() {
		showDirSelectionDialog(win, &ui.OutputDir, ui.LabelOutputDir)
	})
	ui.BtnOutputDir.Resize(fyne.NewSize(210, 40))
	ui.BtnOutputDir.Move(fyne.NewPos(30, 240))

	ui.LabelOutputDir = widget.NewLabel("")
	ui.LabelOutputDir.Resize(fyne.NewSize(400, 40))
	ui.LabelOutputDir.Move(fyne.NewPos(270, 240))

	ui.ProgressBar = widget.NewProgressBar()
	ui.ProgressBar.Resize(fyne.NewSize(660, 30))
	ui.ProgressBar.Move(fyne.NewPos(30, 440))

	return ui
}

func (ui *UI) progress(v float64) {
	ui.ProgressBar.SetValue(v)
}

func showFileSelectionDialog(window fyne.Window, filepath *string, label *widget.Label) {
	fileDialog := dialog.NewFileOpen(func(result fyne.URIReadCloser, err error) {
		if err == nil && result != nil {
			*filepath = result.URI().Path()
			log.Println("Выбранный файл:", *filepath)

			label.SetText(*filepath)
		}
	}, window)

	fileDialog.Show()
}

func showDirSelectionDialog(window fyne.Window, filepath *string, label *widget.Label) {
	fileDialog := dialog.NewFolderOpen(func(result fyne.ListableURI, err error) {
		if err == nil && result != nil {
			*filepath = result.Path()
			log.Println("Выбранная папка:", *filepath)

			label.SetText(*filepath)
		}
	}, window)

	fileDialog.Show()
}
