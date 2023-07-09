package main

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

type INIParser interface {
	GetSectionNames() ([]string, error)
	GetSections() (map[string]map[string]string, error)
	Get(section_name string, key any) (string, error)
	Set(section_name string, key any, value any) error
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
	f := INIFile{path: path}
	file, err := os.Open(f.path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	if err := checkSyntax(scanner); err != nil {
		return &f, err
	}
	return &f, nil
}

func LoadFromString(data string) (*INIString, error) {
	s := INIString{inputData: data}
	scanner := bufio.NewScanner(strings.NewReader(s.inputData))
	if err := checkSyntax(scanner); err != nil {
		return &s, err
	}
	return &s, nil
}

func (f INIFile) GetSectionNames() ([]string, error) {
	str, _ := f.ToString()

	return getSectionNames(str)
}

func (s INIString) GetSectionNames() ([]string, error) {
	str, _ := s.ToString()
	return getSectionNames(str)
}
func (f *INIFile) GetSections() (map[string]map[string]string, error) {
	if f.iniData != nil {
		return f.iniData, nil
	} else {
		file, err := os.Open(f.path)
		if (err) != nil {
			return nil, (err)
		}

		defer file.Close()
		scanner := bufio.NewScanner(file)
		sections, err := getSections(scanner)
		if (err) != nil {
			return nil, (err)
		}
		return sections, nil
	}
}

func (s *INIString) GetSections() (map[string]map[string]string, error) {
	if s.iniData != nil {
		return s.iniData, nil
	} else {
		scanner := bufio.NewScanner(strings.NewReader(s.inputData))
		sections, err := getSections(scanner)
		if (err) != nil {
			return nil, (err)
		}

		return sections, nil
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

func (f INIFile) Get(sectionName string, key any) (string, error) {
	sections, err := f.GetSections()
	if (err) != nil {
		return "", (err)
	}
	val, err := get(sections, sectionName, key)
	if (err) != nil {
		return "", (err)
	}

	return val, nil

}

func (s INIString) Get(sectionName string, key any) (string, error) {
	sections, err := s.GetSections()
	if (err) != nil {
		return "", (err)
	}
	val, err := get(sections, sectionName, key)
	if (err) != nil {
		return "", (err)
	}
	return val, nil

}
func (f *INIFile) Set(sectionName string, key any, val any) error {

	sections, err := f.GetSections()
	if (err) != nil {
		return (err)
	}
	updatedSections := set(sections, sectionName, key, val)
	if (err) != nil {
		return (err)
	}

	f.iniData = updatedSections
	return nil
}

func (s *INIString) Set(sectionName string, key any, val any) error {
	sections, err := s.GetSections()
	if (err) != nil {
		return (err)
	}
	updatedSections := set(sections, sectionName, key, val)

	s.iniData = updatedSections
	return nil

}
func (f INIFile) SaveToFile(path string) error {
	var data string
	var err error
	if f.iniData != nil {
		data = iniToString(f.iniData)
	} else {
		data, err = f.ToString()
		if (err) != nil {
			return err
		}

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
	path := "./iniFiles/correctFile.ini"

	f, err := LoadFromFile(path)
	checkErrors(err)

	testInput(f, "Database", "port", "Email", "hanya", "hanya@mail.com", "./iniFiles/file.ini")
	str, err := LoadFromString("[Person] \n  age= 1\n phone = 123 \n ] \n type = cat \n age = 3 \n  \n 123 = onetwothree \n email=tst")
	checkErrors(err)

	testInput(str, "Person", "age", "Person", "email", "john@mail.com", "./iniFiles/string.ini")

}
