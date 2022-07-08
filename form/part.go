package form

import (
	"github.com/dtylman/gowd"
	"github.com/dtylman/gowd/bootstrap"
	jww "github.com/spf13/jwalterweatherman"
)

// ValidateFunc is the function used to validate an input value when a form is
// submitted. If the validation fails, an error describing why is returned. If
// validation succeeds, nil is returned.
type ValidateFunc func(str string) error

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
func (p *Part) Validate() bool {
	err := p.v(p.f.GetValue())
	if err != nil {
		jww.ERROR.Printf("Failed to validate input %q: %+v", p.caption, err)
		p.SetHelpText(err.Error())
		return false
	}

	p.ClearHelpText()

	return true
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
	p.f.Kids[1].SetAttribute(key, val)
}

func (p *Part) RemoveAttribute(key string) {
	p.f.Kids[1].RemoveAttribute(key)
}

func (p *Part) GetAttribute(key string) (string, bool) {
	return p.f.Kids[1].GetAttribute(key)
}
