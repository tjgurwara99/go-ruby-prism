package main

import (
	"context"
	"fmt"
	"os"

	parser "github.com/tjgurwara99/go-ruby-prism/parser"
)

type visitor struct{}

func newVisitor() *visitor {
	return &visitor{}
}

func (v *visitor) Visit(node parser.Node) {
	fmt.Printf("%T\n", node)
}

func (v *visitor) traverse(node parser.Node) {
	node.Accept(v)
	for _, child := range node.Children() {
		v.traverse(child)
	}
}

func main() {
	ctx := context.Background()

	p, _ := parser.NewParser(ctx)
	defer p.Close(ctx)

	source := "puts 'Hello, World!'"
	result, err := p.Parse(ctx, []byte(source))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	visitor := newVisitor()
	visitor.traverse(result.Value)
}
