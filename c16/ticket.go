package c16

type Ticket struct {
	Values []int
}

func (t *Ticket) Add(val int) {
	t.Values = append(t.Values, val)
}

func (t Ticket) IsValid(fieldRanges *map[string]TicketFieldRanges) bool {
	for _, val := range t.Values {
		match := false
		for _, fieldRange := range *fieldRanges {
			if fieldRange.MatchesAny(val) {
				match = true
				break
			}
		}

		if !match {
			return false
		}
	}

	return true
}
