package parser

import (
	"bytes"
	"errors"
	"fmt"
)

const prismHeader = "PRISM"
const majorVersion = 0
const minorVersion = 24
const patchVersion = 0

func deserialize(serialized []byte, source []byte) (*ParseResult, error) {
	buff := newBuffer(serialized)

	// check header
	header := make([]byte, 5)
	_, err := buff.read(header)
	if err != nil {
		return nil, fmt.Errorf("error reading header: %w", err)
	}

	if !bytes.Equal(header, []byte(prismHeader)) {
		return nil, errors.New("invalid prism header")
	}

	// check version
	version := make([]byte, 3)
	_, err = buff.read(version)
	if err != nil {
		return nil, fmt.Errorf("error reading version: %w", err)
	}

	if !bytes.Equal(version, []byte{majorVersion, minorVersion, patchVersion}) {
		return nil, errors.New("invalid version number")
	}

	// check location
	var location = make([]byte, 1)
	_, err = buff.read(location)
	if err != nil {
		return nil, fmt.Errorf("error reading location: %w", err)
	}

	if !bytes.Equal(location, []byte{0}) {
		return nil, errors.New("requires no location fields in the serialized output")
	}

	// reading encoding and discard it is always UTF-8
	encodingLen, err := loadVarUInt(buff)
	if err != nil {
		return nil, fmt.Errorf("error reading encoding length: %w", err)
	}

	var encoding = make([]byte, encodingLen)
	_, err = buff.read(encoding)
	if err != nil {
		return nil, fmt.Errorf("error reading encoding: %w", err)
	}

	// skip start line and line offsets
	_, err = loadVarSInt(buff)
	if err != nil {
		return nil, fmt.Errorf("error reading start line: %w", err)
	}

	_, err = loadLineOffsets(buff)
	if err != nil {
		return nil, fmt.Errorf("error reading line offsets: %w", err)
	}

	// load magic comments
	comments, err := loadComments(buff)
	if err != nil {
		return nil, fmt.Errorf("error reading comments: %w", err)
	}

	// load magic comments
	magicComments, err := loadMagicComments(buff)
	if err != nil {
		return nil, fmt.Errorf("error reading magic comments: %w", err)
	}

	// load optional locations
	dataLocation, err := loadOptionalLocation(buff)
	if err != nil {
		return nil, fmt.Errorf("error reading data location: %w", err)
	}

	// load syntax errors
	synErrors, err := loadSynErrors(buff)
	if err != nil {
		return nil, fmt.Errorf("error reading syntax errors: %w", err)
	}

	// load syntax warnings
	synWarnings, err := loadSynWarnings(buff)
	if err != nil {
		return nil, fmt.Errorf("error reading syntax warnings: %w", err)
	}

	// build constant pool
	constantPoolBufferOffset := buff.readUInt32()

	constantPoolLength, err := loadVarUInt(buff)
	if err != nil {
		return nil, fmt.Errorf("error reading constant pool length: %w", err)
	}

	constantPool := newConstantPool(
		source,
		serialized,
		constantPoolBufferOffset,
		constantPoolLength,
	)

	// load first node
	node, err := loadNode(buff, source, constantPool)
	if err != nil {
		return nil, fmt.Errorf("error reading first node: %w", err)
	}

	// build parse result
	return NewParseResult(
		node,
		comments,
		magicComments,
		dataLocation,
		synErrors,
		synWarnings,
	), nil
}
