package main

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/kudrykv/bearkeek"
	"github.com/kudrykv/bearkeek/alfred"
)

func doSearch(bear bearkeek.Bear, a *alfred.Alfred) {
	search := getSearchTerm()
	parse := bearkeek.Parse(search)

	if parse.IsTagLast {
		listTags(bear, a, parse)

		return
	}

	listNotes(bear, a, parse, search)
}

func listNotes(bear bearkeek.Bear, a *alfred.Alfred, parse bearkeek.ParseResult, search string) {
	notes, err := bear.Notes(context.Background(), createNoteQuery(parse))
	if erroredToAlfred(err) {
		return
	}

	for _, note := range notes {
		a.AddItem(alfred.NewItem(note.Title, note.Subtitle).Opts(alfred.Arg(getArg(note, search))))
	}

	if len(notes) == 0 {
		a.AddItem(alfred.NewItem("Nothing was found", "").Opts(alfred.Valid(false)))
	}

	fmt.Println(string(a.MustJSON()))
}

func createNoteQuery(parse bearkeek.ParseResult) bearkeek.NotesQuery {
	return bearkeek.NotesQuery{Tags: parse.Tags, Terms: parse.Terms, Limit: limit}
}

func listTags(bear bearkeek.Bear, a *alfred.Alfred, parse bearkeek.ParseResult) {
	tags, err := bear.Tags(context.Background(), bearkeek.TagsQuery{Term: parse.LastTag})
	if erroredToAlfred(err) {
		return
	}

	for _, tag := range tags {
		tagname := "#" + tag.Name
		item := alfred.
			NewItem(tagname, "").
			Opts(
				alfred.UID(strconv.FormatInt(tag.ID, 10)),
				alfred.Autocomplete(parse.RawButTag+tagname+" "),
				alfred.Valid(false),
			)
		a.AddItem(item)
	}

	fmt.Println(string(a.MustJSON()))
}

func getArg(n bearkeek.Note, st string) string {
	return url.Values{
		"id":   {n.UUID},
		"term": {st},
	}.Encode()
}
