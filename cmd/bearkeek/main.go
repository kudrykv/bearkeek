package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/kudrykv/bearkeek"
)

func main() {
	bear, err := bearkeek.NewDefault()
	if err != nil {
		fmt.Println(err)

		return
	}

	flag.Parse()

	search := ""
	if len(flag.Args()) > 0 {
		search = flag.Args()[len(flag.Args())-1]
	}

	parse := bearkeek.Parse(search)

	fmt.Println(parse)

	if parse.IsTagLast {
		tags, err := bear.Tags(context.Background(), bearkeek.TagsQuery{Term: parse.LastTag})
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(tags)

		return
	}

	bear.Notes(context.Background(), bearkeek.NotesQuery{
		Tags:  parse.Tags,
		Terms: parse.Terms,
	})
}
