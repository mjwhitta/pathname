# pathname

A minimal Golang port of Ruby's `Pathname` class.

## How to install

Open a terminal and run the following:

```
$ go get -u gitlab.com/mjwhitta/pathname
```

## How to use

Below is a sample usage to expand file paths accounting for `~` use:

```
package main

import (
    "fmt"

    "gitlab.com/mjwhitta/pathname"
)

func main() {
    fmt.Println(pathname.ExpandPath("~/bin"))
    fmt.Println(pathname.ExpandPath("~user/bin"))
}
```
