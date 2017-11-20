package main

import (
	"./glassdoor"
	"./parser"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"os"
)

func saveCompany(company string, coll *mgo.Collection, gd glassdoor.Client) {
	result := glassdoor.Employer{}
	err := coll.Find(bson.M{"name": company}).One(&result)
	if err != nil {
		empl, err := gd.SearchEmployer(company)
		if err != nil {
			fmt.Println("No employer found:", company)
			err = coll.Insert(bson.M{"name": company})
			if err != nil {
				panic(company)
			}
		} else {
			fmt.Println("Glassdoor:", empl)
			err = coll.Insert(result)
			if err != nil {
				panic(company)
			}
		}
	} else {
		fmt.Println("Cached:", result)
	}
}

func saveJob(job parser.Job, coll *mgo.Collection) {
	var result interface{}
	err := coll.Find(bson.M{
		"company": job.Company,
		"title":   job.Title,
		"link":    job.Link,
	}).One(&result)

	if err != nil {
		fmt.Println("Saving job:", job)
		coll.Insert(job)
	} else {
		fmt.Println("Job exists:", job)
	}
}

func main() {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	companiesColl := session.DB("seeker").C("companies")
	jobsColl := session.DB("seeker").C("jobs")

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
			saveCompany(job.Company, companiesColl, gd)
			saveJob(job, jobsColl)
		}
	}
}
