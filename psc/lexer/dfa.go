package lexer

type State int
type DFA struct {
	Start  State
	Finals map[State]string
	Trans  map[State]map[rune]State
	curr   State
}

func (d *DFA) Reset() {
	d.curr = d.Start
}

func (d *DFA) State() State {
	return d.curr
}

func (d *DFA) IsFinal(s State) (string, bool) {
	lab, ok := d.Finals[s]
	return lab, ok
}

func (d *DFA) Step(r rune) (State, bool) { // Buat ngecek next statenya apa
	next, ok := d.Trans[d.curr][r]
	return next, ok
}

func (d *DFA) Advance(r rune) bool { // Pindah State
	if next, ok := d.Trans[d.curr][r]; ok {
		d.curr = next
		return true
	}
	return false
}
