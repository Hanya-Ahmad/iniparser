package iniparser

import (
	"os"
	"reflect"
	"sort"
	"testing"
)

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
		want := "error non ini file parsed"
		if (err.Error() != want) && err != nil {
			t.Errorf("expected error %q but got error %q", want, err)
		} else if err == nil {
			t.Errorf("expected error %q but received no error", want)
		}
	})
	t.Run("repeated section name error when passing repeatedSection.ini", func(t *testing.T) {
		path := "./iniFiles/repeatedSection.ini"
		ini := NewIniParser()
		err := ini.LoadFromFile(path)
		want := "error repeated section name Database"
		if (err.Error() != want) && err != nil {
			t.Errorf("expected error %q but got error %q", want, err)
		} else if err == nil {
			t.Errorf("expected error %q but received no error", want)
		}
	})
	t.Run("empty key error when passing emptyKey.ini", func(t *testing.T) {
		path := "./iniFiles/emptyKey.ini"
		ini := NewIniParser()
		err := ini.LoadFromFile(path)
		want := "error empty key at key key3"
		if (err.Error() != want) && err != nil {
			t.Errorf("expected error %q but got error %q", want, err)
		} else if err == nil {
			t.Errorf("expected error %q but received no error", want)
		}
	})
	t.Run("repeated key error when passing repeatedKey.ini", func(t *testing.T) {
		path := "./iniFiles/repeatedKey.ini"
		ini := NewIniParser()
		err := ini.LoadFromFile(path)
		want := "error repeated key name key1 in section Database the last assigned value is applied"
		if (err.Error() != want) && err != nil {
			t.Errorf("expected error %q but got error %q", want, err)
		} else if err == nil {
			t.Errorf("expected error %q but received no error", want)
		}
	})
	t.Run("missing assignment operator '='  error when passing missingEqual.ini", func(t *testing.T) {
		path := "./iniFiles/missingEqual.ini"
		ini := NewIniParser()
		err := ini.LoadFromFile(path)
		want := "error missing key value operator at key key3"
		if (err.Error() != want) && err != nil {
			t.Errorf("expected error %q but got error %q", want, err)
		} else if err == nil {
			t.Errorf("expected error %q but received no error", want)
		}
	})
}

