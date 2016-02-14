package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path"

	"github.com/jessevdk/go-flags"
	"gopkg.in/yaml.v2"
)

func Parse(cfg interface{}) error {
	_, err := ParseCommand(cfg)
	return err
}

func ParseCommand(cfg interface{}) (*flags.Command, error) {
	file, err := getConfigFileName()
	if err != nil {
		return nil, err
	}
	if file != "" {
		if err := parseConfigFile(file, cfg); err != nil {
			return nil, err
		}
	}
	parser := flags.NewParser(cfg, flags.HelpFlag|flags.PassDoubleDash)
	if _, err := parser.Parse(); err != nil {
		return nil, err
	}
	return parser.Command.Active, nil
}

func getConfigFileName() (string, error) {
	var fileConfig struct {
		ConfigFile string `long:"config"`
	}
	if _, err := flags.NewParser(&fileConfig, flags.IgnoreUnknown).Parse(); err != nil {
		return "", err
	}
	if fileConfig.ConfigFile != "" {
		return fileConfig.ConfigFile, nil
	}
	for _, file := range []string{
		"config.yml",
		"config.yaml",
		"config.json",
	} {
		if fileExists(file) {
			return file, nil
		}
	}
	return "", nil
}

func parseConfigFile(file string, cfg interface{}) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	switch path.Ext(file) {
	case ".json":
		return json.NewDecoder(f).Decode(cfg)
	case ".yml", ".yaml":
		in, err := ioutil.ReadAll(f)
		if err != nil {
			return err
		}
		return yaml.Unmarshal(in, cfg)
	}
	return errors.New("unsupported config file format: " + file)
}

func fileExists(file string) bool {
	f, err := os.Open(file)
	if err != nil {
		return false
	}
	f.Close()
	return true
}
