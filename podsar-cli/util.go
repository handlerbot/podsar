package main

import (
	"fmt"
	"strconv"
)

func prettyPrint(lines [][2]string) {
	titleMax := 0
	for _, l := range lines {
		thisLen := len(l[0])
		if thisLen > titleMax {
			titleMax = thisLen
		}
	}
	numMax := len(strconv.Itoa(len(lines)))
	for i, l := range lines {
		fmt.Printf("%[1]*d) %-[3]*s  %s\n", numMax, i+1, titleMax, l[0], l[1])
	}
}
