package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type INIParser interface {
	GetSectionNames(string) []string
	GetSections() map[string]map[string]string
	Get(section_name string, key any) string
	Set(section_name string, key any, value any)
	ToString() string
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

func LoadFromFile(path string) *INIFile {
	fmt.Println("load from file function started")
	f := INIFile{path: path}
	return &f
}

func LoadFromString(data string) *INIString {
	fmt.Println("load from string function started")
	s := INIString{inputData: data}
	return &s
}

func (f INIFile) GetSectionNames(str string) []string {
	return getSectionNames(str)
}

func (s INIString) GetSectionNames(str string) []string {
	return getSectionNames(str)
}
func (f *INIFile) GetSections() map[string]map[string]string {
	if f.iniData != nil {
		return f.iniData
	} else {
		file, err := os.Open(f.path)
		check(err)
		defer file.Close()
		scanner := bufio.NewScanner(file)
		sections := getSections(scanner)
		return sections
	}
}

func (s *INIString) GetSections() map[string]map[string]string {
	if s.iniData != nil {
		return s.iniData
	} else {
		scanner := bufio.NewScanner(strings.NewReader(s.inputData))
		sections := getSections(scanner)
		return sections
	}
}
func (f INIFile) ToString() string {
	file, err := os.Open(f.path)
	check(err)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var contents string
	for scanner.Scan() {
		contents += scanner.Text() + "\n"
	}
	err = scanner.Err()
	check(err)
	return contents
}

func (s INIString) ToString() string {
	if s.iniData != nil {
		return iniToString(s.iniData)
	} else {
		scanner := bufio.NewScanner(strings.NewReader(s.inputData))
		var contents string
		for scanner.Scan() {
			line := scanner.Text()
			line = strings.Replace(line, " ", "", -1)
			contents += line + "\n"
		}
		return contents
	}
}

func (f INIFile) Get(sectionName string, key any) string {
	sections := f.GetSections()
	return get(sections, sectionName, key)

}

func (s INIString) Get(sectionName string, key any) string {
	sections := s.GetSections()
	return get(sections, sectionName, key)

}
func (f *INIFile) Set(sectionName string, key any, val any) {
	sections := f.GetSections()
	updatedSections := set(sections, sectionName, key, val)
	f.iniData = updatedSections
}

func (s *INIString) Set(sectionName string, key any, val any) {
	sections := s.GetSections()
	updatedSections := set(sections, sectionName, key, val)
	s.iniData = updatedSections

}
func (f INIFile) SaveToFile(path string) error {
	var data string
	if f.iniData != nil {
		data = iniToString(f.iniData)
	} else {
		data = f.ToString()
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
	f := LoadFromFile(path)
	testInput(f, "Database", "port", "Email", "hanya", "hanya@mail.com")
	// str := LoadFromString("[Person] \n name = John Doe \n  age = \n phone = 123 \n [Pet] \n type = cat \n age = 3 \n [Numbers] \n 123 = onetwothree ")
	// testInput(str, "Person", "name", "Person", "email", "john@mail.com")

}
