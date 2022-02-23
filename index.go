package main

import (
	"fmt"
	"git.xx.network/elixxir/mainnet-commitments/client"
	"git.xx.network/elixxir/mainnet-commitments/utils"
	"github.com/dtylman/gowd"
	"github.com/dtylman/gowd/bootstrap"
	jww "github.com/spf13/jwalterweatherman"
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
	keyPathInput := NewFormInput("keyPath", &inputs.keyPath)
	idfPathInput := NewFormInput("idfPath", &inputs.idfPath)
	nominatorWalletInput := bootstrap.NewFormInput("text", "nominatorWallet")
	validatorWalletInput := bootstrap.NewFormInput("text", "validatorWallet")
	serverAddressInput := bootstrap.NewFormInput("text", "serverAddress")
	serverCertPathInput := NewFormInput("serverCertPath", &inputs.serverCertPath)
	agreeInput := bootstrap.NewCheckBox("Agree", false)
	submit := bootstrap.NewButton(bootstrap.ButtonPrimary, "Submit")

	errBox := bootstrap.NewElement("span", "errorBox")
	errBox.Hidden = true

	submit.OnEvent(gowd.OnClick, func(_ *gowd.Element, event *gowd.EventElement) {
		var errs int
		if len(inputs.keyPath) == 0 {
			keyPathInput.SetHelpText("Required.")
			errs++
		} else {
			keyPathInput.HideHelpText()
		}
		if len(inputs.idfPath) == 0 {
			idfPathInput.SetHelpText("Required.")
			errs++
		} else {
			idfPathInput.HideHelpText()
		}
		inputs.nominatorWallet = nominatorWalletInput.GetValue()
		if len(inputs.nominatorWallet) == 0 {
			nominatorWalletInput.SetHelpText("Required.")
			errs++
		} else {
			if len(nominatorWalletInput.Kids) > 2 {
				nominatorWalletInput.RemoveElement(nominatorWalletInput.Kids[2])
			}
		}
		inputs.validatorWallet = validatorWalletInput.GetValue()
		if len(inputs.validatorWallet) == 0 {
			validatorWalletInput.SetHelpText("Required.")
			errs++
		} else {
			if len(validatorWalletInput.Kids) > 2 {
				validatorWalletInput.RemoveElement(validatorWalletInput.Kids[2])
			}
		}
		inputs.serverAddress = serverAddressInput.GetValue()
		if len(inputs.serverAddress) == 0 {
			serverAddressInput.SetHelpText("Required.")
			errs++
		} else {
			if len(serverAddressInput.Kids) > 2 {
				serverAddressInput.RemoveElement(serverAddressInput.Kids[2])
			}
		}
		if len(inputs.serverCertPath) == 0 {
			serverCertPathInput.SetHelpText("Required.")
			errs++
		} else {
			serverCertPathInput.HideHelpText()
		}
		inputs.agree = agreeInput.Checked()
		if inputs.agree == false {
			// TODO: print error
			errs++
		} else {
			if len(serverAddressInput.Kids) > 2 {
				serverAddressInput.RemoveElement(serverAddressInput.Kids[2])
			}
		}
		jww.INFO.Printf("Inputs set: %+v", inputs)

		if errs == 0 {
			err := client.SignAndTransmit(
				inputs.keyPath,
				inputs.idfPath,
				inputs.nominatorWallet,
				inputs.validatorWallet,
				inputs.serverAddress,
				inputs.serverCertPath,
				utils.Contract)
			if err != nil {
				errBox.SetText(err.Error())
				errBox.Hidden = false
			}
		}
	})

	contract := bootstrap.NewElement("div", "contractContainer")
	_, err := contract.AddHTML(utils.Contract, nil)
	if err != nil {
		return err
	}

	form := bootstrap.NewFormGroup(
		keyPathInput.Element,
		idfPathInput.Element,
		nominatorWalletInput.Element,
		validatorWalletInput.Element,
		serverAddressInput.Element,
		serverCertPathInput.Element,
		contract,
		agreeInput.Element,
		submit,
		errBox,
	)

	h1 := bootstrap.NewElement("h1", "")
	h1.SetText("MainNet Commitments")
	logo := bootstrap.NewElement("img", "logo")
	logo.SetAttribute("src", "img/xx-logo.svg")
	h1.AddElement(logo)
	p := bootstrap.NewElement("p", "blurb")
	p.SetText(blurbText)
	row := bootstrap.NewRow(bootstrap.NewElement("div", "well", h1, p, form))
	body.AddElement(row)

	// Start the UI loop
	err = gowd.Run(body)
	if err != nil {
		return err
	}

	return nil
}

type Inputs struct {
	keyPath, idfPath, nominatorWallet, validatorWallet, serverAddress, serverCertPath string
	agree                                                                             bool
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
