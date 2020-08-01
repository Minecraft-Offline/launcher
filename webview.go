package main

import (
	"github.com/webview/webview"

	//std necessities
	"fmt"
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
	w.view.SetTitle(fmt.Sprintf("Minecraft Offline (%s-%s-%s-%s-%s)", GitBranch, GitState, GitCommit, runtime.GOOS, runtime.GOARCH))
	w.view.SetSize(650, 320, webview.HintNone)
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
