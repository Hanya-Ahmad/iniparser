package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

type INIParser interface {
	GetSectionNames(string) ([]string, error)
	GetSections() map[string]map[string]string
	Get(section_name string, key any) string
	Set(section_name string, key any, value any)
	ToString() (string, error)
	SaveToFile(string) error
}

type INIFile struct {
	path    string
	iniData map[string]map[string]string
}

type INIString struct {
	inputData string
	iniData   map[string]map[string]string
}

func LoadFromFile(path string) (*INIFile, error) {
	fmt.Println("load from file function started")
	f := INIFile{path: path}
	file, _ := os.Open(f.path)
	// check(err)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	if err := checkSyntax(scanner); err != nil {
		return &f, err
	}
	return &f, nil
}

func LoadFromString(data string) (*INIString, error) {
	fmt.Println("load from string function started")
	s := INIString{inputData: data}
	scanner := bufio.NewScanner(strings.NewReader(s.inputData))
	if err := checkSyntax(scanner); err != nil {
		return &s, err
	}
	return &s, nil
}

func (f INIFile) GetSectionNames(str string) ([]string, error) {
	return getSectionNames(str)
}

func (s INIString) GetSectionNames(str string) ([]string, error) {
	return getSectionNames(str)
}
func (f *INIFile) GetSections() map[string]map[string]string {
	if f.iniData != nil {
		return f.iniData
	} else {
		file, _ := os.Open(f.path)
		// check(err)
		defer file.Close()
		scanner := bufio.NewScanner(file)
		sections, err := getSections(scanner)
		checkErrors(err)
		return sections
	}
}

func (s *INIString) GetSections() map[string]map[string]string {
	if s.iniData != nil {
		return s.iniData
	} else {
		scanner := bufio.NewScanner(strings.NewReader(s.inputData))
		sections, err := getSections(scanner)
		checkErrors(err)

		return sections
	}
}
func (f INIFile) ToString() (string, error) {
	file, err := os.Open(f.path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var contents string
	for scanner.Scan() {
		contents += scanner.Text() + "\n"
	}

	return contents, nil
}

func (s INIString) ToString() (string, error) {
	if s.iniData != nil {
		return iniToString(s.iniData), nil
	} else if s.inputData == "" {
		return "", errors.New("error: no INI data found")
	} else {
		scanner := bufio.NewScanner(strings.NewReader(s.inputData))
		var contents string
		for scanner.Scan() {
			line := scanner.Text()
			line = strings.Replace(line, " ", "", -1)
			contents += line + "\n"
		}
		return contents, nil
	}
}

func (f INIFile) Get(sectionName string, key any) string {
	sections := f.GetSections()
	val, err := get(sections, sectionName, key)
	checkErrors(err)

	return val

}

func (s INIString) Get(sectionName string, key any) string {
	sections := s.GetSections()
	val, err := get(sections, sectionName, key)
	checkErrors(err)

	return val

}
func (f *INIFile) Set(sectionName string, key any, val any) {
	sections := f.GetSections()
	updatedSections, err := set(sections, sectionName, key, val)
	checkErrors(err)

	f.iniData = updatedSections
}

func (s *INIString) Set(sectionName string, key any, val any) {
	sections := s.GetSections()
	updatedSections, err := set(sections, sectionName, key, val)
	checkErrors(err)

	s.iniData = updatedSections

}
func (f INIFile) SaveToFile(path string) error {
	var data string
	var err error
	if f.iniData != nil {
		data = iniToString(f.iniData)
	} else {
		data, err = f.ToString()
		checkErrors(err)

	}

	return createFile(path, data)
}

func (s INIString) SaveToFile(path string) error {
	var data string
	if s.iniData != nil {
		data = iniToString(s.iniData)
	} else {
		data = s.inputData
	}
	return createFile(path, data)
}
func main() {
	path := "./iniFiles/test.ini"

	f, err := LoadFromFile(path)
	checkErrors(err)

	testInput(f, "Database", "port", "Email", "hanya", "hanya@mail.com","./iniFiles/file.ini")
	str, err := LoadFromString("[Person] \n  age= 1\n phone = 123 \n ] \n type = cat \n age = 3 \n  \n 123 = onetwothree \n email=tst")
	checkErrors(err)

	testInput(str, "Person", "age", "Person", "email", "john@mail.com","./iniFiles/string.ini")

}
