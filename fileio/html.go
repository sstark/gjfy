package fileio

import "embed"

var (
	//go:embed *.tmpl
	HtmlTemplates embed.FS
	//go:embed favicon.ico
	Favicon []byte
	//go:embed gjfy-logo-small.png
	GjfyLogoSmall []byte
)

const (
	UserMessageViewDefaultText = `
		The link you invoked contains a secret (a password for example)
		somebody wants to share with you. It will be valid only for a short
		time and you may not be able to invoke it again. Please make sure
		you memorise the secret or write it down in an appropriate way.
		`
	UserMessageViewFilename = "userMessageView.txt"
	CssFileName  = "custom.css"
	LogoFileName = "logo.png"
)
