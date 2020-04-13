package postprocesser

import (
	. "downloader/src/opts"
	. "downloader/src/utils"
	"fmt"
	"os"
	"os/exec"
)

// Converts the downloaded file through command invocations for "ffmpeg".
func Convert(fileInfo FileInfo, opts *Opts) {
	fmt.Println("Converting...")
	downloadedFile := fileInfo.FileName()
	// Get selected audio format and codec
	var codec string
	if !opts.Manual {
		switch opts.Quality {
		case "high":
			codec = "flac"
			fileInfo.FileExtension = "flac"
		case "medium":
			codec = "libvorbis"
			fileInfo.FileExtension = "ogg"
		case "low":
			codec = "mp3"
			fileInfo.FileExtension = "mp3"
		}
	} else {
		codec = opts.Codec
		fileInfo.FileExtension = opts.AudioFormat
	}
	fileInfo.SearchFreePath()
	// Convert download with "ffmpeg"
	cmd := exec.Command("ffmpeg", "-i", downloadedFile, "-c:a", codec, fileInfo.FileName())
	_, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	// Remove downloaded file
	remErr := os.Remove(downloadedFile)
	if remErr != nil {
		panic(remErr)
	}
}
