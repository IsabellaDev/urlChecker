package main

import (
	//"flag"
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"mvdan.cc/xurls/v2"
)

// extract all urls from the string and return them as a slice (array)
func extractURL(str string) []string {
	rxStrict := xurls.Strict()
	foundUrls := rxStrict.FindAllString(str, -1)
	return foundUrls
}

//Function to parse ignore URLs from provided ignore file path
func parseIgnoreURL(ignoreFilePath string) []string {
	var ignoreURLs []string
	//Read the content of file given by ignoreFilePath
	content, err := ioutil.ReadFile(ignoreFilePath)
	if err != nil {
		log.Fatal(err)
	}
	//Create a scanner for file content
	scanner := bufio.NewScanner(strings.NewReader(string(content)))
	re := regexp.MustCompile("^(#|https?://)")
	//Scan the ignore URL file line by line
	for scanner.Scan() {
		line := scanner.Text()

		//Check if the ignore link file is invalid
		if !re.Match([]byte(line)) {
			fmt.Println("Ignore Link File is invalid")
			fmt.Println("Exit with status 1")
			os.Exit(1)
		}

		firstChar := string(line[0])

		//Only look at lines that don't start with #
		if firstChar != "#" {
			URLsFoundFromLine := extractURL(line)
			ignoreURLs = append(ignoreURLs, URLsFoundFromLine...)
		}
	}

	//If there is error during scan, log the error
	if scanErr := scanner.Err(); scanErr != nil {
		log.Fatal(scanErr)
	}

	return ignoreURLs
}

// remove duplicate strings from a slice of strings
func removeDuplicate(urls []string) []string {
	result := make([]string, 0, len(urls))
	temp := map[string]struct{}{}
	for _, item := range urls {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func parseUniqueURLsFromFile(filepath string) []string {
	//open file and read it
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatal(err)
	}
	textContent := string(content)

	//call functions to check the availability of each url
	return removeDuplicate(extractURL(textContent))
}

//check if urls passed reachable or not
func checkURL(urls []string) {

	//use multi-threads to make the process execute faster
	var wg sync.WaitGroup
	wg.Add(len(urls))

	//loop through the urls to check one by one
	for _, v := range urls {
		//annonymous function to make the wg.Done() work
		go func(v string) {
			defer wg.Done()

			client := http.Client{
				Timeout: 8 * time.Second,
			}
			//check if the url is reachable or not
			resp, err := client.Head(v)
			//deal with errors
			if err != nil {

				fmt.Println(v + ": NO RESPONCE!")
			} else {

				//allow environment variables to determine the colors of the output
				clicolor := os.Getenv("CLICOLOR")

				if clicolor == "1" {

					//set different colors and reponse according to different status code
					var (
						greenC = "\033[1;32m%s\033[0m"
						redC   = "\033[1;31m%s\033[0m"
						grayC  = "\033[1;30m%s\033[0m"
					)
					switch code := resp.StatusCode; code {
					case 200:
						fmt.Printf(greenC, v+": GOOD!\n")

					case 400, 404:
						fmt.Printf(redC, v+": BAD!\n")

					default:
						fmt.Printf(grayC, v+": UNKNOWN!\n")

					}
				} else {
					switch code := resp.StatusCode; code {
					case 200:
						fmt.Println(v + ": GOOD!")

					case 400, 404:
						fmt.Println(v + ": BAD!")

					default:
						fmt.Println(v + ": UNKNOWN!")
					}

				}
			}

		}(v)
	}

	wg.Wait()
}

//json output structure
type UrlJson struct {
	//[ { "url": 'https://www.google.com', "status": 200 }, { "url": 'https://bad-link.com', "status": 404 } ]
	Url    string
	Status int
}

//if json output required, check urls and output json
func checkURLJson(urls []string) {

	//use multi-threads to make the process execute faster
	var wg sync.WaitGroup
	wg.Add(len(urls))

	urlsJ := make([]UrlJson, 0)

	//loop through the urls to check one by one
	for _, v := range urls {
		go func(v string) {
			defer wg.Done()

			client := http.Client{
				Timeout: 8 * time.Second,
			}
			//check if the url is reachable or not
			resp, err := client.Head(v)
			//deal with errors
			if err != nil {

				j := UrlJson{v, -1}
				urlsJ = append(urlsJ, j)
			} else {
				u := UrlJson{v, resp.StatusCode}
				urlsJ = append(urlsJ, u)

			}
		}(v)
	}
	wg.Wait()

	urlsInJson, err := json.Marshal(urlsJ)
	if err != nil {
		log.Fatalf("Something is going wrong with json Marshalling: %s", err)
	}
	//fmt.Println(urlsInJson)
	os.Stdout.WriteString(string(urlsInJson))

}
