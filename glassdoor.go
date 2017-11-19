package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Client struct {
	id  string
	key string
}

func (c *Client) SearchEmployer(employer string) interface{} {
	var parsed map[string]interface{}

	url := c.buildUrl("employers", employer)
	resp, _ := http.Get(url)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(body)
	json.Unmarshal(body, &parsed)

	employers := parsed["response"].(map[string]interface{})["employers"].([]interface{})
	var result map[string]interface{}
	for _, employer := range employers {
		if employer.(map[string]interface{})["exactMatch"].(bool) {
			result = employer.(map[string]interface{})
			break
		}
	}

	return result
}

func (c *Client) buildUrl(action, q string) string {
	host := "http://api.glassdoor.com/api/api.htm"
	return fmt.Sprintf("%s?=1&format=json&t.p=%s&t.k=%s&action=%s&q=%s&userip=%s&useragent=%s", host, c.id, c.key, action, q, "192.168.1.1", "Mozilla")
}

func main() {
	gd := Client{id: os.Getenv("GLASSDOOR_ID"), key: os.Getenv("GLASSDOOR_KEY")}
	res := gd.SearchEmployer("symantec")
	fmt.Println(res)
}
