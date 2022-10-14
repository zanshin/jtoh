package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"sync"
)

var (
	inPath  string
	outPath string
	convert bool
	handles []int
	counter int
	postCtr int
	parseCtr int
	// quotesCtr int
	// timeCtr int
	// ytCtr int
)

// For limiting threads.
var tokens = make(chan struct{}, 30)

func main() {
	input := flag.String("i", "testdata/posts", "Path/to/posts")
	output := flag.String("o", "testdata/newposts", "Path/to/newposts")
	toml := flag.Bool("t", false, "Convert YAML to TOML")

	flag.Parse()

	inPath = *input
	outPath = *output
	convert = *toml

	run()
}

func run() {

	postDir, err := os.Open(inPath)
	if err != nil {
		log.Fatalf("Failed opening directory: %s", err)
	}

	defer func() {
		if err = postDir.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	var wg sync.WaitGroup

	postList, _ := postDir.Readdirnames(0)
	fmt.Printf("Files to process: %d\n", len(postList))

	postCtr = 0
	parseCtr = 0

	for _, postFile := range postList {

		// process the post using goroutines to provide concurrency
		// and to keep total file handle count down
		wg.Add(1)
		postCtr++
		go processPost(postFile, &wg)

	}

	wg.Wait()

	// fmt.Printf("\n\nhandle history: %v", handles)
	fmt.Printf("\nPost counter : %d", postCtr)
	fmt.Printf("\nParse counter: %d", postCtr)

}

// postProcess encapsulates the post processing process.
func processPost(postFile string, wg *sync.WaitGroup) {

	// wait until all trheads are done
	defer wg.Done()

	// acquire token
	// fmt.Println("token acquired")
	tokens <- struct{}{}

	inFile := inPath + "/" + postFile

	// open the file to be migrated
	file, err := os.Open(inFile)
	if err != nil {
		log.Fatal(err)
	}

	counter++
	handles = append(handles, counter)

	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	// read all the contents of the file
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("Failed reading file %s with error: %s", data, err)
	}

	// Convert to string, parse, convert result to []byte
	original := string(data[:])
	converted := postParser(original)
	postBytes := []byte(converted)

	// Write out new file
	outFile := outPath + "/" + postFile

	newFile, err := os.OpenFile(
		outFile,
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0666,
	)

	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Write bytes to file
	// bytesWritten, err := newFile.Write(postBytes)
	_, err = newFile.Write(postBytes)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Printf("Wrote %s containing %d bytes.\n", postFile, bytesWritten)

	counter--
	handles = append(handles, counter)

	// release token
	// fmt.Println("token released")
	<-tokens

}

// postParser converts post front matter from YAML to TOML
func postParser(post string) string {

	parseCtr++

	// Regex pattern to match YAML --- and <element>:
	// Initial --- with be at the start of the string (^---)
	// The closing --- will follow a new line (\n---)
	// Hence the OR in the regex pattern
	var reYAML = regexp.MustCompile(`(^|\n)(---)`)

	// YAML elements in Jekyll are: element: value
	// TOML elements in Hugo are: element = value
	var reElement = regexp.MustCompile(`(\n[a-z]*)(:)(.*)?`)

	// YouTube shortcode
	// Jekyll format: {% youtube JdxkVQy7QLM %}
	// Hugo   format: {{ youtube(id="JdxkVQy7QLM") }}
	// var reYT = regexp.MustCompile(`({% )(youtube)(\s)(.{11})( %})`)

	var reYT = regexp.MustCompile(`({% )(youtube)(\s)(.*)( %})`)

	// Posts have a variety of date formats. Some dates are enclosed in
	// double-quotes, some have 2-digits for month or day, while others do not.
	// Some dates include a time, others do not.
	// Goal: have all dates in the form YYYY-MM-DDTHH:MM
	// There are three regex expressions in use:
	// reQuotes   - strips any double quotes found from around dates
	// reDate     - matches dates without a time, spams T00:01 at the time
	// reDateTime - matches dates that have times, preserves the time
	//
	// All of this is complicated by the treatment of each post as a string,
	// and not a file of individual lines. Therefore, even though "date:"
	// appears at the start of a new line, we cannot use the ^ to indicate
	// beginning of string. Instead we must us the new line indicator (\n).
	// To bound the regex expression, a trailing new line character is also
	// used.
	//
	var reQuotes = regexp.MustCompile(`\n(date:).*(")(.*)(")\n`)
	var reDate = regexp.MustCompile(`\n(date:)\s*((19|20)[0-9][0-9])-([0|1]?[0-9])-([0|1|2|3]?[0-9])\n`)
	var reDateTime = regexp.MustCompile(`\n(date:)\s*((19|20)[0-9][0-9])-([0|1]?[0-9])-([0|1|2|3]?[0-9])\s{1}([0-1]?[0-9]|2[0-3]):([0-5][0-9])(.*)?\n`)


	// Strip quotes from dates
	post= reQuotes.ReplaceAllString(post, "\n${1} ${3}\n")

	// Format timeless dates and dates with times
	post= reDateTime.ReplaceAllString(post, "\n${1} ${2}-${4}-${5}T${6}:${7}\n")
	post= reDate.ReplaceAllString(post, "\n${1} ${2}-${4}-${5}T00:01\n")

	// TOML Conversion, if requested
	// Replace matched items with +++ and <element> =
	// MUST HAPPEN AFTER THE DATE REGEX, otherwise it won't work...
	if convert {
		post = reYAML.ReplaceAllString(post, "${1}+++")
		post = reElement.ReplaceAllString(post, "${1} =${3}")
	}

	// Convert YoutTube shortcode
	post= reYT.ReplaceAllString(post, "{{ $2(id=\"$4\") }}")

	return post

}
