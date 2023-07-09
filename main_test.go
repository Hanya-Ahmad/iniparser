package main

import (
	"os"
	"reflect"
	"testing"
)

func TestLoadFromFile(t *testing.T) {
	t.Run("error is nil when ini file path is correct", func(t *testing.T) {
		path := "./iniFiles/correctFile.ini"
		_, err := LoadFromFile(path)
		if err != nil {
			t.Errorf("correct file path %q but error %q occured", path, err)
		}
	})
	t.Run("os.Open error occurs when path is incorrect", func(t *testing.T) {
		path := "./iniFiles/incorrectPath.ini"
		_, err := LoadFromFile(path)
		want := "open ./iniFiles/incorrectPath.ini: no such file or directory"
		if (err.Error() != want) && err != nil {
			t.Errorf("expected error %q but got error %q", want, err)
		} else if err == nil {
			t.Errorf("expected error %q but received no error", want)
		}
	})
	t.Run("repeated section name error when passing repeatedSection.ini", func(t *testing.T) {
		path := "./iniFiles/repeatedSection.ini"
		_, err := LoadFromFile(path)
		want := "error: repeated section name: Database"
		if (err.Error() != want) && err != nil {
			t.Errorf("expected error %q but got error %q", want, err)
		} else if err == nil {
			t.Errorf("expected error %q but received no error", want)
		}
	})
	t.Run("empty key error when passing emptyKey.ini", func(t *testing.T) {
		path := "./iniFiles/emptyKey.ini"
		_, err := LoadFromFile(path)
		want := "error: empty key at key: key3"
		if (err.Error() != want) && err != nil {
			t.Errorf("expected error %q but got error %q", want, err)
		} else if err == nil {
			t.Errorf("expected error %q but received no error", want)
		}
	})
	t.Run("repeated key error when passing repeatedKey.ini", func(t *testing.T) {
		path := "./iniFiles/repeatedKey.ini"
		_, err := LoadFromFile(path)
		want := "error: repeated key name: key1 in section: Database, the last assigned value is applied"
		if (err.Error() != want) && err != nil {
			t.Errorf("expected error %q but got error %q", want, err)
		} else if err == nil {
			t.Errorf("expected error %q but received no error", want)
		}
	})
	t.Run("missing assignment operator '='  error when passing missingEqual.ini", func(t *testing.T) {
		path := "./iniFiles/missingEqual.ini"
		_, err := LoadFromFile(path)
		want := "error: missing key-value operator '=' at key: key3"
		if (err.Error() != want) && err != nil {
			t.Errorf("expected error %q but got error %q", want, err)
		} else if err == nil {
			t.Errorf("expected error %q but received no error", want)
		}
	})
}
func TestLoadFromString(t *testing.T) {
	t.Run("repeated section name error when passing a string", func(t *testing.T) {
		str := "[Database]\nkey1 = val\n[Email]\nkey2 = val\n[Database]\nkey3 = val"
		_, err := LoadFromString(str)
		want := "error: repeated section name: Database"
		if (err.Error() != want) && err != nil {
			t.Errorf("expected error %q but got error %q", want, err)
		} else if err == nil {
			t.Errorf("expected error %q but received no error", want)
		}
	})
	t.Run("empty key error when passing a string", func(t *testing.T) {
		str := "[Database]\nkey1 = val\nkey2 = val\nkey3 ="
		_, err := LoadFromString(str)
		want := "error: empty key at key: key3"
		if (err.Error() != want) && err != nil {
			t.Errorf("expected error %q but got error %q", want, err)
		} else if err == nil {
			t.Errorf("expected error %q but received no error", want)
		}
	})
	t.Run("repeated key error when passing a string", func(t *testing.T) {
		str := "[Database]\nkey1 = val\nkey2 = val\nkey1 = anotherVal"
		_, err := LoadFromString(str)
		want := "error: repeated key name: key1 in section: Database, the last assigned value is applied"
		if (err.Error() != want) && err != nil {
			t.Errorf("expected error %q but got error %q", want, err)
		} else if err == nil {
			t.Errorf("expected error %q but received no error", want)
		}
	})
	t.Run("missing assignment operator '='  error when passing a string", func(t *testing.T) {
		str := "[Database]\nkey1 = val\nkey2 = val\nkey3"
		_, err := LoadFromString(str)
		want := "error: missing key-value operator '=' at key: key3"
		if (err.Error() != want) && err != nil {
			t.Errorf("expected error %q but got error %q", want, err)
		} else if err == nil {
			t.Errorf("expected error %q but received no error", want)
		}
	})
}
func TestGetSectionNames(t *testing.T) {
	t.Run("error is nil when passing correctFile.ini", func(t *testing.T) {
		path := "./iniFiles/correctFile.ini"
		file, _ := LoadFromFile(path)
		_, err := file.GetSectionNames()
		if err != nil {
			t.Errorf("correct file  %q but error %q occured", path, err)
		}
	})
	t.Run("correct section names returned when passing correctFile.ini", func(t *testing.T) {
		path := "./iniFiles/correctFile.ini"
		file, _ := LoadFromFile(path)
		sectionNames, _ := file.GetSectionNames()
		want := []string{"Credentials", "Numbers", "Database"}
		if !reflect.DeepEqual(sectionNames, want) {
			t.Errorf("expected %q but got %q", want, sectionNames)
		}
	})
	t.Run("no sections found error when passing noSections.ini", func(t *testing.T) {
		path := "./iniFiles/noSections.ini"
		file, _ := LoadFromFile(path)
		_, err := file.GetSectionNames()
		want := "error: no sections found"
		if (err.Error() != want) && err != nil {
			t.Errorf("expected error %q but got error %q", want, err)
		} else if err == nil {
			t.Errorf("expected error %q but received no error", want)
		}

	})
	t.Run("error is nil when passing correct ini string", func(t *testing.T) {
		str := "[Credentials]\nuser = root\npassword = root\nport = 3000\n#also a comment\n[Numbers]\n123 = onetwothree\n3.14 = pi\n[Database]\nname = 'John'\nage = 25.5\nID = 1 = one"
		s, _ := LoadFromString(str)
		_, err := s.GetSectionNames()
		if err != nil {
			t.Errorf("correct string  %q but error %q occured", str, err)
		}
	})
	t.Run("correct section names returned when passing correct string", func(t *testing.T) {
		str := "[Credentials]\nuser = root\npassword = root\nport = 3000\n#also a comment\n[Numbers]\n123 = onetwothree\n3.14 = pi\n[Database]\nname = 'John'\nage = 25.5\nID = 1 = one"

		s, _ := LoadFromString(str)
		sectionNames, _ := s.GetSectionNames()
		want := []string{"Credentials", "Numbers", "Database"}
		if !reflect.DeepEqual(sectionNames, want) {
			t.Errorf("expected %q but got %q", want, sectionNames)
		}
	})
	t.Run("no sections found error when passing string with no sections", func(t *testing.T) {
		str := "key1 = val\nkey2 = val"
		s, _ := LoadFromString(str)
		_, err := s.GetSectionNames()
		want := "error: no sections found"
		if (err.Error() != want) && err != nil {
			t.Errorf("expected error %q but got error %q", want, err)
		} else if err == nil {
			t.Errorf("expected error %q but received no error", want)
		}
	})
}

