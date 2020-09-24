package main

import (
	"io"
	"math"
	"math/rand"
)

const seed = 0

func ReadByRandomAccess(file io.ReadSeeker, fileSize uint32, blockSize uint64) []byte {
	blocksCount := uint64(math.Ceil(float64(fileSize) / float64(blockSize)))
	blocks := make([][]byte, blocksCount, blocksCount)
	blockIndexes := make([]uint64, blocksCount)
	var i uint64
	for i = 0; i < blocksCount; i++ {
		blockIndexes[i] = i
	}
	rand.Seed(seed)
	rand.Shuffle(int(blocksCount), func(i, j int) { blockIndexes[i], blockIndexes[j] = blockIndexes[j], blockIndexes[i] })
	for i = 0; i < blocksCount; i++ {
		blockIndex := blockIndexes[i]
		offset := blockIndex * blockSize
		retOffset, err := file.Seek(int64(offset), 0)
		if err != nil {
			panic(err)
		} else if uint64(retOffset) != offset {
			panic("Returned offset does not equal the required offset")
		}
		readBytes := make([]byte, blockSize)
		bytesRead, err := file.Read(readBytes)
		if err != nil {
			panic(err)
		} else if uint64(bytesRead) != blockSize && blockIndex != blocksCount-1 {
			panic("Read less bytes than block size in a non-final block")
		}
		blocks[blockIndex] = readBytes[:bytesRead]
	}
	finalArr := make([]byte, 0, fileSize)
	for _, block := range blocks {
		finalArr = append(finalArr, block...)
	}
	return finalArr
}
