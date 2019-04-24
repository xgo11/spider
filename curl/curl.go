package curl

import (
	"fmt"
	"strings"
)

var (
	whiteChars = map[byte]int{
		' ':  1,
		'\t': 0,
		'\r': 0,
		'\n': 0,
	}

	quoteChars = map[byte]int{
		'\'': 0,
		'"':  0,
	}
)

func analyzeParts(curlCommandString string) []string {
	var parts []string

	var start, end int
	var isLastWhite = false
	var inQuote = false
	var lastQuote = byte(0)

	for i, c := range []byte(curlCommandString) {
		cFlag, ok := whiteChars[c]
		if !ok {
			if isLastWhite {
				start = i
			}
			isLastWhite = false

			if _, ok = quoteChars[c]; ok { // 是引号分隔符
				if !inQuote {
					inQuote = true
					lastQuote = c
					start = i + 1
				} else if c == lastQuote {
					end = i
					inQuote = false

					if end > start {
						parts = append(parts, curlCommandString[start:end])
					}
					start = end + 1
					end = start
				}
			}
			continue
		}
		if cFlag == 1 && inQuote {
			continue
		}

		if !isLastWhite {
			end = i
			if end > start {
				parts = append(parts, curlCommandString[start:end])
			}
			isLastWhite = true
		}
	}
	return parts
}

// extract curl command , return url and a map data of headers/data/method/use_gzip
func AnalyzeCurl(curlCommandString string) (string, map[string]interface{}, error) {
	var s = strings.ReplaceAll(curlCommandString, "\n", "")

	headers := make(map[string]string)
	kwArgs := make(map[string]interface{})

	var urls []string
	var command, currentOpt string

	for _, part := range analyzeParts(s) {
		if len(command) < 1 {
			command = part
		} else if !strings.HasPrefix(part, "-") && len(currentOpt) < 1 {
			urls = append(urls, part)
		} else if len(currentOpt) < 1 && strings.HasPrefix(part, "-") {
			if part == "--compressed" {
				kwArgs["use_gzip"] = true
			} else {
				currentOpt = part
			}
		} else {
			if len(currentOpt) < 1 {
				return "", nil, fmt.Errorf("unknow curl argument: %v", part)
			} else if currentOpt == "-H" || currentOpt == "--header" {
				var k, v string
				if i := strings.Index(part, ":"); i > 0 {
					k = part[0:i]
					v = part[i+1:]
					headers[strings.Trim(k, " ")] = strings.Trim(v, " ")
				}
			} else if currentOpt == "-d" || currentOpt == "--data" {
				kwArgs["data"] = part
			} else if currentOpt == "--data-binary" {
				if part[0:1] == "$" {
					part = part[1:]
				}
				kwArgs["data"] = part
			} else if currentOpt == "-X" || currentOpt == "--request" {
				kwArgs["method"] = strings.ToUpper(part)
			} else {
				return "", nil, fmt.Errorf("unknow curl option: %v", currentOpt)
			}
			currentOpt = ""
		}

	}

	if nil == urls || len(urls) < 1 {
		return "", nil, fmt.Errorf("curl, no url")
	}
	if len(headers) > 0 {
		kwArgs["headers"] = headers
	}

	return urls[0], kwArgs, nil
}
