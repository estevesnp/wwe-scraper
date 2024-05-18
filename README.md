# WWE Champ Scraper

This is a simple webscraper that scrapes the WWE Wikipedia page and gets the information of all the champions from a specific belt.

It uses the [goquery](https://github.com/PuerkitoBio/goquery) pkg to parse the HTML.

## Installation

You need to have Go installed. Afterwards, run the following command:

`go install github.com/estevesnp/wwe-scraper@latest`

## Usage

To use the program, run the following command:

`wwe-scraper "<url>" [output.json]`

For example:

`wwe-scraper "https://en.wikipedia.org/wiki/List_of_WWE_Women's_Champions" womenchamps.json`

Make sure to have your Go bin directory in your PATH.
