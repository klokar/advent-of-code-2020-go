package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

const fileName = "./2/passwords.txt"

type entry struct {
	num1     int
	num2     int
	letter   string
	password string
}

func main() {
	entries, error := readFile(fileName)
	if error != nil {
		fmt.Println("File could not be opened! -> ", error)
		return
	}

	validCount := 0
	for _, entry := range entries {
		if entry.isValidForSecondPolicy() {
			validCount += 1
		}
	}

	fmt.Println(fmt.Sprintf("Answer: Valid number: %d", validCount))
}

func (e entry) isValidForFirstPolicy() bool {
	matches := regexp.MustCompile(e.letter).FindAllStringIndex(e.password, -1)

	return len(matches) >= e.num1 && len(matches) <= e.num2
}

func (e entry) isValidForSecondPolicy() bool {
	matches1 := string(e.password[e.num1-1]) == e.letter
	matches2 := string(e.password[e.num2-1]) == e.letter

	return matches1 != matches2
}

func readFile(path string) ([]entry, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []entry
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, parseLine(scanner.Text()))
	}

	return lines, scanner.Err()
}

func parseLine(line string) entry {
	num1Regex, _ := regexp.Compile("([0-9]{1,2})-")
	stringNum1 := num1Regex.FindStringSubmatch(line)
	num2Regex, _ := regexp.Compile("-([0-9]{1,2})")
	stringNum2 := num2Regex.FindStringSubmatch(line)
	letterRegex, _ := regexp.Compile(" ([a-z]):")
	letter := letterRegex.FindStringSubmatch(line)
	passwordRegex, _ := regexp.Compile(": (.*[a-z])")
	password := passwordRegex.FindStringSubmatch(line)

	num1, err := strconv.Atoi(stringNum1[1])
	if err != nil {
		fmt.Println("Error parsing number - " + stringNum1[1])
	}

	num2, err := strconv.Atoi(stringNum2[1])
	if err != nil {
		fmt.Println("Error parsing number - " + stringNum2[1])
	}

	return entry{num1, num2, letter[1], password[1]}
}
