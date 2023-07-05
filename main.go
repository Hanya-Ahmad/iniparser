package main

import (
	"bufio"
	"fmt"
	"os"
	// "strings"
)

type INIParser interface {
	GetSectionNames(string) []string
	GetSections(string) map[string]map[string]string
	Get(section_name string, key any) string
	Set(section_name string, key any, value any) map[string]map[string]string
	ToString() string
	// SaveToFile(string) error
}

type INIFile struct {
	path string
}

type INIString struct {
	data string
}

func LoadFromFile(ini INIParser) {
	fmt.Println("load from file function started")
	s := ini.ToString()
	sectionNames := ini.GetSectionNames(s)
	fmt.Println("Section Names:", sectionNames)
	sections := ini.GetSections(s)
	fmt.Println("Sections:", sections)
	val := ini.Get("Database", 322)
	fmt.Println("Value: ", val)
	updatedSection := ini.Set("Email", "hanya", "hanya@mail.com")
	fmt.Println("Updated Section: ", updatedSection)
}

func LoadFromString(ini INIParser) {
	fmt.Println("load from string function started")
	s := ini.ToString()
	sectionNames := ini.GetSectionNames(s)
	fmt.Println("Section Names:", sectionNames)
	sections := ini.GetSections(s)
	fmt.Println("Sections:", sections)
	val := ini.Get("Person", "name")
	fmt.Println("Value: ", val)
	updatedSection := ini.Set("Person", "email", "john@mail.com")
	fmt.Println("Updated Section: ", updatedSection)

}

func (f INIFile) GetSectionNames(str string) []string {
	return getSectionNames(str)
}

func (s INIString) GetSectionNames(str string) []string {
	return getSectionNames(str)
}
func (f INIFile) GetSections(str string) map[string]map[string]string {
	return getSections(str)
}
func (f INIString) GetSections(str string) map[string]map[string]string {
	return getSections(str)
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
	return s.data
}

func (f INIFile) Get(sectionName string, key any) string {
	sections := f.GetSections(f.ToString())
	if section, ok := sections[sectionName]; ok {
		if value, ok := section[fmt.Sprintf("%v", key)]; ok {
			return value
		} else {
			return "Key not found"
		}
	} else {
		return "Section not found"
	}
}

func (s INIString) Get(sectionName string, key any) string {
	sections := s.GetSections(s.ToString())
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
func (f INIFile) Set(sectionName string, key any, val any) map[string]map[string]string {
	sections := f.GetSections(f.ToString())
	if section, ok := sections[sectionName]; ok {
		section[fmt.Sprintf("%v", key)] = fmt.Sprintf("%v", val)
	} else {
		sections[sectionName] = make(map[string]string)
		sections[sectionName][fmt.Sprintf("%v", key)] = fmt.Sprintf("%v", val)
	}
	return map[string]map[string]string{sectionName: sections[sectionName]}
}

func (s INIString) Set(sectionName string, key any, val any) map[string]map[string]string {
	sections := s.GetSections(s.ToString())
	if section, ok := sections[sectionName]; ok {
		section[fmt.Sprintf("%v", key)] = fmt.Sprintf("%v", val)
	} else {
		sections[sectionName] = make(map[string]string)
		sections[sectionName][fmt.Sprintf("%v", key)] = fmt.Sprintf("%v", val)
	}
	return map[string]map[string]string{sectionName: sections[sectionName]}

}

func main() {
	path := "./iniFiles/test.ini"
	f := INIFile{path: path}
	s := INIString{data: "[Person] \n name = John Doe \n  age = \n phone = 123 \n [Pet] \n type = cat \n age = 3 \n [Numbers] \n 123 = onetwothree "}
	LoadFromFile(f)
	LoadFromString(s)

}
