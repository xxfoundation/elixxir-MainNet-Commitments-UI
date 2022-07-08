package main

import (
	"fmt"
	"git.xx.network/elixxir/mainnet-commitments-ui/form"
	"github.com/dtylman/gowd"
	"github.com/dtylman/gowd/bootstrap"
	jww "github.com/spf13/jwalterweatherman"
	"strconv"
)

var body *gowd.Element

const (
	blurbTextPg1 = `This applet will allow you to commit your wallets. For more information, please see the&nbsp;`
	blurbTextPg2 = `Below are the committed validator and nominator addresses. Select the checkbox to modify them.`
)
const serverAddress = "https://18.185.229.39:11420"

type Inputs struct {
	nodeID          string
	nominatorWallet string
	validatorWallet string
	multiplier      float32
	maxMultiplier   float32
	idfPath         string
	agree           bool
}

func buildPage() error {

	// creates a new bootstrap fluid container
	body = bootstrap.NewContainer(false)
	// add some elements using the object model

	// keyPathInput := bootstrap.NewFileButton(bootstrap.ButtonDefault, "keyPath", false)

	row := page1()

	body.AddElement(row)

	// Start the UI loop
	err := gowd.Run(body)
	if err != nil {
		return err
	}

	return nil
}

func page1() *gowd.Element {

	inputs := Inputs{}

	nodeID := form.NewPart("text", "Node ID", form.ValidateNodeID)
	idfPathInput := form.NewFileButton("Node IDF (.json)", &inputs.idfPath, form.ValidateFilePath)

	submit := bootstrap.NewButton(bootstrap.ButtonPrimary, "Submit")
	errBox := bootstrap.NewElement("span", "errorBox")
	errBox.Hidden = true
	spinner := bootstrap.NewElement("div", "spinnerContainer", bootstrap.NewElement("div", "spinner", gowd.NewText("Loading...")))
	spinner.Hidden = true
	submitBox := bootstrap.NewElement("div", "", errBox, spinner, submit)
	submitBox.SetAttribute("style", "text-align:center;")

	formErrors := bootstrap.NewElement("p", "formErrors")
	formErrors.Hidden = true

	divWell := bootstrap.NewElement("div", "well")

	submit.OnEvent(gowd.OnClick, func(_ *gowd.Element, event *gowd.EventElement) {
		jww.DEBUG.Printf("sumbit")
		submit.Disable()
		body.Disable()
		spinner.Hidden = false
		formErrors.Hidden = true
		errBox.Hidden = true
		defer func() {
			body.Enable()
			submit.Enable()
			spinner.Hidden = true
		}()
		var errs int

		if nodeID.Validate() {
			inputs.nodeID = nodeID.GetValue()
		} else {
			errs++
		}

		if idfPathInput.Validate() {
			inputs.idfPath = nodeID.GetValue()
		} else {
			errs++
		}

		jww.INFO.Printf("Inputs set: %+v", inputs)

		if errs == 0 {
			if err := body.Render(); err != nil {
				jww.ERROR.Print(err)
			}

			spinner.Hidden = true

			submitNodeID := func(nodeID string) (
				validatorWallet, nominatorWallet string,
				selectedMultiplier, maxMultiplier float32, err error) {
				return "validatorWallet", "nominatorWallet", 1.3, 3.5, nil
			}

			var err error
			inputs.validatorWallet, inputs.nominatorWallet, inputs.multiplier,
				inputs.maxMultiplier, err = submitNodeID(inputs.nodeID)

			if err != nil {
				jww.ERROR.Printf("Submit error: %+v", err)
				errBox.SetText("An error occurred when submitting the request. Please contact support at nodes@xx.network and provide the following error message:")
				errBox.AddElement(bootstrap.NewElement("span", "errorBoxMessage", gowd.NewText(err.Error())))
				errBox.Hidden = false
				formErrors.SetText("The were errors in the form input. Please correct them to continue.")
				formErrors.Hidden = false
			} else {
				body.RemoveElements()
				body.AddElement(page2(inputs))
			}
		} else {
			formErrors.SetText("There were errors in the form input. Please correct them to continue.")
			formErrors.Hidden = false
		}
	})

	formGrp := bootstrap.NewFormGroup(
		formErrors,
		nodeID.Element(),
		idfPathInput.Element,
		submitBox,
	)

	formGrp.SetAttribute("style", "margin-top:35px")

	h1 := bootstrap.NewElement("h1", "")
	h1.SetText("xx network MainNet Wallet Commitment")
	logo := bootstrap.NewElement("img", "logo")
	logo.SetAttribute("src", "img/xx-logo.svg")
	h1.AddElement(logo)
	p := bootstrap.NewElement("p", "blurb")
	p.AddHTML(blurbTextPg1, nil)
	instructionPageLink := bootstrap.NewLinkButton("instructions page")
	instructionPageLink.RemoveAttribute("href")
	instructionPageLink.SetAttribute("style", "cursor:pointer;")
	instructionPageLink.OnEvent(gowd.OnClick, func(*gowd.Element, *gowd.EventElement) {
		gowd.ExecJSNow("window.nw.Shell.openExternal('https://xx.network/mainnet-commit-wallet/')")
	})
	p.AddElement(instructionPageLink)
	p.AddElement(gowd.NewText("."))
	divWell.AddElement(h1)
	divWell.AddElement(p)
	divWell.AddElement(formGrp)
	row := bootstrap.NewRow(divWell)

	return row
}

