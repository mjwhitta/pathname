# pathname

A minimal Golang port of Ruby's `Pathname` class. This mostly contains
functions I use a lot such as `Basename`, `Dirname`, and `ExpandPath`.
Ruby's `Exist?` method has been renamed `DoesExist` to be more
Golang-like.

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
    if pathname.DoesExist("~/bin") {
        fmt.Println(pathname.ExpandPath("~/bin"))
    }
    if pathname.DoesExist("~user/bin") {
        fmt.Println(pathname.Dirname("~user/bin"))
        fmt.Println(pathname.Basename("~user/bin"))
    }
}
```

## Links

- [Source](https://gitlab.com/mjwhitta/pathname)
