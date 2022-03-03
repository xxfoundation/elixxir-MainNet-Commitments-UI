package main

import (
	"git.xx.network/elixxir/mainnet-commitments-ui/formParts"
	"git.xx.network/elixxir/mainnet-commitments/client"
	"git.xx.network/elixxir/mainnet-commitments/utils"
	"github.com/dtylman/gowd"
	"github.com/dtylman/gowd/bootstrap"
	jww "github.com/spf13/jwalterweatherman"
	utils2 "gitlab.com/xx_network/primitives/utils"
)

var body *gowd.Element

const blurbText = "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vivamus malesuada eleifend ultrices. Pellentesque habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas. Etiam pretium tempor massa, a volutpat orci mattis non. Integer tincidunt tincidunt ante sed cursus. In lacinia pulvinar tempor. Nullam id luctus nibh, vitae iaculis ante. Integer vel sem at augue viverra suscipit vel nec orci. Sed sed ultrices quam."
const serverAddress = "http://localhost:11420"

type Inputs struct {
	keyPath         string
	idfPath         string
	paymentWallet   string
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
	keyPathInput := formParts.NewFileButton("BetaNet Server Key (.key)", &inputs.keyPath)
	idfPathInput := formParts.NewFileButton("BetaNet Server IDF (.json)", &inputs.idfPath)
	paymentWalletInput := bootstrap.NewFormInput("text", "Wallet to receive payment")
	serverCertPathInput := formParts.NewFileButton("BetaNet Server Certificate (.crt)", &inputs.serverCertPath)

	agreeInput := bootstrap.NewCheckBox("I agree to the contract above.", false)
	agreeHelpText := bootstrap.NewElement("p", "help-block")
	agreeHelpText.Hidden = true
	agreeInput.AddElement(gowd.NewElement("br"))
	agreeInput.AddElement(agreeHelpText)
	agreeBox1 := bootstrap.NewElement("div", "form-group", agreeInput.Element)

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
		spinner.Hidden = false
		formErrors.Hidden = true
		errBox.Hidden = true
		defer func() {
			body.Enable()
			submit.Enable()
			spinner.Hidden = true
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
		inputs.paymentWallet = paymentWalletInput.GetValue()
		if len(inputs.paymentWallet) == 0 {
			paymentWalletInput.SetHelpText("Required.")
			errs++
		} else {
			if len(paymentWalletInput.Kids) > 2 {
				paymentWalletInput.Kids[2].Hidden = true
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
			if err := body.Render(); err != nil {
				jww.ERROR.Print(err)
			}

			err := client.SignAndTransmit(
				inputs.keyPath,
				inputs.idfPath,
				inputs.paymentWallet,
				serverAddress,
				inputs.serverCert,
				utils.NovemberContract)

			spinner.Hidden = true

			if err != nil {
				jww.ERROR.Printf("Submit error: %+v", err)
				errBox.SetText(err.Error())
				errBox.Hidden = false
				formErrors.SetText("The were errors in the form input. Please correct them to continue.")
				formErrors.Hidden = false
			} else {
				divWell.RemoveElements()
				success := bootstrap.NewElement("span", "success", gowd.NewText("November BetaNet Compensation Successful."))
				divWell.AddElement(success)
			}
		} else {
			formErrors.SetText("The were errors in the form input. Please correct them to continue.")
			formErrors.Hidden = false
		}
	})

	contractText := bootstrap.NewElement("p", "contractText")
	contractText.SetText("Read through the entire contract below and accept the terms.")

	contract := bootstrap.NewElement("div", "contractBox", contractText)
	contract1 := bootstrap.NewElement("div", "contractContainer")
	_, err := contract1.AddHTML(utils.NovemberContract, nil)
	if err != nil {
		return err
	}
	contractLink := bootstrap.NewLinkButton("Open in new window")
	contractLink.RemoveAttribute("href")
	contractLink.OnEvent(gowd.OnClick, func(*gowd.Element, *gowd.EventElement) {
		gowd.ExecJSNow(`
let prtContent = document.getElementById("` + contract1.GetID() + `");
let WinPrint = window.open('', '', '');
WinPrint.document.write('<title>TERMS AND CONDITIONS FOR MAINNET SUPPORT REIMBURSEMENT</title>');
WinPrint.document.write(prtContent.innerHTML);
WinPrint.document.close();
WinPrint.focus();`)
	})
	contractPrintLink := bootstrap.NewLinkButton("Print")
	contractPrintLink.RemoveAttribute("href")
	contractPrintLink.OnEvent(gowd.OnClick, func(*gowd.Element, *gowd.EventElement) {
		gowd.ExecJSNow(`
let prtContent = document.getElementById("` + contract1.GetID() + `");
let WinPrint = window.open('', '', '');
WinPrint.document.write('<title>TERMS AND CONDITIONS FOR MAINNET SUPPORT REIMBURSEMENT</title>');
WinPrint.document.write(prtContent.innerHTML);
WinPrint.document.close();
WinPrint.focus();
WinPrint.print();
WinPrint.close();`)
	})
	contractLinkDiv := bootstrap.NewElement("div", "contractLink", contractLink, contractPrintLink)

	contract.AddElement(contract1)
	contract.AddElement(contractLinkDiv)
	contract.AddElement(agreeBox1)

	form := bootstrap.NewFormGroup(
		formErrors,
		serverCertPathInput.Element,
		keyPathInput.Element,
		idfPathInput.Element,
		paymentWalletInput.Element,
		contract,
		submitBox,
	)

	form.SetAttribute("style", "margin-top:35px")

	h1 := bootstrap.NewElement("h1", "")
	h1.SetText("November BetaNet Compensation")
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
