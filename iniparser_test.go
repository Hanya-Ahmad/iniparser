package iniparser

import (
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/pkg/errors"
)

func TestLoadFromReader(t *testing.T) {

	t.Run("checks that ini data is parsed correctly with no errors", func(t *testing.T) {
		const testIniData = `
		[section1]
		key1=value1
		key2=value2
		[section2]
		key3=value3
		key4=value4
		`
		parser := NewIniParser()
		err := parser.loadFromReader(strings.NewReader(testIniData))
		if err != nil {
			t.Errorf("loadfromreader returned an error %v", err)
		}

		expectedSections := inidata{
			"section1": map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
			"section2": map[string]string{
				"key3": "value3",
				"key4": "value4",
			},
		}
		if !reflect.DeepEqual(parser.data, expectedSections) {
			t.Errorf("loadfromreader did not load the correct sections expected %v got %v", expectedSections, parser.data)
		}
	})
	t.Run("repeated section name error is returned when passing a string", func(t *testing.T) {
		str := `[Database]
		key1=val
		[Email]
		key2=val
		[Database]
		key3=val`
		reader := strings.NewReader(str)
		ini := NewIniParser()
		err := ini.loadFromReader(reader)
		want := ErrRepeatedSection
		if !errors.Is(err, want) {
			t.Errorf("expected error %q but got %q", want, err)
		}
	})

	t.Run("empty key error is returned when passing a string", func(t *testing.T) {
		str := `[Database]
		key1=val
		key2=val
		= d`
		reader := strings.NewReader(str)
		ini := NewIniParser()
		err := ini.loadFromReader(reader)
		want := ErrEmptyKey
		if !errors.Is(err, want) {
			t.Errorf("expected error %q but got %q", want, err)
		}
	})

	t.Run("repeated key error is returned when passing a string", func(t *testing.T) {
		str := `[Database]
		key1=val
		key2=val
		key1=anotherVal`
		reader := strings.NewReader(str)
		ini := NewIniParser()
		err := ini.loadFromReader(reader)
		want := ErrRepeatedKeyName
		if !errors.Is(err, want) {
			t.Errorf("expected error %q but got %q", want, err)
		}
	})

	t.Run("missing assignment operator '=' error is returned when passing a string", func(t *testing.T) {
		str := `[Database]
		key1=val
		key2=val
		key3`
		reader := strings.NewReader(str)
		ini := NewIniParser()
		err := ini.loadFromReader(reader)
		want := ErrInvalidIniContent
		if !errors.Is(err, want) {
			t.Errorf("expected error %q but got %q", want, err)
		}
	})
}

func TestLoadFromFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "tmp")
	if err != nil {
		t.Errorf("failed to create temporary file")
	}
	defer os.RemoveAll(tempDir)
	t.Run("error is nil when ini file path is correct", func(t *testing.T) {
		data :=
			`;comment
        [Section]
        key1=val1
        key2=val2
        #also a comment
        [Credentials]
        user=root
        password=root123
        port=3000
        `
		filePath := "example.ini"
		path := filepath.Join(tempDir, filePath)
		err := os.WriteFile(path, []byte(data), 0644)
		if err != nil {
			t.Errorf("error %v", err)
		}
		defer os.Remove(path)
		ini := NewIniParser()
		err = ini.LoadFromFile(path)
		if err != nil {
			t.Errorf("correct file path %q but error %q occured", path, err)
		}
	})
	t.Run("os.Open error occurs when path is incorrect", func(t *testing.T) {
		path := "./iniFiles/incorrectPath.ini"
		ini := NewIniParser()
		err := ini.LoadFromFile(path)
		want := "open ./iniFiles/incorrectPath.ini: no such file or directory"
		if (err.Error() != want) && err != nil {
			t.Errorf("expected error %q but got error %q", want, err)
		} else if err == nil {
			t.Errorf("expected error %q but received no error", want)
		}
	})
	t.Run("non ini file passed error", func(t *testing.T) {
		path := "./iniFiles/config.txt"
		ini := NewIniParser()
		err := ini.LoadFromFile(path)
		want := ErrInvalidExtension
		if (!errors.Is(err, want)) && err != nil {
			t.Errorf("expected error %q but got error %q", want, err)
		} else if err == nil {
			t.Errorf("expected error %q but received no error", want)
		}
	})
}

