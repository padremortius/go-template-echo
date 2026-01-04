package main

import (
	"github.com/padremortius/go-template-echo/internal/app"
	"github.com/padremortius/go-template-echo/internal/config"
)

var (
	aBuildNumber    = ""
	aBuildTimeStamp = ""
	aGitBranch      = ""
	aGitHash        = ""
)

// @title go-template-echo
// @version 1.0
// @description This is a example api-server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Example
// @contact.url http://misko.su/support
// @contact.email support@misko.su

// @license.name MIT

func main() {
	ver := *config.InitVersion(aBuildNumber, aBuildTimeStamp, aGitBranch, aGitHash)
	app.Run(ver)
}
