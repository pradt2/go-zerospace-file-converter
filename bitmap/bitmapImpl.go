package bitmap

import (
	"bytes"
	"math"
	"pawel.com/go-bitmaps/fsutils"
)

type pixelGetterBitmap struct {
	width       uint32
	height      uint32
	pixelGetter func(x, y uint32) BGRPixel
}

func (bitmap *pixelGetterBitmap) GetRowByteSize() uint32 {
	return getRowByteSize(bitmap.width, 24)
}

func (bitmap *pixelGetterBitmap) GetByteSize() uint32 {
	return bitmap.GetRowByteSize() * bitmap.height
}

func (bitmap *pixelGetterBitmap) GetPixel(x, y uint32) BGRPixel {
	return bitmap.pixelGetter(x, y)
}

func (bitmap *pixelGetterBitmap) ReadAt(offset, length uint64) []byte {
	byteSize := uint64(bitmap.GetByteSize())
	maxByteIndexPlusOne := offset + length
	if byteSize < maxByteIndexPlusOne {
		maxByteIndexPlusOne = byteSize
	}
	rowSize := bitmap.GetRowByteSize()
	widthPixels := bitmap.width
	widthPixelBytes := uint64(widthPixels * bgrPixelSize)
	b := make([]byte, 0, maxByteIndexPlusOne-offset)
	for byteIndex := offset; byteIndex < maxByteIndexPlusOne; byteIndex++ {
		offset = byteIndex
		rowIndex := uint32(math.Floor(float64(byteIndex) / float64(rowSize)))
		offset -= uint64(rowIndex * rowSize)
		if offset >= widthPixelBytes {
			b = append(b, 0)
			continue
		}
		colIndex := uint32(math.Floor(float64(offset) / float64(bgrPixelSize)))
		offset -= uint64(colIndex * bgrPixelSize)
		rowIndex = bitmap.height - 1 - rowIndex

		bgrPixel := bitmap.GetPixel(colIndex, rowIndex)
		bgrPixelBytes := make([]byte, 0, bgrPixelSize)
		fsutils.WriteStructure(bytes.NewBuffer(bgrPixelBytes), bgrPixel)

		b = append(b, bgrPixelBytes[offset:offset+1]...)
	}
	return b
}

func NewPixelGetterBitmap(width, height uint32, pixelGetter func(x, y uint32) BGRPixel) Bitmap {
	header := newBitmapDibBGRHeader(width, height)
	array := &pixelGetterBitmap{
		width:       width,
		height:      height,
		pixelGetter: pixelGetter,
	}
	return newDibBGRBitmap(header, array)
}

func NewPixelMappingBitmap(bitmap Bitmap, mapper func(x, y uint32, pixel BGRPixel) BGRPixel) Bitmap {
	return NewPixelGetterBitmap(bitmap.GetWidth(), bitmap.GetHeight(), func(x, y uint32) BGRPixel {
		return mapper(x, y, bitmap.GetPixel(x, y))
	})
}

func NewIdentityBitmap(bitmap Bitmap) Bitmap {
	return NewPixelMappingBitmap(bitmap, func(x, y uint32, pixel BGRPixel) BGRPixel {
		return pixel
	})
}

func NewGrayscaleBitmap(bitmap Bitmap) Bitmap {
	return NewPixelMappingBitmap(bitmap, func(x, y uint32, pixel BGRPixel) BGRPixel {
		brightness := uint8(0.3*float32(pixel.Red) + 0.59*float32(pixel.Green) + 0.11*float32(pixel.Blue))
		pixel.Blue = brightness
		pixel.Green = brightness
		pixel.Red = brightness
		return pixel
	})
}

func NewNegativeBitmap(bitmap Bitmap) Bitmap {
	return NewPixelMappingBitmap(bitmap, func(x, y uint32, pixel BGRPixel) BGRPixel {
		pixel.Blue = 255 - pixel.Blue
		pixel.Green = 255 - pixel.Green
		pixel.Red = 255 - pixel.Red
		return pixel
	})
}

func NewBGRWBitmap(tileSize uint32) Bitmap {
	return NewPixelGetterBitmap(tileSize*2, tileSize*2, func(x, y uint32) BGRPixel {
		if x < tileSize {
			if y < tileSize {
				//blue
				return newBGRPixelFromBytes([]byte{255, 0, 0})
			} else {
				//red
				return newBGRPixelFromBytes([]byte{0, 0, 255})
			}
		} else {
			if y < tileSize {
				//green
				return newBGRPixelFromBytes([]byte{0, 255, 0})
			} else {
				//white
				return newBGRPixelFromBytes([]byte{255, 255, 255})
			}
		}
	})
}
