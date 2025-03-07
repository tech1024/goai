Go AI
==

[![build](https://github.com/tech1024/goai/actions/workflows/build.yml/badge.svg)](https://github.com/tech1024/goai/actions/workflows/build.yml)
[![Coverage Status](https://coveralls.io/repos/github/tech1024/goai/badge.svg?branch=main)](https://coveralls.io/github/tech1024/goai?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/tech1024/goai)](https://goreportcard.com/report/github.com/tech1024/goai)
[![Godoc](https://godoc.org/github.com/tech1024/goai?status.svg)](https://pkg.go.dev/github.com/tech1024/goai)
[![Release](https://img.shields.io/github/release/tech1024/goai.svg)](https://github.com/tech1024/goai/releases/latest)

A golang API library for AI Engineering.

This is a high level feature overview.

- Chat Completion
- Embedding

## Installation

```shell
go get -u 'github.com/tech1024/goai'
```

## Getting Started

```go
package main

import (
	"context"
	"log"

	"github.com/tech1024/goai"
	"github.com/tech1024/goai/provider/ollama"
)

func main() {
	ollamaClient, _ := ollama.NewClient("http://127.0.0.1:11434")
	chat := goai.NewChat(ollama.NewNewChatModel(ollamaClient, "deepseek-r1"))
	result, err := chat.Chat(context.Background(), "What can you do for me ?")
	if err != nil {
		log.Fatal(err)
	}

	log.Println(result)
}
```

## License

This project is licensed under the [Apache 2.0 license](LICENSE).

## Contact

If you have any issues or feature requests, please contact us. PR is welcomed.
- https://github.com/tech1024/goai/issues