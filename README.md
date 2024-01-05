[![Go Reference](https://pkg.go.dev/badge/github.com/billysutomo/hnsdk.svg)](https://pkg.go.dev/github.com/billysutomo/hnsdk)

# hnsdk

Golang SDK for hackernews API

## Installation

```bash
go get github.com/billysutomo/hnsdk
```

## Usage

```go
package main

import (
	"context"
	"fmt"

	"github.com/billysutomo/hnsdk"
)

func main() {
	client := hnsdk.NewClient()

	ctx := context.Background()
	topStoryIDs, err := client.GetTopStories(ctx, 100)
	if err != nil {
		panic(err)
	}

	item, err := client.GetItem(ctx, topStoryIDs[0])
	if err != nil {
		panic(err)
	}
	fmt.Println(item.Title)

	topStories, err := client.GetTopStoriesWithData(ctx, 100)
	if err != nil {
		panic(err)
	}

	fmt.Println(topStories[0].Title)
}
```
