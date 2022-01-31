package main

import (
	"fmt"

	"github.com/dimtass/go-multiconfig/pkg/config"
)

func main() {

	conf, err := config.NewConfig("config.yml", "override-config.yml")
	if err != nil {
		panic(err)
	}
	fmt.Printf("config: %v\n", conf)
}
