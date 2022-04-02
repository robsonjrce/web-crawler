# A simple web crawler written in golang

For this task we have some points to focus on:

- It is a mirror web crawler
- Should find and download pages
  - A new entry is represented by `href` on a `<a>` tag
  - A valid url is in a sublevel for the original url call
- The crawler should graceful stop on `CTRL + C`
- Perform work in parallel where reasonable
- Support resume

## How to build

You should have a working `go` installation.

We are using a make task to be able to build the application:
```
$ make build
```

After that, you should be able to find it in: 

```
$ ls bin
web-crawler
```

## How to use it

To run the app, its simple:

```
$ ./bin/web-crawler walk --url url_to_crawl --dest destination_path_save_data
```

## Implementation Steps

- First we have defined the cli interface
- As second step we will work on the finding suitable data
  - Means that we will look for valid `<a>` tags and look for the `href` attribute
