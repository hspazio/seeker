package main

import (
	"../seeker-go/glassdoor"
	"../seeker-go/parser"
	"fmt"
	"os"
)

func main() {
	parsers := []parser.Parser{
		parser.New("remoteok"),
		parser.New("weworkremotely"),
		parser.New("remoteco"),
	}

	for _, parser := range parsers {
		jobs, err := parser.Parse()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(jobs)
	}

	// test Glassdoor api
	gd := glassdoor.Client{Id: os.Getenv("GLASSDOOR_ID"), Key: os.Getenv("GLASSDOOR_KEY")}
	empl, err := gd.SearchEmployer("Google")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(empl)
	}
}
