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

func buildPage() error {

	inputs := Inputs{}

	// creates a new bootstrap fluid container
	body = bootstrap.NewContainer(false)
	// add some elements using the object model

	jww.ERROR.Printf("PAGE LOADED")

	// keyPathInput := bootstrap.NewFileButton(bootstrap.ButtonDefault, "keyPath", false)
	keyPathInput := bootstrap.NewFormInput("file", "keyPath")
	keyPathInput.SetValue("test")
	idfPathInput := bootstrap.NewFormInput("file", "idfPath")
	nominatorWalletInput := bootstrap.NewFormInput("text", "nominatorWallet")
	validatorWalletInput := bootstrap.NewFormInput("text", "validatorWallet")
	serverAddressInput := bootstrap.NewFormInput("text", "serverAddress")
	serverCertPathInput := bootstrap.NewFormInput("file", "serverCertPath")
	submit := bootstrap.NewButton(bootstrap.ButtonPrimary, "button")

	errBox := bootstrap.NewElement("div", "well")

	// keyPathInput.OnEvent(gowd.OnChange, func(sender *gowd.Element, event *gowd.EventElement) {
	// 	jww.INFO.Printf("keyPath sender: %+v", sender)
	// 	jww.INFO.Printf("keyPath event: %+v", event)
	// 	inputs.keyPath = keyPathInput.GetValue()
	// 	jww.INFO.Printf("keyPath set: %q", inputs.keyPath)
	// })
	// idfPathInput.OnEvent(gowd.OnChange, func(_ *gowd.Element, event *gowd.EventElement) {
	// 	inputs.idfPath = event.GetValue()
	// 	jww.INFO.Printf("idfPath set: %q", inputs.idfPath)
	// })
	// nominatorWalletInput.OnEvent(gowd.OnChange, func(_ *gowd.Element, event *gowd.EventElement) {
	// 	inputs.nominatorWallet = event.GetValue()
	// 	jww.INFO.Printf("nominatorWallet set: %q", inputs.nominatorWallet)
	// })
	// validatorWalletInput.OnEvent(gowd.OnChange, func(_ *gowd.Element, event *gowd.EventElement) {
	// 	inputs.validatorWallet = event.GetValue()
	// 	jww.INFO.Printf("validatorWallet set: %q", inputs.validatorWallet)
	// })
	// serverAddressInput.OnEvent(gowd.OnChange, func(_ *gowd.Element, event *gowd.EventElement) {
	// 	inputs.serverAddress = event.GetValue()
	// 	jww.INFO.Printf("serverAddress set: %q", inputs.serverAddress)
	// })
	// serverCertPathInput.OnEvent(gowd.OnChange, func(_ *gowd.Element, event *gowd.EventElement) {
	// 	inputs.serverCertPath = event.GetValue()
	// 	jww.INFO.Printf("serverCertPath set: %q", inputs.serverCertPath)
	// })
	submit.OnEvent(gowd.OnClick, func(_ *gowd.Element, event *gowd.EventElement) {
		var errs int
		inputs.keyPath = keyPathInput.GetValue()
		if len(inputs.keyPath) == 0 {
			keyPathInput.SetHelpText("Required.")
			errs++
		}
		inputs.idfPath = keyPathInput.GetValue()
		if len(inputs.idfPath) == 0 {
			idfPathInput.SetHelpText("Required.")
			errs++
		}
		inputs.nominatorWallet = nominatorWalletInput.GetValue()
		if len(inputs.nominatorWallet) == 0 {
			nominatorWalletInput.SetHelpText("Required.")
			errs++
		}
		inputs.validatorWallet = validatorWalletInput.GetValue()
		if len(inputs.validatorWallet) == 0 {
			validatorWalletInput.SetHelpText("Required.")
			errs++
		}
		inputs.serverAddress = serverAddressInput.GetValue()
		if len(inputs.serverAddress) == 0 {
			serverAddressInput.SetHelpText("Required.")
			errs++
		}
		inputs.serverCertPath = serverCertPathInput.GetValue()
		if len(inputs.serverCertPath) == 0 {
			serverCertPathInput.SetHelpText("Required.")
			errs++
		}
		jww.INFO.Printf("Inputs set: %+v", inputs)

		if errs == 0 {
			err := client.SignAndTransmit(inputs.keyPath,
				inputs.idfPath,
				inputs.nominatorWallet,
				inputs.validatorWallet,
				inputs.serverAddress,
				inputs.serverCertPath,
				utils.Contract)
			if err != nil {
				errBox.SetText(err.Error())
			}
		}
	})

	form := bootstrap.NewFormGroup(
		keyPathInput.Element,
		idfPathInput.Element,
		nominatorWalletInput.Element,
		validatorWalletInput.Element,
		serverAddressInput.Element,
		serverCertPathInput.Element,
		submit)

	row := bootstrap.NewRow(bootstrap.NewColumn(bootstrap.ColumnSmall, 3, bootstrap.NewColumn(bootstrap.ColumnSmall, 3, bootstrap.NewElement("div", "well", form), errBox)))
	body.AddElement(row)
	// Start the UI loop
	err := gowd.Run(body)
	if err != nil {
		return err
	}

	return nil
}

type Inputs struct {
	keyPath, idfPath, nominatorWallet, validatorWallet, serverAddress, serverCertPath string
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
