package dfa

import (
	"encoding/json"
	"fmt"
	"os"
)

type rawTransition struct {
	From  int    `json:"from"`
	Input string `json:"input"`
	To    int    `json:"to"`
}

type rawFinal struct {
	State  int    `json:"state"`
	Output string `json:"output"`
}

type rawSpec struct {
	States      int             `json:"states"`
	Start       int             `json:"start"`
	FinalArray  []rawFinal      `json:"final"`
	Transitions []rawTransition `json:"transitions"`
}

func LoadJSON(path string) (*DFA, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var spec rawSpec
	if err := json.Unmarshal(b, &spec); err != nil {
		return nil, err
	}

	d := &DFA{
		Start:  State(spec.Start),
		Finals: make(map[State]string),
		Trans:  make(map[State]map[rune]State),
	}

	for _, f := range spec.FinalArray {
		d.Finals[State(f.State)] = f.Output
	}

	for _, t := range spec.Transitions {
		runes := []rune(t.Input)
		if len(runes) != 1 {
			return nil, fmt.Errorf("input '%s' harus 1 rune", t.Input)
		}
		if d.Trans[State(t.From)] == nil {
			d.Trans[State(t.From)] = make(map[rune]State)
		}
		d.Trans[State(t.From)][runes[0]] = State(t.To)
	}

	d.Reset()
	return d, nil
}
