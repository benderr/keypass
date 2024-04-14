package recordform

import (
	"os"

	"github.com/rivo/tview"
)

type RecordValues struct {
	ID       string
	Meta     string
	DataType string
	Login    string
	Password string
	Text     string
	Number   string
	CVV      string
	Expire   string
	FilePath string
	Binary   []byte
}

var (
	dataTypes         = []string{CREDENTIALS, TEXT, CREDIT, BINARY}
	dynamicFormFields = []string{"Login", "Password", "Text", "Number", "CVV", "Expire", "Filepath", "Info"}
)

func getSelectedIndex(dataType string) int {
	if len(dataType) == 0 {
		return 0
	}
	switch dataType {
	case CREDENTIALS:
		return 0
	case TEXT:
		return 1
	case CREDIT:
		return 2
	case BINARY:
		return 3
	default:
		return 0
	}
}

type recordFormBuilder struct {
	*tview.Form
	record *RecordValues
}

func New(initialValues RecordValues, onSubmit func(values RecordValues) bool, onCancel func()) tview.Primitive {

	record := &initialValues

	inst := &recordFormBuilder{
		Form:   tview.NewForm(),
		record: record,
	}

	inst.init()

	inst.AddButton("Save", func() {
		inst.showErrorMessage("")

		if record.DataType == BINARY {
			if len(record.FilePath) == 0 {
				inst.showErrorMessage("Enter file path")
				return
			}

			file, err := os.ReadFile(record.FilePath)
			if err != nil {
				inst.showErrorMessage(err.Error())
				return
			}

			record.Binary = file
		}

		onSubmit(*record)
	})

	inst.AddButton("Cancel", onCancel)

	return inst
}

func (r *recordFormBuilder) init() {
	r.AddInputField("Meta info", r.record.Meta, 20, nil, func(meta string) {
		r.record.Meta = meta
	})

	initialDataType := getSelectedIndex(r.record.DataType)

	r.AddDropDown("Type", dataTypes, initialDataType, func(dataType string, index int) {
		r.record.DataType = dataType
		r.rebuildByType(dataType)
	})

	r.rebuildByType(dataTypes[initialDataType])
}

func (r *recordFormBuilder) clearDynamicFields() {
	for _, f := range dynamicFormFields {
		i := r.GetFormItemIndex(f)
		if i > -1 {
			r.RemoveFormItem(i)
		}
	}
}

func (r *recordFormBuilder) showErrorMessage(message string) {
	item := r.Form.GetFormItemByLabel("Info")
	if item == nil {
		return
	}
	if v, ok := item.(*tview.TextView); ok {
		v.SetText(message)
	}
}

func (r *recordFormBuilder) rebuildByType(dataTypeLabel string) {
	r.clearDynamicFields()

	switch dataTypeLabel {
	case CREDENTIALS:
		r.AddInputField("Login", r.record.Login, 20, nil, func(login string) {
			r.record.Login = login
		})

		r.AddInputField("Password", r.record.Password, 20, nil, func(pass string) {
			r.record.Password = pass
		})
	case TEXT:
		r.AddInputField("Text", r.record.Text, 20, nil, func(text string) {
			r.record.Text = text
		})

	case CREDIT:
		r.AddInputField("Number", r.record.Number, 20, nil, func(text string) {
			r.record.Number = text
		})
		r.AddInputField("CVV", r.record.CVV, 20, nil, func(text string) {
			r.record.CVV = text
		})
		r.AddInputField("Expire", r.record.Expire, 20, nil, func(text string) {
			r.record.Expire = text
		})

	case BINARY:
		r.AddInputField("Filepath", "", 20, nil, func(text string) {
			r.record.FilePath = text

		})
	}

	r.AddTextView("Info", "", 20, 4, true, true)

}
