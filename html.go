package main

import (
//std necessities
)

func buildHTML() string {
	html := &HTML{}

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
	html.Input("server", "server")
	html.NewLine()
	html.FormSubmit("Run Minecraft")

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

func (html *HTML) Input(id, name string) {
	html.raw += "<input type=\"text\" id=\"" + id + "\" name=\"" + name + "\">"
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
