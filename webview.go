package main

import (
	"github.com/webview/webview"

	//std necessities
	"net/url"
	"runtime"
)

type Webview struct {
	view  webview.WebView
	debug bool
}

func NewWebview(webURL string) (*Webview, error) {
	if _, err := url.ParseRequestURI(webURL); err != nil {
		return nil, err
	}

	w := &Webview{}
	if verbosity > 0 {
		w.debug = true
	}

	w.view = webview.New(w.debug)
	w.view.SetTitle("Minecraft Offline v0.0.0 (" + runtime.GOOS + "/" + runtime.GOARCH + ")")
	w.view.SetSize(800, 600, webview.HintNone)
	w.view.Navigate(webURL)

	return w, nil
}

func (w *Webview) Run() {
	if w.view == nil {
		return
	}

	w.view.Run()
}

func (w *Webview) Destroy() {
	if w.view == nil {
		return
	}

	w.view.Destroy()
	w.view = nil
}
