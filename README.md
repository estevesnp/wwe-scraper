# WWE Champ Scraper

This is a simple webscraper that scrapes the WWE Wikipedia page and gets the information of all the champions from a specific belt.

## Installation

You need to have Go installed. Afterwards, run the following command:

`go install github.com/estevesnp/scraperwwe@latest`

## Usage

To use the program, run the following command:

`scraperwwe "<url>" [output.json]`

For example:

`scraperwwe "https://en.wikipedia.org/wiki/List_of_WWE_Women's_Champions" womenchamps.json`

Make sure to have your Go bin directory in your PATH.
