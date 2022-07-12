package form

import (
	"github.com/dtylman/gowd"
	"github.com/dtylman/gowd/bootstrap"
	jww "github.com/spf13/jwalterweatherman"
)

// ValidateFunc is the function used to validate an input value when a form is
// submitted. If the validation fails, an error describing why is returned. If
// validation succeeds, nil is returned.
type ValidateFunc func(str string) (interface{}, string, error)

type Part struct {
	caption string
	f       *bootstrap.FormInput
	v       ValidateFunc
}

func NewPart(inputType string, caption string, v ValidateFunc) *Part {
	return &Part{
		caption: caption,
		f:       bootstrap.NewFormInput(inputType, caption),
		v:       v,
	}
}

func (p *Part) SetHelpText(help string) {
	p.f.SetHelpText(help)
}

func (p *Part) ClearHelpText() {
	p.f.Kids[2].Hidden = true
}

// Validate checks the value of the form input against the validator function.
// If validation fails, the error is set as the help text and returns true.
// If validations succeeds, it returns true.
func (p *Part) Validate() (interface{}, bool) {
	var validated interface{}
	var err error
	var helpText string
	if t, _ := p.GetInputAttribute("type"); t == "checkbox" {
		val := ""
		if p.Checked() {
			val = "true"
		}
		validated, helpText, err = p.v(val)
	} else {
		validated, helpText, err = p.v(p.f.GetValue())
	}

	if err != nil {
		jww.ERROR.Printf("Failed to validate input %q: %+v", p.caption, err)
		p.SetHelpText(helpText)
		return nil, false
	}

	p.ClearHelpText()

	return validated, true
}

func (p *Part) GetValue() string {
	return p.f.GetValue()
}

func (p *Part) SetValue(value string) {
	p.f.SetValue(value)
}

func (p *Part) Disable() {
	p.f.Kids[1].Disable()
}

func (p *Part) Enable() {
	p.f.Kids[1].Enable()
}

func (p *Part) Hide() {
	p.f.Hide()
}

func (p *Part) Show() {
	p.f.Show()
}

func (p *Part) Check() {
	p.f.Kids[1].SetAttribute("checked", "")
}

func (p *Part) Uncheck() {
	p.f.Kids[1].RemoveAttribute("checked")
}

func (p *Part) Checked() bool {
	_, exists := p.f.Kids[1].GetAttribute("checked")
	return exists
}

func (p *Part) Element() *gowd.Element {
	return p.f.Element
}

func (p *Part) OnEvent(event string, handler gowd.EventHandler) {
	p.f.OnEvent(event, handler)
}

func (p *Part) SetAttribute(key, val string) {
	p.f.SetAttribute(key, val)
}

func (p *Part) RemoveAttribute(key string) {
	p.f.RemoveAttribute(key)
}

func (p *Part) GetAttribute(key string) (string, bool) {
	return p.f.GetAttribute(key)
}

func (p *Part) SetInputAttribute(key, val string) {
	p.f.Kids[1].SetAttribute(key, val)
}

func (p *Part) RemoveInputAttribute(key string) {
	p.f.Kids[1].RemoveAttribute(key)
}

func (p *Part) GetInputAttribute(key string) (string, bool) {
	return p.f.Kids[1].GetAttribute(key)
}

func (p *Part) SetLabelAttribute(key, val string) {
	p.f.Kids[0].SetAttribute(key, val)
}

func (p *Part) RemoveLabelAttribute(key string) {
	p.f.Kids[0].RemoveAttribute(key)
}

func (p *Part) GetLabelAttribute(key string) (string, bool) {
	return p.f.Kids[0].GetAttribute(key)
}

func (p *Part) SetHelpTxtAttribute(key, val string) {
	p.f.Kids[2].SetAttribute(key, val)
}

func (p *Part) RemoveHelpTxtAttribute(key string) {
	p.f.Kids[2].RemoveAttribute(key)
}

func (p *Part) GetHelpTxtAttribute(key string) (string, bool) {
	return p.f.Kids[2].GetAttribute(key)
}

func (p *Part) AddElement(elem *gowd.Element) {
	p.f.AddElement(elem)
}

func (p *Part) SwapKids(i, j int) {
	p.f.Kids[i], p.f.Kids[j] = p.f.Kids[j], p.f.Kids[i]
}
