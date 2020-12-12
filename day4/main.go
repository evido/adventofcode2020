package main

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

type Passport struct {
	Attributes map[string]string
}

func (passport *Passport) HasAttribute(name string) bool {
	_, ok := passport.Attributes[name]
	return ok
}

func validateYear(value string, from, to int) bool {
	year, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return false
	}

	return int(year) >= from && int(year) <= to
}

func validateBirthYear(value string) bool {
	return validateYear(value, 1920, 2002)
}

func validateIssueYear(value string) bool {
	return validateYear(value, 2010, 2020)
}

func validateExpirationYear(value string) bool {
	return validateYear(value, 2020, 2030)
}

func validateHeight(value string) bool {
	if len(value) < 3 {
		return false
	}

	unit := value[len(value)-2:]
	measure, err := strconv.ParseInt(value[:len(value)-2], 10, 32)

	if err != nil {
		return false
	}

	switch unit {
	case "cm":
		return measure >= 150 && measure <= 193
	case "in":
		return measure >= 59 && measure <= 76
	default:
		return false
	}

	return false
}

func validateHairColor(value string) bool {
	if len(value) != 7 || value[0] != '#' {
		return false
	}

	_, err := strconv.ParseInt(value[1:], 16, 32)
	return err == nil
}

func validateEyeColor(value string) bool {
	acceptedValues := []string{
		"amb", "blu", "brn", "gry", "grn", "hzl", "oth",
	}

	for _, acceptedValue := range acceptedValues {
		if acceptedValue == value {
			return true
		}
	}

	return false
}

func validatePassportId(value string) bool {
	if len(value) != 9 {
		return false
	}

	for _, c := range value {
		if c < '0' || c > '9' {
			return false
		}
	}

	return true
}

func (passport *Passport) Validate() bool {
	requiredAttributes := map[string]func(string) bool{
		"byr": validateBirthYear,
		"iyr": validateIssueYear,
		"eyr": validateExpirationYear,
		"hgt": validateHeight,
		"hcl": validateHairColor,
		"ecl": validateEyeColor,
		"pid": validatePassportId,
	}

	for requiredAttribute, validate := range requiredAttributes {
		if !passport.HasAttribute(requiredAttribute) ||
			!validate(passport.Attributes[requiredAttribute]) {
			return false
		}
	}

	return true
}

func countValidPassports(passports []Passport) int {
	validPassports := 0
	for _, passport := range passports {
		if passport.Validate() {
			validPassports += 1
		}
	}
	return validPassports
}

func readPassport(lines []string) Passport {
	attributes := make(map[string]string)
	for _, line := range lines {
		properties := strings.Split(line, " ")
		for _, property := range properties {
			parts := strings.Split(property, ":")
			attributes[parts[0]] = parts[1]
		}
	}
	return Passport{
		Attributes: attributes,
	}
}

func readPassports(filename string) ([]Passport, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	passports := make([]Passport, 0)
	lines := strings.Split(string(bytes), "\n")
	passportLines := make([]string, 0)
	for _, line := range lines {
		if len(line) == 0 && len(passportLines) > 0 {
			passports = append(passports, readPassport(passportLines))
			passportLines = make([]string, 0)
		} else {
			passportLines = append(passportLines, line)
		}
	}

	return passports, nil
}

func main() {
	passports, err := readPassports(os.Args[1])
	if err != nil {
		log.Fatalf("Unable to read passports: %s\n", err)
	}

	count := countValidPassports(passports)
	log.Printf("Valid passports: %d\n", count)
}