func TestLoadFromString(t *testing.T) {
	t.Run("error is nil when correct ini string is parsed", func(t *testing.T) {
		str := ";comment\n[Credentials]\nuser = root\npassword = root\nport = 3000\n#also a comment\n[Numbers]\n123 = onetwothree\n3.14 = pi\n[Database]\nname = 'John'\nage = 25.5\nID = 1 = one"
		ini := NewIniParser()
		err := ini.LoadFromString(str)
		if err != nil {
			t.Errorf("correct ini string parsed but error %q occured", err)

		}
	})
	t.Run("repeated section name error when passing a string", func(t *testing.T) {
		str := "[Database]\nkey1 = val\n[Email]\nkey2 = val\n[Database]\nkey3 = val"
		ini := NewIniParser()
		err := ini.LoadFromString(str)
		want := "error repeated section name Database"
		if (err.Error() != want) && err != nil {
			t.Errorf("expected error %q but got error %q", want, err)
		} else if err == nil {
			t.Errorf("expected error %q but received no error", want)
		}
	})
	t.Run("empty key error when passing a string", func(t *testing.T) {
		str := "[Database]\nkey1 = val\nkey2 = val\nkey3 ="
		ini := NewIniParser()
		err := ini.LoadFromString(str)
		want := "error empty key at key key3"
		if (err.Error() != want) && err != nil {
			t.Errorf("expected error %q but got error %q", want, err)
		} else if err == nil {
			t.Errorf("expected error %q but received no error", want)
		}
	})
	t.Run("repeated key error when passing a string", func(t *testing.T) {
		str := "[Database]\nkey1 = val\nkey2 = val\nkey1 = anotherVal"
		ini := NewIniParser()
		err := ini.LoadFromString(str)
		want := "error repeated key name key1 in section Database the last assigned value is applied"
		if (err.Error() != want) && err != nil {
			t.Errorf("expected error %q but got error %q", want, err)
		} else if err == nil {
			t.Errorf("expected error %q but received no error", want)
		}
	})
	t.Run("missing assignment operator '='  error when passing a string", func(t *testing.T) {
		str := "[Database]\nkey1 = val\nkey2 = val\nkey3"
		ini := NewIniParser()
		err := ini.LoadFromString(str)
		want := "error missing key value operator at key key3"
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
		ini := NewIniParser()
		_ = ini.LoadFromFile(path)
		_, err := ini.GetSectionNames()
		if err != nil {
			t.Errorf("correct file  %q but error %q occured", path, err)
		}
	})
	t.Run("correct section names returned when passing correctFile.ini", func(t *testing.T) {
		path := "./iniFiles/correctFile.ini"
		ini := NewIniParser()
		_ = ini.LoadFromFile(path)
		sectionNames, _ := ini.GetSectionNames()
		want := []string{"Credentials", "Numbers", "Database"}
		sort.Strings(sectionNames)
		sort.Strings(want)
		if !reflect.DeepEqual(sectionNames, want) {
			t.Errorf("expected %q but got %q", want, sectionNames)
		}
	})
	t.Run("no sections found error when passing noSections.ini", func(t *testing.T) {
		path := "./iniFiles/noSections.ini"
		ini := NewIniParser()
		_ = ini.LoadFromFile(path)
		_, err := ini.GetSectionNames()
		want := "error no sections found"
		if (err.Error() != want) && err != nil {
			t.Errorf("expected error %q but got error %q", want, err)
		} else if err == nil {
			t.Errorf("expected error %q but received no error", want)
		}

	})
	t.Run("error is nil when passing correct ini string", func(t *testing.T) {
		str := "[Credentials]\nuser = root\npassword = root\nport = 3000\n#also a comment\n[Numbers]\n123 = onetwothree\n3.14 = pi\n[Database]\nname = 'John'\nage = 25.5\nID = 1 = one"
		ini := NewIniParser()
		_ = ini.LoadFromString(str)
		_, err := ini.GetSectionNames()
		if err != nil {
			t.Errorf("correct string  %q but error %q occured", str, err)
		}
	})
	t.Run("correct section names returned when passing correct string", func(t *testing.T) {
		str := "[Credentials]\nuser = root\npassword = root\nport = 3000\n#also a comment\n[Numbers]\n123 = onetwothree\n3.14 = pi\n[Database]\nname = 'John'\nage = 25.5\nID = 1 = one"

		ini := NewIniParser()
		_ = ini.LoadFromString(str)
		sectionNames, _ := ini.GetSectionNames()
		want := []string{"Credentials", "Numbers", "Database"}
		sort.Strings(sectionNames)
		sort.Strings(want)
		if !reflect.DeepEqual(sectionNames, want) {
			t.Errorf("expected %q but got %q", want, sectionNames)
		}
	})
	t.Run("no sections found error when passing string with no sections", func(t *testing.T) {
		str := "key1 = val\nkey2 = val"
		ini := NewIniParser()
		_ = ini.LoadFromString(str)
		_, err := ini.GetSectionNames()
		want := "error no sections found"
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
		ini := NewIniParser()
		_ = ini.LoadFromFile(path)
		_, err := ini.GetSections()
		want := "error no sections found"
		if (err.Error() != want) && err != nil {
			t.Errorf("expected error %q but got error %q", want, err)
		} else if err == nil {
			t.Errorf("expected error %q but received no error", want)
		}
	})
	t.Run("no sections found error when passing string with no sections", func(t *testing.T) {
		str := "key1 = val\nkey2 = val"
		ini := NewIniParser()
		_ = ini.LoadFromString(str)
		_, err := ini.GetSections()
		want := "error no sections found"
		if (err.Error() != want) && err != nil {
			t.Errorf("expected error %q but got error %q", want, err)
		} else if err == nil {
			t.Errorf("expected error %q but received no error", want)
		}
	})
}

