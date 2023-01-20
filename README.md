# pathname

<a href="https://www.buymeacoffee.com/mjwhitta">🍪 Buy me a cookie</a>

[![Go Report Card](https://goreportcard.com/badge/github.com/mjwhitta/pathname)](https://goreportcard.com/report/github.com/mjwhitta/pathname)
[![Pipeline](https://github.com/mjwhitta/pathname/badges/master/pipeline.svg)](https://github.com/mjwhitta/pathname/-/pipelines/latest)
[![Coverage](https://github.com/mjwhitta/pathname/badges/master/coverage.svg)](https://github.com/mjwhitta/pathname)

## What is this?

A minimal Go port of Ruby's `Pathname` class. This mostly contains
functions I use a lot such as `Basename`, `Dirname`, and `ExpandPath`.
Ruby's `Exist?` method has been renamed `DoesExist` to be more
Go-like.

## How to install

Open a terminal and run the following:

```
$ go get --ldflags="-s -w" --trimpath -u github.com/mjwhitta/pathname
```

## How to use

Below is a sample usage to expand file paths accounting for `~` use:

```
package main

import (
    "fmt"

    "github.com/mjwhitta/pathname"
)

func main() {
    if ok, e := pathname.DoesExist("~/bin"); e != nil {
        panic(e)
    } else if ok {
        fmt.Println(pathname.ExpandPath("~/bin"))
    }

    if ok, e := pathname.DoesExist("~user/bin"); e != nil {
        panic(e)
    } else if ok {
        fmt.Println(pathname.Dirname("~user/bin"))
        fmt.Println(pathname.Basename("~user/bin"))
    }
}
```

## Links

- [Source](https://github.com/mjwhitta/pathname)
