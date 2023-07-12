package iniparser

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type sections map[string]map[string]string

type IniParser struct {
	sections sections
}

// Function to initialize a new IniParser object
func NewIniParser() *IniParser {
	return &IniParser{
		sections: make(sections),
	}
}

// Syntax errors
var (
	ErrRepeatedSectionName     = fmt.Errorf("repeated section name")
	ErrEmptyKey                = fmt.Errorf("empty key")
	ErrRepeatedKeyName         = fmt.Errorf("repeated key name")
	ErrMissingKeyValueOperator = fmt.Errorf("missing key value operator")
	ErrNonIniFileParsed        = fmt.Errorf("non ini file parsed")
	ErrNoSectionsFound         = fmt.Errorf("no sections found")
	ErrSectionNotFound         = fmt.Errorf("section not found")
	ErrKeyNotFound             = fmt.Errorf("key not found")
	ErrNonIniFilePath          = fmt.Errorf("non ini file path")
)

// Function to parse data from reader and checks for syntax errors
func (p *IniParser) LoadFromReader(r io.Reader) error {
	sectionMap := make(map[string]bool)
	keyMap := make(map[string]map[string]bool)
	var sectionName string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			sectionName = line[1 : len(line)-1]
			if sectionMap[sectionName] {
				return fmt.Errorf("%w repeated section name at section %s", ErrRepeatedSectionName, sectionName)
			}
			sectionMap[sectionName] = true
			if keyMap[sectionName] == nil {
				keyMap[sectionName] = make(map[string]bool)
			}
			p.sections[sectionName] = make(map[string]string)
			continue
		}
		keyValPair := strings.Split(line, "=")
		if len(keyValPair) == 2 {
			key := keyValPair[0]
			if key == "" {
				return fmt.Errorf("%w empty key in section %s", ErrEmptyKey, sectionName)
			}
			if sectionName == "" {
				continue
			}
			if keyMap[sectionName][key] {
				return fmt.Errorf("%w repeated key %s in section %s the last assigned value is applied", ErrRepeatedKeyName, key, sectionName)

			}
			keyMap[sectionName][key] = true
			p.sections[sectionName][key] = keyValPair[1]
		} else if !(strings.HasPrefix(line, "[") || strings.HasSuffix(line, "]")) && !strings.Contains(line, "=") {
			return fmt.Errorf("%w missing value at key %s in section %s", ErrMissingKeyValueOperator, line, sectionName)

		}
	}
	return nil
}

// Function to read ini data from file
func (p *IniParser) LoadFromFile(path string) error {
	if !strings.HasSuffix(path, ".ini") {
		return ErrNonIniFileParsed
	}
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return p.LoadFromReader(file)
}

// Function to read ini data from string
func (p *IniParser) LoadFromString(str string) error {
	r := strings.NewReader(str)
	err := p.LoadFromReader(r)

	if err != nil {
		return err
	}
	return nil
}

// Function to get all section names of ini data
func (p *IniParser) GetSectionNames() ([]string, error) {
	var sectionNames []string
	for section := range p.sections {
		sectionNames = append(sectionNames, section)
	}
	if len(sectionNames) == 0 {
		return nil, ErrNoSectionsFound
	}
	return sectionNames, nil
}

// Function to get entire ini data
func (p *IniParser) GetSections() (sections, error) {
	if len(p.sections) == 0 {
		return nil, ErrNoSectionsFound
	}
	return p.sections, nil
}

// Function to convert ini data to string
func (p *IniParser) String() string {
	var sb strings.Builder
	for sectionName, section := range p.sections {
		sb.WriteString("[" + sectionName + "]\n")
		for key, value := range section {
			sb.WriteString(key + "=" + value + "\n")
		}
	}
	return sb.String()
}

// Function to get value of key in a specific section
func (p *IniParser) Get(sectionName string, key string) (string, bool, error) {
	section, ok := p.sections[sectionName]
	if !ok {
		return "", false, fmt.Errorf("%w section %s not found", ErrSectionNotFound, sectionName)
	}
	value, ok := section[key]
	if !ok {
		return "", false, fmt.Errorf("%w key %s not found in section %s", ErrKeyNotFound, key, sectionName)
	}
	if value == "" {
		return "", true, nil
	}
	return value, true, nil
}

// Function to set value of key in a specific section
func (p *IniParser) Set(sectionName string, key string, val string) {
	section, ok := p.sections[sectionName]
	if !ok {
		section = make(map[string]string)
		p.sections[sectionName] = section
	}
	section[key] = val
}

// Function to save ini data to a file
func (p *IniParser) SaveToFile(path string) error {
	if !strings.HasSuffix(path, ".ini") {
		return ErrNonIniFilePath
	}
	str := p.String()
	err := os.WriteFile(path, []byte(str), 0644)
	if err != nil {
		return err
	}
	return nil
}
