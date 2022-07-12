package main

import (
	"encoding/json"
	"git.xx.network/elixxir/mainnet-commitments-ui/form"
	"git.xx.network/elixxir/mainnet-commitments/client"
	"git.xx.network/elixxir/mainnet-commitments/utils"
	"github.com/dtylman/gowd"
	"github.com/dtylman/gowd/bootstrap"
	jww "github.com/spf13/jwalterweatherman"
	"strconv"
	"time"
)

const test = false

var body *gowd.Element

const (
	blurbTextPg1 = `This applet will allow you to commit your wallets. For more information, please see the&nbsp;`
	blurbTextPg2 = `Below are the committed validator and nominator addresses. Select the checkbox to modify them.`
	blurbTextPg3 = `Below is the selected team stake. Select the checkbox to modify it.`
)
const serverAddress = "http://localhost:11420"

type Inputs struct {
	certPath            string
	keyPath             string
	idfPath             string
	nodeID              string
	nominatorWallet     string
	validatorWallet     string
	origNominatorWallet string
	origValidatorWallet string
	agree               bool

	cert      []byte
	key       []byte
	nodeIdHex string
	email     string

	origMultiplier uint64
	multiplier     uint64
	maxMultiplier  uint64

	walletModifyCheck     bool
	multiplierModifyCheck bool
}

func buildPage() error {

	// creates a new bootstrap fluid container
	body = bootstrap.NewContainer(false)
	// add some elements using the object model

	// keyPathInput := bootstrap.NewFileButton(bootstrap.ButtonDefault, "keyPath", false)

	row := page1(Inputs{})
	// row = page1(Inputs{
	// 	certPath: "C:\\Users\\Jono\\Go\\src\\git.xx.network\\elixxir\\mainnet-commitments-ui\\tmp\\server.crt",
	// 	keyPath:  "C:\\Users\\Jono\\Go\\src\\git.xx.network\\elixxir\\mainnet-commitments-ui\\tmp\\commitmenttestkey.key",
	// 	idfPath:  "C:\\Users\\Jono\\Go\\src\\git.xx.network\\elixxir\\mainnet-commitments-ui\\tmp\\testidf.json",
	// })
	// row = page3(Inputs{maxMultiplier: 1500, multiplier: 543, origMultiplier: 543})

	body.AddElement(row)

	// Start the UI loop
	err := gowd.Run(body)
	if err != nil {
		return err
	}

	return nil
}

