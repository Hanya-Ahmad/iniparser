package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)
func checkErrors(err error){
	if err != nil {
		fmt.Println(err)
	}
}
func checkSyntax(scanner *bufio.Scanner) error {
	sectionMap := make(map[string]bool)
	keyMap := make(map[string]map[string]bool)
	var currentSection string
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.Replace(line, " ", "", -1)
		if line == "" || strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") {
			continue
		} else if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			sectionName := strings.TrimPrefix(strings.TrimSuffix(line, "]"), "[")
			if sectionMap[sectionName] {
				return errors.New("error: repeated section name: " + sectionName)
			}
			sectionMap[sectionName] = true
			currentSection = sectionName
			if keyMap[currentSection] == nil {
				keyMap[currentSection] = make(map[string]bool)
			}
		} else if strings.Contains(line, "=") {
			keyValPair := strings.Split(line, "=")
			if len(keyValPair) == 2 && keyValPair[1] == "" {
				return errors.New("error: empty key at line: " + line)
			}
			key := strings.TrimSpace(keyValPair[0])
			if currentSection == "" || keyMap[currentSection] == nil {
				continue
			}
			if keyMap[currentSection][key] {
				return errors.New("error: repeated key name: " + key + " in section: " + currentSection+", the last assigned value is returned")
			}
			keyMap[currentSection][key] = true
		} else if !(strings.HasPrefix(line, "[") || strings.HasSuffix(line, "]")) && (!strings.Contains(line, "=")) {
			return errors.New("error: missing key-value operator '=' at line: " + line)
		}
	}
	return nil
}

// this function is called in both implementations of INIFile's and INIString's GetSectionNames
func getSectionNames(str string) ([]string, error) {
	var sectionNames []string
	scanner := bufio.NewScanner(strings.NewReader(str))
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.Replace(line, " ", "", -1)
		if line == "" || strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			sectionNames = append(sectionNames, line[1:len(line)-1])
		}

	}
	if len(sectionNames) == 0 {
		return nil, errors.New("error: no sections found")
	}
	return sectionNames, nil
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
func getSections(scanner *bufio.Scanner) (map[string]map[string]string, error) {
	sections := make(map[string]map[string]string)
	currentSection := ""
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.Replace(line, " ", "", -1)
		if line == "" || strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") {
			continue
		}
		//If line is a section, create a new section map
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			currentSection = line[1 : len(line)-1]
			sections[currentSection] = make(map[string]string)
		} else {
			//If line is a key-value pair, append it to the current section
			keyValPair := strings.Split(line, "=")
			if len(keyValPair) >= 2 {
				key := keyValPair[0]
				value := strings.Join(keyValPair[1:], "=")
				if currentSection == "" {
					continue
				}
				sections[currentSection][key] = value
			}
		}
	}
	if len(sections) == 0 {
		return sections, errors.New("error: no sections found")
	}
	return sections, nil
}
func get(sections map[string]map[string]string, sectionName string, key any) (string, error) {
	if section, ok := sections[sectionName]; ok {
		if value, ok := section[fmt.Sprintf("%v", key)]; ok {
			if value == "" {
				return "", errors.New("error: value is an empty string")
			}
			return value, nil
		} else {
			return "", errors.New("key " + fmt.Sprintf("%v ", key) + "not found in section " + sectionName)
		}
	} else {
		return "", errors.New("section " + sectionName + " not found")
	}
}
func set(sections map[string]map[string]string, sectionName string, key any, val any) (map[string]map[string]string,error) {
	if sections == nil {
		sections = make(map[string]map[string]string)
	}
	if section, ok := sections[sectionName]; ok {
		if _, ok := section[fmt.Sprintf("%v", key)]; ok {
			section[fmt.Sprintf("%v", key)] = fmt.Sprintf("%v", val)
		} else {
			// initiate a new key-value pair in sectionName
			section[fmt.Sprintf("%v", key)] = fmt.Sprintf("%v", val)
			return sections,errors.New("warning: key "+fmt.Sprintf("%v",key)+" not found in section "+sectionName+" but has been created")
		}
	} else {
		// initiate a new section and key-value pair
		sections[sectionName] = make(map[string]string)
		sections[sectionName][fmt.Sprintf("%v", key)] = fmt.Sprintf("%v", val)
		return sections,errors.New("warning: section "+sectionName+" not found but has been created")
	}
	return sections,nil
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

func testInput(p INIParser, getSection string, getKey any, setSection string, setKey any, setValue any, path string) {
	strParsed, _ := p.ToString()
	fmt.Println(strParsed)
	sectionNames, err := p.GetSectionNames(strParsed)
	checkErrors(err)
	fmt.Println("Section Names:", sectionNames)
	sections := p.GetSections()
	fmt.Println("Sections:", sections)
	val := p.Get(getSection, getKey)
	fmt.Println("Value: ", val)
	p.Set(setSection, setKey, setValue)
	updatedSection := p.GetSections()
	fmt.Println("Updated Section: ", updatedSection)
	// p.Set(setSection, 123, 456)
	p.SaveToFile(path)

}
