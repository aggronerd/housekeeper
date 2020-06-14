// Copyright 2019 Gregory Doran <greg@gregorydoran.co.uk>. 
// All rights reserved.

package time

import "log"

func FatalIfErr(err error) {
	if err != nil {
		log.Fatalf("Got fatal error: %s", err)
	}
}

func TruncateString(str string, length int) string {
	output := str
	if len(str) > length {
		if length > 3 {
			length -= 3
		}
		output = str[0:length] + "..."
	}
	return output
}
