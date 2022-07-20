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

const test = true

const version = "version: 7/20/2022"

var body *gowd.Element

const (
	blurbTextPg1 = `<p>This applet allows you to commit to and configure your team stake from the MainNet transition program.</p>
<p>To learn more about the MainNet transition program, as well as its rules and details on its implementation, refer to&nbsp;<a onclick="window.nw.Shell.openExternal('https://xx.network/archive/mainnet-transition/')">About MainNet Transition</a>.</p>
<p>This applet will allow you to commit your wallets, select how much team stake you would like on your node, and sign the contracts required for the MainNet transition program. For more information, please see the&nbsp;<a onclick="window.nw.Shell.openExternal('https://xx.network/mainnet-transition-program-configuration-and-commitment-applet-instructions/')">instructions page</a>.</p>
<p>To confirm your ownership of your node, use your node keys from BetaNet. Please enter these keys below. They will be used locally to sign your commitment and will not leave this machine.</p>`
	blurbTextPg2 = `Below are the committed validator and nominator addresses. Select the checkbox to modify them.`
	blurbTextPg3 = `Use the following field to select the amount of team stake you would like to receive. You may receive up to a maximum determined by both the network stake and how much of your BetaNet rewards you currently have staked. Optionally, you will want to stake the minimum amount that will keep you in the active set to maximize your and the network as a whole's rewards.`
)
const serverAddress = "https://18.185.229.39:11420"

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

	origMultiplier int
	multiplier     int
	maxMultiplier  int

	walletModifyCheck     bool
	multiplierModifyCheck bool
	autoCalcCheck         bool
	selectStateCheck      bool
}

