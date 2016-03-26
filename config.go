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

type HelpError struct {
	Message string
}

func (e *HelpError) Error() string { return e.Message }

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
		if err := ParseFile(file, cfg); err != nil {
			return nil, err
		}
	}
	parser := flags.NewParser(cfg, flags.HelpFlag|flags.PassDoubleDash|flags.IgnoreUnknown)
	if _, err := parser.Parse(); err != nil {
		if e, ok := err.(*flags.Error); ok && e.Type == flags.ErrHelp {
			return nil, &HelpError{e.Message}
		}
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
	app := path.Base(os.Args[0])
	for _, dir := range []string{
		"",
		path.Join(os.Getenv("HOME"), "."+app),
		"/etc/" + app,
	} {
		for _, file := range []string{
			"config.yaml",
			"config.yml",
			"config.json",
		} {
			fileName := path.Join(dir, file)
			if fileExists(fileName) {
				return fileName, nil
			}
		}
	}
	return "", nil
}

func ParseFile(file string, cfg interface{}) error {
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
