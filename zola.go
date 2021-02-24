package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
)

func main() {
	input := flag.String("i", "testdata/posts", "Path/to/posts")
	output := flag.String("o", "testdata/newposts", "Path/to/newposts")

	flag.Parse()

	run(*input, *output)
}

func run(input string, output string) {

	postDir, err := os.Open(input)
	if err != nil {
		log.Fatalf("Failed opening directory: %s", err)
	}

	defer func() {
		if err = postDir.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	postList, _ := postDir.Readdirnames(0)
	for _, postFile := range postList {
		fmt.Printf("Processing: %s\t", postFile)

		// open the file to be migrated
		filePath := postDir.Name() + "/" + postFile
		file, err := os.Open(filePath)
		if err != nil {
			log.Fatal(err)
		}

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

		yamlPost := string(data[:])

		// Parse the YAML post and make it TOML
		tomlPost := postParser(yamlPost)

		// Convert tomlPost string to postBytes []bytes
		postBytes := []byte(tomlPost)

		// Write out new file
		outFile := output + "/" + postFile

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
		log.Printf("Wrote %d bytes.\n", bytesWritten)

	}
}

// postParser converts post front matter from YAML to TOML
func postParser(yamlPost string) string {

	// Regex pattern to match YAML --- and <element>:
	// Initial --- with be at the start of the string (^---)
	// The closing --- will follow a new line (\n---)
	// Hence the OR in the regex pattern
	var reYAML = regexp.MustCompile(`(^|\n)(---)`)
	var reElement = regexp.MustCompile(`(\n[a-z]*)(:)(.*)?`)

	// Replace matched items with +++ and <element> =
	yamlPost = reYAML.ReplaceAllString(yamlPost, "${1}+++")
	yamlPost = reElement.ReplaceAllString(yamlPost, "${1} =${3}")

	return yamlPost

}
