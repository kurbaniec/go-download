package main

import (
	"downloader/src/download"
	"downloader/src/parser"
)

func main() {
	testUrl := "https://www.youtube.com/watch?v=THRDQmJSBs4"
	//metaUrl := parser.GetMetaUrl(testUrl)
	//metaUrl := parser.Testo(testUrl)

	//fmt.Println(metaUrl)
	//parser.DownloadMetaData(metaUrl)
	//parser.ReadMetaData()

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
