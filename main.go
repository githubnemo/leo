package main

import (
	"flag"
	"log"
	"fmt"
	"strings"
	"github.com/PuerkitoBio/goquery"
	"net/url"
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

func main() {
	flag.Parse()

	if flag.NArg() != 1 {
		log.Fatal("Word missing.")
	}

	doc, err := goquery.NewDocument("http://dict.leo.org/englisch-deutsch/" + url.QueryEscape(flag.Arg(0)))

	if err != nil {
		log.Fatal(err)
	}

	sections := []struct{ header, parent, en, de string }{
		{"Substantive", "div[data-dz-name=subst]", "tr [lang=en]", "tr [lang=de]"},
		{"Verben", "div[data-dz-name=verb]", "tr [lang=en]", "tr [lang=de]"},
		{"Adj./Adv.", "div[data-dz-name=adjadv]", "tr [lang=en]", "tr [lang=de]"},
		{"Phrasen", "div[data-dz-name=phrase]", "tr [lang=en]", "tr [lang=de]"},
		{"Beispiele", "div[data-dz-name=example]", "tr [lang=en]", "tr [lang=de]"},
	}

	for _, section := range sections {
		parent := doc.Find(section.parent)

		if parent.Length() == 0 {
			continue
		}

		fmt.Println("## " + section.header)

		en := parent.Find(section.en)
		de := parent.Find(section.de)

		for i := 0; i < min(en.Length(),de.Length()); i++ {
			enText := sanitize(en.Eq(i).Text())
			deText := sanitize(de.Eq(i).Text())
			fmt.Printf("| %30s | %30s |\n", enText, deText)
		}
		fmt.Println("\n\n")
	}
}
