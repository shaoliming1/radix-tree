# radix tree

it's a radix tree implemented by Golang.

## example
```go
package main

import (
	"fmt"
	radix_tree "github.com/shaoliming1/radix-tree"
)

func main()  {
	t :=radix_tree.NewRadixTree()
	t.Insert("foo", 1)
	t.Insert("foobar", 2)
	t.Insert("make",3)
	ret := t.Find("foo")
	t.Delete("foo")
	fmt.Print(ret)
}
``` 
