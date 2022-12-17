# Soleadify Assignment
## Software Engineer

Attempt at implementing a web scrapper that handles data extraction and indexing.

The assignment description is [here](https://soleadify.notion.site/Assignment-Software-Engineer-0c0cd6c11b1e47ea8ccc677a10160e7b).


<!-- @import "[TOC]" {cmd="toc" depthFrom=1 depthTo=6 orderedList=false} -->

<!-- code_chunk_output -->

- [Soleadify Assignment](#soleadify-assignment)
  - [Software Engineer](#software-engineer)
  - [Technical Decisions](#technical-decisions)
    - [Use the Go programming language](#use-the-gohttpsgodev-programming-language)
  - [Tasks](#tasks)
    - [Parse and validate urls from CSV](#parse-and-validate-urls-from-csv)
    - [Download HTML using HTTP client](#download-html-using-http-client)
    - [Extract information from raw HTML using CSS selectors and textContent](#extract-information-from-raw-html-using-css-selectors-and-textcontent)
    - [Process websites concurrently using goroutines](#process-websites-concurrently-using-goroutines)
    - [Extra goals:](#extra-goals)

<!-- /code_chunk_output -->

## Technical Decisions

Here I'm going to document the technical decisions I made, and
the reasoning behind them.

### Use the [Go](https://go.dev/) programming language

Reasons:
- I want to use this programming language professionally and learn more about it
- It compiles to a native code binary which is easy to deploy
- It makes good use of concurrency, which should help out with
    both IO and CPU parallelism.


## Tasks
### Parse and validate urls from CSV
The tool should parse a CSV file and validate the URLs contained.
It should return both the valid URLs, as well as a list of invalid URLs,
so they can be logged.

###  Download HTML using HTTP client
###  Extract information from raw HTML using CSS selectors and textContent
### Process websites concurrently using goroutines

### Extra goals:
- Use selenium and chromedriver to run JS on websites, so we can
      extract JS rendered content


## Examples:
```
$ ./scrappy scrape phone https://mendiolagardening.com
1 invalid phone number(s)
Invalid phone number: "666 888 0000" ("invalid phone number")

Domain: "https://mendiolagardening.com"
 0. "+1 510-575-7324" (a[href="tel:< phone number >"])
 1. "+1 888-999-0000" (a[href="tel:< phone number >"])

$ ./scrappy scrape phone https://mendiolagardening.com --raw
Domain: "https://mendiolagardening.com"
 0. "(510) 575-7324" (regex match)
 1. "666 888 0000" (regex match)
 2. "(510) 575 7324" (regex match)
 3. "(510) 575 7324" (regex match)
 4. "+15105757324" (a[href="tel:< phone number >"])
 5. "5105757324" (a[href="tel:< phone number >"])
 6. "5105757324" (a[href="tel:< phone number >"])
 7. "5105757324" (a[href="tel:< phone number >"])
 8. "888-999-0000" (a[href="tel:< phone number >"])
 9. "+15105757324" (a[href="tel:< phone number >"])
10. "+15105757324" (a[href="tel:< phone number >"])
```


Look into:
```
$ ./scrappy scrape phone https://verdantporch.com
Error: Get "https://www.verdantporch.com/": Not following redirect to www.verdantporch.com because its not in AllowedDomains
```
