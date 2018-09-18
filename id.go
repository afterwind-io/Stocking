package stocking

import (
	"fmt"
	"math"
)

// id is a channel for fetching an auto-increment, int64 based string key
var id = make(chan string)

func generateID() string {
	i := int64(0)

	for {
		i++
		if i == math.MaxInt64 {
			i = 0
		}

		id <- fmt.Sprint(i)
	}
}
