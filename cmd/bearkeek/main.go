package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/kudrykv/bearkeek"
	"github.com/kudrykv/bearkeek/alfred"
)

func main() {
	bear, err := bearkeek.NewDefault()
	if err != nil {
		fmt.Println(err)

		return
	}

	flag.Parse()

	a := alfred.New()

	search := ""
	if len(flag.Args()) > 0 {
		search = flag.Args()[len(flag.Args())-1]
	}

	parse := bearkeek.Parse(search)

	if parse.IsTagLast {
		tags, err := bear.Tags(context.Background(), bearkeek.TagsQuery{Term: parse.LastTag})
		if err != nil {
			fmt.Println(err)
		}

		for _, tag := range tags {
			tagname := "#" + tag.Name
			item := alfred.
				NewItem(tagname, "").
				Opts(
					alfred.Autocomplete(parse.RawButTag+tagname+" "),
					alfred.IsValid(false),
				)
			a.AddItem(item)
		}

		fmt.Println(string(a.MustJSON()))

		return
	}

	notes, err := bear.Notes(context.Background(), bearkeek.NotesQuery{
		Tags:  parse.Tags,
		Terms: parse.Terms,
		Limit: 100,
	})
	if err != nil {
		fmt.Println(err)

		return
	}

	for _, note := range notes {
		a.AddItem(alfred.NewItem(note.Title, note.Subtitle))
	}

	fmt.Println(string(a.MustJSON()))
}
