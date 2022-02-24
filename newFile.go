package main

import (
	jww "github.com/spf13/jwalterweatherman"
	utils2 "gitlab.com/xx_network/primitives/utils"
)

func saveHTML(title, name, content string) {
	head := `<!DOCTYPE html>
<html lang="en">

<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<link rel="stylesheet" type="text/css" href="css/bootstrap.min.css"/>
	<link rel="stylesheet" type="text/css" href="css/style.css"/>
	<style>body{padding:2em; background:#fff;}</style>
	<title>` + title + `</title>
</head><body>`

	foot := "</body></html>"

	page := head + content + foot

	err := utils2.WriteFileDef(name, []byte(page))
	if err != nil {
		jww.ERROR.Print(err)
	}
}
