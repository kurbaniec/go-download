package main

import (
	"downloader/src/parser"
)

func main() {
	testUrl := "https://www.youtube.com/watch?v=ADlGkXAz1D0"
	//metaUrl := parser.GetMetaUrl(testUrl)
	//metaUrl := parser.Testo(testUrl)

	//fmt.Println(metaUrl)
	//parser.DownloadMetaData(metaUrl)
	//parser.ReadMetaData()

	parser.GetVideoEmbedPage(testUrl)
	parser.GetVideoInfo(testUrl)
	parser.GetVideoWatchPage(testUrl)
	parser.ReadMetaData()
}