func page1(inputs Inputs) *gowd.Element {

	certInput := form.NewFileButton("Node Certificate (.cert or .crt)", form.ValidateFilePath)
	certInput.SetValue(inputs.certPath)
	keyInput := form.NewFileButton("Node Key (.key)", form.ValidateFilePath)
	keyInput.SetValue(inputs.keyPath)
	idfInput := form.NewFileButton("Node IDF (.json)", form.ValidateIdfPath)
	idfInput.SetValue(inputs.idfPath)

	submit := bootstrap.NewButton(bootstrap.ButtonPrimary, "Submit")
	errBox := bootstrap.NewElement("span", "errorBox")
	errBox.Hidden = true
	spinner := bootstrap.NewElement("div", "spinnerContainer", bootstrap.NewElement("div", "spinner", gowd.NewText("Loading...")))
	spinner.Hidden = true
	submitBox := bootstrap.NewElement("div", "submitBox", errBox, spinner, submit)

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

		// Validate inputs
		if validated, valid := certInput.Validate(); valid {
			inputs.certPath = certInput.GetValue()
			inputs.cert = validated.([]byte)
		} else {
			errs++
		}
		if validated, valid := keyInput.Validate(); valid {
			inputs.keyPath = keyInput.GetValue()
			inputs.key = validated.([]byte)
		} else {
			errs++
		}
		if validated, valid := idfInput.Validate(); valid {
			inputs.idfPath = idfInput.GetValue()
			inputs.nodeIdHex = validated.(string)
		} else {
			errs++
		}

		jww.INFO.Printf("Inputs set: %+v", inputs)

		if errs == 0 {
			if err := body.Render(); err != nil {
				jww.ERROR.Print(err)
			}

			spinner.Hidden = true

			getInfoTest := func(nid, serverCert, serverAddress string) ([]byte, error) {
				if test {
					time.Sleep(500 * time.Millisecond)
					return []byte(`{
			    "validator-wallet": "6YLQDuXq2PkgPBQXwPPYnUyiQfJbE5Xfyu8JpmjySFz2T4sP",
			    "nominator-wallet": "6YaEntt2HKZ3ZunAZyzqBfD1xGsoRnwBdWd6Zd4yLnmwHgsg",
			    "selected-multiplier": 573,
			    "max-multiplier": 1500,
				"email": "johnDoe@email.com"
			}`), nil
				} else {
					return client.GetInfo(nid, serverCert, serverAddress)
				}
			}

			jsonData, err := getInfoTest(inputs.nodeIdHex, string(inputs.cert), serverAddress)

			if err != nil {
				jww.ERROR.Printf("Submit error: %+v", err)
				errBox.SetText("An error occurred when submitting the request. Please contact support at nodes@xx.network and provide the following error message:")
				errBox.AddElement(bootstrap.NewElement("span", "errorBoxMessage", gowd.NewText(err.Error())))
				errBox.Hidden = false
				formErrors.SetText("There were errors in the form input. Please correct them to continue.")
				formErrors.Hidden = false
			} else {
				type jsonInfo struct {
					ValidatorWallet string `json:"validator-wallet"`
					NominatorWallet string `json:"nominator-wallet"`
					SelectedStake   uint64 `json:"selected-stake"`
					MaxStake        uint64 `json:"max-stake"`
					Email           string `json:"email"`
				}

				var info jsonInfo

				err = json.Unmarshal(jsonData, &info)
				if err != nil {
					jww.ERROR.Printf("JSON unmarshal error: %+v", err)
					errBox.SetText("An error occurred when submitting the request. Please contact support at nodes@xx.network and provide the following error message:")
					errBox.AddElement(bootstrap.NewElement("span", "errorBoxMessage", gowd.NewText(err.Error())))
					errBox.Hidden = false
					formErrors.SetText("There were errors in the form input. Please correct them to continue.")
					formErrors.Hidden = false
				} else {

					inputs.validatorWallet = info.ValidatorWallet
					inputs.nominatorWallet = info.NominatorWallet
					inputs.origValidatorWallet = info.ValidatorWallet
					inputs.origNominatorWallet = info.NominatorWallet
					inputs.multiplier = info.SelectedStake
					inputs.origMultiplier = info.SelectedStake
					inputs.maxMultiplier = info.MaxStake
					inputs.email = info.Email

					body.RemoveElements()
					body.AddElement(page2(inputs))

					// divWell.RemoveElements()
					// success := bootstrap.NewElement("span", "success", gowd.NewText("MainNet Commitments Successful."))
					// result := gowd.NewText(fmt.Sprintf("%+v", inputs))
					// divWell.AddElement(success)
					// divWell.AddElement(result)
				}

			}
		} else {
			formErrors.SetText("There were errors in the form input. Please correct them to continue.")
			formErrors.Hidden = false
		}
	})

	formGrp := bootstrap.NewFormGroup(
		formErrors,
		certInput.Element,
		keyInput.Element,
		idfInput.Element,
		submitBox,
	)

	formGrp.SetAttribute("style", "margin: 2.5em 0 0")

	h1 := bootstrap.NewElement("h1", "")
	h1.SetText("Update Team Stake")
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

	emailInput := form.NewPart("email", "Email to receive notification on changes to the state of your validator (optional)", form.ValidateEmail)
	emailInput.SetValue(inputs.email)

	validatorWallet := form.NewPart("text", "Validator Wallet Address", form.ValidateXXNetworkAddress)
	validatorWallet.SetValue(inputs.validatorWallet)
	validatorWallet.Disable()
	nominatorWallet := form.NewPart("text", "Nominator Wallet Address", form.ValidateXXNetworkAddress)
	nominatorWallet.SetValue(inputs.nominatorWallet)
	nominatorWallet.Disable()

	modifyCheck := form.NewPart("checkbox", "Modify Wallet Addresses", nil)
	modifyCheck.SetAttribute("style", "margin-top:3em;")
	if inputs.walletModifyCheck {
		modifyCheck.Check()
		validatorWallet.Enable()
		nominatorWallet.Enable()
	}
	modifyCheck.OnEvent(gowd.OnClick, func(sender *gowd.Element, event *gowd.EventElement) {
		if modifyCheck.Checked() {
			inputs.walletModifyCheck = true
			validatorWallet.Enable()
			nominatorWallet.Enable()
		} else {
			inputs.walletModifyCheck = false
			validatorWallet.Disable()
			nominatorWallet.Disable()
			validatorWallet.SetValue(inputs.origValidatorWallet)
			nominatorWallet.SetValue(inputs.origNominatorWallet)
		}
	})

	back := bootstrap.NewButton(bootstrap.ButtonPrimary, "Back")
	back.SetAttribute("style", "margin-right:1em;")
	submit := bootstrap.NewButton(bootstrap.ButtonPrimary, "Next")
	errBox := bootstrap.NewElement("span", "errorBox")
	errBox.Hidden = true
	spinner := bootstrap.NewElement("div", "spinnerContainer", bootstrap.NewElement("div", "spinner", gowd.NewText("Loading...")))
	spinner.Hidden = true
	submitBox := bootstrap.NewElement("div", "submitBox", errBox, spinner, back, submit)

	formErrors := bootstrap.NewElement("p", "formErrors")
	formErrors.Hidden = true

	divWell := bootstrap.NewElement("div", "well")

	back.OnEvent(gowd.OnClick, func(sender *gowd.Element, event *gowd.EventElement) {
		body.RemoveElements()
		body.AddElement(page1(inputs))
	})

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

		if validated, valid := emailInput.Validate(); valid {
			inputs.email = validated.(string)
		} else {
			errs++
		}
		if validated, valid := validatorWallet.Validate(); valid {
			inputs.validatorWallet = validated.(string)
		} else {
			errs++
		}
		if validated, valid := nominatorWallet.Validate(); valid {
			inputs.nominatorWallet = validated.(string)
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
		emailInput.Element(),
		modifyCheck.Element(),
		validatorWallet.Element(),
		nominatorWallet.Element(),
		submitBox,
	)

	formGrp.SetAttribute("style", "margin: 2.5em 0 0")

	h1 := bootstrap.NewElement("h1", "")
	h1.SetText("Update Team Stake")
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

	multiplier := form.NewPart("number", "Selected Stake (Max "+strconv.FormatUint(inputs.maxMultiplier, 10)+"xx): ", form.ValidateMultiplier(inputs.maxMultiplier))
	multiplier.SetValue(strconv.FormatUint(inputs.multiplier, 10))
	multiplier.SetInputAttribute("class", "multiplier modifier")
	multiplier.SetInputAttribute("step", "1")
	multiplier.SetInputAttribute("min", "0")
	multiplier.SetInputAttribute("max", strconv.FormatUint(inputs.maxMultiplier, 10))
	multiplier.SetLabelAttribute("for", "number")
	multiplier.Disable()

	multiplier.AddElement(bootstrap.NewElement("p", "inputXX", gowd.NewText("xx")))
	multiplier.SetHelpTxtAttribute("style", "display:table;")
	multiplier.SwapKids(2, 3)

	multiplier.SetInputAttribute("id", "number")
	multiplier.SetInputAttribute("onkeyup", "changeRangeValue(this.value, "+strconv.FormatUint(inputs.maxMultiplier, 10)+")")
	multiplier.SetInputAttribute("onclick", "changeRangeValue(this.value, "+strconv.FormatUint(inputs.maxMultiplier, 10)+")")

	slider := bootstrap.NewElement("input", "stakeRange")
	slider.SetAttribute("type", "range")
	slider.SetAttribute("min", "0")
	slider.SetAttribute("max", strconv.FormatUint(inputs.maxMultiplier, 10))
	slider.SetValue(strconv.FormatUint(inputs.multiplier, 10))
	slider.SetAttribute("id", "range")
	slider.SetAttribute("oninput", "changeInputValue(this.value)")
	slider.Disable()
	multiplier.AddElement(slider)
	multiplier.SwapKids(3, 4)

	modifyCheck := form.NewPart("checkbox", "Modify the selected stake", nil)
	if inputs.multiplierModifyCheck {
		modifyCheck.Check()
		multiplier.Enable()
		slider.Enable()
	}
	modifyCheck.OnEvent(gowd.OnClick, func(sender *gowd.Element, event *gowd.EventElement) {
		if modifyCheck.Checked() {
			inputs.multiplierModifyCheck = true
			multiplier.Enable()
			slider.Enable()
		} else {
			inputs.multiplierModifyCheck = false
			multiplier.Disable()
			slider.Disable()
			multiplier.SetValue(strconv.FormatUint(inputs.origMultiplier, 10))
		}
	})

	back := bootstrap.NewButton(bootstrap.ButtonPrimary, "Back")
	back.SetAttribute("style", "margin-right:1em;")
	submit := bootstrap.NewButton(bootstrap.ButtonPrimary, "Next")
	errBox := bootstrap.NewElement("span", "errorBox")
	errBox.Hidden = true
	spinner := bootstrap.NewElement("div", "spinnerContainer", bootstrap.NewElement("div", "spinner", gowd.NewText("Loading...")))
	spinner.Hidden = true
	submitBox := bootstrap.NewElement("div", "submitBox", errBox, spinner, back, submit)

	formErrors := bootstrap.NewElement("p", "formErrors")
	formErrors.Hidden = true

	divWell := bootstrap.NewElement("div", "well")

	back.OnEvent(gowd.OnClick, func(sender *gowd.Element, event *gowd.EventElement) {
		inputs.multiplierModifyCheck = modifyCheck.Checked()
		inputs.multiplier, _ = strconv.ParseUint(multiplier.GetValue(), 10, 64)
		body.RemoveElements()
		body.AddElement(page2(inputs))
	})

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

		if validated, valid := multiplier.Validate(); valid {
			inputs.multiplier = validated.(uint64)
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
			body.AddElement(page4(inputs))
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

	formGrp.SetAttribute("style", "margin: 2.5em 0 0")

	h1 := bootstrap.NewElement("h1", "")
	h1.SetText("Update Team Stake")
	logo := bootstrap.NewElement("img", "logo")
	logo.SetAttribute("src", "img/xx-logo.svg")
	h1.AddElement(logo)
	p := bootstrap.NewElement("p", "blurb")
	p.AddHTML(blurbTextPg3, nil)
	p.AddElement(gowd.NewText("."))
	divWell.AddElement(h1)
	divWell.AddElement(p)
	divWell.AddElement(formGrp)
	row := bootstrap.NewRow(divWell)

	return row
}
func page4(inputs Inputs) *gowd.Element {

	agreeInput := form.NewPart("checkbox", "I agree to the contract above.", form.ValidateCheckbox)
	agreeInput.SetAttribute("style", "margin-top:1em;")
	back := bootstrap.NewButton(bootstrap.ButtonPrimary, "Back")
	back.SetAttribute("style", "margin-right:1em;")
	submit := bootstrap.NewButton(bootstrap.ButtonPrimary, "Sign and Submit")
	errBox := bootstrap.NewElement("span", "errorBox")
	errBox.Hidden = true
	spinner := bootstrap.NewElement("div", "spinnerContainer", bootstrap.NewElement("div", "spinner", gowd.NewText("Loading...")))
	spinner.Hidden = true
	submitBox := bootstrap.NewElement("div", "submitBox", errBox, spinner, back, submit)

	multiplier := form.NewPart("text", "Your Selected Stake: ", form.ValidateMultiplier(inputs.maxMultiplier))
	multiplier.SetValue(strconv.FormatUint(inputs.multiplier, 10))
	multiplier.SetInputAttribute("class", "multiplier")
	multiplier.SetInputAttribute("step", "1")
	multiplier.SetInputAttribute("min", "0")
	multiplier.SetInputAttribute("max", strconv.FormatUint(inputs.maxMultiplier, 10))
	multiplier.SetInputAttribute("style", "margin-left:-0.75em;padding-right:1em;")
	multiplier.SetAttribute("style", "margin:1em 0 2em;text-align:center;")
	multiplier.SetLabelAttribute("style", "display:block;")
	multiplier.Disable()

	xx := bootstrap.NewElement("p", "inputXX", gowd.NewText("xx"))
	xx.SetAttribute("style", "margin-left:-1.5em;")
	multiplier.AddElement(xx)
	multiplier.SetHelpTxtAttribute("style", "display:table;")
	multiplier.SwapKids(2, 3)

	multiplier.SetInputAttribute("id", "number")
	multiplier.SetInputAttribute("onkeyup", "changeRangeValue(this.value, "+strconv.FormatUint(inputs.maxMultiplier, 10)+")")
	multiplier.SetInputAttribute("onclick", "changeRangeValue(this.value, "+strconv.FormatUint(inputs.maxMultiplier, 10)+")")

	formErrors := bootstrap.NewElement("p", "formErrors")
	formErrors.Hidden = true

	divWell := bootstrap.NewElement("div", "well")

	back.OnEvent(gowd.OnClick, func(sender *gowd.Element, event *gowd.EventElement) {
		body.RemoveElements()
		body.AddElement(page3(inputs))
	})

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

		if validated, valid := agreeInput.Validate(); valid {
			inputs.agree = validated.(bool)
		} else {
			errs++
		}

		jww.INFO.Printf("Inputs set: %+v", inputs)

		if errs == 0 {
			if err := body.Render(); err != nil {
				jww.ERROR.Print(err)
			}

			spinner.Hidden = true

			signAndTransmitTest := func(keyPath, idfPath, nominatorWallet,
				validatorWallet, serverAddress, serverCert, contract,
				email string, selectedMultiplier float32) error {
				if test {
					time.Sleep(500 * time.Millisecond)
					return nil
				} else {
					return client.SignAndTransmit(keyPath, idfPath,
						nominatorWallet, validatorWallet, serverAddress,
						serverCert, contract, email, selectedMultiplier)
				}
			}

			err := signAndTransmitTest(
				inputs.keyPath,
				inputs.idfPath,
				inputs.nominatorWallet,
				inputs.validatorWallet,
				serverAddress,
				string(inputs.cert),
				utils.Contract,
				inputs.email,
				float32(inputs.multiplier),
			)

			if err != nil {
				jww.ERROR.Printf("Submit error: %+v", err)
				errBox.SetText("An error occurred when submitting the request. Please contact support at nodes@xx.network and provide the following error message:")
				errBox.AddElement(bootstrap.NewElement("span", "errorBoxMessage", gowd.NewText(err.Error())))
				errBox.Hidden = false
				formErrors.SetText("The were errors in the form input. Please correct them to continue.")
				formErrors.Hidden = false
			} else {
				divWell.RemoveElements()
				success := bootstrap.NewElement("span", "success", gowd.NewText("Successful submitted commitment."))
				// 				// result := gowd.NewText(fmt.Sprintf("%+v", inputs))
				// 				result2 := gowd.NewElement("result2")
				// 				_, _ = result2.AddHTML(`
				// <table style=" font-family: "Roboto Mono", monospace;">
				//   <tr>
				//     <td><strong>keyPath</strong></td>
				//     <td>`+inputs.keyPath+`</td>
				//   </tr>
				//   <tr>
				//     <td><strong>idfPath</strong></td>
				//     <td>`+inputs.idfPath+`</td>
				//   </tr>
				//   <tr>
				//     <td><strong>nominatorWallet</strong></td>
				//     <td>`+inputs.nominatorWallet+`</td>
				//   </tr>
				//   <tr>
				//     <td><strong>validatorWallet</strong></td>
				//     <td>`+inputs.validatorWallet+`</td>
				//   </tr>
				//   <tr>
				//     <td><strong>serverAddress</strong></td>
				//     <td>`+serverAddress+`</td>
				//   </tr>
				//   <tr>
				//     <td><strong>cert</strong></td>
				//     <td>`+string(inputs.cert)+`</td>
				//   </tr>
				//   <tr>
				//     <td><strong>email</strong></td>
				//     <td>`+inputs.email+`</td>
				//   </tr>
				//   <tr>
				//     <td><strong>multiplier</strong></td>
				//     <td>`+strconv.FormatUint(inputs.multiplier, 10)+`</td>
				//   </tr>
				// </table>`, nil)
				divWell.AddElement(success)
				// divWell.AddElement(result2)
			}
		} else {
			formErrors.SetText("There were errors in the form input. Please correct them to continue.")
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
		jww.FATAL.Panic(err)
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
	contract.AddElement(agreeInput.Element())

	formGrp := bootstrap.NewFormGroup(
		formErrors,
		contract,
		multiplier.Element(),
		submitBox,
	)

	formGrp.SetAttribute("style", "margin: 2.5em 0 0")

	h1 := bootstrap.NewElement("h1", "")
	h1.SetText("Update Team Stake")
	logo := bootstrap.NewElement("img", "logo")
	logo.SetAttribute("src", "img/xx-logo.svg")
	h1.AddElement(logo)
	p := bootstrap.NewElement("p", "blurb")
	p.AddHTML(blurbTextPg3, nil)
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
