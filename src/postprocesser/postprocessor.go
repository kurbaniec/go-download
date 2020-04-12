package postprocesser

import (
	. "downloader/src/opts"
	. "downloader/src/utils"
	"os"
	"os/exec"
)

func Convert(fileInfo FileInfo, opts *Opts) {
	downloadedFile := fileInfo.FileName()
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
	cmd := exec.Command("ffmpeg", "-i", downloadedFile, "-c:a", codec, fileInfo.FileName())
	_, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	remErr := os.Remove(downloadedFile)
	if remErr != nil {
		panic(remErr)
	}
}
