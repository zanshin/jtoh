package main

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func TestMain(m *testing.M) {

	binName := "main"
	fmt.Println("Building tool...")
	build := exec.Command("go", "build", "-o", binName)

	if err := build.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Cannot build tool %s: %s", binName, err)
		os.Exit(1)
	}

	fmt.Println("Running tests...")
	m.Run()

	fmt.Println("Cleaning up...")
	os.Remove(binName)
	os.Exit(0)
}

func TestPostParser(t *testing.T) {
	// yamlPost := "---\ntitle: \"My Title\"\ncategory: blog\n---\nThis is test content"
	// tomlPost := "+++\ntitle = \"My Title\"\ncategory = blog\n+++\nThis is test content"

	yamlPost := `---
title: "My Title"
category: blog
layout: post
date: 2021-02-18 17:05
link:
---
Now is the winter: of our discontent.
{% youtube JdxkVQy7QLM %}`

	tomlPost := `+++
title = "My Title"
category = blog
layout = post
date = 2021-02-18 17:05
link =
+++
Now is the winter: of our discontent.
{{ youtube(id="JdxkVQy7QLM") }}`

	result := postParser(yamlPost)

	if result != tomlPost {
		fmt.Fprintf(os.Stderr, "Parser failed. Expected: \n%q. \nGot: \n%q\n", tomlPost, result)
	}

}
