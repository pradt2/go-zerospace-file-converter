package bitmap

import (
	"bytes"
	"encoding/binary"
	"io"
	"math"
	"pawel.com/go-bitmaps/fsutils"
)

type BGRPixel struct {
	Blue  byte
	Green byte
	Red   byte
}

func newBGRPixelFromBytes(b []byte) BGRPixel {
	pixel := &BGRPixel{}
	fsutils.ReadIntoStructure(bytes.NewBuffer(b), pixel)
	return *pixel
}

var bgrPixelSize = uint32(binary.Size(BGRPixel{}))

type PixelArray interface {
	GetRowByteSize() uint32
	GetByteSize() uint32
	GetPixel(x, y uint32) BGRPixel
	fsutils.ReaderAt
}

type fileBackedBGRPixelArr struct {
	width  uint32
	height uint32
	offset uint32
	file   io.ReadSeeker
}

func (f *fileBackedBGRPixelArr) GetRowByteSize() uint32 {
	size := uint32(getRowByteSize(f.width, bgrPixelSize*8))
	return size
}

func (f *fileBackedBGRPixelArr) GetByteSize() uint32 {
	size := uint32(f.GetRowByteSize() * f.height)
	return size
}

func (f *fileBackedBGRPixelArr) GetPixel(x, y uint32) BGRPixel {
	y = f.height - y - 1
	offset := f.offset + y*f.GetRowByteSize() + x*bgrPixelSize
	bgrPixel := &BGRPixel{}
	fsutils.ReadIntoStructureAtOffset(f.file, uint64(offset), bgrPixel)
	return *bgrPixel
}

func (f *fileBackedBGRPixelArr) ReadAt(offset, length uint64) []byte {
	b := fsutils.ReadFile(f.file, uint64(f.offset)+offset, length)
	return b
}

func newBGRPixelArrFromFile(bitmapHeader DibBGRBitmapHeader, file io.ReadSeeker) PixelArray {
	return &fileBackedBGRPixelArr{
		width:  bitmapHeader.dibHeader.Width,
		height: bitmapHeader.dibHeader.Height,
		offset: bitmapHeader.commonHeader.PixelArrOffset,
		file:   file,
	}
}

func getRowByteSize(pixelRowSize, pixelSizeBits uint32) uint32 {
	pixelByteSize := pixelSizeBits / 8
	pixelRowByteSize := uint32(math.Ceil(float64(pixelRowSize*pixelByteSize)/float64(4)) * 4)
	return pixelRowByteSize
}
