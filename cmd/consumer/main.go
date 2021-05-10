package main

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"log"
)

type Specification struct {
	Debug bool `envconfig:"DEBUG" default:"false"`
	Port  int  `envconfig:"PORT" default:"8080"`
	Key string `envconfig:"KEY" default:"empty"`
}

func main() {
	var s Specification
	err := envconfig.Process("", &s)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%#v", s)
}
