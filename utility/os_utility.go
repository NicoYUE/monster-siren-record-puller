package utility

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"
)

func Mkdir(path string, perm os.FileMode) error {

	err := os.Mkdir(path, perm)
	if err != nil {
		if os.IsExist(err) {
			return fmt.Errorf("directory %s already exists", path)
		}
		return fmt.Errorf("failed to create directory: %v", err)
	}
	return nil
}

func WriteResponse2File(filename string, destDirectory string, response *http.Response) string {

	contentType := response.Header.Get("content-type")
	outputPath := destDirectory + filename

	switch contentType {
	case "audio/mpeg":
		outputPath = outputPath + ".mp3"
	case "audio/wav":
		outputPath = outputPath + ".wav"
	case "image/jpeg":
		outputPath = outputPath + ".jpg"
	case "image/png":
		outputPath = outputPath + ".png"
	case "application/octet-stream":
		outputPath = outputPath + ".lrc"
	}

	file, err := os.Create(outputPath)
	if err != nil {
		print(err)
	}

	_, err = io.Copy(file, response.Body)
	if err != nil {
		fmt.Printf("ERROR Copying")
	}

	file.Close()
	response.Body.Close()

	// Special case for wav files (which needs to be created beforehand)
	if contentType == "audio/wav" {
		return WavToMp3(outputPath)
	}

	return outputPath
}

// WaitForExit attend une entrée utilisateur pour garder le terminal ouvert
// Fonctionne même quand stdin n'est pas disponible (double-clic sur .exe Windows)
func WaitForExit() {
	fmt.Println("\nAppuyez sur Entrée pour quitter...")
	
	// Vérifier si stdin est disponible
	stat, err := os.Stdin.Stat()
	if err != nil || (stat.Mode()&os.ModeCharDevice) == 0 {
		// stdin n'est pas disponible (double-clic sur .exe)
		if runtime.GOOS == "windows" {
			// Utiliser la commande pause de Windows
			cmd := exec.Command("cmd", "/c", "pause")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Run()
		} else {
			// Sur Unix, attendre quelques secondes
			time.Sleep(5 * time.Second)
		}
	} else {
		// stdin est disponible, lire une entrée
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}
}