func TestGet(t *testing.T) {
	t.Run("empty string value returned and exists is true when passing emptyString.ini file", func(t *testing.T) {
		path := "./iniFiles/emptyString.ini"
		ini := NewIniParser()
		_ = ini.LoadFromFile(path)
		_, exists, err := ini.Get("Database", "key2")

		if !(exists == true && err == nil) {
			t.Errorf("expected empty string and exist to be true but received error %q", err)

		}
	})
	t.Run("section not found error when passing file", func(t *testing.T) {
		path := "./iniFiles/correctFile.ini"
		ini := NewIniParser()
		_ = ini.LoadFromFile(path)
		_, _, err := ini.Get("Names", "key2")
		want := "error section Names not found"
		if (err.Error() != want) && err != nil {
			t.Errorf("expected error %q but got error %q", want, err)
		} else if err == nil {
			t.Errorf("expected error %q but received no error", want)
		}
	})
	t.Run("key not found in section error when passing file", func(t *testing.T) {
		path := "./iniFiles/correctFile.ini"
		ini := NewIniParser()
		_ = ini.LoadFromFile(path)
		_, _, err := ini.Get("Database", "status")
		want := "error key status not found in section Database"
		if (err.Error() != want) && err != nil {
			t.Errorf("expected error %q but got error %q", want, err)
		} else if err == nil {
			t.Errorf("expected error %q but received no error", want)
		}
	})
	t.Run("empty string value returned and exists is true when passing an ini string with empty string", func(t *testing.T) {
		str := "[Database]\nkey1=val\nkey2= "
		ini := NewIniParser()
		_ = ini.LoadFromString(str)
		_, exists, err := ini.Get("Database", "key2")
		if !(exists == true && err == nil) {
			t.Errorf("expected empty string and exist to be true but received error %q", err)

		}
	})

	t.Run("section not found error when passing string", func(t *testing.T) {
		str := "[Credentials]\nuser = root\npassword = root\nport = 3000\n#also a comment\n[Numbers]\n123 = onetwothree\n3.14 = pi\n[Database]\nname = 'John'\nage = 25.5\nID = 1 = one"
		ini := NewIniParser()
		_ = ini.LoadFromString(str)
		_, _, err := ini.Get("Names", "key2")
		want := "error section Names not found"
		if (err.Error() != want) && err != nil {
			t.Errorf("expected error %q but got error %q", want, err)
		} else if err == nil {
			t.Errorf("expected error %q but received no error", want)
		}
	})

	t.Run("key not found in section error when passing string", func(t *testing.T) {
		str := "[Credentials]\nuser = root\npassword = root\nport = 3000\n#also a comment\n[Numbers]\n123 = onetwothree\n3.14 = pi\n[Database]\nname = 'John'\nage = 25.5\nID = 1 = one"
		ini := NewIniParser()
		_ = ini.LoadFromString(str)
		_, _, err := ini.Get("Database", "status")
		want := "error key status not found in section Database"
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
		ini := NewIniParser()
		_ = ini.LoadFromFile(path)
		err := ini.SaveToFile(createdPath)
		if err != nil {
			t.Errorf("error saving file %v", err)
		}
		_, err = os.Stat(createdPath)
		if os.IsNotExist(err) {
			t.Errorf("error file not created %v", err)
		} else {
			// Clean up the created file
			err = os.Remove(createdPath)
			if err != nil {
				t.Errorf("error removing file %v", err)
			}
		}
	})
	t.Run("check if file has been created or not when called by an ini string struct", func(t *testing.T) {
		str := "[Credentials]\nuser = root\npassword = root\nport = 3000\n#also a comment\n[Numbers]\n123 = onetwothree\n3.14 = pi\n[Database]\nname = 'John'\nage = 25.5\nID = 1 = one"
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

// func TestToString(t *testing.T) {
// 	t.Run("check if the string representation of the ini file is correct", func(t *testing.T) {
// 		path := "./iniFiles/correctFile.ini"
// 		ini := NewIniParser()
// 		_ = ini.LoadFromFile(path)
// 		expectedStr := "[Credentials]\nuser = root\npassword = root\nport = 3000\n[Numbers]\n123 = onetwothree\n3.14 = pi\n[Database]\nname = 'John'\nage = 25.5\nID = 1 = one"
// 		if iniStr := ini.ToString(); iniStr != expectedStr {
// 			t.Errorf("expected %v but got %v", expectedStr, iniStr)
// 		}
// 	})
// }
