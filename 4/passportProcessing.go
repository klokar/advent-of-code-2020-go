package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const path = "./4/passports.txt"

type passport struct {
	id             string
	countryId      string
	birthYear      int
	issueYear      int
	expirationYear int
	height         string
	hairColor      string
	eyeColor       string
}

func main() {
	passports, _ := parsePassports(path)

	validCount := 0
	for _, passport := range *passports {
		if passport.isValid() {
			validCount += 1
		}
	}

	fmt.Println(validCount)
}

func parsePassports(path string) (*[]passport, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	passports := make([]passport, 0)
	scanner := bufio.NewScanner(file)

	psp := passport{}
	for scanner.Scan() {
		if scanner.Text() == "" {
			passports = append(passports, psp)
			psp = passport{}
		} else {
			entries := strings.Fields(scanner.Text())
			for _, entry := range entries {
				data := strings.Split(entry, ":")
				psp.addData(data)
			}
		}
	}

	passports = append(passports, psp)

	return &passports, scanner.Err()
}

func (psp *passport) addData(data []string) {
	switch data[0] {
	case "byr":
		psp.birthYear, _ = strconv.Atoi(data[1])
	case "iyr":
		psp.issueYear, _ = strconv.Atoi(data[1])
	case "eyr":
		psp.expirationYear, _ = strconv.Atoi(data[1])
	case "hgt":
		psp.height = data[1]
	case "hcl":
		psp.hairColor = data[1]
	case "ecl":
		psp.eyeColor = data[1]
	case "pid":
		psp.id = data[1]
	case "cid":
		psp.countryId = data[1]
	}
}

func (psp passport) isValid() bool {
	return psp.birthYear >= 1920 && psp.birthYear <= 2002 &&
		psp.issueYear >= 2010 && psp.issueYear <= 2020 &&
		psp.expirationYear >= 2020 && psp.expirationYear <= 2030 &&
		psp.hasValidHeight() &&
		psp.hasValidHairColor() &&
		psp.hasValidEyeColor() &&
		psp.hasValidId()
}

func (psp passport) hasValidHeight() bool {
	regex, _ := regexp.Compile("([0-9]{2,3})(in|cm)")
	split := regex.FindStringSubmatch(psp.height)

	if len(split) != 3 {
		return false
	}

	height, _ := strconv.Atoi(split[1])

	if split[2] == "in" {
		return height >= 59 && height <= 76
	}

	return height >= 150 && height <= 193
}

func (psp passport) hasValidHairColor() bool {
	matches, _ := regexp.MatchString("^#[0-9a-z]{6}\\b", psp.hairColor)

	return matches
}

func (psp passport) hasValidEyeColor() bool {
	switch psp.eyeColor {
	case "amb", "blu", "brn", "gry", "grn", "hzl", "oth":
		return true
	}

	return false
}

func (psp passport) hasValidId() bool {
	matches, _ := regexp.MatchString("^[0-9]{9}\\b", psp.id)

	return matches
}
