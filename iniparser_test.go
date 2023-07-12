package iniparser

import (
	"strings"

	"github.com/pkg/errors"

	"os"
	"reflect"
	"sort"
	"testing"
)

func TestLoadFromReader(t *testing.T) {

	t.Run("checks that ini data is parsed correctly with no errors", func(t *testing.T) {
		const testIniData = `
		[section1]
		key1 = value1
		key2 = value2

		[section2]
		key3 = value3
		key4 = value4
		`
		parser := NewIniParser()
		err := parser.LoadFromReader(strings.NewReader(testIniData))
		if err != nil {
			t.Errorf("loadfromreader returned an error %v", err)
		}

		expectedSections := sections{
			"section1": map[string]string{
				"key1 ": " value1",
				"key2 ": " value2",
			},
			"section2": map[string]string{
				"key3 ": " value3",
				"key4 ": " value4",
			},
		}
		if !reflect.DeepEqual(parser.sections, expectedSections) {
			t.Errorf("loadfromreader did not load the correct sections expected %v got %v", expectedSections, parser.sections)
		}
	})
	t.Run("repeated section name error is returned when passing a string", func(t *testing.T) {
		str := `[Database]
		key1 = val
		[Email]
		key2 = val
		[Database]
		key3 = val`
		reader := strings.NewReader(str)
		ini := NewIniParser()
		err := ini.LoadFromReader(reader)
		want := ErrRepeatedSectionName
		if !errors.Is(err, want) {
			t.Errorf("expected error %q but got %q", want, err)
		}
	})

	t.Run("empty key error is returned when passing a string", func(t *testing.T) {
		str := `[Database]
		key1 = val
		key2 = val
		= d`
		reader := strings.NewReader(str)
		ini := NewIniParser()
		err := ini.LoadFromReader(reader)
		want := ErrEmptyKey
		if !errors.Is(err, want) {
			t.Errorf("expected error %q but got %q", want, err)
		}
	})

	t.Run("repeated key error is returned when passing a string", func(t *testing.T) {
		str := `[Database]
		key1 = val
		key2 = val
		key1 = anotherVal`
		reader := strings.NewReader(str)
		ini := NewIniParser()
		err := ini.LoadFromReader(reader)
		want := ErrRepeatedKeyName
		if !errors.Is(err, want) {
			t.Errorf("expected error %q but got %q", want, err)
		}
	})

	t.Run("missing assignment operator '=' error is returned when passing a string", func(t *testing.T) {
		str := `[Database]
		key1 = val
		key2 = val
		key3`
		reader := strings.NewReader(str)
		ini := NewIniParser()
		err := ini.LoadFromReader(reader)
		want := ErrMissingKeyValueOperator
		if !errors.Is(err, want) {
			t.Errorf("expected error %q but got %q", want, err)
		}
	})
}

