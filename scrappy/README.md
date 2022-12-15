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
