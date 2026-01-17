package utility

import (
	"fmt"
	"github.com/NicoYUE/godub/converter"
	"github.com/bogem/id3v2/v2"
	"log"
	"monster-siren-record-puller/domain/model"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func SetID3Tags(filename string, metadata model.SongMetadata) error {
	tag, err := id3v2.Open(filename, id3v2.Options{Parse: false})
	if err != nil {
		log.Println("Error when creating Metadata Tags: ", err)
	}
	defer tag.Close()

	tag.SetDefaultEncoding(id3v2.EncodingUTF8)
	tag.SetTitle(metadata.Title)
	tag.SetAlbum(metadata.AlbumName)
	tag.SetArtist(strings.Join(metadata.Artists, ","))
	tag.AddAttachedPicture(metadata.PictureFrame)

	return tag.Save()
}

func PictureMime(filename string) (string, error) {
	extension := filepath.Ext(filename)

	switch extension {
	case ".jpg":
		return "image/jpeg", nil
	case ".png":
		return "image/png", nil
	}
	return "", fmt.Errorf("currently unhandled image type %s", extension)
}

// checkFFmpeg vérifie si ffmpeg est disponible dans le PATH
func checkFFmpeg() error {
	_, err := exec.LookPath("ffmpeg")
	if err != nil {
		return fmt.Errorf("ffmpeg n'est pas installé ou pas dans le PATH. Veuillez installer ffmpeg depuis https://ffmpeg.org/download.html")
	}
	return nil
}

func WavToMp3(filename string) string {
	// Vérifier si ffmpeg est disponible avant d'essayer de convertir
	if err := checkFFmpeg(); err != nil {
		log.Printf("ERREUR: %v", err)
		log.Printf("La conversion WAV vers MP3 nécessite ffmpeg.")
		log.Printf("Téléchargez ffmpeg depuis: https://ffmpeg.org/download.html")
		log.Printf("Ou utilisez un gestionnaire de paquets comme Chocolatey: choco install ffmpeg")
		WaitForExit()
		os.Exit(1)
	}

	mp3Filename := filename[0:len(filename)-4] + ".mp3"

	wavFile, err := os.Open(filename)
	if err != nil {
		log.Printf("Erreur lors de l'ouverture du fichier WAV '%s': %v", filename, err)
		WaitForExit()
		os.Exit(1)
	}
	defer wavFile.Close()

	// Create if not exist, write only, truncate content
	mp3File, err := os.OpenFile(mp3Filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		log.Printf("Erreur lors de la création du fichier MP3 '%s': %v", mp3Filename, err)
		WaitForExit()
		os.Exit(1)
	}
	defer mp3File.Close()

	err = converter.NewConverter(mp3File).WithBitRate(128000).WithDstFormat("mp3").Convert(wavFile)
	if err != nil {
		log.Printf("Erreur lors de la conversion WAV vers MP3: %v", err)
		log.Printf("Assurez-vous que ffmpeg est correctement installé et dans votre PATH.")
		WaitForExit()
		os.Exit(1)
	}
	
	// Supprimer le fichier WAV temporaire après conversion réussie
	os.Remove(filename)

	return mp3Filename
}
