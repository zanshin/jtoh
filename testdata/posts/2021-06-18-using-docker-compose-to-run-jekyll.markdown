---
layout: post
title: "Using Docker Compose to Run Jekyll"
date: 2021-06-18 13:38
categories: nerdliness
link:
---
While setting up my new computer, one of the tasks I was faced with was recreating my
[Jekyll](https://jekyllrb.com "Jekyll") environment. Jekyll is the Ruby based static site generator
I use for this site. MacOS has not always kept up with Ruby, and so I've had to employ
[rbenv](https://github.com/rbenv/rbenv "rbenv") to have an up to date Ruby installed. And to
segregate the Gem dependencies Jekyll has from any other Ruby project that might have conflicting
dependencies.

I kept putting off the task as I wasn't looking forward to the possibility of some problems in Ruby
Gem land. And then I thought about Docker—how hard would it be to use Docker as my Jekyll
environment? As it turns out, not very hard at all. Eight lines of YAML total, not very hard.

A quick search lead me to this article, [How to run Jekyll locally with Docker and Docker
Compose](https://caioteixeira.dev/blog/jekyll-docker/ "ow to run Jekyll locally with Docker and
Docker Compose"). Using `docker-compose` and a Docker image that has Jekyll installed on it, you can
quickly spin up a fully functioning Jekyll environment.

By putting the `docker-compose.yml` file in the root directory of my blog, and running

    docker-compose up

I get Jekyll without having to install Jekyll, manage Gems, or use rbenv to control which version of
Ruby and which collection of Gems is in use. Best of all, since Ruby **is** installed by default
(version 2.6.3p62 as of this writing), I have `rake` at my disposal, so the `Rakefile` I have still
works. This file has tasks for creating new postings, or new draft postings, for publishing drafts,
and for deploying my site using `rsync`.

The only changes I made to the `docker-compose.yml` file described in the article, is that I used
different flags on the `jekyll serve` command. I'm using

    jekyll serve --watch --drafts

which is what I've had for a long time. The performance of the Docker image is a bit slower than it
would be if it were natively installed. I may try `--incremental` at some point to see if that
speeds things up.

Here's my `docker-compose.yml` file.

    services:
      jekyll:
        image: jekyll/jekyll:latest
        command: jekyll serve --watch --drafts
        ports:
          - 4000:4000
        volumes:
          - .:/srv/jekyll

Briefly, this file says, compose a Docker image using the latest jekyll base image. Once it's ready
run the Jekyll server, watching for new drafts, and expose the server on port 4000. Finally, map the
current directory `.`, to `/srv/jekyll` in the image. Done.

Docker and docker-compose saved me a lot of time and effort. Now I wonder where else I could use it.
