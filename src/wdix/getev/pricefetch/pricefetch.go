package pricefetch

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type Card struct {
	name  string
	price float64
}

func CardUrl(name string) string {
	s := []string{"http://store.tcgplayer.com/magic/return-to-ravnica/", name}
	return strings.Join(s, "")

}

type myRegexp struct {
	*regexp.Regexp
}

func (r *myRegexp) FindStringSubmatchMap(s string) map[string]string {
	captures := make(map[string]string)

	match := r.FindStringSubmatch(s)
	if match == nil {
		return captures
	}

	for i, name := range r.SubexpNames() {
		// Ignore the whole regexp match and unnamed groups
		if i == 0 || name == "" {
			continue
		}

		captures[name] = match[i]

	}
	return captures
}

func FetchCardPrice(name string) (price string) {
	price = "0.0"
	url := CardUrl(name)
	res, _ := http.Get(url)
	response, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()
	avgRegex := myRegexp{regexp.MustCompile(`<td class=\"avg\">(?P<price>.*?)</td>`)}
	captures := avgRegex.FindStringSubmatchMap(string(response))
	elem, ok := captures["price"]
	if ok {
		price = elem
	}
	return
}

func LookupCard(returnChannel chan Card, name string) {
	price := FetchCardPrice(name)
	floatPrice := parsePriceString(price)
	fmt.Println("completed: ", name)
	returnChannel <- Card{name, floatPrice}
}

func parsePriceString(price string) (cost float64) {
	replacer := strings.NewReplacer("$", "")
	stripped := replacer.Replace(price)

	cost, _ = strconv.ParseFloat(stripped, 64)

	return
}
