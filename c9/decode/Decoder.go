package main

import (
	. "avc20/c9"
	"bufio"
	"fmt"
	"os"
	"strconv"
)

const (
	path           = "./c9/decode/transmission.txt"
	preambleLength = 25
)

func main() {
	lines, err := loadTransmission(path)
	if err != nil {
		fmt.Println("Error parsing file: ", err)
		return
	}

	trn := Transmission{Data: *lines}
	prm := PreambleTransmission{Transmission: trn, Length: preambleLength}
	searchedValue := 0
	for prm.ReadNext() {
		if !prm.IsValid() {
			fmt.Println("Non matching: ", prm.Current())
			searchedValue = prm.Current()
			break
		}
	}

	con := ContiguousTransmission{Transmission: trn, SearchedValue: searchedValue}
	for con.ReadNext() {
		found, values, min, max := con.Find()
		if found {
			fmt.Println("Discovered range: ", values)
			fmt.Println("Found pair: ", min, max)
			fmt.Println("Pair summary: ", min+max)
			break
		}
	}
}

func loadTransmission(path string) (*[]int, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	lines := make([]int, 0, 1000)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return nil, err
		}

		lines = append(lines, line)
	}

	return &lines, scanner.Err()
}
