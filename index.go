package main

import (
	"github.com/dtylman/gowd"
	"github.com/dtylman/gowd/bootstrap"
	jww "github.com/spf13/jwalterweatherman"
)

func buildPage() {

	var body *gowd.Element

	// creates a new bootstrap fluid container
	body = bootstrap.NewContainer(false)

	// add some elements using the object model
	div := bootstrap.NewElement("div", "well")
	div.SetAttribute("style", "font-size:1.5em;margin-top:25px;")
	body.AddElement(div)

	logo := bootstrap.NewElement("img", "")
	logo.SetAttribute("src", "img/xx_logo.svg")
	logo.SetAttribute("style", "float:right;margin: -10px -10px 0 0;")
	logo.SetAttribute("id", "logo")
	div.AddElement(logo)

	// Start the ui loop
	err = gowd.Run(body)
	if err != nil {
		jww.ERROR.Printf("Failed to start ui loop: %+v", err)
	}
}
