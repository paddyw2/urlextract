package urlextract

import (
	"regexp"
)

// IP pattern explanation
// 1. allow 250-255 or 200-249 or 0-9 with optionally 0 or 1 at the front or 0-9 at the back
// 2. repeat same logic as 1, but with period at start
// 3. expect thep pattern from 2 to be repeated three times
const IpPattern string = `(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`

// hostname pattern explanation:
// 1. in a web page, we only care about matches starting with one of "'/ (i.e http://, or "www...)
// 2. next comes the sub domains pattern, which allows 1 or more valid subdomains (i.e. www.my-site, or just my-site)
// 3. all hostnames must end with a period followed by valid tld characters (i.e. www.my-site.com)
const HostnamePattern string = `["'/]([a-zA-Z0-9]+(-[a-zA-Z0-9]+)*\.)+[a-zA-Z]+`

// urlPattern explanation:
// 1. url must start with hostname
// 2. a url can then have optional paths or options, so allow a /followed by any combination of
// legal url characters
// 3. the urls we care about must target some file, so this pattern of legal url characters must
// end with . followed by lower case alphabet characters to mark a file extension (i.e. .js)
// 4. lastly, the url may have 0 or 1 options after it (i.e. ?v=3411234)
const UrlPattern string = HostnamePattern + `((/[a-zA-Z0-9-_&=\.%?/]*)*(\.[a-z]+)(\?[a-zA-Z0-9-_&=\.%?]*){0,1}){0,1}`

func NewExtractor(options ...interface{}) *Extractor {
	validateTlds := true
	// use switch to validate
	for index, value := range options {
		switch index {
		case 0: // validateTlds
			if validateTldsParam, ok := value.(bool); ok {
				validateTlds = validateTldsParam
			} else {
				panic("First parameter must be of type boolean")
			}
		}
	}
	return &Extractor{ValidateTlds: validateTlds}
}

type Url struct {
	Url      string
	Hostname string
	Tld      string
}

type Extractor struct {
	Ips          []string
	Urls         []*Url
	ValidateTlds bool
}

func (extractor *Extractor) ExtractHostnamesIps(inputString string) {
	extractor.ExtractIps(inputString)
	extractor.ExtractHostnames(inputString)
}

func (extractor *Extractor) ExtractIps(inputString string) {
	ipPattern := regexp.MustCompile(IpPattern)
	submatchall := ipPattern.FindAllString(inputString, -1)
	for _, element := range submatchall {
		extractor.Ips = append(extractor.Ips, element)
	}
}

func (extractor *Extractor) ExtractHostnames(inputString string) {
	hostnameRegex := regexp.MustCompile(HostnamePattern)
	urlRegex := regexp.MustCompile(UrlPattern)
	submatchall := urlRegex.FindAllString(inputString, -1)

	for _, rawUrl := range submatchall {
		rawHostname := hostnameRegex.FindString(rawUrl)
		tldPattern := regexp.MustCompile(`\.([a-zA-Z]+(-[a-zA-Z]+)*)$`)
		tldMatch := tldPattern.FindString(rawHostname)

		// skip if not a valid tld
		if _, ok := ValidTlds[tldMatch[1:]]; !ok && extractor.ValidateTlds {
			continue
		}

		var url, hostname string

		// clean up values before adding to discovered list
		if rawUrl[0] == '"' || rawUrl[0] == '\'' || rawUrl[0] == '/' {
			url = rawUrl[1:]
			hostname = rawHostname[1:]
		} else {
			url = rawUrl
			hostname = rawHostname
		}

		newUrl := Url{Url: url, Hostname: hostname, Tld: tldMatch[1:]}
		extractor.Urls = append(extractor.Urls, &newUrl)
	}
}
