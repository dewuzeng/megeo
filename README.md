# megeo
Distance calculation based on Redis-GEO.

```go
package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/somecodeio/megeo"
)

func main() {

	conn, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		panic(err)
	}
	geo := megeo.NewMeGEO(conn, "stores:locations")
	items := megeo.GEOItems{
		{
			Longitude: 113.33752,
			Latitude:  23.09787,
			Member:    441,
		},
		{
			Longitude: 113.321316,
			Latitude:  23.095648,
			Member:    453,
		},
	}
	gotNum, err := geo.BatchAdd(items)
	if err != nil {
		panic(err)
	}
	fmt.Printf("update %d items \n", gotNum)
	gotItems, err := geo.Radius(113.3172, 23.08331, 10)
	if err != nil {
		panic(err)
	}
	for _, item := range gotItems {
		fmt.Printf("distince = %f, member = %d \n", item.Distance, item.Member)
	}
}

```