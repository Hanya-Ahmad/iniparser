package iniparser

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type section map[string]string

type inidata map[string]section

// IniParser contains the ini data
type IniParser struct {
	data inidata
}

// NewIniParser initializes a new IniParser object
func NewIniParser() *IniParser {
	return &IniParser{
		data: make(inidata),
	}
}

// Syntax errors
var (
	ErrRepeatedSection   = fmt.Errorf("repeated section name")
	ErrEmptyKey          = fmt.Errorf("empty key")
	ErrRepeatedKeyName   = fmt.Errorf("repeated key name")
	ErrInvalidIniContent = fmt.Errorf("invalid ini content")
	ErrInvalidExtension  = fmt.Errorf("invalid extension")
	ErrNoSectionsFound   = fmt.Errorf("no sections found")
	ErrSectionNotFound   = fmt.Errorf("section not found")
	ErrKeyNotFound       = fmt.Errorf("key not found")
)

// loadFromReader checks for syntax errors
func (p *IniParser) loadFromReader(r io.Reader) error {
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
			if p.data[sectionName] != nil {
				return fmt.Errorf("%w: at section %s", ErrRepeatedSection, sectionName)
			}
			p.data[sectionName] = make(section)
			continue
		}

		keyValPair := strings.Split(line, "=")
		if len(keyValPair) == 2 && sectionName != "" {
			key := keyValPair[0]
			if key == "" {
				return fmt.Errorf("%w: in section %s", ErrEmptyKey, sectionName)
			}
			for k := range p.data[sectionName] {
				if k == key {
					return fmt.Errorf("%w: %s in section %s the last assigned value is applied", ErrRepeatedKeyName, key, sectionName)
				}
			}
			p.data[sectionName][key] = keyValPair[1]
			continue
		}
		return fmt.Errorf("%w: at key %s in section %s", ErrInvalidIniContent, line, sectionName)

	}
	return nil
}

// LoadFromFile parses ini data from file
func (p *IniParser) LoadFromFile(path string) error {
	if !strings.HasSuffix(path, ".ini") {
		return ErrInvalidExtension
	}
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return p.loadFromReader(file)
}

// LoadFromString parses ini data from string
func (p *IniParser) LoadFromString(str string) error {
	r := strings.NewReader(str)
	return p.loadFromReader(r)
}

// GetSectionNames retrieves all section names of parsed ini data
func (p *IniParser) GetSectionNames() []string {
	var sectionNames []string
	for section := range p.data {
		sectionNames = append(sectionNames, section)
	}
	return sectionNames
}

// GetSections retrieves all sections of parsed ini data
func (p *IniParser) GetSections() inidata {
	return p.data
}

func (p *IniParser) String() string {
	var result string
	for sectionName, section := range p.data {
		result += fmt.Sprintf("[%s]\n", sectionName)
		for key, value := range section {
			result += fmt.Sprintf("%s=%s\n", key, value)
		}
	}
	return result
}

// Get retrieves the value of a key in a given section name
func (p *IniParser) Get(sectionName string, key string) (string, error) {
	section, ok := p.data[sectionName]
	if !ok {

		return "", fmt.Errorf("%w: at section %s", ErrSectionNotFound, sectionName)
	}
	value, ok := section[key]
	if !ok {
		return "", fmt.Errorf("%w: at key %s", ErrKeyNotFound, key)
	}
	fmt.Println(value, "vak")
	return value, nil
}

// Set updates the value of a key in a given section name
func (p *IniParser) Set(sectionName string, key string, val string) {
	section, ok := p.data[sectionName]
	if !ok {
		section = make(map[string]string)
		p.data[sectionName] = section
	}
	section[key] = val
}

// SaveToFile saves ini data to a given path
func (p *IniParser) SaveToFile(path string) error {
	if !strings.HasSuffix(path, ".ini") {
		return ErrInvalidExtension
	}
	str := p.String()
	return os.WriteFile(path, []byte(str), 0644)

}
