package stash

import "errors"

type StashFile struct {
	Stashes []Stash `json:"stashes"`
	Demo    bool    `json:"demo"`
}

func (h *StashFile) Validate() error {
	if len(h.Stashes) > 50 {
		return errors.New("Too many stashes")
	}
	return nil
}

type Stash struct {
	ID       string   `json:"id"`
	Lat      float64  `json:"lat"`
	Lon      float64  `json:"lon"`
	Location string   `json:"location"`
	Contents []string `json:"contents,omitempty"`
	Type     string   `json:"type"`
	W3W      string   `json:"w3w"`
	Points   int      `json:"points"`
	Hide     bool     `json:"hide"`
}
