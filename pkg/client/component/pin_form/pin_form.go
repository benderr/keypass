package pinform

import "github.com/rivo/tview"

type PinForm struct {
	Pin string
}

func New(text string, onSubmit func(a PinForm) bool) *tview.Form {
	auth := &PinForm{}

	form := tview.NewForm()

	form.AddInputField(text, "", 10, nil, func(pin string) {
		auth.Pin = pin
	})

	form.AddButton("Enter", func() {
		if ok := onSubmit(*auth); ok {
			auth.Pin = ""
		}
	})

	return form
	// grid := tview.NewGrid().
	// 	SetColumns(0, 40, 0).
	// 	SetRows(0, 40, 0).
	// 	AddItem(form, 1, 1, 1, 1, 0, 0, true)

	// return grid
}
