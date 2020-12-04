package alfred

type Output struct {
	Rerun     float64           `json:"rerun,omitempty"`
	Variables map[string]string `json:"variables,omitempty"`
	Items     []Item            `json:"items"`
}

type Icon struct {
	Type string `json:"type,omitempty"`
	Path string `json:"path,omitempty"`
}

type Mod struct {
	Variables map[string]string `json:"variables,omitempty"`
	Valid     *bool             `json:"valid,omitempty"`
	Arg       string            `json:"arg,omitempty"`
	Subtitle  string            `json:"subtitle,omitempty"`
	Icon      *Icon             `json:"icon,omitempty"`
}

type Text struct {
	Copy      string `json:"copy,omitempty"`
	Largetype string `json:"largetype,omitempty"`
}

// Item specifies the members of the "items" array.
type Item struct {
	Variables    map[string]string `json:"variables,omitempty"`
	UID          string            `json:"uid,omitempty"`
	Title        string            `json:"title"`
	Subtitle     string            `json:"subtitle,omitempty"`
	Arg          string            `json:"arg,omitempty"`
	Icon         *Icon             `json:"icon,omitempty"`
	Autocomplete string            `json:"autocomplete,omitempty"`
	Type         string            `json:"type,omitempty"`
	Valid        *bool             `json:"valid,omitempty"`
	Match        string            `json:"match,omitempty"`
	Mods         map[string]Mod    `json:"mods,omitempty"`
	Text         *Text             `json:"text,omitempty"`
	QuicklookURL string            `json:"quicklookurl,omitempty"`
}

func NewItem(title, subtitle string) Item {
	return Item{
		Title:    title,
		Subtitle: subtitle,
	}
}

type Option func(Item) Item

func (r Item) Opts(opts ...Option) Item {
	item := r
	for _, opt := range opts {
		item = opt(item)
	}

	return item
}

func Autocomplete(autocomplete string) Option {
	return func(item Item) Item {
		item.Autocomplete = autocomplete
		return item
	}
}

func Arg(arg string) Option {
	return func(item Item) Item {
		item.Arg = arg
		return item
	}
}

func UID(uid string) Option {
	return func(item Item) Item {
		item.UID = uid
		return item
	}
}

func Valid(flag bool) Option {
	return func(item Item) Item {
		ptr := &flag
		item.Valid = ptr

		return item
	}
}
