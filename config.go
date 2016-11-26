package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
)

// Initialize global variables
var Config config

// Initialize config options
type config struct {
	PortDir    string
	Order      []string
	Alias      [][]string
	IndentChar string
	Pull       map[string]pull
}

type pull struct {
	Url    string
	Branch string
}

func init() {
	// Read out config
	f, err := ioutil.ReadFile("./config/config.toml")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not read config!")
		os.Exit(1)
	}

	// Decode config
	if _, err := toml.Decode(string(f), &Config); err != nil {
		fmt.Fprintln(os.Stderr, "Could not decode config!")
		os.Exit(1)
	}
}
