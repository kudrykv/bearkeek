package main

import (
	"fmt"
	"net/url"
)

func doPrepXCallback() {
	query, err := url.ParseQuery(getSearchTerm())
	if erroredToAlfred(err) {
		return
	}

	var inNewWindow string
	if noteNewWindow == "yes" {
		inNewWindow = "&new_window=yes"
	}

	var inEdit string
	if noteEdit == "yes" {
		inEdit = "&edit=yes"
	}

	fmt.Println("bear://x-callback-url/open-note?id=" + query.Get("id") + inNewWindow + inEdit)
}