func TestGetSections(t *testing.T) {
	t.Run("no sections found error when passing noSections.ini", func(t *testing.T) {
		path := "./iniFiles/noSections.ini"
		file, _ := LoadFromFile(path)
		_, err := file.GetSections()
		want := "error: no sections found"
		if (err.Error() != want) && err != nil {
			t.Errorf("expected error %q but got error %q", want, err)
		} else if err == nil {
			t.Errorf("expected error %q but received no error", want)
		}
	})
	t.Run("no sections found error when passing string with no sections", func(t *testing.T) {
		str := "key1 = val\nkey2 = val"
		s, _ := LoadFromString(str)
		_, err := s.GetSections()
		want := "error: no sections found"
		if (err.Error() != want) && err != nil {
			t.Errorf("expected error %q but got error %q", want, err)
		} else if err == nil {
			t.Errorf("expected error %q but received no error", want)
		}
	})
}
func TestGet(t *testing.T) {
	t.Run("empty string value error when passing emptyString.ini file", func(t *testing.T) {
		path := "./iniFiles/emptyString.ini"
		file, _ := LoadFromFile(path)
		_, err := file.Get("Database", "key2")
		want := "error: value is an empty string"
		if (err.Error() != want) && err != nil {
			t.Errorf("expected error %q but got error %q", want, err)
		} else if err == nil {
			t.Errorf("expected error %q but received no error", want)
		}
	})
	t.Run("section not found error when passing file", func(t *testing.T) {
		path := "./iniFiles/correctFile.ini"
		file, _ := LoadFromFile(path)
		_, err := file.Get("Names", "key2")
		want := "section Names not found"
		if (err.Error() != want) && err != nil {
			t.Errorf("expected error %q but got error %q", want, err)
		} else if err == nil {
			t.Errorf("expected error %q but received no error", want)
		}
	})
	t.Run("key not found in section error when passing file", func(t *testing.T) {
		path := "./iniFiles/correctFile.ini"
		file, _ := LoadFromFile(path)
		_, err := file.Get("Database", "status")
		want := "key status not found in section Database"
		if (err.Error() != want) && err != nil {
			t.Errorf("expected error %q but got error %q", want, err)
		} else if err == nil {
			t.Errorf("expected error %q but received no error", want)
		}
	})
	t.Run("empty string value error when getting a key whose value is an empty string from an ini string", func(t *testing.T) {
		str := "[Database]\nkey1=val\nkey2= "
		s, _ := LoadFromString(str)
		_, err := s.Get("Database", "key2")
		want := "error: value is an empty string"
		if (err.Error() != want) && err != nil {
			t.Errorf("expected error %q but got error %q", want, err)
		} else if err == nil {
			t.Errorf("expected error %q but received no error", want)
		}
	})

	t.Run("section not found error when passing string", func(t *testing.T) {
		str := "[Credentials]\nuser = root\npassword = root\nport = 3000\n#also a comment\n[Numbers]\n123 = onetwothree\n3.14 = pi\n[Database]\nname = 'John'\nage = 25.5\nID = 1 = one"
		s, _ := LoadFromString(str)
		_, err := s.Get("Names", "key2")
		want := "section Names not found"
		if (err.Error() != want) && err != nil {
			t.Errorf("expected error %q but got error %q", want, err)
		} else if err == nil {
			t.Errorf("expected error %q but received no error", want)
		}
	})

	t.Run("key not found in section error when passing string", func(t *testing.T) {
		str := "[Credentials]\nuser = root\npassword = root\nport = 3000\n#also a comment\n[Numbers]\n123 = onetwothree\n3.14 = pi\n[Database]\nname = 'John'\nage = 25.5\nID = 1 = one"
		s, _ := LoadFromString(str)
		_, err := s.Get("Database", "status")
		want := "key status not found in section Database"
		if (err.Error() != want) && err != nil {
			t.Errorf("expected error %q but got error %q", want, err)
		} else if err == nil {
			t.Errorf("expected error %q but received no error", want)
		}

	})
}
func TestSaveToFile(t *testing.T) {
	t.Run("check if file has been created or not when called by a file struct", func(t *testing.T) {
		path := "./iniFiles/correctFile.ini"
		createdPath := "./iniFiles/createdFile.ini"
		file, _ := LoadFromFile(path)
		err := file.SaveToFile(createdPath)
		if err != nil {
			t.Errorf("Error saving file: %v", err)
		}
		_, err = os.Stat(createdPath)
		if os.IsNotExist(err) {
			t.Errorf("File not created: %v", err)
		}
	})
	t.Run("check if file has been created or not when called by an ini string struct", func(t *testing.T) {
		str := "[Credentials]\nuser = root\npassword = root\nport = 3000\n#also a comment\n[Numbers]\n123 = onetwothree\n3.14 = pi\n[Database]\nname = 'John'\nage = 25.5\nID = 1 = one"
		s, _ := LoadFromString(str)
		createdPath := "./iniFiles/createdFileFromString.ini"
		err := s.SaveToFile(createdPath)
		if err != nil {
			t.Errorf("Error saving file: %v", err)
		}
		_, err = os.Stat(createdPath)
		if os.IsNotExist(err) {
			t.Errorf("File not created: %v", err)
		}
	})
}
