package utility

import (
	"io"
	"log"
	"os"

	"github.com/bogem/id3v2/v2"
)

func PictureFrame(filename string) id3v2.PictureFrame {
	mime, _ := PictureMime(filename)
	file, err := os.Open(filename)
	if err != nil {
		log.Printf("Failed to open picture file '%s': %v", filename, err)
		WaitForExit()
		os.Exit(1)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Printf("Failed to read picture from file '%s': %v", filename, err)
		WaitForExit()
		os.Exit(1)
	}

	return id3v2.PictureFrame{
		Encoding:    id3v2.EncodingISO,
		MimeType:    mime,
		PictureType: id3v2.PTFrontCover,
		Description: "Cover",
		Picture:     data,
	}
}
