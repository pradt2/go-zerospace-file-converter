package main

import (
	"os"
	"testing"
)

const filename = "TESTFILE1.randomAccessRead_test.bin"

var bytes = []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}

func createFile() *os.File {
	testFile, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	_, err = testFile.Write(bytes)
	if err != nil {
		panic(err)
	}
	if err = testFile.Sync(); err != nil {
		panic(err)
	}
	return testFile
}

func TestRead(t *testing.T) {
	file := createFile()
	for blockSize := 1; blockSize < 3; blockSize++ {
		out := ReadByRandomAccess(file, uint32(len(bytes)), uint64(blockSize))
		if len(bytes) != len(out) {
			t.Errorf("[Blocksize: %d] Output length is incorrect, expected %d , got %d", blockSize, len(bytes), len(out))
		}
		for i := 0; i < len(bytes); i++ {
			if bytes[i] == out[i] {
				continue
			}
			t.Errorf("[Blocksize: %d] Output content is incorrect at byte %d, expected %d , got %d", blockSize, i, bytes[i], out[i])
		}
	}
	file.Close()
	os.Remove(file.Name())
}
