package iniparser

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

type sections map[string]map[string]string

type IniParser struct {
	sections sections
}

func NewIniParser() *IniParser {
	return &IniParser{
		sections: make(sections),
	}
}

// TODO: how to export errors and embed variables into them as well?
var (
	ErrRepeatedSectionName = errors.New("error repeated section name")
	ErrEmptyKey            = errors.New("empty key")
	ErrRepeatedKeyName     = errors.New("repeated key name")
	ErrMissingKeyValue     = errors.New("missing key value operator")
	ErrNonIniFileParsed    = errors.New("non ini file parsed")
	ErrNoSectionsFound     = errors.New("no sections found")
	ErrSectionNotFound     = errors.New("section not found")
	ErrKeyNotFound         = errors.New("key not found")
	ErrNonIniFilePath      = errors.New("non ini file path")
)

func (p *IniParser) checkSyntax(contents string) error {
	sectionMap := make(map[string]bool)
	keyMap := make(map[string]map[string]bool)
	var currentSection string
	scanner := bufio.NewScanner(strings.NewReader(contents))
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.Replace(line, " ", "", -1)
		if line == "" || strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") {
			continue
		} else if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			sectionName := strings.TrimPrefix(strings.TrimSuffix(line, "]"), "[")
			if sectionMap[sectionName] {
				return errors.New("error repeated section name " + sectionName)
			}
			sectionMap[sectionName] = true
			currentSection = sectionName
			if keyMap[currentSection] == nil {
				keyMap[currentSection] = make(map[string]bool)
			}
		} else if strings.Contains(line, "=") {
			keyValPair := strings.Split(line, "=")
			if len(keyValPair) == 2 && keyValPair[1] == "" {
				return errors.New("error empty key at key " + keyValPair[0])
			}
			key := strings.TrimSpace(keyValPair[0])
			if currentSection == "" || keyMap[currentSection] == nil {
				continue
			}
			if keyMap[currentSection][key] {
				return errors.New("error repeated key name " + key + " in section " + currentSection + " the last assigned value is applied")
			}
			keyMap[currentSection][key] = true
		} else if !(strings.HasPrefix(line, "[") || strings.HasSuffix(line, "]")) && (!strings.Contains(line, "=")) {
			return errors.New("error missing key value operator at key " + line)
		}
	}
	return nil
}

func (p *IniParser) LoadFromFile(path string) error {
	if !strings.HasSuffix(path, ".ini") {
		return errors.New("error non ini file parsed")
	}
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var contents string
	for scanner.Scan() {
		contents += scanner.Text() + "\n"
	}
	err = p.LoadFromString(contents)
	if err != nil {
		return err
	}
	return nil
}

func (p *IniParser) LoadFromString(str string) error {
	scanner := bufio.NewScanner(strings.NewReader(str))
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
			p.sections[currentSection] = make(map[string]string)
		} else {
			//If line is a key-value pair, append it to the current section
			keyValPair := strings.Split(line, "=")
			if len(keyValPair) >= 2 {
				key := keyValPair[0]
				value := strings.Join(keyValPair[1:], "=")
				if currentSection == "" {
					continue
				}
				p.sections[currentSection][key] = value
			}
		}
	}
	if len(p.sections) == 0 {
		return errors.New("error no sections found")
	}
	err := p.checkSyntax(str)
	if err != nil {
		return err
	}
	return nil
}
func (p *IniParser) GetSectionNames() ([]string, error) {
	var sectionNames []string
	for section := range p.sections {
		sectionNames = append(sectionNames, section)
	}
	if len(sectionNames) == 0 {
		return nil, errors.New("error no sections found")
	}
	return sectionNames, nil
}

func (p *IniParser) GetSections() (sections, error) {
	if len(p.sections) == 0 {
		return nil, errors.New("error no sections found")
	}
	return p.sections, nil
}
func (p *IniParser) ToString() string {
	var sb strings.Builder
	for sectionName, section := range p.sections {
		sb.WriteString("[" + sectionName + "]\n")
		for key, value := range section {
			sb.WriteString(key + " = " + value + "\n")
		}
		// sb.WriteString("\n")
	}
	return sb.String()
}

func (p *IniParser) Get(sectionName string, key any) (string, bool, error) {
	exists := true
	section, ok := p.sections[sectionName]
	if !ok {
		exists = false
		return "", exists, errors.New("error section " + sectionName + " not found")
	}
	value, ok := section[fmt.Sprintf("%v", key)]
	if !ok {
		exists = false
		return "", exists, errors.New("error key " + fmt.Sprintf("%v", key) + " not found in section " + sectionName)
	}
	if value == "" {
		exists = true
		return "", exists, nil
	}
	return value, exists, nil
}
func (p *IniParser) Set(sectionName string, key any, val any) {
	keyStr, _ := key.(string)
	valStr, _ := val.(string)
	section, ok := p.sections[sectionName]
	if !ok {
		section = make(map[string]string)
		p.sections[sectionName] = section
	}
	section[keyStr] = valStr
}

func (p *IniParser) SaveToFile(path string) error {
	if !strings.HasSuffix(path, ".ini") {
		return errors.New("error non ini file path")
	}
	str := p.ToString()
	err := os.WriteFile(path, []byte(str), 0644)
	if err != nil {
		return err
	}
	return nil
}
