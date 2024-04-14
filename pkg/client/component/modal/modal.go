package modal

import (
	"github.com/rivo/tview"
)

func New(text string, onClose func()) *tview.Modal {
	return tview.NewModal().
		SetText(text).
		AddButtons([]string{"Close"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			onClose()
		})
}
