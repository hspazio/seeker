package main

import (
	"./glassdoor"
	"./parser"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"os"
)

func main() {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	companiesColl := session.DB("seeker").C("companies")

	gd := glassdoor.Client{Id: os.Getenv("GLASSDOOR_ID"), Key: os.Getenv("GLASSDOOR_KEY")}

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
		for _, job := range jobs {
			result := glassdoor.Employer{}
			err = companiesColl.Find(bson.M{"name": job.Company}).One(&result)
			if err != nil {
				empl, err := gd.SearchEmployer(job.Company)
				if err != nil {
					// TODO: save empty record to db so we won't check glassdoor again
					fmt.Println("No employer found:", job.Company)
				} else {
					// TODO: save new record to db
					fmt.Println("Glassdoor:", empl)
				}
			} else {
				fmt.Println("Cached:", result)
			}

			// TODO: save job if it doesn't exist
			fmt.Println("Saving job:", job)
		}
	}
}
