# Introduction
This project exists solely to convert blog posting files that were created for Jekyll, to ones that
will work with [Hugo](https://gohugo.io "Hugo").

## Particulars
A number of things need to happen to prepare my old blog postings to work with Hugo. My site has
approximately 2100 postings, spanning 22 years and several blogging back ends. Originally hand
coded, the site went through Blogger, MoveableType, WordPress, and Octopress, before ending up on
Jekyll.

The items needing attention are:
- YouTube shortcodes
- Date formatting
- TOML vs. YAML front matter

### TODO
In addition to YouTube shortcodes, code hightlight needs to be addressed. Many code examples are
simply indented 4 spaces as per Markdown rules. However, some leverage a Jekyll feature and use a
highlight tag, with a language specified, to colorize the code. All instances of these highlight
tags will need to be updated.

### YouTube Shortcodes
Hugo has a slightly different format for their YouTube shortcode than Jekyll's format. Via regular
expressions `{% youtube JdxkVQy7QLM %}` becomes `{{ youtube(id="JdxkVQy7QLM") }}`.

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

Processing 2181 files takes only a few seconds.

    Files to process: 2181

    real    0m3.013s
    user    0m0.379s
    sys     0m0.453s
