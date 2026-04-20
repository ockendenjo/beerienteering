package stash

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type StashFile struct {
	Stashes          []*Stash `json:"stashes"`
	Demo             bool     `json:"demo"`
	ShowInstructions bool     `json:"showInstructions,omitempty"`
}

func (h *StashFile) Validate() error {
	if len(h.Stashes) > 50 {
		return errors.New("too many stashes")
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

func (s *Stash) Validate() error {
	_, err := uuid.Parse(s.ID)
	if err != nil {
		return fmt.Errorf("invalid stash id %w", err)
	}
	if s.Lat > 90 || s.Lat < -90 {
		return errors.New("invalid stash latitude")
	}
	if s.Lon > 180 || s.Lon < -180 {
		return errors.New("invalid stash longitude")
	}
	if len(s.Location) > 100 {
		return errors.New("stash location too long")
	}
	if len(s.Contents) > 10 {
		return errors.New("stash contents has too many items")
	}
	for _, c := range s.Contents {
		if len(c) > 100 {
			return errors.New("stash content item is too long")
		}
	}
	if len(s.Type) > 20 {
		return errors.New("stash type is too long")
	}
	if len(s.W3W) > 100 {
		return errors.New("stash w3w is too long")
	}
	if s.Points < 0 {
		return errors.New("stash points must be positive")
	}
	return nil
}
