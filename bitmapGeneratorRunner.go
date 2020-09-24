package main

import (
	"os"
	"pawel.com/go-bitmaps/bitmap"
)

func main() {
	const size = 200
	f, err := os.Open("coronavirus.bmp")
	if err != nil {
		panic(err)
	}
	bmap := bitmap.NewBitmapFromFile(f)
	bmap = bitmap.NewNegativeBitmap(bmap)
	b := bmap.ReadAt(0, uint64(bmap.GetFileSize()))
	f, err = os.OpenFile("test.bmp", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	if _, err := f.Write(b); err != nil {
		panic(err)
	}
	if err := f.Close(); err != nil {
		panic(err)
	}
}
