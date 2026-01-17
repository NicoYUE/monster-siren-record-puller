package main

import (
	"fmt"
	"monster-siren-record-puller/cache"
	"monster-siren-record-puller/domain/model"
	"monster-siren-record-puller/domain/service"
	"monster-siren-record-puller/infra/ms/repo"
	"monster-siren-record-puller/utility"
	"net/http"
	"os"
	"sync"
)

const BaseDirectory = "./MSR/"
const DefaultConcurrency = 5

var client = http.Client{}
var msrService = service.NewMonsterSirenService(client)
var mediaRepository = repo.NewMonsterSirenMediaRepository(client)
var albumCache = cache.NewAlbumCache()

func main() {

	fmt.Printf("Initiating MSR downloading \n")

	err := utility.Mkdir(BaseDirectory, 0755)
	if err != nil {
		fmt.Println("Warn:", err)
	} else {
		fmt.Println("Created Directory: ", BaseDirectory)
	}

	albums := msrService.RetrieveAlbums()
	AsyncAlbumDownload(albums)

	// Garder le terminal ouvert pour voir les erreurs
	utility.WaitForExit()
}

func AsyncAlbumDownload(albums []model.Album) {
	var group sync.WaitGroup
	channel := make(chan struct{}, DefaultConcurrency)

	for _, album := range albums {
		group.Add(1)
		channel <- struct{}{}

		// intermediate var
		album := album
		go func() {
			defer func() { <-channel }()
			defer group.Done()
			DownloadAlbumSongs(album)
		}()
	}
	group.Wait()
}

func DownloadAlbumSongs(album model.Album) {
	fmt.Printf("Album: %s\n", album.Name)

	songs := msrService.RetrieveAlbumSongs(album)

	albumPath := fmt.Sprintf(BaseDirectory+"%s/", album.Name)

	if !albumCache.AlbumExists(album.Name) {
		utility.Mkdir(albumPath, 0755)
	}

	coverResp, _ := mediaRepository.RetrieveImage(album.CoverUrl)
	coverPath := utility.WriteResponse2File(album.Name, albumPath, coverResp)
	pictureFrame := utility.PictureFrame(coverPath)
	os.Remove(coverPath)

	for _, song := range songs {
		if song.SongId == "048709" {
			fmt.Printf("Special case - Skipping %s", song.Name)
			continue
		}

		fmt.Printf("Downloading : %s ---- Album: %s\n", song.Name, album.Name)

		if !albumCache.SongExists(album.Name, song.Name) {
			audioResp, _ := mediaRepository.RetrieveAudio(song.SourceUrl)
			audioPath := utility.WriteResponse2File(song.Name, albumPath, audioResp)

			//if song.LyricUrl != "" {
			//	lyricResp, _ := mediaRepository.RetrieveLyric(song.LyricUrl)
			//	utility.WriteResponse2File(song.Name, albumPath, lyricResp)
			//}

			err := utility.SetID3Tags(audioPath, model.SongMetadata{
				Title:        song.Name,
				AlbumName:    album.Name,
				Artists:      song.Artists,
				AlbumArtists: album.Artists,
				PictureFrame: pictureFrame,
			})

			if err != nil {
				fmt.Printf("error! %s when trying to set ID tags for Song: %s", err, song.Name)
				return
			}

			fmt.Printf("---- Finished : %s\n", song.Name)
		} else {
			fmt.Printf("Song was already download, skipping \n")
		}
		albumCache.Cache(album.Name, song.Name)
	}
}
