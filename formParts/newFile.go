package formParts

import (
	jww "github.com/spf13/jwalterweatherman"
	utils2 "gitlab.com/xx_network/primitives/utils"
)

func SaveHTML(title, name, content string) {
	head := `<!DOCTYPE html>
<html lang="en">

<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<title>` + title + `</title>
</head><body>`

	foot := "</body></html>"

	page := head + content + foot

	err := utils2.WriteFileDef(name, []byte(page))
	if err != nil {
		jww.ERROR.Print(err)
	}
}
