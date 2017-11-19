package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strings"
)

type Job struct {
	company string
	title   string
	link    string
}

type parser struct {
	name string
}

func (p parser) Parse() ([]Job, error) {
	var jobs []Job
	var err error
	switch p.name {
	case "remoteok":
		jobs, err = parseRemoteok()
	case "weworkremotely":
		jobs, err = parseWeworkremotely()
	case "remoteco":
		jobs, err = parseRemoteco()
	}

	return jobs, err
}

func parseRemoteok() ([]Job, error) {
	site := "https://remoteok.io/remote-dev-jobs"
	var jobs []Job

	res, err := http.Get(site)
	if err != nil {
		return nil, err
	}

	doc, _ := goquery.NewDocumentFromResponse(res)
	jobList := doc.Find("#jobsboard").Find(".job")

	for i := range jobList.Nodes {
		company, _ := jobList.Eq(i).Attr("data-company")
		title, _ := jobList.Eq(i).Attr("data-search")
		link, _ := jobList.Eq(i).Attr("data-url")
		job := Job{
			company: strings.TrimSpace(company),
			title:   strings.TrimSpace(title),
			link:    strings.TrimSpace(link),
		}

		if company != "" && title != "" {
			jobs = append(jobs, job)
		}
	}

	return jobs, err
}

func parseWeworkremotely() ([]Job, error) {
	site := "https://weworkremotely.com/categories/2-programming/jobs#intro"
	var jobs []Job

	res, err := http.Get(site)
	if err != nil {
		return nil, err
	}

	doc, _ := goquery.NewDocumentFromResponse(res)
	jobList := doc.Find(".jobs").Find("li > a")

	for i := range jobList.Nodes {
		company := jobList.Eq(i).Find(".company").Text()
		title := jobList.Eq(i).Find(".title").Text()
		link, _ := jobList.Eq(i).Attr("href")
		job := Job{
			company: strings.TrimSpace(company),
			title:   strings.TrimSpace(title),
			link:    strings.TrimSpace(link),
		}

		if company != "" && title != "" {
			jobs = append(jobs, job)
		}
	}

	return jobs, err
}

func parseRemoteco() ([]Job, error) {
	site := "https://remote.co/remote-jobs/developer"
	var jobs []Job

	res, err := http.Get(site)
	if err != nil {
		return nil, err
	}

	doc, _ := goquery.NewDocumentFromResponse(res)
	jobList := doc.Find(".job_listings").Find(".job_listing")

	for i := range jobList.Nodes {
		company := jobList.Eq(i).Find(".company > strong").Text()
		title := jobList.Eq(i).Find(".position > h3").Text()
		link, _ := jobList.Eq(i).Find("a").Attr("href")
		job := Job{
			company: strings.TrimSpace(company),
			title:   strings.TrimSpace(title),
			link:    strings.TrimSpace(link),
		}

		if company != "" && title != "" {
			jobs = append(jobs, job)
		}
	}

	return jobs, err
}

func main() {
	parsers := []parser{
		// parser{name: "remoteok"},
		// parser{name: "weworkremotely"},
		parser{name: "remoteco"},
	}

	for _, parser := range parsers {
		jobs, err := parser.Parse()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(jobs)
	}
}
