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

	"github.com/mb0/glob"
	flag "github.com/spf13/pflag"
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
	ignoreURLs := []string{}
	//Read the content of file given by ignoreFilePath
	content, readErr := ioutil.ReadFile(ignoreFilePath)
	if readErr != nil {
		log.Fatal(readErr)
	}
	//Create a scanner for file content
	scanner := bufio.NewScanner(strings.NewReader(string(content)))
	//Scan the ignore URL file line by line
	for scanner.Scan() {
		line := scanner.Text()
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

func removeLinkFromList(s []string, i int) []string {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
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

func main() {
	globFlag := flag.Bool("g", false, "Glob pattern")

	//add flags of -j, --jason to enable the program to output JSON
	jflag := flag.BoolP("json", "j", false, "json output")

	// checking version flag
	vflag := flag.BoolP("version", "v", false, "version check")

	// ignore url flag
	ignoreFlag := flag.BoolP("ignore", "i", false, "ignore url patterns")

	flag.Parse()
	//deal with non-file path, giving usage message
	if len(os.Args) == 1 {
		fmt.Println("help/usage message: To run this program, please pass an argument to it,i.e.: go run urlChecker.go urls.txt")

	} else {
		//feature of checking version
		if *vflag {
			fmt.Println("  *****  urlChecker Version 0.2  *****  ")
			return
		}

		if *globFlag {
			//Assign the glob pattern provided to a local variable
			pattern := flag.Args()[0]
			//Read all files in the current directory
			files, _ := ioutil.ReadDir(".")
			//Create a globber object
			globber, _ := glob.New(glob.Default())
			//Loop through all files
			for _, file := range files {
				//Check if the file name match the glob pattern provided
				matched, _ := globber.Match(pattern, file.Name())
				//If matched then run the url check on that file
				if matched {
					//open file and read it
					content, err := ioutil.ReadFile(file.Name())
					if err != nil {
						log.Fatal(err)
					}
					textContent := string(content)

					fmt.Println(">>  ***** UrlChecker is working now...... *****  <<")
					fmt.Println("--------------------------------------------------------------------------------------------------")
					//call functions to check the availability of each url
					checkURL(extractURL(textContent))
				}
			}
			return
		}

		if *ignoreFlag {
			combineStrIgnorePatterns := ""
			ignoreFilePath := flag.Args()[0]

			ignoreList := parseIgnoreURL(ignoreFilePath)

			filepath := flag.Arg(1)

			//If the user did not provide a second arg, exit with status code 1
			if filepath == "" {
				fmt.Println("A filepath as a second arg is required")
				os.Exit(1)
			}

			//Extract the URLs from filepath provided
			URLList := parseUniqueURLsFromFile(filepath)

			//If there are links to ignore, then filter them out from URLs list extracted above, else check regularly
			if len(ignoreList) != 0 {
				URLListWithoutIgnores := []string{}
				for index, pattern := range ignoreList {
					if index != len(ignoreList)-1 {
						combineStrIgnorePatterns += pattern + "|"
					} else {
						combineStrIgnorePatterns += pattern
					}
				}

				//Create regex object to match any ignore links in list
				re := regexp.MustCompile("^(" + combineStrIgnorePatterns + ")")

				//Filter out the ignored links
				for _, link := range URLList {
					if !re.Match([]byte(link)) {
						URLListWithoutIgnores = append(URLListWithoutIgnores, link)
					}
				}

				//Check with URL List that has no ignored links
				checkURL(URLListWithoutIgnores)
			} else {
				//Check regularly if there is nothing to ignore
				checkURL(URLList)
			}
			return
		}

		//use for loop to deal with multiple file paths
		i := 1
		for i+1 <= len(os.Args) {

			var urls []string
			if os.Args[i][0] != '-' {

				//call functions to check the availability of each url
				urls = parseUniqueURLsFromFile(os.Args[i])

				//check if there are flags for JSON output or not
				if *jflag {

					checkURLJson(urls)
				} else {

					fmt.Println()
					fmt.Println(">>  ***** UrlChecker is working now...... *****  <<")
					fmt.Println("--------------------------------------------------------------------------------------------------")
					checkURL(urls)
				}
			}
			i++

		}

	}
}
