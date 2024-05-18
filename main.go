package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Reign struct {
	Name     string `json:"name"`
	Date     string `json:"date"`
	Event    string `json:"event"`
	Location string `json:"location"`
	Days     int    `json:"days"`
}

func scrape(url string) ([]Reign, error) {
	reigns := []Reign{}

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	tables := doc.Find("table.mw-datatable")
	if tables.Length() != 1 {
		return nil, fmt.Errorf("expected 1 table, found %d", tables.Length())
	}
	t := tables.First()

	rows := t.Find("tr:not([style])")
	if rows.Length() == 0 {
		return nil, fmt.Errorf("no rows found")
	}

	rows.Each(func(i int, row *goquery.Selection) {
		if err != nil {
			return
		}

		// Skip the header rows
		if i <= 1 {
			return
		}

		var r Reign

		row.Find("td").Each(func(j int, c *goquery.Selection) {
			switch j {
			case 0:
				r.Name = cleanString(c.Text())
			case 1:
				r.Date = cleanString(c.Text())
			case 2:
				r.Event = cleanString(c.Text())
			case 3:
				r.Location = cleanString(c.Text())
			case 6:
				r.Days, err = cleanDays(c.Text())
			}
		})
		if err != nil {
			return
		}

		reigns = append(reigns, r)
	})
	if err != nil {
		return nil, err
	}

	return reigns, nil
}

func cleanString(s string) string {
	if idx := strings.Index(s, "["); idx != -1 {
		s = s[:idx]
	}

	return strings.TrimSpace(s)
}

func cleanDays(s string) (int, error) {
	if strings.HasPrefix(s, "<1") {
		return 0, nil
	}

	mapping := func(r rune) rune {
		switch r {
		case ',', '+', '\n':
			return ' '
		default:
			return r
		}
	}

	s = strings.Map(mapping, s)
	s = strings.ReplaceAll(s, " ", "")

	days, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("failed to convert %q to int: %v", s, err)
	}

	return days, nil
}

func createJson(reigns []Reign, fileName string) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	b, err := json.MarshalIndent(reigns, "", "    ")
	if err != nil {
		return err
	}

	_, err = f.Write(b)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	if len(os.Args) == 1 {
		fmt.Fprintln(os.Stderr, "Usage: scrapewwe \"<wiki_url>\" [file.json]")
		os.Exit(1)
	}

	url := os.Args[1]

	var fileName string
	if len(os.Args) > 2 {
		fileName = os.Args[2]
	} else {
		fileName = "reigns.json"
	}

	reigns, err := scrape(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error scraping %s: %s", url, err)
		os.Exit(1)
	}

	if err := createJson(reigns, fileName); err != nil {
		fmt.Fprintf(os.Stderr, "error creating file %q: %v\n", fileName, err)
		os.Exit(1)
	}
}
