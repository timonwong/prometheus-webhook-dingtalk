# Gocov HTML export

This is a simple helper tool for generating HTML output from
[axw/gocov](https://github.com/axw/gocov/)

Here is a screenshot:

![HTML coverage report screenshot](https://github.com/matm/gocov-html/blob/master/gocovh-html.png)

## Installation

Just type the following to install the program and its dependencies:
```
$ go get github.com/axw/gocov/gocov
$ go get github.com/matm/gocov-html
```

## Usage

`gocov-html` can read a JSON file or read from standard input:
```
$ gocov test strings | gocov-html > strings.html
ok      strings 0.700s  coverage: 98.1% of statements
```

The generated HTML content comes along with a default embedded CSS. Use the `-s` 
flag to use a custom stylesheet:
```
$ gocov test net/http | gocov-html -s mystyle.css > http.html
```
