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

	before:= `---
title: "My Title"
category: blog
layout: post
date: 2021-02-18 17:05
link:
---
Now is the winter: of our discontent.
{% youtube JdxkVQy7QLM %}`

	after := `---
title: "My Title"
category: blog
layout: post
date: 2021-02-18T17:05:00
link:
---
Now is the winter: of our discontent.
{{ youtube(id="JdxkVQy7QLM") }}`

	result := postParser(before)

	if result != after {
		fmt.Fprintf(os.Stderr, "Parser failed. Expected: \n%q. \nGot: \n%q\n", after, result)
	}
}

func TestDateQuoteParser(t *testing.T) {
	before:= `---
title: "My Title"
category: blog
layout: post
date: "2021-02-18"
link:
---
Now is the winter: of our discontent.
{% youtube JdxkVQy7QLM %}`

	after := `---
title: "My Title"
category: blog
layout: post
date: 2021-02-18T03:02:00
link:
---
Now is the winter: of our discontent.
{{ youtube(id="JdxkVQy7QLM") }}`

	result := postParser(before)

	if result != after {
		fmt.Fprintf(os.Stderr, "Parser failed. Expected: \n%q. \nGot: \n%q\n", after, result)
	}

}

func TestDateTimeParser(t *testing.T) {
	before:= `---
title: "My Title"
category: blog
layout: post
date: 2021-02-18 12:34
link:
---
Now is the winter: of our discontent.
{% youtube JdxkVQy7QLM %}`

	after := `---
title: "My Title"
category: blog
layout: post
date: 2021-02-18T12:34:00
link:
---
Now is the winter: of our discontent.
{{ youtube(id="JdxkVQy7QLM") }}`

	result := postParser(before)

	if result != after {
		fmt.Fprintf(os.Stderr, "Parser failed. Expected: \n%q. \nGot: \n%q\n", after, result)
	}

}

func TestDateNoTimeParser(t *testing.T) {
	before:= `---
title: "My Title"
category: blog
layout: post
date: 2021-02-18
link:
---
Now is the winter: of our discontent.
{% youtube JdxkVQy7QLM %}`

	after := `---
title: "My Title"
category: blog
layout: post
date: 2021-02-18T03:02:00
link:
---
Now is the winter: of our discontent.
{{ youtube(id="JdxkVQy7QLM") }}`

	result := postParser(before)

	if result != after {
		fmt.Fprintf(os.Stderr, "Parser failed. Expected: \n%q. \nGot: \n%q\n", after, result)
	}

}

func TestHighlight(t *testing.T) {
	before:= `---
title: "My Title"
category: blog
layout: post
date: 2021-02-18
link:
---
Now is the winter: of our discontent.
{% highlight bash %}
blah
blah
blah
{% endhighlight %}`

	after := `---
title: "My Title"
category: blog
layout: post
date: 2021-02-18T03:02:00
link:
---
Now is the winter: of our discontent.
{{< highlight bash >}}
blah
blah
blah
{{< / highlight >}}`

	result := postParser(before)

	if result != after {
		fmt.Fprintf(os.Stderr, "Parser failed. Expected: \n%q. \nGot: \n%q\n", after, result)
	}

}
