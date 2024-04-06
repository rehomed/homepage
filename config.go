package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

var config Configuration

type Configuration struct {
	DefaultPage string              `yaml:"default_page"`
	Pages       []ConfigurationPage `yaml:"pages"`
}

type ConfigurationPage struct {
	Title   string                    `yaml:"title" json:"title"`     // page title
	Path    string                    `yaml:"path" json:"path"`       // the pathname of the page (prefixed with /)
	Inject  ConfigurationPageInject   `yaml:"inject" json:"inject"`   // content to inject to the page
	Search  ConfigurationPageSearch   `yaml:"search" json:"search"`   // search engine config
	Widgets []ConfigurationPageWidget `yaml:"widgets" json:"widgets"` // widgets injected in the page
	Links   []ConfigurationPageLink   `yaml:"links" json:"links"`
}

type ConfigurationPageInject struct {
	CSS        string `yaml:"css" json:"css"`
	HTML       string `yaml:"html" json:"html"`
	JavaScript string `yaml:"js" json:"js"`
}

type ConfigurationPageSearch struct {
	Enabled bool `yaml:"enabled" json:"enabled"` // whether to show the search bar
	Label   bool `yaml:"label" json:"label"`     // "Search with [label]" as search bar placeholder
	URL     bool `yaml:"url" json:"url"`         // the placeholder URL for the service (i.e https://google.com/search?q=%s)
}

type ConfigurationPageWidget struct {
	Bultin string            `yaml:"builtin"` // ID of built-in widget
	Inject string            `yaml:"inject"`  // HTML metadata to inject
	KV     map[string]string // extra data passed in query params
}

type ConfigurationPageLink struct {
	Title string `yaml:"title" json:"title"`
	Icon  string `yaml:"icon" json:"icon"`
	URL   string `yaml:"url" json:"url"`
	// StatusURL TODO!!
}

func init() {
	// if in devmode, load dev.env
	if os.Getenv("ENV") != "production" {
		godotenv.Load("dev.env")
	}

	cwd, err := os.Getwd()
	if err != nil {
		slog.Error("failed to get cwd", "err", err)
		os.Exit(1)
	}

	// data path
	dataPath := os.Getenv("DATA_PATH")
	if dataPath == "" {
		dataPath = cwd + "/data"
		slog.Warn(fmt.Sprintf("DATA_PATH is unset. Using %q (default)", dataPath))

		os.Setenv("DATA_PATH", dataPath)
	}

	slog.Debug("using data path", "name", dataPath)

	_, err = os.Stat(dataPath)
	if os.IsNotExist(err) {
		slog.Warn("DATA_PATH is missing. Creating it now...", "name", dataPath)

		err = os.MkdirAll(dataPath, 0770)
		if err != nil {
			slog.Error("failed creating data directory:", "err", err)
			os.Exit(1)
		}
	}

	// http listener
	port := os.Getenv("PORT")
	if port == "" {
		slog.Debug("PORT is unset, defaulting to 8080")
		port = "8080"
		os.Setenv("PORT", port)
	}

	listen := os.Getenv("LISTEN")
	if listen == "" {
		listen = "0.0.0.0:" + port
		os.Setenv("LISTEN", listen)
		slog.Debug("overriding LISTEN to", "val", listen, "why", "unset")
	} else if !strings.Contains(listen, ":") {
		listen += ":" + port
		os.Setenv("LISTEN", listen)
		slog.Debug("overriding LISTEN to", "val", listen, "why", "no port in value")
	}

	// config data
	configContent := os.Getenv("CONFIG")
	configFile := os.Getenv("CONFIG_FILE")
	if configContent == "" && configFile == "" {
		slog.Error("unable to start: no config found (variable CONFIG or CONFIG_FILE)")
		os.Exit(1)
	}

	if configContent == "" && configFile != "" {
		slog.Debug("loading config from file", "name", configFile)
		configFileData, err := os.ReadFile(configFile)
		if err != nil {
			slog.Error("unable to start: config unreadable:", "err", err)
			os.Exit(1)
		}

		configContent = string(configFileData)
	} else {
		slog.Debug("loading config from env")
	}

	// cfg parsing
	var cfg Configuration
	err = yaml.Unmarshal([]byte(configContent), &cfg)
	if err != nil {
		slog.Error("failed unmarshaling config yaml", "err", err)
		os.Exit(1)
	}

	spew.Dump(cfg)
	config = cfg
}

func resolvePage(path string, cfg Configuration) *ConfigurationPage {
	for _, p := range cfg.Pages {
		if p.Path == path {
			return &p
		}
	}

	return nil
}
