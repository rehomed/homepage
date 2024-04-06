package main

import (
	"fmt"
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Configuration struct {
	DefaultPage string              `yaml:"default_page"`
	Pages       []ConfigurationPage `yaml:"pages"`
}

type ConfigurationPage struct {
	Title   string                    `yaml:"title"`   // page title
	Path    string                    `yaml:"path"`    // the pathname of the page (prefixed with /)
	Inject  ConfigurationPageInject   `yaml:"inject"`  // content to inject to the page
	Search  ConfigurationPageSearch   `yaml:"search"`  // search engine config
	Widgets []ConfigurationPageWidget `yaml:"widgets"` // widgets injected in the page
	Links   []ConfigurationPageLink   `yaml:"links"`
}

type ConfigurationPageInject struct {
	CSS        string `yaml:"css"`
	HTML       string `yaml:"html"`
	JavaScript string `yaml:"js"`
}

type ConfigurationPageSearch struct {
	Enabled bool `yaml:"enabled"` // whether to show the search bar
	Label   bool `yaml:"label"`   // "Search with [label]" as search bar placeholder
	URL     bool `yaml:"url"`     // the placeholder URL for the service (i.e https://google.com/search?q=%s)
}

type ConfigurationPageWidget struct {
	Bultin string            `yaml:"builtin"` // ID of built-in widget
	Inject string            `yaml:"inject"`  // HTML metadata to inject
	KV     map[string]string // extra data passed in query params
}

type ConfigurationPageLink struct {
	Title string `yaml:"title"`
	Icon  string `yaml:"icon"`
	URL   string `yaml:"url"`
	// StatusURL TODO!!
}

func init() {
	// if in devmode, load dev.env
	if os.Getenv("ENV") != "production" {
		godotenv.Load("dev.env")
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("failed to get cwd:", err)
	}

	// test/validate env config
	dataPath := os.Getenv("DATA_PATH")
	if dataPath == "" {
		dataPath = cwd + "/data"
		fmt.Printf("WARNING: DATA_PATH is unset. Using %q (default)\n", dataPath)
		os.Setenv("DATA_PATH", dataPath)
	}

	_, err = os.Stat(dataPath)
	if os.IsNotExist(err) {
		fmt.Printf("WARNING: DATA_PATH is missing. Creating it now... (%q)\n", dataPath)

		err = os.MkdirAll(dataPath, 0770)
		if err != nil {
			log.Fatal("failed creating data directory:", dataPath, err)
		}
	}

	// config data
	configContent := os.Getenv("CONFIG")
	configFile := os.Getenv("CONFIG_FILE")
	if configContent == "" && configFile == "" {
		log.Fatal("unable to start: no config found (variable CONFIG or CONFIG_FILE)")
	}

	if configContent == "" && configFile != "" {
		configFileData, err := os.ReadFile(configFile)
		if err != nil {
			log.Fatal("unable to start: config unreadable:", err)
		}

		configContent = string(configFileData)
	}

	// config parsing
	var config Configuration
	err = yaml.Unmarshal([]byte(configContent), &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	spew.Dump(config)
}
