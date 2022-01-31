package config

import (
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type ConfigA struct {
	A1 string `yaml:"a1"`
	A2 string `yaml:"a2"`
}

type ConfigB struct {
	B1 string `yaml:"b1"`
	B2 string `yaml:"b2"`
}

type Config struct {
	A ConfigA
	B ConfigB
}

func NewConfig(filenames ...string) (*Config, error) {
	if len(filenames) <= 0 {
		log.WithFields(log.Fields{"module": "config"}).Fatalf("You must provide at least one filename for reading Values")
	}
	var resultValues map[string]interface{}
	// Create config structure
	config := &Config{}
	for _, filename := range filenames {
		err := ValidateConfigPath(filename)
		if err != nil {
			log.WithFields(log.Fields{"module": "config"}).Fatalf("%v", err)
		}
		var override map[string]interface{}
		bs, err := ioutil.ReadFile(filename)
		if err != nil {
			log.WithFields(log.Fields{"module": "config"}).Fatalf("%v", err)
			continue
		}
		if err := yaml.Unmarshal(bs, &override); err != nil {
			log.WithFields(log.Fields{"module": "config"}).Fatalf("%v", err)
			continue
		}

		//check if is nil. This will only happen for the first filename
		if resultValues == nil {
			resultValues = override
		} else {
			for k, v := range override {
				resultValues[k] = v
			}
		}

	}
	confContent, err := yaml.Marshal(resultValues)
	// expand environment variables
	confContent = []byte(os.ExpandEnv(string(confContent)))
	if err != nil {
		log.WithFields(log.Fields{"module": "config"}).Fatalf("%v", err)
	}
	if err := yaml.Unmarshal(confContent, config); err != nil {
		log.WithFields(log.Fields{"module": "config"}).Fatalf("%v", err)
	}

	return config, nil
}

// ValidateConfigPath just makes sure, that the path provided is a file,
// that can be read
func ValidateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	return nil
}
