package parser

import (
	"fmt"
	"math/big"
)

func loadVarUInt(buff *buffer) (uint32, error) {
	result := uint32(0)
	shift := 0

	for {
		byte, err := buff.readByte()
		if err != nil {
			return 0, fmt.Errorf("error reading byte: %w", err)
		}

		result += uint32(byte&0x7F) << shift
		shift += 7

		if (byte & 0x80) == 0 {
			break
		}
	}

	return result, nil
}

func loadVarSInt(buff *buffer) (int32, error) {
	x, err := loadVarUInt(buff)
	if err != nil {
		return 0, fmt.Errorf("error reading VarUInt: %w", err)
	}

	return int32((x >> 1) ^ (-(x & 1))), nil
}

func loadLineOffsets(buff *buffer) ([]uint32, error) {
	count, err := loadVarUInt(buff)
	if err != nil {
		return nil, fmt.Errorf("error reading line offsets count: %w", err)
	}

	lineOffsets := make([]uint32, count)
	for i := range count {
		lineOffsets[i], err = loadVarUInt(buff)
		if err != nil {
			return nil, fmt.Errorf("error reading line offset: %w", err)
		}
	}

	return lineOffsets, nil
}

func loadLocation(buff *buffer) (*Location, error) {
	startOffset, err := loadVarUInt(buff)
	if err != nil {
		return nil, fmt.Errorf("error reading location startOffset: %w", err)
	}

	length, err := loadVarUInt(buff)
	if err != nil {
		return nil, fmt.Errorf("error reading location length: %w", err)
	}

	return NewLocation(startOffset, length), nil
}

func loadComments(buff *buffer) ([]*Comment, error) {
	count, err := loadVarUInt(buff)
	if err != nil {
		return nil, fmt.Errorf("error reading comments count: %w", err)
	}

	comments := make([]*Comment, count)

	for i := range count {
		typpe, err := loadVarUInt(buff)
		if err != nil {
			return nil, fmt.Errorf("error reading comment type: %w", err)
		}

		loc, err := loadLocation(buff)
		if err != nil {
			return nil, fmt.Errorf("error reading comment value location: %w", err)
		}

		comment := NewComment(typpe, loc)
		comments[i] = comment
	}

	return comments, nil
}

func loadMagicComments(buff *buffer) ([]*MagicComment, error) {
	count, err := loadVarUInt(buff)
	if err != nil {
		return nil, fmt.Errorf("error reading magic comments count: %w", err)
	}

	comments := make([]*MagicComment, count)

	for i := range count {
		keyLocation, err := loadLocation(buff)
		if err != nil {
			return nil, fmt.Errorf("error reading magic comment key location: %w", err)
		}

		valueLocation, err := loadLocation(buff)
		if err != nil {
			return nil, fmt.Errorf("error reading magic comment value location: %w", err)
		}

		comment := NewMagicComment(keyLocation, valueLocation)
		comments[i] = comment
	}

	return comments, nil
}

func loadOptionalLocation(buff *buffer) (*Location, error) {
	nextByte, err := buff.readByte()
	if err != nil {
		return nil, fmt.Errorf("error reading optional location: %w", err)
	}

	if nextByte == 0 {
		return nil, nil
	}

	location, err := loadLocation(buff)
	if err != nil {
		return nil, fmt.Errorf("error reading location: %w", err)
	}

	return location, nil
}

func loadSynErrors(buff *buffer) ([]*SyntaxError, error) {
	count, err := loadVarUInt(buff)
	if err != nil {
		return nil, fmt.Errorf("error reading errors count: %w", err)
	}

	synErrs := make([]*SyntaxError, count)
	for i := range count {
		errorType, err := buff.readByte()
		if err != nil {
			return nil, fmt.Errorf("error reading error type: %w", err)
		}

		messageBytes, err := loadEmbeddedStr(buff)
		if err != nil {
			return nil, fmt.Errorf("error reading error message: %w", err)
		}

		location, err := loadLocation(buff)
		if err != nil {
			return nil, fmt.Errorf("error reading error location: %w", err)
		}

		errorLevel, err := buff.readByte()
		if err != nil {
			return nil, fmt.Errorf("error reading error level: %w", err)
		}

		synErrs[i] = NewSyntaxError(
			string(messageBytes),
			location,
			SyntaxErrorLevel(errorLevel),
			SyntaxErrorTypes[(errorType&0xFF)],
		)
	}

	return synErrs, nil
}

func loadSynWarnings(buff *buffer) ([]*SyntaxWarning, error) {
	count, err := loadVarUInt(buff)
	if err != nil {
		return nil, fmt.Errorf("error reading warnings count: %w", err)
	}

	synWarnings := make([]*SyntaxWarning, count)
	for i := range count {
		warningType, err := buff.readByte()
		if err != nil {
			return nil, fmt.Errorf("error reading warning type: %w", err)
		}

		messageBytes, err := loadEmbeddedStr(buff)
		if err != nil {
			return nil, fmt.Errorf("error reading warning message: %w", err)
		}

		location, err := loadLocation(buff)
		if err != nil {
			return nil, fmt.Errorf("error reading warning location: %w", err)
		}

		warningLevel, err := buff.readByte()
		if err != nil {
			return nil, fmt.Errorf("error reading warning level: %w", err)
		}

		synWarnings[i] = NewSyntaxWarning(
			string(messageBytes),
			location,
			SyntaxWarningLevel(warningLevel),
			SyntaxWarningTypes[(warningType&0xFF)-224],
		)
	}

	return synWarnings, nil
}

