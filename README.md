# query [![Build Status](https://travis-ci.org/inappcloud/query.svg?branch=master)](https://travis-ci.org/inappcloud/query)

SQL Query Builder in Go.

# MongoDB-like where

``` go
package main

import (
  "github.com/inappcloud/query"
  "github.com/inappcloud/query/where"
)

func main() {
  c := where.Parse(`{"$or":[{"price":{"$gt":10},{"name":"phoenix"}]}`)
  c.String() // "(price > 10 OR name = ?)"
  c.Params() // []interface{}{"phoenix"}

  q := query.Select("teas").Where(where.And(c, where.Eq("id", 1)))
  q.String() // "SELECT * FROM teas WHERE (price > 10 OR name = ?) AND id = 1"
  q.Params() // []interface{}{"phoenix", 1}
}
```
