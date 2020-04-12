package main

import (
	"downloader/src/download"
	. "downloader/src/opts"
	"downloader/src/parser"
	"downloader/src/postprocesser"
	"downloader/src/utils"
	"fmt"
	"github.com/jessevdk/go-flags"
	"os"
)

func main() {
	// Parse arguments
	var opts Opts
	args, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		panic(err)
	}
	if len(args) == 1 {
		fmt.Println("URL parameter is missing. Please try again with it.")
		os.Exit(0)
	}
	// Parse meta data to find song url and tags
	cipherStore := map[string]*parser.CipherOperations{}
	audioStream := parser.GetAudioStream(args[1], cipherStore, &opts)
	// Create file and download stream to it
	fileInfo := utils.FileLocation(audioStream)
	download.DownloadStream(audioStream, fileInfo)
	// Postprocess file through command invocations for "ffmpeg
	postprocesser.Convert(fileInfo, &opts)
}
