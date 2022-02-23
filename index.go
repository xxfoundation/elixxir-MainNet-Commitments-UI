package main

import (
	"fmt"
	"git.xx.network/elixxir/mainnet-commitments/client"
	"git.xx.network/elixxir/mainnet-commitments/utils"
	"github.com/dtylman/gowd"
	"github.com/dtylman/gowd/bootstrap"
	jww "github.com/spf13/jwalterweatherman"
	utils2 "gitlab.com/xx_network/primitives/utils"
	"time"
)

var body *gowd.Element

const blurbText = "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vivamus malesuada eleifend ultrices. Pellentesque habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas. Etiam pretium tempor massa, a volutpat orci mattis non. Integer tincidunt tincidunt ante sed cursus. In lacinia pulvinar tempor. Nullam id luctus nibh, vitae iaculis ante. Integer vel sem at augue viverra suscipit vel nec orci. Sed sed ultrices quam. Vestibulum hendrerit tellus justo, id tempus enim fringilla quis. Vivamus nunc ante, tincidunt et tellus eget, varius lobortis tortor. Maecenas in porta erat. Nam in dolor turpis. Aliquam a tristique arcu, vitae pulvinar tellus. Sed nec imperdiet mi, vel maximus est. Nunc aliquet eros arcu, quis faucibus nibh feugiat a. Fusce id orci nunc.\n\n"

func buildPage() error {

	inputs := Inputs{}

	// creates a new bootstrap fluid container
	body = bootstrap.NewContainer(false)
	// add some elements using the object model

	// keyPathInput := bootstrap.NewFileButton(bootstrap.ButtonDefault, "keyPath", false)
	keyPathInput := NewFileButton("keyPath", &inputs.keyPath)
	idfPathInput := NewFileButton("idfPath", &inputs.idfPath)
	nominatorWalletInput := bootstrap.NewFormInput("text", "nominatorWallet")
	validatorWalletInput := bootstrap.NewFormInput("text", "validatorWallet")
	serverAddressInput := bootstrap.NewFormInput("text", "serverAddress")
	serverCertPathInput := NewFileButton("serverCertPath", &inputs.serverCertPath)

	agreeInput := bootstrap.NewCheckBox("agree", false)
	agreeHelpText := bootstrap.NewElement("p", "help-block")
	agreeHelpText.Hidden = true
	agreeBox := bootstrap.NewElement("div", "form-group", agreeInput.Element, agreeHelpText)

	submit := bootstrap.NewButton(bootstrap.ButtonPrimary, "Submit")
	errBox := bootstrap.NewElement("span", "errorBox")
	errBox.Hidden = true
	submitBox := bootstrap.NewElement("div", "", submit, errBox)
	submitBox.SetAttribute("style", "text-align:center;")

	formErrors := bootstrap.NewElement("p", "formErrors")
	formErrors.Hidden = true

	divWell := bootstrap.NewElement("div", "well")

	submit.OnEvent(gowd.OnClick, func(_ *gowd.Element, event *gowd.EventElement) {
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
		inputs.serverAddress = serverAddressInput.GetValue()
		if len(inputs.serverAddress) == 0 {
			serverAddressInput.SetHelpText("Required.")
			errs++
		} else {
			if len(serverAddressInput.Kids) > 2 {
				serverAddressInput.Kids[2].Hidden = true
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

		if errs == 0 {
			formErrors.Hidden = true
			err := client.SignAndTransmit(
				inputs.keyPath,
				inputs.idfPath,
				inputs.nominatorWallet,
				inputs.validatorWallet,
				inputs.serverAddress,
				inputs.serverCert,
				utils.Contract)
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
		keyPathInput.Element,
		idfPathInput.Element,
		nominatorWalletInput.Element,
		validatorWalletInput.Element,
		serverAddressInput.Element,
		serverCertPathInput.Element,
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

type Inputs struct {
	keyPath, idfPath, nominatorWallet, validatorWallet, serverAddress, serverCert, serverCertPath string
	agree                                                                                         bool
}

// happens when the 'start' button is clicked
func btnClicked(sender *gowd.Element, event *gowd.EventElement) {
	// adds a text and progress bar to the body
	sender.SetText("Working...")
	text := body.AddElement(gowd.NewStyledText("Working...", gowd.BoldText))
	progressBar := bootstrap.NewProgressBar()
	body.AddElement(progressBar.Element)

	// makes the body stop responding to user events
	body.Disable()

	// clean up - remove the added elements
	defer func() {
		sender.SetText("Start")
		body.RemoveElement(text)
		body.RemoveElement(progressBar.Element)
		body.Enable()
	}()

	// render the progress bar
	for i := 0; i <= 123; i++ {
		progressBar.SetValue(i, 123)
		text.SetText(fmt.Sprintf("Working %v", i))
		time.Sleep(time.Millisecond * 20)
		// this will cause the body to be refreshed
		body.Render()
	}

}