func TestGetSectionNames(t *testing.T) {
	t.Run("correct section names returned with no error when passing correct ini string", func(t *testing.T) {
		str := `[Credentials]
		user=root
		password=root
		port=3000
		#also a comment
		[Numbers]
		[Database]
		name='John'
		age=25.5
		ID=two=one`
		ini := NewIniParser()
		_ = ini.LoadFromString(str)
		sectionNames := ini.GetSectionNames()
		expectedNames := []string{"Credentials", "Numbers", "Database"}
		sort.Strings(sectionNames)
		sort.Strings(expectedNames)
		if !reflect.DeepEqual(sectionNames, expectedNames) {
			t.Errorf("expected %q but got %q", expectedNames, sectionNames)
		}

	})
}
func TestGetSections(t *testing.T) {
	// Define test data
	testData := `
        [section1]
        key1=value1
        key2=value2
        [section2]
        key3=value3
        key4=value4
    `
	iniparser := NewIniParser()
	_ = iniparser.LoadFromString(testData)
	sections := iniparser.GetSections()
	expected := inidata{
		"section1": map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
		"section2": map[string]string{
			"key3": "value3",
			"key4": "value4",
		},
	}
	if !reflect.DeepEqual(sections, expected) {
		t.Errorf("got %v expected %v", sections, expected)
	}
}

func TestGet(t *testing.T) {
	str := `[Credentials]
	user= 
	password=root
	port=3000
	key=
	[Database]
	name='John'
	age=25.5
	`
	t.Run("section not found error when passing string", func(t *testing.T) {
		ini := NewIniParser()
		_ = ini.LoadFromString(str)
		_, err := ini.Get("Names", "key2")
		want := ErrSectionNotFound
		if (!errors.Is(err, want)) && err != nil {
			t.Errorf("expected error %q but got error %q", want, err)
		} else if err == nil {
			t.Errorf("expected error %q but received no error", want)
		}
	})

	t.Run("key not found in section error when passing string", func(t *testing.T) {
		ini := NewIniParser()
		_ = ini.LoadFromString(str)
		_, err := ini.Get("Database", "status")
		want := ErrKeyNotFound
		if (!errors.Is(err, want)) && err != nil {
			t.Errorf("expected error %q but got error %q", want, err)
		} else if err == nil {
			t.Errorf("expected error %q but received no error", want)
		}
	})
	t.Run("value retrieved is correct", func(t *testing.T) {
		ini := NewIniParser()
		_ = ini.LoadFromString(str)
		key, err := ini.Get("Credentials", "key")
		want := ""
		if !(err == nil && key == want) {
			t.Errorf("got %s expected %s", key, want)
		}
	})
}

func TestSaveToFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "tmp")
	str := `[Credentials]
	user=root
	password=root
	port=3000
	[Database]
	name='John'
	age=25.5
	`
	if err != nil {
		t.Errorf("failed to create temporary file")
	}
	defer os.RemoveAll(tempDir)
	t.Run("check if file has been created or not when called by an ini string struct", func(t *testing.T) {
		ini := NewIniParser()
		_ = ini.LoadFromString(str)
		createdPath := "createdFileFromString.ini"
		path := filepath.Join(tempDir, createdPath)
		err := ini.SaveToFile(path)
		if err != nil {
			t.Errorf("failed to save file %v", err)
		}
		_, err = os.Stat(path)
		if os.IsNotExist(err) {
			t.Errorf("failed to create file %v", err)
		} else {
			err = os.Remove(path)
			if err != nil {
				t.Errorf("failed to remove file %v", err)
			}
		}
	})
	t.Run("invalid extension error", func(t *testing.T) {
		ini := NewIniParser()
		_ = ini.LoadFromString(str)
		createdPath := "createdFileFromString.txt"
		path := filepath.Join(tempDir, createdPath)
		err := ini.SaveToFile(path)
		want := ErrInvalidExtension
		if !errors.Is(err, want) {
			t.Errorf("got %s expected %s", err, want)
		}
	})

}

func TestSet(t *testing.T) {
	ini := NewIniParser()
	ini.data["Section1"] = map[string]string{
		"key1": "value1",
		"key2": "value2",
	}
	ini.Set("Section1", "key1", "new value")
	if ini.data["Section1"]["key1"] != "new value" {
		t.Errorf("set did not update key value, expected %q got %q", "newvalue", ini.data["Section1"]["key1"])
	}
	ini.Set("Section1", "key3", "value3")
	if ini.data["Section1"]["key3"] != "value3" {
		t.Errorf("set did not add new key with correct value expected %q  got %q", "value3", ini.data["Section1"]["key3"])
	}
	ini.Set("Section2", "key1", "value1")
	if ini.data["Section2"]["key1"] != "value1" {
		t.Errorf("set did not add new section with correct key value pair expected %q got %q", "value1", ini.data["Section2"]["key1"])
	}
}
