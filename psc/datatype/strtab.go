package datatype

import "fmt"

type StrTabEntry struct {
	Length    int
	String    string
	Reference int
}

type StrTab []StrTabEntry

func (t StrTab) FindString(value string) (int, *StrTabEntry) {
	for i, v := range t {
		if v.String == value {
			return i, &v
		}
	}

	return -1, nil
}

func (t StrTab) String() string {
	if t == nil || len(t) == 0 {
		return "<empty string table>"
	}

	const maxStringWidth = 30

	trunc := func(s string, max int) string {
		if len(s) <= max {
			return s
		}
		return s[:max-1] + "…"
	}

	out := "Idx  Len   String\n"
	out += "---- ----- --------------------------------\n"

	for i, e := range t {
		s := trunc(e.String, maxStringWidth)

		out += fmt.Sprintf(
			"%-4d %-5d %-30s\n",
			i,
			e.Length,
			s,
		)
	}

	return out
}
