package main

import (
	"avc20/c16"
	ticketCreator "avc20/c16/parser"
	"fmt"
	"math"
	"strings"
)

const (
	path            = "./c16/parser/findings.txt"
	ticketFields    = 20
	multiplyingWord = "departure"
)

func main() {
	rules, myTicket, tickets, _ := ticketCreator.Input{Path: path}.Parse()

	validTickets := validateTickets(tickets, rules)
	fmt.Println("Valid tickets:", len(validTickets))

	mappings := obtainMappings(&validTickets, rules)
	fmt.Println("Mappings:", mappings)

	outcome := multiplyDepartures(myTicket, &mappings)
	fmt.Println("Outcome:", outcome)
}

func validateTickets(tickets *[]c16.Ticket, fieldRanges *map[string]c16.TicketFieldRanges) []c16.Ticket {
	validTickets := make([]c16.Ticket, 0, 200)
	for _, ticket := range *tickets {
		if ticket.IsValid(fieldRanges) {
			validTickets = append(validTickets, ticket)
		}
	}

	return validTickets
}

func obtainMappings(validTickets *[]c16.Ticket, rules *map[string]c16.TicketFieldRanges) map[int]string {
	// Get possible field mappings from rules
	// Rule passes if all this field in all tickets match
	possibleMappings := make(map[string][]int)
	for name, rule := range *rules {
		// Check ticket field
		for field := 0; field < ticketFields; field++ {

			// Iterate over all tickets for this field and find non-matching
			matches := true
			for _, ticket := range *validTickets {
				if !rule.MatchesAny(ticket.Values[field]) {
					matches = false
				}
			}

			// If all match assign mapping
			if matches {
				possibleMappings[name] = append(possibleMappings[name], field)
			}
		}
	}

	return narrowMapping(&possibleMappings, make(map[int]string))
}

func narrowMapping(possibleMappings *map[string][]int, exactMapping map[int]string) map[int]string {
	// End where there are not more fields to check
	if len(*possibleMappings) == 0 {
		return exactMapping
	}

	// Set highest possible number for matches
	minMatches := math.MaxInt
	minName := ""

	// Get field that has the lowest possible mappings
	for rule, matches := range *possibleMappings {
		if len(matches) < minMatches {
			minMatches = len(matches)
			minName = rule
		}
	}

	// Assign first not taken field
	for _, field := range (*possibleMappings)[minName] {
		if exactMapping[field] == "" {
			exactMapping[field] = minName
			delete(*possibleMappings, minName)
			break
		}
	}

	return narrowMapping(possibleMappings, exactMapping)
}

func multiplyDepartures(myTicket *c16.Ticket, results *map[int]string) int {
	outcome := 1

	for field, name := range *results {
		if strings.Contains(name, multiplyingWord) {
			outcome *= myTicket.Values[field]
		}
	}

	return outcome
}