func page2(inputs Inputs) *gowd.Element {

	validatorWallet := form.NewPart("text", "Validator Wallet Address", form.ValidateXXNetworkAddress)
	validatorWallet.SetValue(inputs.validatorWallet)
	validatorWallet.Disable()
	nominatorWallet := form.NewPart("text", "Nominator Wallet Address", form.ValidateXXNetworkAddress)
	nominatorWallet.SetValue(inputs.nominatorWallet)
	nominatorWallet.Disable()

	modifyCheck := form.NewPart("checkbox", "Modify Wallet Addresses", nil)
	modifyCheck.OnEvent(gowd.OnClick, func(sender *gowd.Element, event *gowd.EventElement) {
		if modifyCheck.Checked() {
			validatorWallet.Enable()
			nominatorWallet.Enable()
		} else {
			validatorWallet.Disable()
			nominatorWallet.Disable()
		}
	})

	submit := bootstrap.NewButton(bootstrap.ButtonPrimary, "Submit")
	errBox := bootstrap.NewElement("span", "errorBox")
	errBox.Hidden = true
	spinner := bootstrap.NewElement("div", "spinnerContainer", bootstrap.NewElement("div", "spinner", gowd.NewText("Loading...")))
	spinner.Hidden = true
	submitBox := bootstrap.NewElement("div", "", errBox, spinner, submit)
	submitBox.SetAttribute("style", "text-align:center;")

	formErrors := bootstrap.NewElement("p", "formErrors")
	formErrors.Hidden = true

	divWell := bootstrap.NewElement("div", "well")

	submit.OnEvent(gowd.OnClick, func(_ *gowd.Element, event *gowd.EventElement) {
		jww.DEBUG.Printf("sumbit")
		submit.Disable()
		body.Disable()
		spinner.Hidden = false
		formErrors.Hidden = true
		errBox.Hidden = true
		defer func() {
			body.Enable()
			submit.Enable()
			spinner.Hidden = true
		}()
		var errs int

		if validatorWallet.Validate() {
			inputs.validatorWallet = validatorWallet.GetValue()
		} else {
			errs++
		}

		if validatorWallet.Validate() {
			inputs.nominatorWallet = nominatorWallet.GetValue()
		} else {
			errs++
		}
		jww.INFO.Printf("Inputs set: %+v", inputs)

		if errs == 0 {
			if err := body.Render(); err != nil {
				jww.ERROR.Print(err)
			}

			spinner.Hidden = true

			body.RemoveElements()
			body.AddElement(page3(inputs))
		} else {
			formErrors.SetText("There were errors in the form input. Please correct them to continue.")
			formErrors.Hidden = false
		}
	})

	formGrp := bootstrap.NewFormGroup(
		formErrors,
		modifyCheck.Element(),
		validatorWallet.Element(),
		nominatorWallet.Element(),
		submitBox,
	)

	formGrp.SetAttribute("style", "margin-top:35px")

	h1 := bootstrap.NewElement("h1", "")
	h1.SetText("xx network MainNet Wallet Commitment")
	logo := bootstrap.NewElement("img", "logo")
	logo.SetAttribute("src", "img/xx-logo.svg")
	h1.AddElement(logo)
	p := bootstrap.NewElement("p", "blurb")
	p.AddHTML(blurbTextPg2, nil)
	p.AddElement(gowd.NewText("."))
	divWell.AddElement(h1)
	divWell.AddElement(p)
	divWell.AddElement(formGrp)
	row := bootstrap.NewRow(divWell)

	return row
}

