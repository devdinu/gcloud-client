package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/devdinu/gcloud-client/logger"
)

type Defaults struct {
	User         string `json:"user"`
	DBFile       string `json:"db_file"`
	SSHFile      string `json:"ssh_key"`
	LogLevel     string `json:"log_level"`
	TemplatesDir string `json:"templates_dir"`
}

var conf Defaults

func loadDefaults(configFile string) (*Defaults, error) {
	homeDir := os.Getenv("HOME")
	if configFile == "" {
		return nil, errors.New("No configuration file mentioned")
	}

	configDir := filepath.Dir(configFile)
	appConfig := Defaults{
		User:         os.Getenv("USER"),
		SSHFile:      homeDir + string(os.PathSeparator) + "ssh" + string(os.PathSeparator) + "id_rsa.pub",
		DBFile:       configDir + string(os.PathSeparator) + "hosts.db",
		TemplatesDir: getTemplatesDir(),
		LogLevel:     "info",
	}

	if _, err := os.Stat(configFile); err != nil {
		if os.IsNotExist(err) {
			if err = setDefaultConfig(configFile, appConfig); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
		return &appConfig, err
	} else {
		f, err := os.Open(configFile)
		if err != nil {
			return nil, fmt.Errorf("config file %s open failed with error %v", configFile, err)
		}

		err = json.NewDecoder(f).Decode(&appConfig)
		if err != nil {
			logger.Debugf("Try removing config file %s", configFile)
			return nil, fmt.Errorf("reading config file %s failed with error %v", configFile, err)
		}
	}
	logger.Debugf("[Config] Loaded default config: %+v from file %s", appConfig, configFile)

	return &appConfig, nil
}

func setDefaultConfig(configFile string, appConfig Defaults) error {
	err := os.MkdirAll(filepath.Dir(configFile), os.ModePerm)
	if err != nil {
		return err
	}
	f, err := os.Create(configFile)
	if err != nil {
		return err
	}
	err = json.NewEncoder(f).Encode(appConfig)
	if err != nil {
		return err
	}
	logger.Infof("[Config] created default configuration %s", configFile)
	return nil
}

func getTemplatesDir() string {
	defaultTemplatesDir := strings.Join([]string{os.Getenv("HOME"), ".config", "tmuxinator"}, string(os.PathSeparator))
	templatesDir := strings.Join([]string{"/usr", "local", "Cellar", "gcloud-client"}, string(os.PathSeparator))
	files, err := ioutil.ReadDir(templatesDir)
	if err != nil {
		logger.Debugf("[Config] Couldn't figure out brew cellar directory for templates: %v", err)
	}
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".yml") {
			return templatesDir
		}
	}
	return defaultTemplatesDir
}
