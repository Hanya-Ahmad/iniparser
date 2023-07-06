package main

import (
	"bufio"
	"fmt"
	"os"
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
func iniToString(iniData map[string]map[string]string) string {
	var b strings.Builder

	for sectionName, section := range iniData {
		b.WriteString("[" + sectionName + "]\n")
		for key, val := range section {
			b.WriteString(key + "=" + val + "\n")
		}
		b.WriteString("\n")
	}

	return b.String()
}

// this function is called in both implementations of INIFile's and INIString's GetSections
func getSections(scanner *bufio.Scanner) map[string]map[string]string {
	sections := make(map[string]map[string]string)
	currentSection := ""
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

func get(sections map[string]map[string]string, sectionName string, key any) string {
	// sections := s.GetSections(s.ToString())
	if section, ok := sections[sectionName]; ok {
		if value, ok := section[fmt.Sprintf("%v", key)]; ok {
			return value
		} else {
			return "Key not found"
		}
	} else {
		return "Section Not found"
	}
}
func set(sections map[string]map[string]string, sectionName string, key any, val any) map[string]map[string]string {
	if sections == nil {
		sections = make(map[string]map[string]string)
	}
	if section, ok := sections[sectionName]; ok {
		if _, ok := section[fmt.Sprintf("%v", key)]; ok {
			section[fmt.Sprintf("%v", key)] = fmt.Sprintf("%v", val)
		} else {
			// initiate a new key-value pair in sectionName
			section[fmt.Sprintf("%v", key)] = fmt.Sprintf("%v", val)
			return sections
		}
	} else {
		// initiate a new section and key-value pair
		sections[sectionName] = make(map[string]string)
		sections[sectionName][fmt.Sprintf("%v", key)] = fmt.Sprintf("%v", val)
		return sections
	}
	return sections
}

func createFile(path, data string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := file.WriteString(data); err != nil {
		return err
	}
	return nil
}

func testInput(p INIParser, getSection string, getKey any, setSection string, setKey any, setValue any) {
	fmt.Println(p.ToString())
	sectionNames := p.GetSectionNames(p.ToString())
	fmt.Println("Section Names:", sectionNames)
	sections := p.GetSections()
	fmt.Println("Sections:", sections)
	val := p.Get(getSection, getKey)
	fmt.Println("Value: ", val)
	p.Set(setSection, setKey, setValue)
	updatedSection := p.GetSections()
	fmt.Println("Updated Section: ", updatedSection)
	p.Set(setSection,123,456)
	p.SaveToFile("./iniFiles/test.ini")


}
