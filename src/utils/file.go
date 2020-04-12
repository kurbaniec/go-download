package utils

import (
	"downloader/src/parser"
	"os"
	"strconv"
)

type FileInfo struct {
	Title         string
	FileExtension string
}

func (fileInfo *FileInfo) FileName() string {
	return fileInfo.Title + "." + fileInfo.FileExtension
}

func (fileInfo *FileInfo) SearchFreePath() {
	fileExists := true
	songIndex := 1
	for fileExists {
		_, err := os.Stat(fileInfo.FileName())
		if os.IsNotExist(err) {
			fileExists = false
		} else {
			fileInfo.Title = fileInfo.Title + "(" + strconv.Itoa(songIndex) + ")"
		}
	}
}

// Creates and returns the location of the file
// which is used for the download
func FileLocation(stream parser.AudioStream) FileInfo {
	// Create and open output file
	fileInfo := FileInfo{
		Title:         stream.Title,
		FileExtension: string(stream.Container),
	}
	fileInfo.SearchFreePath()
	_, err := os.Create(fileInfo.FileName())
	if err != nil {
		panic(err)
	}
	return fileInfo
}
