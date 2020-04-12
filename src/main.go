package main

import (
	"downloader/src/download"
	. "downloader/src/opts"
	"downloader/src/parser"
	"downloader/src/postprocesser"
	"fmt"
	"github.com/jessevdk/go-flags"
	"os"
)

func main() {
	//testUrl := "https://www.youtube.com/watch?v=THRDQmJSBs4"
	//testUrl := "https://www.youtube.com/watch?v=ADlGkXAz1D0"

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

	cipherStore := map[string]*parser.CipherOperations{}
	audioStreams := make([]parser.AudioStream, 0, 10)
	parser.GetStreams(args[1], cipherStore, &audioStreams)

	for _, stream := range audioStreams {
		if stream.Container == parser.Webm {
			fileInfo := download.FileLocation(stream)
			download.DownloadStream(stream, fileInfo)
			postprocesser.Convert(fileInfo)
			break
		}
	}
}
