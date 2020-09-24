package bitmap

import (
	"io"
	"pawel.com/go-bitmaps/fsutils"
)

type Bitmap interface {
	GetWidth() uint32
	GetHeight() uint32
	GetFileSize() uint32
	GetPixel(x, y uint32) BGRPixel
	fsutils.ReaderAt
}

type dibBGRbitmapFile struct {
	bitmapHeader DibBGRBitmapHeader
	pixelArr     PixelArray
}

func (f *dibBGRbitmapFile) GetWidth() uint32 {
	return f.bitmapHeader.dibHeader.Width
}

func (f *dibBGRbitmapFile) GetHeight() uint32 {
	return f.bitmapHeader.dibHeader.Height
}

func (f *dibBGRbitmapFile) GetFileSize() uint32 {
	return f.bitmapHeader.commonHeader.FileSize
}

func (f *dibBGRbitmapFile) GetPixel(x, y uint32) BGRPixel {
	return f.pixelArr.GetPixel(x, y)
}

func (f *dibBGRbitmapFile) ReadAt(offset, length uint64) []byte {
	if uint32(offset+length) <= dibBGRHeaderSize {
		b := f.bitmapHeader.ReadAt(offset, length)
		return b
	}
	if offset > uint64(dibBGRHeaderSize) && uint32(offset+length) > dibBGRHeaderSize {
		offset -= uint64(dibBGRHeaderSize)
		b := f.pixelArr.ReadAt(offset, length)
		return b
	}
	b := f.bitmapHeader.ReadAt(offset, length)
	length -= uint64(len(b))
	b = append(b, f.pixelArr.ReadAt(0, length)...)
	return b
}

func newDibBGRBitmap(header DibBGRBitmapHeader, array PixelArray) Bitmap {
	return &dibBGRbitmapFile{
		bitmapHeader: header,
		pixelArr:     array,
	}
}

func NewBitmapFromFile(file io.ReadSeeker) Bitmap {
	bitmapHeader := newBitmapDibBGRHeaderFromFile(file)
	pixelArr := newBGRPixelArrFromFile(bitmapHeader, file)
	return newDibBGRBitmap(bitmapHeader, pixelArr)
}
