package main

import (
	//std necessities
	"fmt"
)

var mainCSS = `
body {
	color: #FAFAFA;
	background-color: #1C1C1C;
}
`

func htmlLogin(loginErr error) string {
	html := &HTML{
		css: mainCSS,
	}

	html.FormStart("/login")
	html.Label("email", "Mojang Email: ")
	html.Input("text", "mcemail", "mcemail")
	html.NewLine()
	html.Label("password", "Mojang Password: ")
	html.Input("password", "mcpwd", "mcpwd")
	html.NewLine()
	html.FormSubmit("Login")
	html.FormEnd()

	if loginErr != nil {
		html.Label("errorMsg", "An error occurred logging into your Mojang account.")
		html.NewLine()
		html.Label("errorDetail", fmt.Sprintf("%v", loginErr))
	}

	return html.String()
}

func htmlLauncher(launchVer string) string {
	html := &HTML{
		css: mainCSS,
	}

	html.FormStart("/launch")
	html.Label("version", "Choose a version: ")
	html.SelectStart("version", "version")
	for i := 0; i < len(versions.Versions); i++ {
		version := versions.Versions[i]
		html.Option(version.ID, version.ID)
	}
	html.SelectEnd()
	html.NewLine()
	html.Label("server", "Enter direct connect (empty to ignore): ")
	html.Input("text", "server", "server")
	html.NewLine()
	html.FormSubmit("Run Minecraft")
	html.FormEnd()

	if launchVer != "" {
		html.Label("launchMsg", "Launching Minecraft "+launchVer+"...")
	}

	return html.String()
}

type HTML struct {
	css string
	js  string
	raw string
}

func (html *HTML) String() string {
	raw := "<html>"

	if html.css != "" || html.js != "" {
		raw += "<head><style>" + html.css + "</style></head>"
	}

	raw += "<body>" + html.raw + "<script>" + html.js + "</script></body></html>"

	return raw
}

func (html *HTML) NewLine() {
	html.raw += "<br />"
}

func (html *HTML) FormStart(action string) {
	html.raw += "<form action=\"" + action + "\">"
}

func (html *HTML) FormSubmit(value string) {
	html.raw += "<input type=\"submit\" value=\"" + value + "\">"
}

func (html *HTML) FormEnd() {
	html.raw += "</form>"
}

func (html *HTML) Label(id, msg string) {
	html.raw += "<label for=\"" + id + "\">" + msg + "</label>"
}

func (html *HTML) Input(inputType, id, name string) {
	html.raw += "<input type=\"" + inputType + "\" id=\"" + id + "\" name=\"" + name + "\">"
}

func (html *HTML) Button(call, msg string) {
	html.raw += "<button onclick=\"" + call + "()\">" + msg + "</button>"
}

func (html *HTML) SelectStart(name, id string) {
	html.raw += "<select name=\"" + name + "\" id=\"" + id + "\">"
}

func (html *HTML) SelectEnd() {
	html.raw += "</select>"
}

func (html *HTML) Option(value, msg string) {
	html.raw += "<option value=\"" + value + "\">" + msg + "</option>"
}
