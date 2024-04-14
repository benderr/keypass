package confirm

import (
	"github.com/rivo/tview"
)

func New(text string, onConfirm func(), onCancel func()) *tview.Modal {
	return NewWithLabels(text, "Confirm", "Cancel", onConfirm, onCancel)
}

func NewWithLabels(text, labelConfirm, labelCancel string, onConfirm func(), onCancel func()) *tview.Modal {
	return tview.NewModal().
		SetText(text).
		AddButtons([]string{labelConfirm, labelCancel}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonIndex == 0 {
				onConfirm()
			} else {
				onCancel()
			}
		})
}
