package main

import (
	"fmt"
	"github.com/br36b/blog-aggregator/internal/config"
)

func main() {
	dbConfig, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}

	dbConfig.SetUser("quack")

	dbConfig, err = config.Read()
	if err != nil {
		fmt.Println(err)
	}

}
