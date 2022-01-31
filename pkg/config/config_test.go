package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"gopkg.in/yaml.v3"
)

type TestA struct {
	A1 string `yaml:"a1"`
	A2 string `yaml:"a2"`
}

type TestB struct {
	B1 string `yaml:"b1"`
	B2 string `yaml:"b2"`
}

type TestConfigAB struct {
	A TestA
	B TestB
}

type TestConfigB struct {
	B TestB
}

func WriteFile(filename string, content []byte) error {
	err := ioutil.WriteFile(filename, content, 0644)
	if err != nil {
		panic("Unable to write data into the file")
	}
	return nil
}

func CreateFiles() {
	TestConfig1 := &TestConfigAB{
		A: TestA{
			A1: "${ENV_VALUE_A1}",
			A2: "Value A2",
		},
		B: TestB{
			B1: "XXXXXXX",
			B2: "XXXXXXX",
		},
	}
	yamlData, err := yaml.Marshal(&TestConfig1)
	if err != nil {
		panic(err)
	}
	if WriteFile("file1.yml", yamlData) != nil {
		panic(err)
	}

	TestConfig2 := &TestConfigB{
		B: TestB{
			B1: "Value B1",
			B2: "${ENV_VALUE_B2}",
		},
	}
	yamlData, err = yaml.Marshal(&TestConfig2)
	if err != nil {
		panic(err)
	}
	if WriteFile("file2.yml", yamlData) != nil {
		panic(err)
	}
}

func RemoveFiles() {
	os.Remove("file1.yml")
	os.Remove("file2.yml")
}

func TestNewConfig(t *testing.T) {

	os.Setenv("ENV_VALUE_A1", "Qwerty1234")
	os.Setenv("ENV_VALUE_B2", "Qwerty1234")

	CreateFiles()

	conf, err := NewConfig("file1.yml", "file2.yml")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v", conf)

	RemoveFiles()

	if conf.A.A1 != "Qwerty1234" {
		t.Errorf("conf.A.A1 failed")
	}
	if conf.A.A2 != "Value A2" {
		t.Errorf("conf.A.A1 failed")
	}
	if conf.B.B1 != "Value B1" {
		t.Errorf("conf.A.A1 failed")
	}
	if conf.B.B2 != "Qwerty1234" {
		t.Errorf("conf.A.A1 failed")
	}
}
