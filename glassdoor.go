package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Client struct {
	id  string
	key string
}
type Employer struct {
	Id                            int    `json:"id"`
	Name                          string `json:"name"`
	Website                       string `json:"website"`
	ExactMatch                    bool   `json:"exactMatch"`
	Industry                      string `json:"industry"`
	NumberOfRatings               int    `json:"NumberOfRatings"`
	OverallRating                 string `json:"overallRating"`
	RatingDescription             string `json:"ratingDescription"`
	CultureAndValuesRating        string `json:"cultureAndValuesRating"`
	SeniorLeadershipRating        string `json:"seniorLeadershipRating"`
	CompensationAndBenefitsRating string `json:"compensationAndBenefitsRating"`
	CareerOpportunitiesRating     string `json:"careerOpportunitiesRating"`
	WorkLifeBalanceRating         string `json:"workLifeBalanceRating"`
	RecommentToFriendRating       int    `json:"recommendToFriendRating"`
	SectorName                    string `json:"sectorName"`
	IndustryName                  string `json:"industryName"`
}
type gdResponse struct {
	Employers []Employer `json:"employers"`
}
type gdData struct {
	Success  string     `json:"success"`
	Status   string     `json:"status"`
	Response gdResponse `json:"response"`
}

func (c *Client) SearchEmployer(employer string) (Employer, error) {
	var parsed gdData
	var result Employer

	url := c.buildUrl("employers", employer)
	resp, err := http.Get(url)
	if err != nil {
		return result, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}
	// fmt.Println(string(body))
	json.Unmarshal(body, &parsed)

	if parsed.Status != "OK" {
		return result, errors.New(parsed.Status)
	}

	employers := parsed.Response.Employers
	found := false
	for _, employer := range employers {
		if employer.ExactMatch {
			result = employer
			found = true
			break
		}
	}

	if found {
		return result, nil
	}
	return result, errors.New("no exact match found")
}

func (c *Client) buildUrl(action, q string) string {
	host := "http://api.glassdoor.com/api/api.htm"
	return fmt.Sprintf("%s?=1&format=json&t.p=%s&t.k=%s&action=%s&q=%s&userip=%s&useragent=%s", host, c.id, c.key, action, q, "192.168.1.1", "Mozilla")
}

func main() {
	gd := Client{id: os.Getenv("GLASSDOOR_ID"), key: os.Getenv("GLASSDOOR_KEY")}
	empl, err := gd.SearchEmployer("Google")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(empl)
	}
}
