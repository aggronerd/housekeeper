// Copyright 2019 Gregory Doran <greg@gregorydoran.co.uk>. 
// All rights reserved.

package main

import (
	"github.com/aggronerd/housekeeper/cmd"
	"log"
)

func main() {
	err := cmd.Execute()

	if err != nil {
		log.Fatal(err)
	}
}


