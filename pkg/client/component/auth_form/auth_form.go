package authform

import "github.com/rivo/tview"

type AuthForm struct {
	Login    string
	Password string
}

func New(title string, onSubmit func(a AuthForm) bool) *tview.Form {
	auth := &AuthForm{}

	form := tview.NewForm()

	form.SetTitle(title)

	form.AddInputField("Login", "", 20, nil, func(login string) {
		auth.Login = login
	})

	form.AddPasswordField("Password", "", 20, rune('*'), func(password string) {
		auth.Password = password
	})

	form.AddButton("Submit", func() {
		if ok := onSubmit(*auth); ok {
			auth.Login = ""
			auth.Password = ""
		}
	})

	return form
}