func TestLoadFromFile(t *testing.T) {
	t.Run("error is nil when ini file path is correct", func(t *testing.T) {
		path := "./iniFiles/correctFile.ini"
		ini := NewIniParser()
		err := ini.LoadFromFile(path)
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
		want := ErrNonIniFileParsed
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
		user = root
		password = root
		port = 3000
		#also a comment
		[Numbers]
		[Database]
		name = 'John'
		age = 25.5
		ID = two = one`
		ini := NewIniParser()
		_ = ini.LoadFromString(str)
		sectionNames, err := ini.GetSectionNames()
		if err != nil {
			t.Errorf("correct string  %q but error %q occured", str, err)
		}
		expectedNames := []string{"Credentials", "Numbers", "Database"}
		sort.Strings(sectionNames)
		sort.Strings(expectedNames)
		if !reflect.DeepEqual(sectionNames, expectedNames) {
			t.Errorf("expected %q but got %q", expectedNames, sectionNames)
		}

	})
	t.Run("no sections found error when passing string with no sections", func(t *testing.T) {
		str := `key1 = val
		key2 = val`
		ini := NewIniParser()
		_ = ini.LoadFromString(str)
		_, err := ini.GetSectionNames()
		want := ErrNoSectionsFound
		if (!errors.Is(err, want)) && err != nil {
			t.Errorf("expected error %q but got error %q", want, err)
		} else if err == nil {
			t.Errorf("expected error %q but received no error", want)
		}
	})
}

func TestGetSections(t *testing.T) {
	t.Run("no sections found error when passing string with no sections", func(t *testing.T) {
		str := `key1 = val
		key2 = val`
		ini := NewIniParser()
		_ = ini.LoadFromString(str)
		_, err := ini.GetSections()
		want := ErrNoSectionsFound
		if (!errors.Is(err, want)) && err != nil {
			t.Errorf("expected error %q but got error %q", want, err)
		} else if err == nil {
			t.Errorf("expected error %q but received no error", want)
		}
	})
}

func TestGet(t *testing.T) {
	t.Run("empty string value returned and exists is true when passing an ini string with empty string", func(t *testing.T) {
		str := `[Database]
		key1=val
		key2= `
		ini := NewIniParser()
		_ = ini.LoadFromString(str)
		_, exists, err := ini.Get("Database", "key2")
		if !(exists == true && err == nil) {
			t.Errorf("expected empty string and exist to be true but received error %q and exists is %t", err, exists)
		}
	})

	t.Run("section not found error when passing string", func(t *testing.T) {
		str := `[Credentials]
		user = root
		password = root
		port = 3000
		[Database]
		name = 'John'
		age = 25.5
		`
		ini := NewIniParser()
		_ = ini.LoadFromString(str)
		_, _, err := ini.Get("Names", "key2")
		want := ErrSectionNotFound
		if (!errors.Is(err, want)) && err != nil {
			t.Errorf("expected error %q but got error %q", want, err)
		} else if err == nil {
			t.Errorf("expected error %q but received no error", want)
		}
	})

	t.Run("key not found in section error when passing string", func(t *testing.T) {
		str := `[Credentials]
		user = root
		password = root
		port = 3000
		[Database]
		name = 'John'
		age = 25.5
		`
		ini := NewIniParser()
		_ = ini.LoadFromString(str)
		_, _, err := ini.Get("Database", "status")
		want := ErrKeyNotFound
		if (!errors.Is(err, want)) && err != nil {
			t.Errorf("expected error %q but got error %q", want, err)
		} else if err == nil {
			t.Errorf("expected error %q but received no error", want)
		}
	})
}

func TestSaveToFile(t *testing.T) {
	t.Run("check if file has been created or not when called by an ini string struct", func(t *testing.T) {
		str := `[Credentials]
		user = root
		password = root
		port = 3000
		[Database]
		name = 'John'
		age = 25.5
		`
		ini := NewIniParser()
		_ = ini.LoadFromString(str)
		createdPath := "./iniFiles/createdFileFromString.ini"
		err := ini.SaveToFile(createdPath)
		if err != nil {
			t.Errorf("error saving file %v", err)
		}
		_, err = os.Stat(createdPath)
		if os.IsNotExist(err) {
			t.Errorf("error file not created %v", err)
		} else {
			err = os.Remove(createdPath)
			if err != nil {
				t.Errorf("error removing file %v", err)
			}
		}
	})

}

func TestSet(t *testing.T) {
	ini := NewIniParser()
	ini.sections["Section1"] = map[string]string{
		"key1": "value1",
		"key2": "value2",
	}
	ini.Set("Section1", "key1", "new value")
	if ini.sections["Section1"]["key1"] != "new value" {
		t.Errorf("set did not update key value, expected %q got %q", "newvalue", ini.sections["Section1"]["key1"])
	}
	ini.Set("Section1", "key3", "value3")
	if ini.sections["Section1"]["key3"] != "value3" {
		t.Errorf("set did not add new key with correct value expected %q  got %q", "value3", ini.sections["Section1"]["key3"])
	}
	ini.Set("Section2", "key1", "value1")
	if ini.sections["Section2"]["key1"] != "value1" {
		t.Errorf("set did not add new section with correct key value pair expected %q got %q", "value1", ini.sections["Section2"]["key1"])
	}
}
