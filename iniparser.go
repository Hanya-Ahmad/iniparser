package iniparser

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/pkg/errors"
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

// Exported errors to library user
var (
	ErrRepeatedSectionName     = errors.New("error repeated section name")
	ErrEmptyKey                = errors.New("error empty key")
	ErrRepeatedKeyName         = errors.New("error repeated key name")
	ErrMissingKeyValueOperator = errors.New("error missing key value operator")
	ErrNonIniFileParsed        = errors.New("error non ini file parsed")
	ErrNoSectionsFound         = errors.New("error no sections found")
	ErrSectionNotFound         = errors.New("error section not found")
	ErrKeyNotFound             = errors.New("error key not found")
	ErrNonIniFilePath          = errors.New("error non ini file path")
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
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			sectionName := strings.TrimPrefix(strings.TrimSuffix(line, "]"), "[")
			if sectionMap[sectionName] {
				return errors.Wrapf(ErrRepeatedSectionName, "error repeated section name at section %s", sectionName)
			}
			sectionMap[sectionName] = true
			currentSection = sectionName
			if keyMap[currentSection] == nil {
				keyMap[currentSection] = make(map[string]bool)
			}
		}
		if strings.Contains(line, "=") {
			keyValPair := strings.Split(line, "=")
			if len(keyValPair) == 2 && keyValPair[1] == "" {
				return errors.Wrapf(ErrEmptyKey, "error empty %s key", keyValPair[0])
			}
			key := strings.TrimSpace(keyValPair[0])
			if currentSection == "" || keyMap[currentSection] == nil {
				continue
			}
			if keyMap[currentSection][key] {
				return errors.Wrapf(ErrRepeatedKeyName, "error repeated key %s in section %s the last assigned value is applied", key, currentSection)
			}
			keyMap[currentSection][key] = true
		}
		if !(strings.HasPrefix(line, "[") || strings.HasSuffix(line, "]")) && (!strings.Contains(line, "=")) {
			return errors.Wrapf(ErrMissingKeyValueOperator, "error missing value at key %s", line)

		}
	}
	return nil
}

// Function to read ini data
func (p *IniParser) LoadFromReader(r io.Reader) {
	scanner := bufio.NewScanner(r)
	currentSection := ""
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.Replace(line, " ", "", -1)
		if line == "" || strings.HasPrefix(line, ";") || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			currentSection = line[1 : len(line)-1]
			p.sections[currentSection] = make(map[string]string)
		} else {
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

// Function to read ini data from string
func (p *IniParser) LoadFromString(str string) error {
	r := strings.NewReader(str)
	p.LoadFromReader(r)
	if len(p.sections) == 0 {
		return ErrNoSectionsFound
	}
	err := p.checkSyntax(str)
	if err != nil {
		return err
	}
	return nil
}

// Function to get section names of ini data
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
func (p *IniParser) ToString() string {
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
func (p *IniParser) Get(sectionName string, key any) (string, bool, error) {
	exists := true
	section, ok := p.sections[sectionName]
	if !ok {
		exists = false

		return "", exists, errors.Wrapf(ErrSectionNotFound, "error section %s not found", sectionName)
	}
	value, ok := section[fmt.Sprintf("%v", key)]
	if !ok {
		exists = false

		return "", exists, errors.Wrapf(ErrKeyNotFound, "error key %s not found in section %s", fmt.Sprintf("%v", key), sectionName)
	}
	if value == "" {
		exists = true
		return "", exists, nil
	}
	return value, exists, nil
}

// Function to set value of key in a specific section
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

// Function to save ini data to a file
func (p *IniParser) SaveToFile(path string) error {
	if !strings.HasSuffix(path, ".ini") {
		return ErrNonIniFilePath
	}
	str := p.ToString()
	err := os.WriteFile(path, []byte(str), 0644)
	if err != nil {
		return err
	}
	return nil
}
