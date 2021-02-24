# Introduction
This project exists solely to convert blog posting files that were created for Jekyll, to ones that
will work with [Zola](https://www.getxola.org "Zola").

## Particulars
### YAML vs. TOML
The postings in Jekyll all have YAML front matter. This block of lines is at the top of every
posting. It is delimited by a line containing `---` both before and after the block. The elements
within the block — metadata about the post — are all YAML format, e.g., `title: "My Title"`.

These YAML specific elements will be converted to TOML specific elements.

* The `---` will become `+++`
* The `title:` will become `title =`

The remainder of the file is left unchanged. The updated version of the post is written to a new
directory, preserving the original files.
