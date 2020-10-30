package leo

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func sanitize(a string) string {
	return strings.Trim(a, "\n")
}

func Query(leoUrl string, langCode string) {
	flag.Parse()

	if flag.NArg() != 1 {
		log.Fatal("Please specify the word you wish to translate")
	}

	doc, err := goquery.NewDocument(leoUrl + url.QueryEscape(flag.Arg(0)))

	if err != nil {
		log.Fatal(err)
	}

	selector := fmt.Sprintf("tr [lang=%s]", langCode)
	sections := []struct{ header, parent, other, german string }{
		{"Substantive", "div[data-dz-name=subst]", selector, "tr [lang=de]"},
		{"Verben", "div[data-dz-name=verb]", selector, "tr [lang=de]"},
		{"Adj./Adv.", "div[data-dz-name=adjadv]", selector, "tr [lang=de]"},
		{"Phrasen", "div[data-dz-name=phrase]", selector, "tr [lang=de]"},
		{"Beispiele", "div[data-dz-name=example]", selector, "tr [lang=de]"},
	}

	for _, section := range sections {
		parent := doc.Find(section.parent)

		if parent.Length() == 0 {
			continue
		}

		fmt.Println("## " + section.header)

		other := parent.Find(section.other)
		german := parent.Find(section.german)

		for i := 0; i < min(other.Length(), german.Length()); i++ {
			otherText := sanitize(other.Eq(i).Text())
			germanText := sanitize(german.Eq(i).Text())
			fmt.Printf("| %30s | %30s |\n", otherText, germanText)
		}
		fmt.Println("\n\n")
	}
}
