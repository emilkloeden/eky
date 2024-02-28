package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

type Attributes map[string]interface{}

type Eky struct {
	FilePath string
	Data     Attributes
}

// NewEky creates a new instance of Eky with a file path determined by the current user's home directory.
func NewEky() *Eky {
	filePath, err := getFilePath()
	if err != nil {
		fmt.Println("Error getting file path... Exiting")
		os.Exit(1)
	}
	data := Attributes{}
	return &Eky{FilePath: filePath, Data: data}
}

func (e *Eky) load() error {
	file, err := os.Open(e.FilePath)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	var data map[string]interface{}
	err = json.NewDecoder(file).Decode(&data)
	if err != nil {
		return fmt.Errorf("error decoding JSON: %v", err)
	}
	e.Data = data
	return nil
}

func (e *Eky) save() error {
	file, err := os.OpenFile(e.FilePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()
	err = json.NewEncoder(file).Encode(e.Data)
	if err != nil {
		return fmt.Errorf("error encoding JSON: %v", err)
	}
	return nil
}

// Get prints the value associated with the given key in the Eky instance.
func (e *Eky) List() {
	for key := range e.Data {
		fmt.Println(key)
	}
}

// Get prints the JSON representation of the value associated with the given key in the Eky instance,
// if the value can be encoded to JSON. Otherwise, it prints the value as is.
func (e *Eky) Get(key string) {
	value := e.Data[key]
	if value == nil {
		return
	}

	jsonValue, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		fmt.Println(value)
	} else {
		fmt.Println(string(jsonValue))
	}
}

// Set sets the value for the given key in the Eky instance.
// If the provided value is valid JSON, it is decoded and stored.
// Otherwise, the value is stored as a string.
func (e *Eky) Set(key, value string) error {
	var decodedValue interface{}
	err := json.Unmarshal([]byte(value), &decodedValue)
	if err != nil {
		// Unable to decode, store as a string
		e.Data[key] = value
	} else {
		e.Data[key] = decodedValue
	}
	err = e.save()
	if err != nil {
		return fmt.Errorf("error setting key: %v", key)
	}
	return nil
}

// Remove removes the values associated with the given keys from the Eky instance.
func (e *Eky) Remove(args ...string) {
	for _, key := range args {
		value := e.Data[key]
		if value != nil {
			delete(e.Data, key)
		}
	}
	e.save()
}

// Clear removes all data from the Eky instance.
func (e *Eky) Clear(args ...string) {
	e.Data = Attributes{}
	e.save()
}

func getFilePath() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("error getting current user: %v", err)
	}
	homeDir := currentUser.HomeDir
	return filepath.Join(homeDir, ".eky.json"), nil
}

func main() {
	getCmd := flag.NewFlagSet("get", flag.ExitOnError)
	setCmd := flag.NewFlagSet("set", flag.ExitOnError)
	rmCmd := flag.NewFlagSet("rm", flag.ExitOnError)

	if len(os.Args) < 2 {
		fmt.Println("Please specify a subcommand : `get`, `set`, `list`, `rm` or `clear`.")
		os.Exit(1)
	}
	e := NewEky()
	err := e.load()
	if err != nil {
		log.Fatalf("Error loading data: %v", err)
	}

	subCommandArgs := os.Args[2:]
	switch os.Args[1] {
	case "clear":
		e.Clear()
	case "list":
		e.List()
	case "get":
		getCmd.Parse(subCommandArgs)
		args := getCmd.Args()
		if len(args) < 1 {
			log.Fatal("Please provide a key to get.")
		}
		e.Get(args[0])
	case "set":
		setCmd.Parse(subCommandArgs)
		args := setCmd.Args()
		if len(args) < 2 {
			log.Fatal("Please provide key and value to set.")
		}
		e.Set(args[0], args[1])
	case "rm":
		rmCmd.Parse(subCommandArgs)
		args := rmCmd.Args()
		if len(args) < 1 {
			log.Fatal("Please provide at least one key to remove.")
		}
		e.Remove(args...)

	default:
		log.Fatalf("Invalid subcommand, please specify one of : `get`, `set`, `list`, `rm` or `clear`.")
		os.Exit(1)
	}

}
