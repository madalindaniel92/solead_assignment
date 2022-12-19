# Soleadify Assignment
## Software Engineer

Attempt at implementing a web scrapper that handles data extraction and indexing.

The assignment description is [here](https://soleadify.notion.site/Assignment-Software-Engineer-0c0cd6c11b1e47ea8ccc677a10160e7b).

## Technical Decisions

Here I'm going to document the technical decisions I made, and
the reasoning behind them.

### Use the [Go](https://go.dev/) programming language

Reasons:
- I want to use this programming language professionally and learn more about it
- It compiles to a native code binary which is easy to deploy
- It makes good use of concurrency, which should help out with IO and CPU parallelism.

I've chosen to implement a CLI helper tool, since it allows me to run various  
tasks from the command line and divide the problem into smaller pieces.

### Approach
- Check website reachability using http HEAD requests

- Extract links for each website using HTML anchor tags,
  as well as parsing the website's sitemap.

- Extract phone numbers using 2 different strategies:
  - finding anchor tags with a `href` of `tel:`
  - parsing the textContent of the page and matching it
    against a [regular expression for US numbers](/scrappy/internal/phone/phone.go#L55)  
    (only **US** numbers are supported for now, since using a single phone number format is easier;   
    to add more formats multiple regex matches for various formats are needed,  
    or one really complicated and imprecise regex to rule them all...)

- Validate and normalize the extracted phone numbers using a [Golang port of libphonenumbers](https://github.com/nyaruka/phonenumbers#phonenumbers).  
  this allows us to check if a number is valid within a specific number plan.

- Setup ElasticSearch and Kibana for local development using a [docker-compose](/elastic_support/docker-compose.yml) file.    
  The credentials are loaded into the CLI using a [config file](/scrappy/.scrappy.yaml)   
  The cluster uses a self-signed CA certificate copied into the repo [here](elasticsearch_ca.crt).   
  The `setup` part of the docker-compose.yml file automatically generates this file.

- Import company info into elasticsearch using a the official [go-elasticsearch client](https://github.com/elastic/go-elasticsearch)

## Todo (WIP)

- Implement ES search using `multi_match` query. 

- Scrape websites concurrently (only HEAD requests are currently handled concurrently)

- Write Dockerfile that builds the CLI tool and shoves it into an alpine Linux image.

- Write small http server returning JSON serialized data

## Local setup

Currently golang is required in order to build the tool.   
My goal is to use a Dockerfile to create a small alpine Linux container with the binary,  
then push said binary to docker hub, allowing it to be used without Go installed.


```sh
# Clone repo
git clone git@github.com:madalindaniel92/solead_assignment.git

# Navigate to golang project folder
cd solead_assignment/scrappy

# Build binary
go build

# See help information for tool
./scrappy help

# Check which domains are reachable
./scrappy check domains testdata/sample-websites.csv
```

There's also setup involved to get Elasticsearch and Kibana up and running using  
the provided docker-compose.yml.

```sh
# Navigate to ES support folder
cd solead_assignment/elastic_support

# Pull images and start ES cluster
docker-compose up -d
```

### Grabbing certificates from Docker:

The docker-compose ElasticSearch setup uses TLS, and we need  
to copy (or otherwise make accessible) the certificate used by the  
local development cluster.

The option I currently settled for is copying the certificate from

```sh
# Need either sudo, or current user as part of the docker group
sudo su

# Read Mountpoint value to figure out where docker mounts said volume
docker volume inspect elastic_support_certs | grep Mountpoint

# Copy certificate from volume into current work directory
cp <volume mountpoint path>/elastic_support_certs/_data/ca/ca.crt elasticsearch_ca.crt

# Make sure current user has file permissions.
chown $USER:$USER elasticsearch_ca.crt
```


## Tasks

This section documents each subproblem identified, the CLI subcommand used  
to tackle said problem, and example output when running the command.


### Parse and validate website domains from a CSV file.
The tool should parse a CSV file with website domains and validate the URLs contained.  
It should return both the valid URLs, as well as a list of invalid URLs, so they can be logged.

The CLI subcommand for parsing the domains CSV and checking the domains  
using http HEAD requests is:
```sh
$ ./scrappy check domains <csv file to load domain names from>
```

Example:
```sh
$ ./scrappy check domains testdata/small-sample.csv
```

Example output of running this command on the [small-sample.csv](./scrappy/testdata/small-sample.csv):
```
2022/12/19 08:54:14 Failed request to domain "https://coffee-homemachines.club": "Head \"https://coffee-homemachines.club\": dial tcp: lookup coffee-homemachines.club: no such host"
2022/12/19 08:54:14 Failed request to domain "https://maddux.pro": "Head \"https://maddux.pro\": x509: certificate is valid for *.secureserversites.net, secureserversites.net, not maddux.pro"
2022/12/19 08:54:14 HEAD "https://takapartners.com" - 200
2022/12/19 08:54:14 HEAD "https://thestonenc.com" - 200
2022/12/19 08:54:14 HEAD "https://ohanaconsulting.net" - 200
2022/12/19 08:54:14 HEAD "https://mendiolagardening.com" - 200
2022/12/19 08:54:14 HEAD "https://techbarstore.com" - 200
2022/12/19 08:54:14 Failed request to domain "https://tlalocrivas.com": "Head \"https://tlalocrivas.com\": x509: certificate is valid for *.weebly.com, *.weeblysite.com, weebly.com, weeblysite.com, not tlalocrivas.com"
2022/12/19 08:54:14 HEAD "https://postmodern-strings.com" - 200
2022/12/19 08:54:14 HEAD "https://cumberland-river.com" - 200
2022/12/19 08:54:14 HEAD "https://kkcger.com" - 200
2022/12/19 08:54:15 HEAD "https://workitstudio.com" - 200
2022/12/19 08:54:15 HEAD "https://mazautoglass.com" - 200
2022/12/19 08:54:15 HEAD "https://verdantporch.com" - 200
2022/12/19 08:54:15 HEAD "https://timent.com" - 200
2022/12/19 08:54:15 HEAD "https://melatee.com" - 200
2022/12/19 08:54:16 HEAD "https://bostonzen.org" - 200
2022/12/19 08:54:18 HEAD "https://creativebusinessassociates.com" - 200
2022/12/19 08:54:22 HEAD "https://kansaslimousin.org" - 200

Successful requests: 16
Failed to connect: 3
Bad requests: 0
```

From the output we can see that from these 19 domains, 16 returned success (200 OK),  
while 3 failed to connect.

The implementation for this command is in [cmd/domains.go](/scrappy/cmd/domains.go), in the [domainAction](/scrappy/cmd/domains.go#L56) function.  
The logic is split up into two parts:
- [LoadDomainsFromFile](/scrappy/internal/csv/csv.go#L39) which calls     [ParseDomainsCSV](scrappy/internal/csv/csv.go#L69) to import the domain urls from the CSV file.  
  Tests for this funcion are [csv_test.go](/scrappy/internal/csv/csv_test.go)

- [CheckURLs](scrappy/internal/web/web.go#L58) which uses worker goroutines to
  issue HEAD requests to each domain and check that it is reachable.


### Parse and validate company info from a CSV file
The tool should parse a CSV file with company information and display it.

The CLI subcommand for parsing the domains CSV printing company information is:
```sh
./scrappy check companies <csv file to load company info from>
```

Example:
```sh
./scrappy check companies testdata/small-company-names.csv
```

Example output of running this command on the [small-company-names.csv](/scrappy/testdata/small-company-names.csv):
```
Domain: https://bostonzen.org
Commercial name: Greater Boston Zen Center
Legal name: GREATER BOSTON ZEN CENTER INC.
Other names:
    - Greater Boston Zen Center
    - Boston Zen
    - GREATER BOSTON ZEN CENTER INC.

Domain: https://mazautoglass.com
Commercial name: MAZ Auto Glass
Other names:
    - MAZ Auto Glass

Domain: https://melatee.com
Commercial name: Melatee
Other names:
    - Melatee

Domain: https://timent.com
Commercial name: Timent Technologies
Other names:
    - Timent Technologies
    - Timent
```

The implementation for this subcommand is in the [cmd/companies.go](/scrappy/cmd/companies.go#L41) file.


### Extract links from domain using sitemap
The tool should extract the links associated with a domain that it could crawl,
using the domain's sitemap.

The CLI subcommand for getting links from a domains sitemap is:
```sh
./scrappy links sitemap <domain with scheme information>
```

Example:
```sh
./scrappy links sitemap https://cumberland-river.com
```

Output:
```
2022/12/19 09:14:48 Visiting "https://cumberland-river.com/wp-sitemap.xml"
2022/12/19 09:14:48 Visiting "https://cumberland-river.com/wp-sitemap-posts-post-1.xml"
2022/12/19 09:14:49 Visiting "https://cumberland-river.com/wp-sitemap-posts-page-1.xml"
2022/12/19 09:14:49 Visiting "https://cumberland-river.com/wp-sitemap-taxonomies-category-1.xml"
2022/12/19 09:14:49 Visiting "https://cumberland-river.com/wp-sitemap-taxonomies-post_tag-1.xml"
2022/12/19 09:14:50 Visiting "https://cumberland-river.com/wp-sitemap-users-1.xml"
        0 "https://cumberland-river.com/cumberland-river-nashville/"
        1 "https://cumberland-river.com/cumberland-river-cruises/"
        2 "https://cumberland-river.com/cumberland-river-fishing/"
        3 "https://cumberland-river.com/cumberland-river-has-more-than-fishing/"
        4 "https://cumberland-river.com/7-fun-things-to-do-in-kentucky/"
        5 "https://cumberland-river.com/best-hotels-in-kentucky/"
        6 "https://cumberland-river.com/best-restaurants-in-louisville-kentucky/"
        7 "https://cumberland-river.com/"
        8 "https://cumberland-river.com/about/"
        9 "https://cumberland-river.com/contact/"
        10 "https://cumberland-river.com/links/"
        11 "https://cumberland-river.com/category/uncategorized/"
        12 "https://cumberland-river.com/category/nashville/"
... output trimmed>
```

The implementation of this command is in [cmd/sitemap.go](/scrappy/cmd/sitemap.go#L39).  
It uses the [GetSitemapLinks]() function to fetch `robots.txt` from the domain,  
extracts the `Sitemap:` entry from `robots.txt` if present, then extracts  
the links from the sitemap using a [scraping framework for Go](https://github.com/gocolly/colly).
The heart of the implementation is in the [CollectSitemapLinks](/scrappy/internal/web/sitemap.go#50) function.


### Extract links from domain HTML anchor tags
The tool should extract all URL links from a website's anchor tags.

The CLI subcommand getting the links is:
```sh
./scrappy links nav <domain with scheme information>
```

Example:
```sh
./scrappy links nav https://cumberland-river.com
```

Output:
```
2022/12/19 09:24:56 Visiting "https://cumberland-river.com"
        0 "https://cumberland-river.com/"
        1 "https://cumberland-river.com/about/"
        2 "https://cumberland-river.com/contact/"
        3 "https://cumberland-river.com/links/"
        4 "https://cumberland-river.com/category/recommendations/"
        5 "https://cumberland-river.com/category/activities/"
        6 "https://cumberland-river.com/category/hotels/"
        7 "https://cumberland-river.com/category/dining/"
```

The implementation for this subcommand is in the [cmd/nav.go.](/scrappy/cmd/nav.go#L39) file.

It relies on the [GetLinks](/scrappy/internal/web/nav_links.go#L11) function.


### Scrape phone numbers from an URL
The tool should extract phone numbers from the URL of an HTML document.

The CLI subcommand for extracting phone number information is:
```sh
./scrappy scrape phone <domain with scheme information>
```

Example:
```sh
./scrappy scrape phone https://mazautoglass.com/
```

Output:
```
Domain: "https://mazautoglass.com/"
 0. "+1 415-626-4474" (a[href="tel:< phone number >"])
```

We can see that the number has been extracted from an anchor tag,
whose `href` attribute is of type `tel:`.

Another example:
```sh
./scrappy scrape phone https://verdantporch.com
```

Output:
```
1 invalid phone number(s)
Invalid phone number: "1635271811" ("invalid phone number")

Domain: "https://verdantporch.com"
 0. "+1 910-639-7205" (regex match)
```

First, we see that we found a value that is an invalid phone number.   
Next, we can see that a number has been extracted using a regex match,

Yet another example:
```sh
$ ./scrappy scrape phone https://mendiolagardening.com
1 invalid phone number(s)
Invalid phone number: "666 888 0000" ("invalid phone number")

Domain: "https://mendiolagardening.com"
 0. "+1 510-575-7324" (a[href="tel:< phone number >"])
 1. "+1 888-999-0000" (a[href="tel:< phone number >"])

Using the `--raw` flag we get the results without validation and deduplication:
```sh
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

The implementation for scraping getting phone numbers is in:
- [internal/web/phone.go](/scrappy/internal/web/phone.go)
  `GetPhoneNums` implements the webscraping part, using 2 strategies,
  looking for anchor tags that have a `tel:` href attribute, and also
  matching the text agains a regular expression for US phone numbers.

- [internal/phone/phone.go](/scrappy/internal/phone/phone.go) - implements text scraping and
  number validation using the [usPhoneNumberRegex](/scrappy/internal/phone/phone.go#L55) to find matches and [ValidatePhoneNumber](/scrappy/internal/phone/phone.go#L121) to validate phone numbers


The command only extracts US phone numbers currently since it's easier  
to focus on a single phone number format.

In order to validate the phone numbers, a [Golang port of libphonenumbers](https://github.com/nyaruka/phonenumbers#phonenumbers) is used.

## Bits and pieces to sort out

### Extra goals:
- Use selenium and chromedriver to run JS on websites, so we can
      extract JS rendered content
