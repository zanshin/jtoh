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
	handles []int
	counter int
)

func main() {
	input := flag.String("i", "testdata/posts", "Path/to/posts")
	output := flag.String("o", "testdata/newposts", "Path/to/newposts")

	flag.Parse()

	inPath = *input
	outPath = *output

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
	for _, postFile := range postList {
		// fmt.Printf("Processing: %s\t", postFile)

		// process the post using goroutines to provide concurrency
		// and to keep total file handle count down
		wg.Add(1)
		go func(postFile string) {
			defer wg.Done()
			postProcess(postFile)
		}(postFile)
	}

	wg.Wait()

	fmt.Printf("\n\nhandle history: %v", handles)

}

// postProcess encapsulates the post processing process.
func postProcess(postFile string) {

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
	yamlPost := string(data[:])
	tomlPost := postParser(yamlPost)
	postBytes := []byte(tomlPost)

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
	bytesWritten, err := newFile.Write(postBytes)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Wrote %s containing %d bytes.\n", postFile, bytesWritten)

	counter--
	handles = append(handles, counter)

}

// postParser converts post front matter from YAML to TOML
func postParser(yamlPost string) string {

	// Regex pattern to match YAML --- and <element>:
	// Initial --- with be at the start of the string (^---)
	// The closing --- will follow a new line (\n---)
	// Hence the OR in the regex pattern
	var reYAML = regexp.MustCompile(`(^|\n)(---)`)

	// YAML elements in Jekyll are: element: value
	// TOML elements in Zola   are: element = value
	var reElement = regexp.MustCompile(`(\n[a-z]*)(:)(.*)?`)

	// YouTube shortcode
	// Jekyll format: {% youtube JdxkVQy7QLM %}
	// Zola   format: {{ youtube(id="JdxkVQy7QLM") }}
	// var reYT = regexp.MustCompile(`({% )(youtube)(\s)(.{11})( %})`)
	var reYT = regexp.MustCompile(`({% )(youtube)(\s)(.*)( %})`)

	// Replace matched items with +++ and <element> =
	yamlPost = reYAML.ReplaceAllString(yamlPost, "${1}+++")
	yamlPost = reElement.ReplaceAllString(yamlPost, "${1} =${3}")
	yamlPost = reYT.ReplaceAllString(yamlPost, "{{ $2(id=\"$4\") }}")

	return yamlPost

}
