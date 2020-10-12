package main

import (
	//"flag"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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

				//fmt.Println(v + ": NO RESPONCE!")
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

	flag.Parse()
	//deal with non-file path, giving usage message
	if len(os.Args) == 1 {
		fmt.Println("help/usage message: To run this program, please pass an argument to it,i.e.: go run urlChecker.go urls.txt")

	} else {
		//feature of checking version
		if *vflag {
			fmt.Println("  *****  urlChecker Version 0.1  *****  ")
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

		//use for loop to deal with multiple file paths
		i := 1
		for i+1 <= len(os.Args) {

			var urls []string
			if os.Args[i][0] != '-' {

				//open file and read it
				content, err := ioutil.ReadFile(os.Args[i])
				if err != nil {
					log.Fatal(err)
				}
				textContent := string(content)

				//call functions to check the availability of each url
				urls = removeDuplicate(extractURL(textContent))

				//check if there are flags for JSON output or not
				if *jflag {
					//urlsJ := make([]UrlJson, 0)
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
