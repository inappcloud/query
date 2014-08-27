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
  c.String() // "(price > ? OR name = ?)"
  c.Params() // []interface{}{10, "phoenix"}

  q := query.Select("teas").Where(where.And(c, where.Eq("id", 1)))
  q.String() // "SELECT * FROM teas WHERE (price > $1 OR name = $2) AND id = $3"
  q.Params() // []interface{}{10, "phoenix", 1}
}
```
