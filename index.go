package main

import (
	"git.xx.network/elixxir/mainnet-commitments/client"
	"git.xx.network/elixxir/mainnet-commitments/utils"
	"github.com/dtylman/gowd"
	"github.com/dtylman/gowd/bootstrap"
	jww "github.com/spf13/jwalterweatherman"
	utils2 "gitlab.com/xx_network/primitives/utils"
)

var body *gowd.Element

const blurbText = "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vivamus malesuada eleifend ultrices. Pellentesque habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas. Etiam pretium tempor massa, a volutpat orci mattis non. Integer tincidunt tincidunt ante sed cursus. In lacinia pulvinar tempor. Nullam id luctus nibh, vitae iaculis ante. Integer vel sem at augue viverra suscipit vel nec orci. Sed sed ultrices quam. Vestibulum hendrerit tellus justo, id tempus enim fringilla quis. Vivamus nunc ante, tincidunt et tellus eget, varius lobortis tortor. Maecenas in porta erat. Nam in dolor turpis. Aliquam a tristique arcu, vitae pulvinar tellus. Sed nec imperdiet mi, vel maximus est. Nunc aliquet eros arcu, quis faucibus nibh feugiat a. Fusce id orci nunc.\n\n"
const serverAddress = "http://localhost:11420"

type Inputs struct {
	keyPath         string
	idfPath         string
	nominatorWallet string
	validatorWallet string
	serverCert      string
	serverCertPath  string
	agree           bool
}

func buildPage() error {

	inputs := Inputs{}

	// creates a new bootstrap fluid container
	body = bootstrap.NewContainer(false)
	// add some elements using the object model

	// keyPathInput := bootstrap.NewFileButton(bootstrap.ButtonDefault, "keyPath", false)
	keyPathInput := NewFileButton("Server Key (.key)", &inputs.keyPath)
	idfPathInput := NewFileButton("IDF (.json)", &inputs.idfPath)
	nominatorWalletInput := bootstrap.NewFormInput("text", "Nominator Wallet Address")
	validatorWalletInput := bootstrap.NewFormInput("text", "Validator Wallet Address")
	serverCertPathInput := NewFileButton("BetaNet Server Certificate (.crt)", &inputs.serverCertPath)

	agreeInput := bootstrap.NewCheckBox("I agree to the contract above.", false)
	agreeHelpText := bootstrap.NewElement("p", "help-block")
	agreeHelpText.Hidden = true
	agreeBox := bootstrap.NewElement("div", "form-group", agreeInput.Element, agreeHelpText)

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
		submit.Disable()
		body.Disable()
		defer func() {
			body.Enable()
			submit.Enable()
		}()
		var errs int
		if len(inputs.keyPath) == 0 {
			keyPathInput.SetHelpText("Required.")
			errs++
		} else {
			keyPathInput.HideHelpText()

			_, err := utils2.ReadFile(inputs.keyPath)
			if err != nil {
				jww.ERROR.Printf("keyPath error: %+v", err)
				keyPathInput.SetHelpText(err.Error())
				errs++
			}
		}
		if len(inputs.idfPath) == 0 {
			idfPathInput.SetHelpText("Required.")
			errs++
		} else {
			idfPathInput.HideHelpText()
			_, err := utils2.ReadFile(inputs.idfPath)
			if err != nil {
				jww.ERROR.Printf("idfPath error: %+v", err)
				idfPathInput.SetHelpText(err.Error())
				errs++
			}
		}
		inputs.nominatorWallet = nominatorWalletInput.GetValue()
		if len(inputs.nominatorWallet) == 0 {
			nominatorWalletInput.SetHelpText("Required.")
			errs++
		} else {
			if len(nominatorWalletInput.Kids) > 2 {
				nominatorWalletInput.Kids[2].Hidden = true
			}
		}
		inputs.validatorWallet = validatorWalletInput.GetValue()
		if len(inputs.validatorWallet) == 0 {
			validatorWalletInput.SetHelpText("Required.")
			errs++
		} else {
			if len(validatorWalletInput.Kids) > 2 {
				validatorWalletInput.Kids[2].Hidden = true
			}
		}
		if len(inputs.serverCertPath) == 0 {
			serverCertPathInput.SetHelpText("Required.")
			errs++
		} else {
			serverCertPathInput.HideHelpText()

			data, err := utils2.ReadFile(inputs.serverCertPath)
			if err != nil {
				jww.ERROR.Printf("serverCertPath error: %+v", err)
				serverCertPathInput.SetHelpText(err.Error())
				errs++
			} else {
				inputs.serverCert = string(data)
			}
		}
		inputs.agree = agreeInput.Checked()
		if inputs.agree == false {
			agreeHelpText.SetText("Required.")
			agreeHelpText.Hidden = false
			errs++
		} else {
			agreeHelpText.Hidden = true
		}
		jww.INFO.Printf("Inputs set: %+v", inputs)

		if errs >= 0 {
			spinner.Hidden = false
			formErrors.Hidden = true
			errBox.Hidden = true

			if err := body.Render(); err != nil {
				jww.ERROR.Print(err)
			}

			err := client.SignAndTransmit(
				inputs.keyPath,
				inputs.idfPath,
				inputs.nominatorWallet,
				inputs.validatorWallet,
				serverAddress,
				inputs.serverCert,
				utils.Contract)

			spinner.Hidden = true

			if err != nil {
				jww.ERROR.Printf("Submit error: %+v", err)
				errBox.SetText(err.Error())
				errBox.Hidden = false
				formErrors.SetText("The were errors in the form input. Please correct them to continue.")
				formErrors.Hidden = false
			} else {
				divWell.RemoveElements()
				success := bootstrap.NewElement("span", "success", gowd.NewText("MainNet Commitments Successful."))
				divWell.AddElement(success)
			}
		} else {
			formErrors.SetText("The were errors in the form input. Please correct them to continue.")
			formErrors.Hidden = false
		}
	})

	contractText := bootstrap.NewElement("p", "contractText", gowd.NewText("Read through the entire contract below and accept the terms."))
	contract := bootstrap.NewElement("div", "contractContainer")
	_, err := contract.AddHTML(utils.Contract, nil)
	if err != nil {
		return err
	}

	form := bootstrap.NewFormGroup(
		formErrors,
		serverCertPathInput.Element,
		keyPathInput.Element,
		idfPathInput.Element,
		nominatorWalletInput.Element,
		validatorWalletInput.Element,
		contractText,
		contract,
		agreeBox,
		submitBox,
	)

	form.SetAttribute("style", "margin-top:35px")

	h1 := bootstrap.NewElement("h1", "")
	h1.SetText("xx network MainNet Commitments")
	logo := bootstrap.NewElement("img", "logo")
	logo.SetAttribute("src", "img/xx-logo.svg")
	h1.AddElement(logo)
	p := bootstrap.NewElement("p", "blurb")
	p.SetText(blurbText)
	divWell.AddElement(h1)
	divWell.AddElement(p)
	divWell.AddElement(form)
	row := bootstrap.NewRow(divWell)
	body.AddElement(row)

	// Start the UI loop
	err = gowd.Run(body)
	if err != nil {
		return err
	}

	return nil
}
