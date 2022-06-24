package main

import (
	"fmt"
	"os"

	"github.com/jmpa-io/up-go"
)

func main() {

	// setup client.
	c, err := up.New("xxxx")
	if err != nil {
		fmt.Printf("failed to setup client: %v\n", err)
		os.Exit(1)
	}

	// do something with the client..
	// like send a ping to the api to check if the token is valid.
	p, err := c.Ping()
	if err != nil {
		fmt.Printf("failed to ping: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("%+v\n", p)
}
