package c16

type TicketFieldRanges struct {
	From1 int
	To1   int
	From2 int
	To2   int
}

func (fr TicketFieldRanges) MatchesAny(val int) bool {
	return val >= fr.From1 && val <= fr.To1 || val >= fr.From2 && val <= fr.To2
}