func page3(inputs Inputs) *gowd.Element {

	multiplier := form.NewPart("number", "Selected Multiplier", form.ValidateMultiplier(inputs.maxMultiplier))
	multiplier.SetValue(strconv.FormatFloat(float64(inputs.multiplier), 'f', 3, 32))
	multiplier.SetAttribute("min", strconv.FormatFloat(0.0, 'f', 3, 32))
	multiplier.SetAttribute("max", strconv.FormatFloat(float64(inputs.maxMultiplier), 'f', 3, 32))
	multiplier.SetAttribute("step", ".001")
	multiplier.SetAttribute("pattern", `"^\d*(\.\d{0,2})?$"`)
	multiplier.Disable()

	modifyCheck := form.NewPart("checkbox", "Modify the selected multiplier", nil)
	modifyCheck.OnEvent(gowd.OnClick, func(sender *gowd.Element, event *gowd.EventElement) {
		if modifyCheck.Checked() {
			multiplier.Enable()
		} else {
			multiplier.Disable()
		}
	})

	submit := bootstrap.NewButton(bootstrap.ButtonPrimary, "Submit")
	errBox := bootstrap.NewElement("span", "errorBox")
	errBox.Hidden = true
	spinner := bootstrap.NewElement("div", "spinnerContainer", bootstrap.NewElement("div", "spinner", gowd.NewText("Loading...")))
	spinner.Hidden = true
	submitBox := bootstrap.NewElement("div", "", errBox, spinner, submit)
	submitBox.SetAttribute("style", "text-align:center;")

	formErrors := bootstrap.NewElement("p", "formErrors")
	formErrors.Hidden = true

	divWell := bootstrap.NewElement("div", "well")

	submit.OnEvent(gowd.OnClick, func(_ *gowd.Element, event *gowd.EventElement) {
		jww.DEBUG.Printf("sumbit")
		submit.Disable()
		body.Disable()
		spinner.Hidden = false
		formErrors.Hidden = true
		errBox.Hidden = true
		defer func() {
			body.Enable()
			submit.Enable()
			spinner.Hidden = true
		}()
		var errs int

		if multiplier.Validate() {
			inputs.multiplier = getFloat(multiplier.GetValue())
		} else {
			errs++
		}

		jww.INFO.Printf("Inputs set: %+v", inputs)

		if errs == 0 {
			if err := body.Render(); err != nil {
				jww.ERROR.Print(err)
			}

			spinner.Hidden = true

			submitInputs := func(nodeID, validatorWallet, nominatorWallet string,
				selectedMultiplier float32) error {
				return nil
			}
			err := submitInputs(inputs.nodeID, inputs.validatorWallet, inputs.nominatorWallet, inputs.multiplier)

			if err != nil {
				jww.ERROR.Printf("Submit error: %+v", err)
				errBox.SetText("An error occurred when submitting the request. Please contact support at nodes@xx.network and provide the following error message:")
				errBox.AddElement(bootstrap.NewElement("span", "errorBoxMessage", gowd.NewText(err.Error())))
				errBox.Hidden = false
				formErrors.SetText("The were errors in the form input. Please correct them to continue.")
				formErrors.Hidden = false
			} else {
				divWell.RemoveElements()
				success := bootstrap.NewElement("span", "success", gowd.NewText("MainNet Commitments Successful."+fmt.Sprintf("%+v", inputs)))
				divWell.AddElement(success)
			}
		} else {
			formErrors.SetText("There were errors in the form input. Please correct them to continue.")
			formErrors.Hidden = false
		}
	})

	formGrp := bootstrap.NewFormGroup(
		formErrors,
		modifyCheck.Element(),
		multiplier.Element(),
		submitBox,
	)

	formGrp.SetAttribute("style", "margin-top:35px")

	h1 := bootstrap.NewElement("h1", "")
	h1.SetText("xx network MainNet Wallet Commitment")
	logo := bootstrap.NewElement("img", "logo")
	logo.SetAttribute("src", "img/xx-logo.svg")
	h1.AddElement(logo)
	p := bootstrap.NewElement("p", "blurb")
	p.AddHTML(blurbTextPg2, nil)
	p.AddElement(gowd.NewText("."))
	divWell.AddElement(h1)
	divWell.AddElement(p)
	divWell.AddElement(formGrp)
	row := bootstrap.NewRow(divWell)

	return row
}

func getFloat(str string) float32 {
	f, err := strconv.ParseFloat(str, 32)
	if err != nil {
		jww.FATAL.Panicf("Failed to parse string as float %q: %+v", str, err)
	}

	return float32(f)
}
