package main

import (
	"errors"
	"flag"

	"github.com/kudrykv/bearkeek"
	"github.com/kudrykv/bearkeek/alfred"
)

const (
	modeSearch = "search"
	modeXCb    = "xcb"

	limit400 = 400

	openNewWindow = "yes"
	openEdit      = "yes"
)

// nolint:gochecknoglobals
var (
	mode  string
	limit = 400

	noteNewWindow string
	noteEdit      string

	errUnknownMode = errors.New("unknown mode")
)

func init() { // nolint:gochecknoinits
	flag.StringVar(&mode, "mode", modeSearch, "search")
	flag.IntVar(&limit, "limit", limit400, "<num>")

	flag.StringVar(&noteNewWindow, "new-window", openNewWindow, "yes|no")
	flag.StringVar(&noteEdit, "edit", openEdit, "yes|no")
}

func main() {
	bear, err := bearkeek.NewDefault()
	if erroredToAlfred(err) {
		return
	}

	flag.Parse()

	a := alfred.New()

	switch mode {
	case modeSearch:
		doSearch(bear, a)
	case modeXCb:
		doPrepXCallback()
	default:
		erroredToAlfred(errUnknownMode)
	}
}
