package parser

import (
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Job spec
type Job struct {
	Company string
	Title   string
	Link    string
}

// Parser is the component representing the site to parse
type Parser struct {
	Name string
}

// New creates a new parser by name
func New(name string) Parser {
	return Parser{Name: name}
}

// Parse starts parsing the site
func (p Parser) Parse() ([]Job, error) {
	var jobs []Job
	var err error
	switch p.Name {
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
	site := "https://remoteok.io"
	jobsUrl := site + "/remote-dev-jobs"
	var jobs []Job

	res, err := http.Get(jobsUrl)
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
			Company: strings.TrimSpace(company),
			Title:   strings.TrimSpace(title),
			Link:    site + strings.TrimSpace(link),
		}

		if company != "" && title != "" {
			jobs = append(jobs, job)
		}
	}

	return jobs, err
}

func parseWeworkremotely() ([]Job, error) {
	site := "https://weworkremotely.com"
	jobsUrl := site + "/categories/2-programming/jobs#intro"
	var jobs []Job

	res, err := http.Get(jobsUrl)
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
			Company: strings.TrimSpace(company),
			Title:   strings.TrimSpace(title),
			Link:    site + strings.TrimSpace(link),
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
			Company: strings.TrimSpace(company),
			Title:   strings.TrimSpace(title),
			Link:    strings.TrimSpace(link),
		}

		if company != "" && title != "" {
			jobs = append(jobs, job)
		}
	}

	return jobs, err
}
