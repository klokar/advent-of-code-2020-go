package parser

import (
	. "avc20/c16"
	"bufio"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Input struct {
	Path string
}

func (inp Input) Parse() (*map[string]TicketFieldRanges, *Ticket, *[]Ticket, error) {
	file, err := os.Open(inp.Path)
	if err != nil {
		return nil, nil, nil, err
	}
	defer file.Close()

	rules := make(map[string]TicketFieldRanges)
	var myTicket Ticket
	tickets := make([]Ticket, 0, 300)
	scanner := bufio.NewScanner(file)

	typeSequence := 1
	for scanner.Scan() {
		switch scanner.Text() {
		case "":
		case "your ticket:":
			typeSequence = 2
		case "nearby tickets:":
			typeSequence = 3
		default:
			switch typeSequence {
			case 1:
				ruleParts := strings.Split(scanner.Text(), ":")
				regex, _ := regexp.Compile("^ ([0-9]{1,3})-([0-9]{1,3}) or ([0-9]{1,3})-([0-9]{1,3})$")
				split := regex.FindStringSubmatch(ruleParts[1])
				from1, _ := strconv.Atoi(split[1])
				to1, _ := strconv.Atoi(split[2])
				from2, _ := strconv.Atoi(split[3])
				to2, _ := strconv.Atoi(split[4])

				rules[ruleParts[0]] = TicketFieldRanges{from1, to1, from2, to2}
			case 2:
				myTicket = parseTicket(scanner.Text())
			case 3:
				tickets = append(tickets, parseTicket(scanner.Text()))
			}
		}
	}

	return &rules, &myTicket, &tickets, scanner.Err()
}

func parseTicket(text string) Ticket {
	fields := strings.Split(text, ",")

	var ticket Ticket
	for _, stringField := range fields {
		val, _ := strconv.Atoi(stringField)
		ticket.Add(val)
	}

	return ticket
}