func loadEmbeddedStr(buff *buffer) ([]byte, error) {
	length, err := loadVarUInt(buff)
	if err != nil {
		return nil, fmt.Errorf("error reading embedded string length: %w", err)
	}

	strBytes := make([]byte, length)
	_, err = buff.read(strBytes)
	if err != nil {
		return nil, fmt.Errorf("error reading embedded string: %w", err)
	}

	return strBytes, nil
}

func loadStr(buff *buffer, src []byte) ([]byte, error) {
	b, err := buff.readByte()
	if err != nil {
		return nil, fmt.Errorf("error reading string type: %w", err)
	}

	switch b {
	case 1:
		start, err := loadVarUInt(buff)
		if err != nil {
			return nil, fmt.Errorf("error reading string start: %w", err)
		}

		length, err := loadVarUInt(buff)
		if err != nil {
			return nil, fmt.Errorf("error reading string length: %w", err)
		}

		bytes := make([]byte, length)
		copy(bytes, src[start:start+length])
		return bytes, nil
	case 2:
		strBytes, err := loadEmbeddedStr(buff)
		if err != nil {
			return nil, fmt.Errorf("error reading embedded string: %w", err)
		}

		return strBytes, nil
	default:
		return nil, fmt.Errorf("invalid string type: %d: %w", b, err)
	}
}

func loadOptionalNode(buff *buffer, src []byte, pool *constantPool) (Node, error) {
	nextByte, err := buff.readByte()
	if err != nil {
		return nil, fmt.Errorf("error reading optional node: %w", err)
	}

	if nextByte == 0 {
		return nil, nil
	}

	buff.setPosition(buff.position() - 1)
	node, err := loadNode(buff, src, pool)
	if err != nil {
		return nil, fmt.Errorf("error reading node: %w", err)
	}

	return node, nil
}

func loadInteger(buff *buffer) (*big.Int, error) {
	negative, err := buff.readByte()
	if err != nil {
		return nil, fmt.Errorf("error reading integer sign: %w", err)
	}

	length, err := loadVarUInt(buff)
	if err != nil {
		return nil, fmt.Errorf("error reading integer words length: %w", err)
	}

	if length == 0 {
		return nil, fmt.Errorf("invalid integer words length: %d: %w", length, err)
	}

	firstWord, err := loadVarUInt(buff)
	if err != nil {
		return nil, fmt.Errorf("error reading integer first word: %w", err)
	}

	if length == 1 {
		if negative != 0 {
			return big.NewInt(-int64(firstWord)), nil
		} else {
			return big.NewInt(int64(firstWord)), nil
		}
	}

	result := big.NewInt(int64(firstWord))
	for index := 1; index < int(length); index++ {
		word, err := loadVarUInt(buff)
		if err != nil {
			return nil, fmt.Errorf("error reading word: %w", err)
		}

		temp := big.NewInt(int64(word))
		temp = temp.Lsh(temp, uint(index*32))
		result = result.Or(result, temp)
	}

	if negative != 0 {
		result = result.Neg(result)
	} else {
		result = result.Abs(result)
	}

	return result, nil
}

func loadConstant(buff *buffer, pool *constantPool) (string, error) {
	idx, err := loadVarUInt(buff)
	if err != nil {
		return "", fmt.Errorf("error reading constant index: %w", err)
	}

	constant, err := pool.Get(buff, idx)
	if err != nil {
		return "", fmt.Errorf("error getting constant: %w", err)
	}

	return constant, nil
}

func loadOptionalConstant(buff *buffer, pool *constantPool) (*string, error) {
	nextByte, err := buff.readByte()
	if err != nil {
		return nil, fmt.Errorf("error reading optional node: %w", err)
	}

	if nextByte == 0 {
		return nil, nil
	}

	buff.setPosition(buff.position() - 1)
	constant, err := loadConstant(buff, pool)
	if err != nil {
		return nil, fmt.Errorf("error reading node: %w", err)
	}

	return &constant, nil
}

func loadConstants(buff *buffer, pool *constantPool) ([]string, error) {
	count, err := loadVarUInt(buff)
	if err != nil {
		return nil, fmt.Errorf("error reading constants count: %w", err)
	}

	constants := make([]string, count)
	for i := range count {
		constants[i], err = loadConstant(buff, pool)
		if err != nil {
			return nil, fmt.Errorf("error reading constant: %w", err)
		}
	}

	return constants, nil
}

func loadFlags(buff *buffer) (int16, error) {
	flags, err := loadVarUInt(buff)
	if err != nil {
		return 0, fmt.Errorf("error reading flags: %w", err)
	}

	if flags > 0x7FFF {
		return 0, fmt.Errorf("invalid flags: %d: %w", flags, err)
	}

	return int16(flags), nil
}
