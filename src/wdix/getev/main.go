package main

import (
	"fmt"
	"github.com/moovweb/gokogiri/html"
	"io/ioutil"
	"net/http"
	"wdix/getev/pricefetch"
)

const ravURL string = "http://gatherer.wizards.com/Pages/Search/Default.aspx?output=checklist&action=advanced&set=%5b%22Return+to+Ravnica%22%5d"

func waitForCards(responseChannel chan pricefetch.Card, numberOfCards int) (cards []pricefetch.Card) {
	returnedCount := 0
	for {
		cards = append(cards, <-responseChannel)
		returnedCount++

		if returnedCount >= numberOfCards {
			break
		}
	}
	return
}

func fetchCardNames(names *[]string) {
	res, err := http.Get(ravURL)
	if err != nil {
		fmt.Println(err)
	}
	response, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	res.Body.Close()

	doc, err := html.Parse(response, html.DefaultEncodingBytes, nil, html.DefaultParseOption, html.DefaultEncodingBytes)

	if err != nil {
		fmt.Println(err)
	}

	html := doc.Root().FirstChild()
	defer doc.Free()

	results, err := html.Search("//tr[@class='cardItem']")

	for _, row := range results {

		name, err := row.Search("./td[@class='name']")

		if err != nil {
			fmt.Println(err)
			continue
		}

		stringName := name[0].Content()
		*names = append(*names, stringName)
	}

	if err != nil {
		fmt.Println(err)
	}

	return
}

func main() {
	var names = make([]string, 0, 1)
	fetchCardNames(&names)
	fmt.Println(names)

	cardChannel := make(chan pricefetch.Card)

	for _, cardName := range names {
		go pricefetch.LookupCard(cardChannel, cardName)
	}
	cards := waitForCards(cardChannel, len(names))
	fmt.Println(cards)
}
