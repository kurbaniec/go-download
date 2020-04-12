package postprocesser

import (
	"downloader/src/download"
	"os"
	"os/exec"
)

func Convert(fileInfo download.FileInfo) {
	cmd := exec.Command("ffmpeg", "-i", fileInfo.FileName(), "-c:a", "flac", fileInfo.Title+".flac")
	_, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	remErr := os.Remove(fileInfo.FileName())
	if remErr != nil {
		panic(remErr)
	}
}
