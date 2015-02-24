package main

import (
	"fmt"
	"log"
	"os"
)

func bail(s string, m ...interface{}) {
	fmt.Printf("FATAL: "+s, m...)
	os.Exit(1)
}

func debugLogln(m ...interface{}) {
	if *debug {
		log.Println(m...)
	}
}

func debugLogf(s string, m ...interface{}) {
	if *debug {
		log.Printf(s, m...)
	}
}
