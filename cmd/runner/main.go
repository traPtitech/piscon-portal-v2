package main

import (
	"log"

	"github.com/traPtitech/piscon-portal-v2/runner"
)

func main() {
	r := runner.Prepare()

	for {
		if err := r.Run(); err != nil {
			log.Printf("error: %v", err)
		}
	}
}
