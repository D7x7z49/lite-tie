package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/spf13/viper"
)

var (
	registryDir    = filepath.Dir(os.Args[0])
	registryType   = "json"
	registryName   = ".registry"
	registryFile   = filepath.Join(registryDir, registryName+"."+registryType)
	logOutput      = os.Stdout
	registryLogger = log.New(logOutput, "", log.LstdFlags)
	regMu          sync.Mutex
)

type Entry struct {
	Source string `mapstructure:"source" json:"source"`
	Link   string `mapstructure:"link" json:"link"`
	User   string `mapstructure:"user" json:"user"`
}

func InitRegistry() error {
	viper.AddConfigPath(registryDir)
	viper.SetConfigType(registryType)
	viper.SetConfigName(registryName)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			viper.SetDefault("", nil)
			if err := viper.WriteConfigAs(registryFile); err != nil {
				registryLogger.Printf("create failed: %v", err)
				return fmt.Errorf("create failed")
			}
			viper.SetConfigFile(registryFile)
			if err := viper.ReadInConfig(); err != nil {
				registryLogger.Printf("read failed: %v", err)
				return fmt.Errorf("read failed")
			}
		} else {
			registryLogger.Printf("read failed: %v", err)
			return fmt.Errorf("read failed")
		}
	}
	return nil
}

func SaveRegistry() error {
	if err := viper.WriteConfig(); err != nil {
		registryLogger.Printf("save failed: %v", err)
		return fmt.Errorf("save failed")
	}
	return nil
}

func AddEntry(alias, source string) error {
	regMu.Lock()
	defer regMu.Unlock()
	viper.Set(alias, Entry{
		Source: source,
		Link:   filepath.Join(registryDir, alias),
		User:   os.Getenv("USERNAME"),
	})
	return SaveRegistry()
}

func RemoveEntry(alias string) error {
	regMu.Lock()
	defer regMu.Unlock()
	viper.Set(alias, nil)
	return SaveRegistry()
}

func GetEntries() (map[string]Entry, error) {
	regMu.Lock()
	defer regMu.Unlock()
	entries := make(map[string]Entry)
	if err := viper.Unmarshal(&entries); err != nil {
		registryLogger.Printf("parse failed: %v", err)
		return nil, fmt.Errorf("parse failed")
	}
	return entries, nil
}

func SetLogOutput(output *os.File) {
	registryLogger.SetOutput(output)
}
