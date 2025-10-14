package main

import (
	"fmt"
	"github.com/br36b/blog-aggregator/internal/config"
)

func main() {
	fmt.Println("test")

	fmt.Println("Reading file contents of configuration file")
	dbConfig, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Printing file contents of configuration file")
	fmt.Println(dbConfig)

	fmt.Println("Setting user in configuration file")
	dbConfig.SetUser("quack")

	dbConfig, err = config.Read()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Printing file contents of configuration file")
	fmt.Println(dbConfig)
}
