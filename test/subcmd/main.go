package main

import (
	"flag"
	"log"
	"os"
	"strconv"
	"time"
)

var file string
var dur time.Duration

func init() {
	flag.DurationVar(&dur, "duration", time.Millisecond*100, "set running duration")
	flag.Parse()
	file = flag.Arg(0)
	if file == "" {
		panic("pass file argument")
	}
}

func main() {
	time.Sleep(dur)
	var prev int

	if bytes, err := os.ReadFile(file); os.IsNotExist(err) {
		// do nothing
	} else if err != nil {
		log.Fatal(err)
	} else {
		i, err := strconv.Atoi(string(bytes))
		if err != nil {
			log.Fatal(err)
		}
		prev = i
	}

	err := os.WriteFile(file, []byte(strconv.Itoa(prev+1)), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

}
