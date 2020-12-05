package main

import (
	"flag"
	"fmt"

	"github.com/kudrykv/bearkeek/alfred"
)

func getSearchTerm() string {
	if len(flag.Args()) > 0 {
		return flag.Args()[len(flag.Args())-1]
	}

	return ""
}

func erroredToAlfred(err error) bool {
	if err == nil {
		return false
	}

	item := alfred.NewItem(err.Error(), "").Opts(alfred.Valid(false))
	fmt.Println(string(alfred.New().AddItem(item).MustJSON()))

	return true
}
