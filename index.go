package main

import (
	"git.xx.network/elixxir/mainnet-commitments-ui/formParts"
	"git.xx.network/elixxir/mainnet-commitments/client"
	"git.xx.network/elixxir/mainnet-commitments/utils"
	"github.com/dtylman/gowd"
	"github.com/dtylman/gowd/bootstrap"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/xx-labs/sleeve/wallet"
	utils2 "gitlab.com/xx_network/primitives/utils"
	"strconv"
)

var body *gowd.Element

const blurbText = `This applet will allow you to commit your wallets. For more information, please see the&nbsp;`
const serverAddress = "https://18.185.229.39:11420"

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
	keyPathInput := formParts.NewFileButton("BetaNet Server Key (.key)", &inputs.keyPath)
	idfPathInput := formParts.NewFileButton("BetaNet Server IDF (.json)", &inputs.idfPath)
	nominatorWalletInput := bootstrap.NewFormInput("text", "Nominator Wallet Address")
	validatorWalletInput := bootstrap.NewFormInput("text", "Validator Wallet Address")
	serverCertPathInput := formParts.NewFileButton("BetaNet Server Certificate (.cert or .crt)", &inputs.serverCertPath)

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

	shade := bootstrap.NewElement("div", "shade")
	body.AddElement(shade)
	confirmBoxText := bootstrap.NewElement("span", "", gowd.NewText("You did not supply a nominator wallet address. If you sent your November compensation to your validator wallet address, then you may continue. Otherwise, you should cancel and add a validator wallet address."))
	confirmBoxYes := bootstrap.NewButton(bootstrap.ButtonPrimary, "Submit")
	confirmBoxNo := bootstrap.NewButton(bootstrap.ButtonPrimary, "Cancel")
	confirmBox := bootstrap.NewElement("div", "confirmBox", confirmBoxText, confirmBoxYes, confirmBoxNo)

	confirmBoxYes.OnEvent(gowd.OnClick, func(*gowd.Element, *gowd.EventElement) {
		confirmBox.Hide()
		shade.Hide()
		spinner.Show()
		body.Render()
		submitContract(inputs, spinner, errBox, formErrors, divWell)
	})
	confirmBoxNo.OnEvent(gowd.OnClick, func(*gowd.Element, *gowd.EventElement) {
		confirmBox.Hide()
		shade.Hide()
		nominatorWalletInput.SetAttribute("style", "background: #ffffba;border:6px solid #ffffba;margin:-6px -6px 9px;")
	})
	body.AddElement(confirmBox)
	confirmBox.Hide()
	shade.Hide()

	submit.OnEvent(gowd.OnClick, func(_ *gowd.Element, event *gowd.EventElement) {
		nominatorWalletInput.RemoveAttribute("style")
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
		inputs.nominatorWallet = nominatorWalletInput.GetValue()
		if len(inputs.nominatorWallet) == 0 {
			// nominatorWalletInput.SetHelpText("Required.")
			// errs++
		} else {
			ok, err := wallet.ValidateXXNetworkAddress(inputs.nominatorWallet)
			if !ok || err != nil {
				nominatorWalletInput.SetHelpText("Invalid wallet address.")
				jww.ERROR.Printf("Invalid nominator wallet address: %+v", err)
				errs++
			} else if len(nominatorWalletInput.Kids) > 2 {
				nominatorWalletInput.Kids[2].Hidden = true
			}
		}
		inputs.validatorWallet = validatorWalletInput.GetValue()
		if len(inputs.validatorWallet) == 0 {
			validatorWalletInput.SetHelpText("Required.")
			errs++
		} else {
			ok, err := wallet.ValidateXXNetworkAddress(inputs.validatorWallet)
			if !ok || err != nil {
				validatorWalletInput.SetHelpText("Invalid wallet address.")
				jww.ERROR.Printf("Invalid validator wallet address: %+v", err)
				errs++
			} else if len(validatorWalletInput.Kids) > 2 {
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

		if errs == 0 {
			if err := body.Render(); err != nil {
				jww.ERROR.Print(err)
			}

			if inputs.nominatorWallet == "" {
				confirmBox.Show()
				shade.Show()
			} else {
				submitContract(inputs, spinner, errBox, formErrors, divWell)
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
	contractFontSize := 100
	contract1.SetAttribute("style", "font-size:"+strconv.Itoa(contractFontSize)+"%;")
	_, err := contract1.AddHTML(utils.Contract, nil)
	if err != nil {
		return err
	}
	contractLink := bootstrap.NewLinkButton("Open in new window")
	contractLink.RemoveAttribute("href")
	contractLink.OnEvent(gowd.OnClick, func(*gowd.Element, *gowd.EventElement) {
		gowd.ExecJSNow(`
let prtContent = document.getElementById("` + contract1.GetID() + `");
let WinPrint = window.open('', '', '');
WinPrint.document.write('<title>PARTICIPANT TERMS AND CONDITIONS FOR MAINNET TRANSITION PROGRAM</title>');
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
WinPrint.document.write('<title>PARTICIPANT TERMS AND CONDITIONS FOR MAINNET TRANSITION PROGRAM</title>');
WinPrint.document.write(prtContent.innerHTML);
WinPrint.document.close();
WinPrint.focus();
WinPrint.print();
WinPrint.close();`)
	})

	incFontSizeLink := bootstrap.NewLinkButton("+")
	incFontSizeLink.RemoveAttribute("href")
	incFontSizeLink.OnEvent(gowd.OnClick, func(*gowd.Element, *gowd.EventElement) {
		contractFontSize += 5
		contract1.SetAttribute("style", "font-size:"+strconv.Itoa(contractFontSize)+"%;")
	})
	decFontSizeLink := bootstrap.NewLinkButton("-")
	decFontSizeLink.RemoveAttribute("href")
	decFontSizeLink.OnEvent(gowd.OnClick, func(*gowd.Element, *gowd.EventElement) {
		contractFontSize -= 5
		contract1.SetAttribute("style", "font-size:"+strconv.Itoa(contractFontSize)+"%;")
	})
	fontSizeSpan := bootstrap.NewElement("span", "", gowd.NewText("Font size: "), incFontSizeLink, decFontSizeLink)
	fontSizeSpan.SetAttribute("style", "float:right;font-size:92%;")

	contractLinkDiv := bootstrap.NewElement("div", "contractLink", contractLink, contractPrintLink, fontSizeSpan)

	contract.AddElement(contract1)
	contract.AddElement(contractLinkDiv)
	contract.AddElement(agreeBox1)

	form := bootstrap.NewFormGroup(
		formErrors,
		serverCertPathInput.Element,
		keyPathInput.Element,
		idfPathInput.Element,
		nominatorWalletInput.Element,
		validatorWalletInput.Element,
		contract,
		submitBox,
	)

	form.SetAttribute("style", "margin-top:35px")

	h1 := bootstrap.NewElement("h1", "")
	h1.SetText("xx network MainNet Wallet Commitment")
	logo := bootstrap.NewElement("img", "logo")
	logo.SetAttribute("src", "img/xx-logo.svg")
	h1.AddElement(logo)
	p := bootstrap.NewElement("p", "blurb")
	p.AddHTML(blurbText, nil)
	instructionPageLink := bootstrap.NewLinkButton("instructions page")
	instructionPageLink.RemoveAttribute("href")
	instructionPageLink.SetAttribute("style", "cursor:pointer;")
	instructionPageLink.OnEvent(gowd.OnClick, func(*gowd.Element, *gowd.EventElement) {
		gowd.ExecJSNow("window.nw.Shell.openExternal('https://xx.network/november-betanet-compensation-applet-instructions/')")
	})
	p.AddElement(instructionPageLink)
	p.AddElement(gowd.NewText("."))
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

func submitContract(inputs Inputs, spinner, errBox, formErrors, divWell *gowd.Element) {
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
		errBox.SetText("An error occurred when submitting the request. Please contact support at nodes@xx.network and provide the following error message:")
		errBox.AddElement(bootstrap.NewElement("span", "errorBoxMessage", gowd.NewText(err.Error())))
		errBox.Hidden = false
		formErrors.SetText("The were errors in the form input. Please correct them to continue.")
		formErrors.Hidden = false
	} else {
		divWell.RemoveElements()
		success := bootstrap.NewElement("span", "success", gowd.NewText("MainNet Commitments Successful."))
		divWell.AddElement(success)
	}
}
