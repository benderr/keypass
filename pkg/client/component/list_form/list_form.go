package listform

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type RowButton struct {
	btnView   *tview.Button
	btnEdit   *tview.Button
	btnDelete *tview.Button
	*tview.Flex
	label string
}

func NewItem(label string, labelShow, labelEdit, labelDelete string) *RowButton {
	text := tview.NewTextView().SetText(label)

	btnView := tview.NewButton(labelShow)
	btnEdit := tview.NewButton(labelEdit).SetStyle(tcell.StyleDefault.Background(tcell.ColorDarkOrange))
	btnDelete := tview.NewButton(labelDelete).SetStyle(tcell.StyleDefault.Background(tcell.ColorDarkRed))
	flex := tview.NewFlex().AddItem(text, 0, 2, false).
		AddItem(btnView, 0, 1, false).
		AddItem(btnEdit, 0, 1, false).
		AddItem(btnDelete, 0, 1, false)
	flex.SetTitle(label)

	r := &RowButton{
		btnView:   btnView,
		btnEdit:   btnEdit,
		btnDelete: btnDelete,
		Flex:      flex,
	}
	r.SetTitle(label)
	return r
}

func (r *RowButton) GetLabel() string {

	return r.label
}

func (r *RowButton) GetFieldHeight() int {
	return 1
}

func (r *RowButton) GetFieldWidth() int {
	return tview.TaggedStringWidth(r.GetLabel()) + 4
}

func (r *RowButton) SetViewHandler(handler func()) *RowButton {
	r.btnView.SetSelectedFunc(handler)
	return r
}

func (r *RowButton) SetEditHandler(handler func()) *RowButton {
	r.btnEdit.SetSelectedFunc(handler)
	return r
}

func (r *RowButton) SetDeleteHandler(handler func()) *RowButton {
	r.btnDelete.SetSelectedFunc(handler)
	return r
}

func (r *RowButton) SetDisabled(disabled bool) tview.FormItem {
	r.btnView.SetDisabled(disabled)
	r.btnEdit.SetDisabled(disabled)
	r.btnDelete.SetDisabled(disabled)
	return r
}

func (r *RowButton) SetFinishedFunc(handler func(key tcell.Key)) tview.FormItem {
	return r
}

func (r *RowButton) SetFormAttributes(labelWidth int, labelColor, bgColor, fieldTextColor, fieldBgColor tcell.Color) tview.FormItem {
	r.SetBackgroundColor(bgColor)
	return r
}
