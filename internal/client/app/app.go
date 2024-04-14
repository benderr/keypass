package app

import (
	"context"
	"fmt"
	"time"

	"github.com/benderr/keypass/internal/client/dto"
	"github.com/benderr/keypass/internal/client/session"
	authform "github.com/benderr/keypass/pkg/client/component/auth_form"
	"github.com/benderr/keypass/pkg/client/component/confirm"
	listform "github.com/benderr/keypass/pkg/client/component/list_form"
	"github.com/benderr/keypass/pkg/client/component/modal"
	pinform "github.com/benderr/keypass/pkg/client/component/pin_form"
	recordform "github.com/benderr/keypass/pkg/client/component/record_form"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type IAppLogic interface {
	Login(login string, pass string) error
	Register(login string, pass string) error
	Logout() error
	LoadUser() error

	SyncRecords() error
	GetRecords() ([]dto.ClientRecord, error)
	DeleteRecord(ID string) error
	UpdateRecord(record dto.ServerRecord) error
	AddRecord(record dto.ServerRecord) error

	GetSessionState() session.State
	SuspendSession() error
	CheckPin(pin string) (bool, error)
	SetPin(pin string) error

	SaveBinaryFile(r *dto.ClientRecord) (string, error)
}

type appClient struct {
	logic        IAppLogic
	app          *tview.Application
	pages        *tview.Pages
	lastActivity time.Time
}

// New return instance terminal user interface
func New(al IAppLogic) *appClient {
	app := tview.NewApplication()

	ac := &appClient{
		pages:        tview.NewPages(),
		app:          app,
		logic:        al,
		lastActivity: time.Now(),
	}

	ac.init()
	return ac
}

func (a *appClient) init() {
	if err := a.logic.LoadUser(); err != nil {
		a.showModal(err.Error(), "main")
		return
	}
	a.SwitchTo("main")
}

// SwitchTo change page in terminal
func (a *appClient) SwitchTo(page string) {
	if page == "login" || page == "register" {
		a.logic.Logout()
		a.buildPage(page)
		return
	}

	sessionState := a.logic.GetSessionState()

	if sessionState == session.NoSession {
		a.buildPage("login")
		return
	}

	if sessionState == session.NeedPin {
		a.buildPage("setPin")
		return
	}

	if sessionState == session.Suspended {
		a.buildPage("checkPin")
		return
	}

	a.buildPage(page)
}

func (a *appClient) buildPage(page string) {
	switch page {
	case "register":
		a.pages.AddAndSwitchToPage("register", authform.New("Registration", a.onRegister).AddButton("Go to login", func() { a.SwitchTo("login") }), true)
	case "setPin":
		a.pages.AddAndSwitchToPage("setPin", pinform.New("Protect PIN", a.onSetPin).AddButton("Logout", a.onLogout), true)
	case "checkPin":
		a.pages.AddAndSwitchToPage("checkPin", pinform.New("Enter PIN", a.onCheckPin).AddButton("Logout", a.onLogout), true)
	case "login":
		a.pages.AddAndSwitchToPage("login", authform.New("Authorization", a.onLogin).AddButton("Go to registration", func() { a.SwitchTo("register") }), true)
	case "addRecord":
		a.buildAddRecordPage()
	case "logout":
		a.pages.AddAndSwitchToPage("logout", confirm.New("Logout?", a.onLogout, func() { a.SwitchTo("main") }), true)
	case "main":
		a.buildRecordsPage()
	default:
		a.showModal("Page not found", "main")
	}
}

func (a *appClient) buildAddRecordPage() {
	onSubmit := func(v recordform.RecordValues) bool {
		if err := a.logic.AddRecord(dto.ConvertToServerRecord(v)); err != nil {
			goBack := func() {
				a.pages.SwitchToPage("addRecord")
			}
			a.pages.AddAndSwitchToPage("error", modal.New(err.Error(), goBack), true)
			return false
		}
		a.showModal("Record saved", "main")
		return true
	}
	goBack := func() {
		a.SwitchTo("main")
	}
	a.pages.AddAndSwitchToPage("addRecord", recordform.New(recordform.RecordValues{}, onSubmit, goBack), true)
}

func (a *appClient) buildRecordsPage() {
	records, err := a.logic.GetRecords()

	if err != nil {
		a.showModal("Main: "+err.Error(), "main")
		return
	}

	var box tview.Primitive

	if len(records) == 0 {
		text := tview.NewTextView().SetText("Record list is empty")
		box = text
	} else {
		form := tview.NewForm()

		for _, r := range records {
			r := r
			item := listform.NewItem(r.Meta, "Show", "Edit", "Delete")

			item.
				SetViewHandler(func() { a.showRecord(&r) }).
				SetEditHandler(func() { a.editRecord(&r) }).
				SetDeleteHandler(func() { a.deleteRecord(r.ID) })

			form.AddFormItem(item)
		}
		box = form
	}

	btnNew := tview.NewButton("Add new").
		SetStyle(tcell.StyleDefault.Background(tcell.ColorLightGreen)).
		SetSelectedFunc(func() {
			a.SwitchTo("addRecord")
		})
	btnLogout := tview.NewButton("Logout").
		SetStyle(tcell.StyleDefault.Background(tcell.ColorIndianRed)).
		SetSelectedFunc(func() {
			a.SwitchTo("logout")
		})
	footer := tview.NewFlex().AddItem(btnNew, 0, 1, false).AddItem(btnLogout, 0, 1, false)

	flex := tview.NewFlex().AddItem(box, 0, 1, false).
		AddItem(footer, 3, 1, false).
		SetDirection(tview.FlexRow)

	a.pages.AddAndSwitchToPage("main", flex, true)
}

func (a *appClient) checkActivity() {
	sec := time.Since(a.lastActivity).Seconds()
	if sec > session.ActivePeriodSeconds {
		a.logic.SuspendSession()
		a.SwitchTo("checkPin")
		a.app.Draw()
	}
}

func (a *appClient) syncData() {
	err := a.logic.SyncRecords()
	if err == nil {
		p, _ := a.pages.GetFrontPage()
		if p == "main" {
			a.SwitchTo("main")
			a.app.Draw()
		}
		// a.showModal("Sync error: "+err.Error(), "main")
		// a.app.Draw()
	}
}

func (a *appClient) updateActivity() {
	a.lastActivity = time.Now()
}

func (a *appClient) showModal(text string, backRoute string) {
	goBack := func() {
		a.SwitchTo(backRoute)
	}
	a.pages.AddAndSwitchToPage("error", modal.New(text, goBack), true)
}

func (a *appClient) showRecord(r *dto.ClientRecord) {
	goBack := func() {
		a.SwitchTo("main")
	}

	backToRecord := func() {
		a.pages.SwitchToPage("viewRecord")
	}

	if r.DataType == recordform.BINARY {
		content := fmt.Sprintf("%v \n ", r.Meta)

		downloadFile := func() {
			filePath, err := a.logic.SaveBinaryFile(r)
			if err != nil {
				a.pages.AddAndSwitchToPage("error", modal.New(err.Error(), backToRecord), true)
			} else {
				a.pages.AddAndSwitchToPage("success", modal.New(filePath, backToRecord), true)
			}
		}

		a.pages.AddAndSwitchToPage("viewRecord", confirm.NewWithLabels(content, "Download", "Close", downloadFile, goBack), true)

	} else {
		content := fmt.Sprintf("%v \n ", r.Meta)

		for key, val := range r.Info {
			content += fmt.Sprintf("%v: %v\n", key, val)
		}

		a.pages.AddAndSwitchToPage("viewRecord", modal.New(content, goBack), true)
	}

}

func (a *appClient) editRecord(r *dto.ClientRecord) {
	onSubmit := func(v recordform.RecordValues) bool {
		if err := a.logic.UpdateRecord(dto.ConvertToServerRecord(v)); err != nil {
			goBack := func() {
				a.pages.SwitchToPage("editRecord")
			}
			a.pages.AddAndSwitchToPage("error", modal.New(err.Error(), goBack), true)
			return false
		}
		a.showModal("Record saved", "main")
		return true
	}
	goBack := func() {
		a.SwitchTo("main")
	}
	a.pages.AddAndSwitchToPage("editRecord", recordform.New(dto.ConvertToFormValues(*r), onSubmit, goBack), true)
}

func (a *appClient) deleteRecord(ID string) {
	goBack := func() {
		a.SwitchTo("main")
	}

	deleteRecord := func() {
		if err := a.logic.DeleteRecord(ID); err != nil {
			a.pages.AddAndSwitchToPage("error", modal.New(err.Error(), func() {
				a.deleteRecord(ID)
			}), true)
			return
		}
		goBack()
	}

	a.pages.AddAndSwitchToPage("delete-record", confirm.New("Confirm delete record?", deleteRecord, goBack), true)
}

func (a *appClient) onCheckPin(p pinform.PinForm) bool {
	valid, err := a.logic.CheckPin(p.Pin)
	if err != nil {
		a.showModal(err.Error(), "checkPin")
		return false
	}

	if !valid {
		a.showModal("Invalid PIN", "checkPin")
		return false
	}

	a.syncData()
	a.SwitchTo("main")

	return true
}

func (a *appClient) onSetPin(p pinform.PinForm) bool {
	if err := a.logic.SetPin(p.Pin); err != nil {
		a.showModal(err.Error(), "main")
		return false
	}
	a.syncData()
	a.SwitchTo("main")
	return true
}

func (a *appClient) onLogin(p authform.AuthForm) bool {
	err := a.logic.Login(p.Login, p.Password)

	if err != nil {
		a.showModal(err.Error(), "login")
		return false
	}

	a.SwitchTo("main")
	return true
}

func (a *appClient) onLogout() {
	err := a.logic.Logout()
	if err != nil {
		a.showModal(err.Error(), "login")
		return
	}
	a.SwitchTo("login")
}

func (a *appClient) onRegister(p authform.AuthForm) bool {
	err := a.logic.Register(p.Login, p.Password)

	if err != nil {
		a.showModal(err.Error(), "register")
		return false
	}

	a.SwitchTo("main")
	return true
}

// Run start process with TUI app
func (a *appClient) Run(ctx context.Context) error {

	app := a.app.SetRoot(a.pages, true).
		EnableMouse(true).
		EnablePaste(true).
		SetAfterDrawFunc(func(screen tcell.Screen) {
			a.updateActivity()
		})

	go func() {
		<-ctx.Done()
		app.Stop()
	}()

	activityCheck := time.NewTicker(time.Second * 10)

	syncCheck := time.NewTicker(time.Second * 15)

	go func() {
		for {
			select {
			case <-ctx.Done():
				app.Stop()
				activityCheck.Stop()
				return
			case <-activityCheck.C:
				a.checkActivity()
			case <-syncCheck.C:
				a.syncData()
			}

		}
	}()

	return app.Run()
}
