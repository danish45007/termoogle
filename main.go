package main

import (
	"fmt"
	"math/rand"
)

var googleDomains = map[string]string{}

var userAgents = []string{}

type SerachResults struct {
	ResultRank  int
	ResultUrl   string
	ResultTitle string
	ResultDesc  string
}

func randomUserAgent() string {
	randNum := rand.Int() % len(userAgents)
	return userAgents[randNum]
}

func searchOnGoogle(string) ([]string, error) {
	return []string{
		"xyz",
	}, nil
}

func main() {
	res, err := searchOnGoogle("xyz")
	if err != nil {
		panic("unable to fetch results.")
	}
	for _, r := range res {
		fmt.Println(r)
	}
}
