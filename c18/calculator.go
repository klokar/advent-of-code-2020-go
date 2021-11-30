package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

const (
	path = "./c18/expressions.txt"
)

type calculation struct {
	Input  string
	first  expression
	levels int
}

type expression struct {
	Next  *expression
	Value string
	Level int
}

func main() {
	calculations, err := loadCalculations(path)
	if err != nil {
		fmt.Println("Error parsing file: ", err)
		return
	}

	result := 0
	for _, calculation := range *calculations {
		result += calculation.Calculate()
	}

	fmt.Println("Result: ", result)
}

func (c *calculation) Calculate() int {
	c.assembleParts()

	for level := c.levels; level >= 0; level-- {
		current := &c.first
		for current != nil {
			if current.Level == level {
				current.aggregateLevel()
			}
			current = current.Next
		}
	}

	if val, err := strconv.Atoi(c.first.Value); err == nil {
		return val
	}

	return 0
}

func (c calculation) display() {
	current := c.first
	expression := current.Value

	for current.Next != nil {
		current = *current.Next
		expression += current.Value
	}
}

func (e *expression) aggregateLevel() {
	numbers := make([]int, 0, 10)
	operators := make([]string, 0, 10)
	current := e

	for current != nil && current.Level == e.Level {
		if number, err := strconv.Atoi(current.Value); err == nil {
			numbers = append(numbers, number)
		} else {
			operators = append(operators, current.Value)
		}

		current = current.Next
	}

	fmt.Println(fmt.Sprintf("Calculating level %d: ", e.Level), numbers, operators)

	operators, numbers = aggregateParams("+", operators, numbers)
	fmt.Println("Calculated '+': ", numbers, operators)
	operators, numbers = aggregateParams("*", operators, numbers)
	fmt.Println("Calculated '*': ", numbers, operators)

	fmt.Println("Result for level: ", numbers[0])

	e.Value = strconv.Itoa(numbers[0])
	// Lower level as this one is calculated
	e.Level--
	// Set next to as one from lower level
	e.Next = current
}

func aggregateParams(selectedOperator string, operators []string, numbers []int) ([]string, []int) {
	for i := 0; i < len(operators); i++ {
		if operators[i] == selectedOperator {
			fmt.Println("Aggregating", i, operators[i], numbers)
			if selectedOperator == "+" {
				numbers[i] = numbers[i] + numbers[i+1]
			} else {
				numbers[i] = numbers[i] * numbers[i+1]
			}

			// Remove second number
			numbers = append(numbers[:i+1], numbers[i+2:]...)
			// Remove used operator
			operators = append(operators[:i], operators[i+1:]...)
			// Subtract i as 1 entry was removed from operators and numbers
			i--
		}
	}

	return operators, numbers
}

func (c *calculation) assembleParts() {
	level := 0
	maxLevel := 0
	var first, last *expression

	for _, rne := range c.Input {
		char := string(rne)
		switch char {
		case "(":
			level++
			if level > maxLevel {
				maxLevel = level
			}
		case " ":
			continue
		case ")":
			level--
		default:
			current := expression{
				Value: char,
				Level: level,
			}

			// Set pointer of previous expression to current
			if last != nil {
				last.Next = &current

				// Set first element if not set and if last has next
				if first == nil {
					first = last
				}
			}

			// Set current expression as last for next one
			last = &current
		}
	}

	c.first = *first
	c.levels = maxLevel
}

func loadCalculations(path string) (*[]calculation, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	calculations := make([]calculation, 0, 380)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		calculations = append(calculations, calculation{Input: scanner.Text()})
	}

	return &calculations, scanner.Err()
}
