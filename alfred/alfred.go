package alfred

import (
	"encoding/json"
	"fmt"
)

type Alfred struct {
	prep Output
}

func New() *Alfred {
	return &Alfred{}
}

func (a *Alfred) AddItems(items []Item) *Alfred {
	a.prep.Items = append(a.prep.Items, items...)
	return a
}

func (a *Alfred) AddItem(item Item) *Alfred {
	a.prep.Items = append(a.prep.Items, item)
	return a
}

func (a *Alfred) JSON() ([]byte, error) {
	bts, err := json.Marshal(&a.prep)
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}

	return bts, nil
}

func (a *Alfred) MustJSON() []byte {
	bts, err := a.JSON()
	if err != nil {
		panic(fmt.Errorf("json: %w", err))
	}

	return bts
}
