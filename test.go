package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	APP := app.New()
	w1 := APP.NewWindow("Hello World")
	w1.Resize(fyne.NewSize(576, 370))
	w1.ShowAndRun()
}
