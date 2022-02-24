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

const blurbText = "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vivamus malesuada eleifend ultrices. Pellentesque habitant morbi tristique senectus et netus et malesuada fames ac turpis egestas. Etiam pretium tempor massa, a volutpat orci mattis non. Integer tincidunt tincidunt ante sed cursus. In lacinia pulvinar tempor. Nullam id luctus nibh, vitae iaculis ante. Integer vel sem at augue viverra suscipit vel nec orci. Sed sed ultrices quam."
const serverAddress = "http://localhost:11420"
const twoContracts = true

type Inputs struct {
	keyPath         string
	idfPath         string
	nominatorWallet string
	validatorWallet string
	serverCert      string
	serverCertPath  string
	agree1, agree2  bool
}

func buildPage() error {

	inputs := Inputs{}

	// creates a new bootstrap fluid container
	body = bootstrap.NewContainer(false)
	// add some elements using the object model

	// keyPathInput := bootstrap.NewFileButton(bootstrap.ButtonDefault, "keyPath", false)
	keyPathInput := NewFileButton("BetaNet Server Key (.key)", &inputs.keyPath)
	idfPathInput := NewFileButton("BetaNet Server IDF (.json)", &inputs.idfPath)
	nominatorWalletInput := bootstrap.NewFormInput("text", "Nominator Wallet Address")
	validatorWalletInput := bootstrap.NewFormInput("text", "Validator Wallet Address")
	serverCertPathInput := NewFileButton("BetaNet Server Certificate (.crt)", &inputs.serverCertPath)

	agreeInput1 := bootstrap.NewCheckBox("I agree to the contract above.", false)
	agreeHelpText1 := bootstrap.NewElement("p", "help-block")
	agreeHelpText1.Hidden = true
	agreeInput1.AddElement(gowd.NewElement("br"))
	agreeInput1.AddElement(agreeHelpText1)
	agreeBox1 := bootstrap.NewElement("div", "form-group", agreeInput1.Element)

	agreeInput2 := bootstrap.NewCheckBox("I agree to the contract above.", false)
	agreeHelpText2 := bootstrap.NewElement("p", "help-block")
	agreeHelpText2.Hidden = true
	agreeInput2.AddElement(gowd.NewElement("br"))
	agreeInput2.AddElement(agreeHelpText2)
	agreeBox2 := bootstrap.NewElement("div", "form-group", agreeInput2.Element)

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
		inputs.agree1 = agreeInput1.Checked()
		if inputs.agree1 == false {
			agreeHelpText1.SetText("Required.")
			agreeHelpText1.Hidden = false
			errs++
		} else {
			agreeHelpText1.Hidden = true
		}
		if twoContracts {
			inputs.agree2 = agreeInput2.Checked()
			if inputs.agree2 == false {
				agreeHelpText2.SetText("Required.")
				agreeHelpText2.Hidden = false
				errs++
			} else {
				agreeHelpText2.Hidden = true
			}
		}
		jww.INFO.Printf("Inputs set: %+v", inputs)

		if errs == 0 {
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

	contractText := bootstrap.NewElement("p", "contractText")
	if twoContracts {
		contractText.SetText("Read through both contracts below and accept the terms both.")
	} else {
		contractText.SetText("Read through the entire contract below and accept the terms.")
	}

	contract := bootstrap.NewElement("div", "contractBox", contractText)
	contract1 := bootstrap.NewElement("div", "contractContainer")
	_, err := contract1.AddHTML(utils.Contract, nil)
	if err != nil {
		return err
	}
	contractLink1 := bootstrap.NewLinkButton("Open in new window")
	contractLink1.OnEvent(gowd.OnClick, func(*gowd.Element, *gowd.EventElement) {
		saveHTML("Contract", "contract1.html", utils.Contract)
	})
	contractLink1.SetAttribute("href", "contract1.html")
	contractLink1.SetAttribute("target", "_blank")
	contractLinkDiv1 := bootstrap.NewElement("div", "contractLink", contractLink1)

	if twoContracts {
		contract2 := bootstrap.NewElement("div", "contractContainer")
		_, err = contract2.AddHTML(utils.Contract, nil)
		if err != nil {
			return err
		}
		contractLink2 := bootstrap.NewLinkButton("Open in new window")
		contractLink2.OnEvent(gowd.OnClick, func(*gowd.Element, *gowd.EventElement) {
			saveHTML("Contract", "contract2.html", utils.Contract)
		})
		contractLink2.SetAttribute("href", "contract2.html")
		contractLink2.SetAttribute("target", "_blank")
		contractLinkDiv2 := bootstrap.NewElement("div", "contractLink", contractLink2)

		contract.AddElement(contract1)
		contract.AddElement(contractLinkDiv1)
		contract.AddElement(agreeBox1)
		contract.AddElement(contract2)
		contract.AddElement(contractLinkDiv2)
		contract.AddElement(agreeBox2)
	} else {
		contract.AddElement(contract1)
		contract.AddElement(contractLinkDiv1)
		contract.AddElement(agreeBox1)
	}

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