func buildPage() error {

	// creates a new bootstrap fluid container
	body = bootstrap.NewContainer(false)
	if test {
		testMode := bootstrap.NewElement("span", "", gowd.NewText("TEST MODE"))
		testMode.SetAttribute("style", "position:absolute;top:0;left:0;background:red;color:#fff")
		body.AddElement(testMode)
	}

	row := page1(Inputs{})
	row = page3(Inputs{
		certPath:              `C:\Users\Jono\Go\src\git.xx.network\elixxir\mainnet-commitments\client\test\server.crt`,
		keyPath:               `C:\Users\Jono\Go\src\git.xx.network\elixxir\mainnet-commitments\client\test\commitmenttestkey.key`,
		idfPath:               `C:\Users\Jono\Go\src\git.xx.network\elixxir\mainnet-commitments\client\test\testidf.json`,
		nodeID:                "",
		nominatorWallet:       "6a1TiUWcjderApE4876zGH5hbxjTbFV8sAb7sE3Tx2FfEGJt",
		validatorWallet:       "6WtKWycWig29uFfMN1PgGGEuzwzdurvhbf7Qxu9g5RuXQ7cM",
		origNominatorWallet:   "6a1TiUWcjderApE4876zGH5hbxjTbFV8sAb7sE3Tx2FfEGJt",
		origValidatorWallet:   "6WtKWycWig29uFfMN1PgGGEuzwzdurvhbf7Qxu9g5RuXQ7cM",
		agree:                 false,
		email:                 "johndoe@example.com",
		origMultiplier:        -6,
		multiplier:            0,
		maxMultiplier:         8425465,
		walletModifyCheck:     false,
		multiplierModifyCheck: false,
	})

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
			    "selected-stake": 457,
			    "max-stake": 1500,
				"email": "johnDoe@email.com"
			}`), nil
				} else {
					return client.GetInfo(nid, serverCert, serverAddress)
				}
			}

			jsonData, err := getInfoTest(inputs.nodeIdHex, string(inputs.cert), serverAddress)

			jww.DEBUG.Printf("JSON: %s", jsonData)

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
					SelectedStake   int    `json:"selected-stake"`
					MaxStake        int    `json:"max-stake"`
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
	h1.SetText("MainNet Transition Program Configuration and Commitment")
	h1.SetAttribute("style", "padding-right: 50px;")
	logo := bootstrap.NewElement("img", "logo")
	logo.SetAttribute("src", "img/xx-logo.svg")
	h1.AddElement(logo)
	p := bootstrap.NewElement("p", "blurb")
	p.AddHTML(blurbTextPg1, nil)
	divWell.AddElement(h1)
	divWell.AddElement(p)
	divWell.AddElement(formGrp)
	divWell.AddElement(bootstrap.NewElement("span", "version", gowd.NewText(version)))
	row := bootstrap.NewRow(divWell)

	return row
}

func page2(inputs Inputs) *gowd.Element {

	emailInput := form.NewPart("email", "Email to receive notification on changes to the state of your validator (optional)", form.ValidateEmail)
	emailInput.SetValue(inputs.email)

	validatorWallet := form.NewPart("text", "Validator Wallet Address", form.ValidateXXNetworkAddress)
	validatorWallet.SetValue(inputs.validatorWallet)
	validatorWallet.Disable()
	nominatorWallet := form.NewPart("text", "Nominator Wallet Address", form.ValidateXXNetworkAddressNotRequired)
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
	h1.SetText("Wallet Commitment")
	logo := bootstrap.NewElement("img", "logo")
	logo.SetAttribute("src", "img/xx-logo.svg")
	h1.AddElement(logo)
	p := bootstrap.NewElement("p", "blurb")
	p.AddHTML(blurbTextPg2, nil)
	divWell.AddElement(h1)
	divWell.AddElement(p)
	divWell.AddElement(formGrp)
	divWell.AddElement(bootstrap.NewElement("span", "version", gowd.NewText(version)))
	row := bootstrap.NewRow(divWell)

	return row
}

func page3(inputs Inputs) *gowd.Element {

	autoCalcRadio := form.NewPart("radio", "Auto Calculate Stake each Era", nil)
	autoCalcRadio.SetAttribute("style", "margin-left:2em;")
	autoCalcRadio.SetInputAttribute("name", "stake")
	autoCalcRadio.SetInputAttribute("value", "1")
	autoCalcRadio.SetLabelAttribute("style", "font-weight:bold")
	moreInfo := bootstrap.NewElement("img", "moreInfo")
	moreInfo.SetAttribute("src", "img/Info_icon.svg")
	autoCalcRadio.GetKid(0).AddElement(moreInfo)
	autoCalcRadio.Disable()

	autoCalcWarn := bootstrap.NewElement("p", "warn")
	_, _ = autoCalcWarn.AddHTML(`<img src="img/Warning_icon.svg"/><strong>Warning:</strong>&nbsp;For team stake auto calculate to ensure your node gets into the active set, you must also nominate your own node. This is your responsibility as a node operator.`, nil)
	autoCalcWarn.Hide()
	autoCalcNote := bootstrap.NewElement("p", "note")
	_, _ = autoCalcNote.AddHTML(`<img src="img/Note_icon.svg"/><strong>Note:</strong>&nbsp;For more information about team stake auto calculate, refer to&nbsp;<a onclick="window.nw.Shell.openExternal('https://xx.network/archive/mainnet-transition/')">ARTICLE</a>.`, nil)
	autoCalcNote.Hide()
	autoCalcRadio.AddElement(autoCalcWarn)
	autoCalcRadio.AddElement(autoCalcNote)

	popupShadow := bootstrap.NewElement("div", "popupShadow")
	popupText := bootstrap.NewElement("div", "popupText")
	popupBox := bootstrap.NewElement("div", "popupBox", popupText)
	popupContainer := bootstrap.NewElement("div", "popupContainer", popupShadow, popupBox)
	autoCalcWarnP := bootstrap.NewElement("p", "warn")
	_, _ = autoCalcWarnP.AddHTML(`<img src="img/Warning_icon.svg"/><strong>Warning:</strong>&nbsp;For team stake auto calculate to ensure your node gets into the active set, you must also nominate your own node. This is your responsibility as a node operator.`, nil)
	autoCalcNoteP := bootstrap.NewElement("p", "note")
	_, _ = autoCalcNoteP.AddHTML(`<img src="img/Note_icon.svg"/><strong>Note:</strong>&nbsp;For more information about team stake auto calculate, refer to&nbsp;<a onclick="window.nw.Shell.openExternal('https://xx.network/archive/mainnet-transition/')">ARTICLE</a>.`, nil)
	clearBtn := bootstrap.NewElement("img", "clearBtn")
	clearBtn.SetAttribute("src", "img/Clear_icon.svg")
	popupText.AddElement(clearBtn)
	popupText.AddElement(autoCalcWarnP)
	popupText.AddElement(autoCalcNoteP)
	popupContainer.Hide()
	body.AddElement(popupContainer)

	moreInfo.OnEvent(gowd.OnClick, func(*gowd.Element, *gowd.EventElement) {
		popupContainer.Show()
	})

	clearBtn.OnEvent(gowd.OnClick, func(*gowd.Element, *gowd.EventElement) {
		popupContainer.Hide()
	})

	popupShadow.OnEvent(gowd.OnClick, func(*gowd.Element, *gowd.EventElement) {
		popupContainer.Hide()
	})

	selectRadio := form.NewPart("radio", "Select Stake", nil)
	selectRadio.SetAttribute("style", "margin-left:2em;")
	selectRadio.SetInputAttribute("name", "stake")
	selectRadio.SetInputAttribute("value", "1")
	selectRadio.SetLabelAttribute("style", "font-weight:bold")
	selectRadio.Disable()

	selectText := bootstrap.NewElement("p", "selectText", gowd.NewText("Select a stake up to a max of "+strconv.Itoa(inputs.maxMultiplier)+" xx."))
	selectRadio.AddElement(selectText)
	selectText.Hide()

	radioBox := bootstrap.NewElement("div", "radioBox", autoCalcRadio.Element(), selectRadio.Element())

	multiplier := form.NewPart("number", "", form.ValidateMultiplier(inputs.maxMultiplier))
	multiplier.SetAttribute("style", "margin-left:2em;")
	multiplier.SetValue("0")
	multiplier.SetInputAttribute("class", "multiplier modifier")
	multiplier.SetInputAttribute("step", "1")
	multiplier.SetInputAttribute("min", "0")
	multiplier.SetInputAttribute("max", strconv.Itoa(inputs.maxMultiplier))
	multiplier.SetLabelAttribute("for", "number")
	multiplier.SetLabelAttribute("id", "numberLabel")
	multiplier.Disable()

	multiplier.AddElement(bootstrap.NewElement("p", "inputXX", gowd.NewText("xx")))
	multiplier.SetHelpTxtAttribute("style", "display:table;")
	multiplier.SwapKids(2, 3)

	multiplier.SetInputAttribute("id", "number")
	multiplier.SetInputAttribute("onkeyup", "changeRangeValue(this.value, "+strconv.Itoa(inputs.maxMultiplier)+")")
	multiplier.SetInputAttribute("onclick", "changeRangeValue(this.value, "+strconv.Itoa(inputs.maxMultiplier)+")")

	slider := bootstrap.NewElement("input", "stakeRange")
	slider.SetAttribute("type", "range")
	slider.SetAttribute("min", "0")
	slider.SetAttribute("max", strconv.Itoa(inputs.maxMultiplier))
	slider.SetValue("0")
	slider.SetAttribute("id", "range")
	slider.SetAttribute("oninput", "changeInputValue(this.value)")
	slider.Disable()
	multiplier.AddElement(slider)
	multiplier.SwapKids(3, 4)
	// gowd.ExecJS(`window.onload = updateRangeWidth()`)

	multiplier.Hide()
	selectRadio.OnEvent(gowd.OnChange, func(*gowd.Element, *gowd.EventElement) {
		multiplier.Show()
		selectText.Show()
		autoCalcWarn.Hide()
		autoCalcNote.Hide()

		if selectRadio.Checked() {
			inputs.selectStateCheck = true
			inputs.autoCalcCheck = false
		} else {
			inputs.selectStateCheck = false
			inputs.autoCalcCheck = true
		}
	})
	autoCalcRadio.OnEvent(gowd.OnChange, func(*gowd.Element, *gowd.EventElement) {
		multiplier.Hide()
		selectText.Hide()
		autoCalcWarn.Show()
		autoCalcNote.Show()
		if autoCalcRadio.Checked() {
			inputs.autoCalcCheck = true
			inputs.selectStateCheck = false
		} else {
			inputs.autoCalcCheck = false
			inputs.selectStateCheck = false
		}
	})
	selectRadio.AddElement(multiplier.Element())

	modifyCheck := form.NewPart("checkbox", "Modify the selected stake", nil)
	modifyCheck.SetLabelAttribute("style", "font-weight:bold")
	if inputs.multiplierModifyCheck {
		modifyCheck.Check()
		multiplier.Enable()
		slider.Enable()
		autoCalcRadio.Enable()
		selectRadio.Enable()

		if inputs.selectStateCheck {
			selectRadio.Check()
			multiplier.Show()
			selectText.Show()
			autoCalcWarn.Hide()
			autoCalcNote.Hide()
			slider.SetValue(strconv.Itoa(inputs.multiplier))
			multiplier.SetValue(strconv.Itoa(inputs.multiplier))
		} else if inputs.autoCalcCheck {
			autoCalcRadio.Check()
			selectText.Hide()
			autoCalcWarn.Show()
			autoCalcNote.Show()
		}

	} else {
		if inputs.origMultiplier == -2 {
			autoCalcRadio.Check()
			selectText.Hide()
			autoCalcWarn.Show()
			autoCalcNote.Show()

			inputs.selectStateCheck = false
			inputs.autoCalcCheck = true
		} else if inputs.origMultiplier == -3 {
			autoCalcRadio.Uncheck()
			selectRadio.Uncheck()
			selectText.Hide()
			inputs.selectStateCheck = false
			inputs.autoCalcCheck = false
		} else {
			selectRadio.Check()
			multiplier.Show()
			selectText.Show()
			autoCalcWarn.Hide()
			autoCalcNote.Hide()
			slider.SetValue(strconv.Itoa(inputs.origMultiplier))
			multiplier.SetValue(strconv.Itoa(inputs.origMultiplier))

			inputs.selectStateCheck = true
			inputs.autoCalcCheck = false
		}
	}
	modifyCheck.OnEvent(gowd.OnChange, func(*gowd.Element, *gowd.EventElement) {
		if modifyCheck.Checked() {
			inputs.multiplierModifyCheck = true
			multiplier.Enable()
			slider.Enable()
			autoCalcRadio.Enable()
			selectRadio.Enable()
		} else {
			inputs.multiplierModifyCheck = false

			jww.DEBUG.Printf("orig: %d", inputs.origMultiplier)

			multiplier.Disable()
			slider.Disable()
			autoCalcRadio.Disable()
			selectRadio.Disable()

			if inputs.origMultiplier == -2 {
				selectRadio.Uncheck()
				autoCalcRadio.Check()
				multiplier.Hide()
				selectText.Hide()
				autoCalcWarn.Show()
				autoCalcNote.Show()
			} else if inputs.origMultiplier == -3 {
				autoCalcRadio.Uncheck()
				selectRadio.Uncheck()
				inputs.selectStateCheck = false
				inputs.autoCalcCheck = false
			} else {
				selectRadio.Check()
				multiplier.Show()
				selectText.Show()
				autoCalcWarn.Hide()
				autoCalcNote.Hide()
				slider.SetValue(strconv.Itoa(inputs.origMultiplier))
				multiplier.SetValue(strconv.Itoa(inputs.origMultiplier))
			}

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
		inputs.multiplier, _ = strconv.Atoi(multiplier.GetValue())
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

		if autoCalcRadio.Checked() {
			inputs.multiplier = -2
		} else if selectRadio.Checked() {
			if validated, valid := multiplier.Validate(); valid {
				inputs.multiplier = validated.(int)
			} else {
				errs++
			}
		} else {
			modifyCheck.SetHelpText("A team stake must be selected.")
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

	checkboxText := bootstrap.NewElement("p", "", gowd.NewText("The applet has imported your previous team stake. Select the checkbox below to modify it. If you have not previously selected a team stake, you must do so before proceeding."))
	checkboxText.SetAttribute("style", "background: #eaf8fd;padding: 0.75em;margin: 1em -0.75em;")

	formGrp := bootstrap.NewFormGroup(
		formErrors,
		checkboxText,
		modifyCheck.Element(),
		radioBox,
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
	divWell.AddElement(h1)
	divWell.AddElement(p)
	divWell.AddElement(formGrp)
	row := bootstrap.NewRow(divWell)
	divWell.AddElement(bootstrap.NewElement("span", "version", gowd.NewText(version)))

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
	multiplier.SetInputAttribute("class", "multiplier")
	multiplier.SetAttribute("style", "margin:1em 0 2em;text-align:center;")
	multiplier.SetLabelAttribute("style", "display:block;")
	if inputs.multiplier == -2 {
		multiplier.SetValue("Auto Calculated")
		multiplier.SetInputAttribute("style", "text-align:center;font-family: \"Roboto\", \"Franklin Gothic Medium\", Tahoma, sans-serif;")
	} else {
		multiplier.SetValue(strconv.Itoa(inputs.multiplier))
		multiplier.SetInputAttribute("style", "margin-left:-0.75em;padding-right:1em;")
		xx := bootstrap.NewElement("p", "inputXX", gowd.NewText("xx"))
		xx.SetAttribute("style", "margin-left:-1.5em;")
		multiplier.AddElement(xx)
		multiplier.SetHelpTxtAttribute("style", "display:table;")
		multiplier.SwapKids(2, 3)
	}
	multiplier.Disable()

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
				email string, selectedMultiplier int) error {
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
				inputs.multiplier,
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
				success := bootstrap.NewElement("span", "success", gowd.NewText("Successfully submitted commitment."))
				// result := gowd.NewText(fmt.Sprintf("%+v", inputs))
				result2 := gowd.NewElement("result2")
				_, _ = result2.AddHTML(`
				<table style=" font-family: "Roboto Mono", monospace;">
				  <tr>
				    <td><strong>keyPath</strong></td>
				    <td>`+inputs.keyPath+`</td>
				  </tr>
				  <tr>
				    <td><strong>idfPath</strong></td>
				    <td>`+inputs.idfPath+`</td>
				  </tr>
				  <tr>
				    <td><strong>nominatorWallet</strong></td>
				    <td>`+inputs.nominatorWallet+`</td>
				  </tr>
				  <tr>
				    <td><strong>validatorWallet</strong></td>
				    <td>`+inputs.validatorWallet+`</td>
				  </tr>
				  <tr>
				    <td><strong>serverAddress</strong></td>
				    <td>`+serverAddress+`</td>
				  </tr>
				  <tr>
				    <td><strong>cert</strong></td>
				    <td>`+string(inputs.cert)+`</td>
				  </tr>
				  <tr>
				    <td><strong>email</strong></td>
				    <td>`+inputs.email+`</td>
				  </tr>
				  <tr>
				    <td><strong>multiplier</strong></td>
				    <td>`+strconv.Itoa(inputs.multiplier)+`</td>
				  </tr>
				</table>`, nil)
				divWell.AddElement(success)
				divWell.AddElement(result2)
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
	h1.SetText("MainNet Transition Contract")
	logo := bootstrap.NewElement("img", "logo")
	logo.SetAttribute("src", "img/xx-logo.svg")
	h1.AddElement(logo)
	p := bootstrap.NewElement("p", "blurb")
	divWell.AddElement(h1)
	divWell.AddElement(p)
	divWell.AddElement(formGrp)
	divWell.AddElement(bootstrap.NewElement("span", "version", gowd.NewText(version)))
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
