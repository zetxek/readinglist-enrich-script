package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

	f, err := os.Open("my-book-list.md")
	check(err)
	defer f.Close()
	scanner := bufio.NewScanner(f)

	bookRegex := regexp.MustCompile(`^-\s`)
	catRegex := regexp.MustCompile(`^## `)
	authorRegex := regexp.MustCompile(`by\s`)

	category := ""
	title := ""
	author := ""
	goodreads := ""

	for scanner.Scan() {

		line := scanner.Text()
		isCat := catRegex.MatchString(line)

		if isCat {
			category = strings.Replace(line, "## ", "", 1)
		}

		isBook := bookRegex.MatchString(line)
		if isBook {
			title = strings.Replace(line, "- ", "", 1)
			parts := strings.Split(title, ", by ")
			author = parts[1]
			title = parts[0]
			title = authorRegex.ReplaceAllString(title, "")

			goodreads = fmt.Sprintf("https://www.goodreads.com/search?q=%s+%s", strings.Replace(title, " ", "+", -1), strings.Replace(author, " ", "+", -1))
		}

		if !isBook && !isCat {
			continue
		}

		if isBook {
			generateMd(title, author, category, goodreads)
		}

	}
}

func generateMd(title string, author string, category string, goodreads string) {

	filename := strings.ToLower(title)
	filename = strings.Replace(filename, " ", "-", -1)
	filename = strings.Replace(filename, ":", "", -1)
	filename = strings.Replace(filename, ",", "", -1)
	f, err := os.Create(fmt.Sprintf("./files/%s.md", filename))

	check(err)

	f.WriteString(fmt.Sprintf("---\ntitle: \"%s\"\nbook_author: [\"%s\"]\nbook_category: [\"%s\"]\nlink: \"%s\"\n---\n", title, author, category, goodreads))
	f.Sync()

	fmt.Printf("Created file: %s.md\n", filename)
}
