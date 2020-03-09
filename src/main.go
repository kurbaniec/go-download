package main

import (
	"downloader/src/download"
	"downloader/src/parser"
)

func main() {
	//testUrl := "https://www.youtube.com/watch?v=THRDQmJSBs4"
	testUrl := "https://www.youtube.com/watch?v=ADlGkXAz1D0"

	cipherStore := map[string]*parser.CipherOperations{}
	audioStreams := make([]parser.AudioStream, 0, 10)

	parser.GetStreams(testUrl, cipherStore, &audioStreams)

	for _, stream := range audioStreams {
		if stream.Container == parser.Webm {
			download.DownloadStream(stream)
			break
		}
	}
}
