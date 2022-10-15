# Introduction
This project exists solely to convert blog posting files that were created for Jekyll, to ones that
will work with [Hugo](https://gohugo.io "Hugo").

## Particulars
A number of things need to happen to prepare my old blog postings to work with Hugo. My site has
just over 2250 postings, spanning 22 years and several blogging back ends. Originally hand
coded, the site went through Blogger, MoveableType, WordPress, and Octopress, before ending up on
Jekyll.

The items needing attention are:
- Shortcodes: YouTube and Gist
- Code highlighting
- Other Liquid markup
- Date formatting
- TOML vs. YAML front matter

### TODO
- [X] {% highlight %} {% endhighlight %}
- [X] {% gist ###### %}
- [ ] {% raw %} {% endraw %}
- [X] {% if ... %}
- [X] {% elsif %}
- [X] Capture and display conversion event counts
- [ ] Inline {% highlight %}{% endhighlight %} instances

### Gist Shortcodes
The Jekyll gist shortcode works with only the ID number portion of the URL. The Hugo one requires
the gist ID number AND the gist account name in order to function. A grep of my postings shows that
there are only 5 instances of the gist shortcode to be modified. Far easier to do by hand than to
fix via code.

### Code Highlighting
While the Hugo highlight shortcode offers more features than the Jekyll one, they both have the
same data requirements: the word "highlight" and a language name. The brackets surrounding these two
pieces of information changes from `{% ... %}` to `{{< ... >}}`.

The `{% endhightlight %}` in Jekyll becomes `{{< / highlight >}}` in Hugo.

### Raw Shortcodes
The `raw` shortcode works (mostly) in tandem with the `highlight` shortcode. It allows you to put
anything in for the code block, and ignores it. I don't see an equivalent code in Hugo, so I am
going to eliminate the `raw` code occurrences. There are only 11, so if I have to do some manual
editing of posts, it would be manageable.

### YouTube Shortcodes
Hugo has a slightly different format for their YouTube shortcode than Jekyll's format. Via regular
expressions `{% youtube JdxkVQy7QLM %}` becomes `{{ youtube(id="JdxkVQy7QLM") }}`.

### IF and ELSIF
While my  initial grepping showed that there were some `{% if %}` and `{% elsif %}` tags in my
postings, a closer look reveals that they are all in code samples, and not part of the site that
needs converting.

### Date Formatting
Due to the age of my site, and the different blogging systems used, the front matter date is
inconsistently formatted. Some dates are in double quotes, some are not. Some have a time specified,
others do not. All the dates need to be formatted the same.

There are three regex transforms that make this happen. One strips the quotes, the second handles
entries that have both the date and the time, and the last handles entries that only have the date.

The end result is a date in the format YYYY-MM-DDTHH:MM. If the incoming date does not have a time,
the time of `00:01` is used.

### TOML vs. YAML Front Matter
Originally I was going to convert all my posting to use a TOML formatted front matter. Later I
decided that wasn't really necessary. I did add a flag so that TOML front matter can be created, if
so desired.

These YAML specific elements will be converted to TOML specific elements.

* The `---` will become `+++`
* The `title:` will become `title =`

## Processing
Rather than read each posting file line-by-line and process them that way, I treat each posting as a
single, large string. This makes the regex portion of the process easier with one slight
complication.

The order of operations performed is critical. The dates must be re-formatted before the YAML to
TOML conversion, if that option has been selected. The regex matching patterns expect the date line
to have `date:` not `date =`.

A count of each type of conversion performed is kept and displayed.

Processing 2252 files takes only a few seconds.

    Files to process: 2252

    YouTube shortcodes converted:       70
    Date and Time formats converted:    482
    Date only converted:                1770
    Quotes stripped from dates:         282
    Highlight shortcodes converted:     70
    End Highlight shortcodes converted: 70

    Total number of bytes processed:    4241847

    Post counter : 2252
    Parse counter: 2252
    real    0m3.607s
    user    0m0.487s
    sys     0m0.441s

## Running
Clone the project onto your `$GOPATH`. Create a temporary directory that has two sub-directories:
`dest` and `source`. Copy all your postings to `source`. My `$GOPATH` is
`~/code/go` so after cloning the project to `~/code/go/src/github` I run this command, while inside
the temporary directory:

    time go run ~/code/go/src/github.com/zanshin/jtoh/main.go -i source -o dest

Adding `time` to the start of the command will show how long the conversion took.
