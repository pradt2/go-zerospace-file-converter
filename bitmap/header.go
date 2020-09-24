package bitmap

import (
	"encoding/binary"
	"io"
	"pawel.com/go-bitmaps/fsutils"
)

type commonHeader struct {
	MagicBytes     [2]byte
	FileSize       uint32
	Unused1        uint16
	Unused2        uint16
	PixelArrOffset uint32
}

var commonHeaderSize = uint32(binary.Size(commonHeader{}))

func newCommonHeaderFromFile(file io.ReadSeeker) commonHeader {
	bitmapHeader := &commonHeader{}
	fsutils.ReadIntoStructureAtOffset(file, 0, bitmapHeader)
	return *bitmapHeader
}

func writeCommonHeader(header commonHeader, file io.Writer) {
	fsutils.WriteStructure(file, header)
}

func readCommonHeader(header commonHeader, offset int, length int) []byte {
	b := fsutils.ReadStructureBytes(uint64(offset), uint64(length), header)
	return b
}

type dibBGRHeader struct {
	HeaderSize                       uint32
	Width                            uint32
	Height                           uint32
	ColorPlanes                      uint16
	BitsPerPixel                     uint16
	PixelCompression                 uint32
	BitmapSize                       uint32
	PrintResPixelsPerMeterHorizontal uint32
	PrintResPixelsPerMeterVertical   uint32
	ColourPaletteSize                uint32
	ImportantColoursSize             uint32
}

var dibBGRHeaderSize = uint32(binary.Size(dibBGRHeader{}))

func newDibBGRHeaderFromFile(file io.ReadSeeker) dibBGRHeader {
	dibHeader := &dibBGRHeader{}
	fsutils.ReadIntoStructureAtOffset(file, uint64(commonHeaderSize), dibHeader)
	return *dibHeader
}

func writeDibBGRHeader(header dibBGRHeader, file io.Writer) {
	fsutils.WriteStructure(file, header)
}

func readDibHeader(header dibBGRHeader, offset, length uint32) []byte {
	b := fsutils.ReadStructureBytes(uint64(offset), uint64(length), header)
	return b
}

type DibBGRBitmapHeader struct {
	commonHeader commonHeader
	dibHeader    dibBGRHeader
}

var dibBGRBitmapHeaderSize = commonHeaderSize + dibBGRHeaderSize

func (h *DibBGRBitmapHeader) GetRowByteSize() uint32 {
	size := getRowByteSize(h.dibHeader.Width, uint32(h.dibHeader.BitsPerPixel))
	return size
}

func (h *DibBGRBitmapHeader) ReadAt(offset, length uint64) []byte {
	end := offset + length
	if uint32(end) <= commonHeaderSize {
		b := readCommonHeader(h.commonHeader, int(offset), int(length))
		return b
	}
	if uint32(offset) >= commonHeaderSize {
		b := readDibHeader(h.dibHeader, uint32(offset), uint32(length))
		return b
	}
	b1 := readCommonHeader(h.commonHeader, int(offset), int(length))
	b2 := readDibHeader(h.dibHeader, 0, uint32(length)-uint32(len(b1)))
	return append(b1, b2...)
}

func newBitmapDibBGRHeaderFromFile(file io.ReadSeeker) DibBGRBitmapHeader {
	return DibBGRBitmapHeader{
		commonHeader: newCommonHeaderFromFile(file),
		dibHeader:    newDibBGRHeaderFromFile(file),
	}
}

func newBitmapDibBGRHeader(width, height uint32) DibBGRBitmapHeader {
	pixelSizeBits := 24
	bitmapSize := getRowByteSize(width, 24) * height
	h := DibBGRBitmapHeader{
		commonHeader: commonHeader{
			MagicBytes:     [2]byte{'B', 'M'},
			FileSize:       commonHeaderSize + dibBGRHeaderSize + bitmapSize,
			Unused1:        0,
			Unused2:        0,
			PixelArrOffset: commonHeaderSize + dibBGRHeaderSize,
		},
		dibHeader: dibBGRHeader{
			HeaderSize:                       dibBGRHeaderSize,
			Width:                            width,
			Height:                           height,
			ColorPlanes:                      1,
			BitsPerPixel:                     uint16(pixelSizeBits),
			PixelCompression:                 0,
			BitmapSize:                       bitmapSize,
			PrintResPixelsPerMeterHorizontal: 2835,
			PrintResPixelsPerMeterVertical:   2835,
			ColourPaletteSize:                0,
			ImportantColoursSize:             0,
		},
	}
	return h
}
