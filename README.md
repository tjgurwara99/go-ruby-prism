# go-ruby-prism

This is a fork of [go-ruby-prism](https://github.com/danielgatis/go-ruby-prism).
I wanted to use it but I did't need some of the package dependencies, so I've cleaned
them up - I don't think the author wants this change so I'm keeping it in the fork and
not opening a PR upstream.

The go-ruby-prism is package that leverages the Ruby Prism parser compiled to WebAssembly for parsing Ruby code without the need for CGO.

## Features

- **CGO-Free**: Go-Ruby-Prism utilizes the [Ruby Prism parser](https://github.com/ruby/prism) compiled to WebAssembly, eliminating the need for CGO bindings.
- **Simplified Integration**: Seamlessly integrate Ruby code parsing into your Go applications with minimal setup.
- **High Performance**: Harnesses the efficiency of WebAssembly for speedy and efficient parsing of Ruby code.
- **Cross-Platform**: Works across various platforms supported by Go, ensuring compatibility in diverse environments.

## Usage

Here's a basic example demonstrating how to use this package:

```go
package main

import (
	"context"
	"fmt"

	parser "github.com/danielgatis/go-ruby-prism/parser"
)

func main() {
	ctx := context.Background()

	p, _ := parser.NewParser(ctx)
	source := "puts 'Hello, World!'"
	result, _ := p.Parse(ctx, source)
	fmt.Println(result)
}
```

You can find more examples in the examples folder.


## License

Original Copyright Notice would remain in this repository under (c) 2024-present [Daniel Gatis](https://github.com/danielgatis)

Licensed under [MIT License](./LICENSE.txt)
