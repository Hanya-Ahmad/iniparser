package main

import (
	"bufio"
	"strings"
)

// function to check for errors
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// this function is called in both implementations of INIFile's and INIString's GetSectionNames
func getSectionNames(str string) []string {
	var sectionNames []string
	scanner := bufio.NewScanner(strings.NewReader(str))

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.Replace(line, " ", "", -1)
		if line == "" || strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "[") {
			sectionNames = append(sectionNames, line[1:len(line)-1])
		}
	}
	return sectionNames
}

// this function is called in both implementations of INIFile's and INIString's GetSections
func getSections(str string) map[string]map[string]string {
	sections := make(map[string]map[string]string)
	currentSection := ""
	scanner := bufio.NewScanner(strings.NewReader(str))
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.Replace(line, " ", "", -1)
		if line == "" || strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") {
			continue
		}
		//If line is a section, create a new section map
		if strings.HasPrefix(line, "[") {
			currentSection = line[1 : len(line)-1]
			sections[currentSection] = make(map[string]string)
		} else {
			//If line is a key-value pair, append it to the current section
			keyValPair := strings.Split(line, "=")
			if len(keyValPair) == 2 {
				key := keyValPair[0]
				value := keyValPair[1]
				sections[currentSection][key] = value
			}
		}
	}
	return sections
}

