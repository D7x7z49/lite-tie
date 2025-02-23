package config

import (
	"crypto/sha256"
	"fmt"
	"io"
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
	Source    string `mapstructure:"source" json:"source"`
	Link      string `mapstructure:"link" json:"link"`
	User      string `mapstructure:"user" json:"user"`
	Hash      string `mapstructure:"hash" json:"hash"`
	Available bool   `mapstructure:"available" json:"available"`
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

func AddEntry(alias, source, link string) error {
	hash, err := computeFileHash(source)
	if err != nil {
		return fmt.Errorf("hash failed: %v", err)
	}

	regMu.Lock()
	defer regMu.Unlock()
	viper.Set(alias, Entry{
		Source:    source,
		Link:      link,
		User:      os.Getenv("USERNAME"),
		Hash:      hash,
		Available: true,
	})
	return SaveRegistry()
}

func RemoveEntries(aliases []string) error {
	regMu.Lock()
	defer regMu.Unlock()

	all := viper.AllSettings()
	for _, alias := range aliases {
		delete(all, alias)
	}

	viper.Reset()
	viper.AddConfigPath(registryDir)
	viper.SetConfigType(registryType)
	viper.SetConfigName(registryName)

	for k, v := range all {
		viper.Set(k, v)
	}

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

func computeFileHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func UpdateEntries() error {
	entries, err := GetEntries()
	if err != nil {
		return err
	}

	regMu.Lock()
	defer regMu.Unlock()
	for alias, entry := range entries {
		fileInfo, err := os.Stat(entry.Source)
		available := err == nil && !fileInfo.IsDir()
		if available {
			hash, err := computeFileHash(entry.Source)
			if err != nil || hash != entry.Hash {
				available = false
			}
		}
		viper.Set(alias, Entry{
			Source:    entry.Source,
			Link:      entry.Link,
			User:      entry.User,
			Hash:      entry.Hash,
			Available: available,
		})
	}
	return SaveRegistry()
}
