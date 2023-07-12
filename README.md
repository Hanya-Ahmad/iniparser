# Overview
## iniparser is a Go package for parsing, manipulating, and generating INI files. It provides a simple and easy-to-use interface for reading and writing INI files. 

# Installation
## To install iniparser run this command
```
go get github.com/codescalersinternships/INIParser-Hanya
```

# How to use
## 1. Import iniparser in your project
```
import github.com/codescalersinternships/INIParser-Hanya 
```
## 2. Initialize a new iniparser object
```
parser := iniparser.NewIniParser()
```
## 3. Use the package's methods on your parser, for example:
```
err:= parser.LoadFromFile("path")
```

# API
## Methods
- ## NewIniParser() *IniParser 
    ### *Creates a new IniParser object*<br>

- ## LoadFromFile(path string) error 
    ### *Loads ini data from a file at the specified path*<br>

- ## LoadFromString(str string) error 
    ### *Loads ini data from a string*<br>

- ## GetSectionNames() ([]string, error)
    ### *Returns the entire ini data as a map of sections and their key/value pairs* <br>

- ## GetSections() (sections, error)
    ### *Returns the names of all sections in the ini data*<br>

- ## Get(sectionName string, key string) (string, bool, error)
    ### *Gets the value of a key in a specific section*<br>

- ## Set(sectionName string, key string, val string)
    ### *Sets the value of a key in a specific section*<br>

- ## String() string
    ### *Converts the ini data to a string*<br>

- ## SaveToFile(path string) error
    ### *Saves the ini data to a file at the specified path*<br>


## Syntax Errors
### The following syntax errors are defined and tested in the package:
- ### ErrRepeatedSectionName
- ### ErrEmptyKey
- ### ErrRepeatedKeyName
- ### ErrMissingKeyValueOperator
- ### ErrNonIniFileParsed
- ### ErrNoSectionsFound
- ### ErrSectionNotFound
- ### ErrKeyNotFound
- ### ErrNonIniFilePath

