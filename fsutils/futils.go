package fsutils

import (
	"bytes"
	"encoding/binary"
	"io"
)

type ReaderAt interface {
	ReadAt(offset, length uint64) []byte
}

func ReadFile(file io.ReadSeeker, offset uint64, length uint64) []byte {
	if ret, err := file.Seek(int64(offset), 0); err != nil || uint64(ret) != offset {
		panic(err)
	}
	b := make([]byte, length)
	n, err := file.Read(b)
	if err != nil || n == -1 || uint64(n) < length {
		panic(err)
	}
	return b
}

func ReadIntoStructureAtOffset(file io.ReadSeeker, offset uint64, data interface{}) {
	structureSize := binary.Size(data)
	if structureSize == -1 {
		panic("Structure size is negative")
	}
	b := ReadFile(file, offset, uint64(structureSize))
	ReadIntoStructure(bytes.NewReader(b), data)
}

func ReadIntoStructure(file io.Reader, data interface{}) {
	if err := binary.Read(file, binary.LittleEndian, data); err != nil {
		panic(err)
	}
}

func ReadStructureBytes(offset, length uint64, data interface{}) []byte {
	structureSizeInt := binary.Size(data)
	if structureSizeInt == -1 {
		panic("Structure size is negative")
	}
	structureSize := uint64(structureSizeInt)
	if offset >= structureSize {
		panic("Offset larger than header size")
	}
	b := make([]byte, 0, structureSize)
	WriteStructure(bytes.NewBuffer(b), data)
	var end uint64 = 0
	if structureSize <= offset+length {
		end = structureSize
	} else {
		end = offset + length
	}
	b = b[offset:end]
	return b
}

func WriteStructure(file io.Writer, data interface{}) {
	if err := binary.Write(file, binary.LittleEndian, data); err != nil {
		panic(err)
	}
}
