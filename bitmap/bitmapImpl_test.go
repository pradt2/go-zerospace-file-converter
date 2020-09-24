package bitmap_test

import (
	"go-bitmaps/bitmap"
	"testing"
)

func TestPerf(t *testing.T) {
	bmap := bitmap.NewBGRWBitmap(1)
	_ = bmap.ReadAt(0, uint64(bmap.GetFileSize()))
}
