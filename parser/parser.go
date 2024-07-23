package parser

import (
	"context"
	"fmt"

	"github.com/tjgurwara99/go-ruby-prism/wasm"
)

type Parser struct {
	runtime *wasm.Runtime
}

func NewParser(ctx context.Context) (*Parser, error) {
	runtime, err := wasm.NewRuntime(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate wasm runtime: %w", err)
	}

	return &Parser{
		runtime: runtime,
	}, nil
}

func (p *Parser) Close(ctx context.Context) error {
	if err := p.runtime.Close(ctx); err != nil {
		return fmt.Errorf("failed to close the wasm runtime: %w", err)
	}

	return nil
}

func (p *Parser) Parse(ctx context.Context, source []byte) (result *ParseResult, err error) {
	result = nil
	err = nil

	defer func() {
		if r := recover(); r != nil {
			result = nil
			err = fmt.Errorf("recovered from panic: %v", r)
		}
	}()

	result, err = p.parseWithOptions(ctx, source, newParseOptions())

	if err != nil {
		return nil, fmt.Errorf("failed to parse with options: %w", err)
	}

	return result, nil
}

func (p *Parser) parseWithOptions(ctx context.Context, source []byte, opts *parseOptions) (*ParseResult, error) {
	sourcePtr, err := p.runtime.Calloc(ctx, 1, uint64(len(source)))
	if err != nil {
		return nil, fmt.Errorf("failed to allocate memory for source: %w", err)
	}

	if !p.runtime.MemoryWrite(sourcePtr, source) {
		return nil, fmt.Errorf("failed to write the source into memory: %w", err)
	}

	// put option into memory
	optBytes, err := opts.bytes()
	if err != nil {
		return nil, fmt.Errorf("failed to convert options into bytes: %w", err)
	}

	optPtr, err := p.runtime.Calloc(ctx, 1, uint64(len(optBytes)))
	if err != nil {
		return nil, fmt.Errorf("failed to allocate memory for options: %w", err)
	}

	if !p.runtime.MemoryWrite(optPtr, optBytes) {
		return nil, fmt.Errorf("failed to write the options into memory: %w", err)
	}

	// call the serialize parse function
	bufferSizeOf, err := p.runtime.BufferSizeOf(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get the buffer size: %w", err)
	}

	bufferPtr, err := p.runtime.Calloc(ctx, bufferSizeOf, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to get the buffer ptr: %w", err)
	}

	if err := p.runtime.BufferInit(ctx, bufferPtr); err != nil {
		return nil, fmt.Errorf("failed to init the buffer: %w", err)
	}

	if _, err := p.runtime.SerializeParse(ctx, bufferPtr, sourcePtr, uint64(len(source)), optPtr); err != nil {
		return nil, fmt.Errorf("failed to call the parse function: %w", err)
	}

	// read result from memory
	bufferValue, err := p.runtime.BufferValue(ctx, bufferPtr)
	if err != nil {
		return nil, fmt.Errorf("failed to get the buffer value: %w", err)
	}

	bufferLen, err := p.runtime.BufferLength(ctx, bufferPtr)
	if err != nil {
		return nil, fmt.Errorf("failed to get the buffer length: %w", err)
	}

	serializedBytes, ok := p.runtime.MemoryRead(bufferValue, bufferLen)
	if !ok {
		return nil, fmt.Errorf("failed to read the buffer content from memory: %w", err)
	}

	// free memory
	if err := p.runtime.BufferFree(ctx, bufferPtr); err != nil {
		return nil, fmt.Errorf("failed to free memory for buffer ptr: %w", err)
	}

	if err := p.runtime.Free(ctx, sourcePtr); err != nil {
		return nil, fmt.Errorf("failed to free memory for source ptr: %w", err)
	}

	if err := p.runtime.Free(ctx, bufferPtr); err != nil {
		return nil, fmt.Errorf("failed to free memory for buffer ptr: %w", err)
	}

	if err := p.runtime.Free(ctx, optPtr); err != nil {
		return nil, fmt.Errorf("failed to free memory for option ptr: %w", err)
	}

	result, err := deserialize(serializedBytes, source)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize the result: %w", err)
	}

	return result, nil
}
