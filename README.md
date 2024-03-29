# Name
jtoh - Convert Jekyll-formatted blog posts to Hugo-formatting

# Synopsis

    jtoh -i source -o destination [-t true]

# Description
`jtoh` converts a number of front matter elements, as well as some shortcodes from Jekyll formatting
to Hugo formatting. See [DETAILS](DETAILS.md) for more information.

# Options

    -i - The path to the postings to convert

    -o - The path were the converted postings should be written

    -t - OPTIONAL - If true, will convert front matter from YAML to TOML

# Installing
Clone the repository to your $GOPATH.

    git clone https://github.com/zanshin/jtoh.git $GOPATH/src/github/zanshin/jtoh

# Building
The repository does not contain a pre-built executable. One can be made for the current platform
using:

    go build .

# Running

    $ jtoh -i source -o dest
    Files to process: 2248

    Front Matter conversions
    Date and Time formats converted:    479
    Date only converted:                1767
    Quotes stripped from dates:         282
    Leading zero added to M:            1351
    Leading zero added to D:            527
    Categories converted to tags:       2248

    Shortcods converted
    YouTube shortcodes converted:       70
    Highlight shortcodes converted:     66
    End Highlight shortcodes converted: 66
    Image shortcodes converted:         88

    Total number of bytes processed:    4214704

    Post counter : 2248
    Parse counter: 2248
