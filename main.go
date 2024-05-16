package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Reign struct {
	Date     string `json:"date"`
	Event    string `json:"event"`
	Location string `json:"location"`
	Days     int    `json:"days"`
}

func scrape(url string) map[string][]Reign {
	reigns := make(map[string][]Reign)

	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	tables := doc.Find("table.mw-datatable")
	if tables.Length() != 1 {
		log.Fatalf("expected 1 table, found %d", tables.Length())
	}
	t := tables.First()

	rows := t.Find("tr:not([style])")
	if rows.Length() == 0 {
		log.Fatalf("no rows found")
	}

	rows.Each(func(i int, r *goquery.Selection) {
		// Skip the header rows
		if i <= 1 {
			return
		}

		var (
			name     string
			date     string
			event    string
			location string
			days     int
		)

		r.Find("td").Each(func(j int, c *goquery.Selection) {
			switch j {
			case 0:
				name = cleanString(c.Text())
			case 1:
				date = cleanString(c.Text())
			case 2:
				event = cleanString(c.Text())
			case 3:
				location = cleanString(c.Text())
			case 6:
				days = cleanDays(c.Text())
			}
		})

		reign := Reign{
			Date:     date,
			Event:    event,
			Location: location,
			Days:     days,
		}

		reigns[name] = append(reigns[name], reign)
	})

	return reigns
}

func cleanString(s string) string {
	if idx := strings.Index(s, "["); idx != -1 {
		s = s[:idx]
	}

	return strings.TrimRight(s, "\n")
}

func cleanDays(s string) int {
	if strings.HasPrefix(s, "<1") {
		return 0
	}

	mapping := func(r rune) rune {
		switch r {
		case ',', '<', '\n':
			return ' '
		default:
			return r
		}
	}

	s = strings.Map(mapping, s)
	s = strings.ReplaceAll(s, " ", "")

	days, err := strconv.Atoi(s)
	if err != nil {
		log.Fatal("error:", err)
	}

	return days
}

func createJson(reigns map[string][]Reign, fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	b, err := json.Marshal(reigns)
	if err != nil {
		return err
	}

	_, err = file.Write(b)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	url := "https://en.wikipedia.org/wiki/List_of_WWE_Women's_Champions_(1956-2010)"
	reigns := scrape(url)

	if err := createJson(reigns, "data/reigns.json"); err != nil {
		fmt.Fprintf(os.Stderr, "error creating file %q: %v\n", "reigns.json", err)
		os.Exit(1)
	}
}
