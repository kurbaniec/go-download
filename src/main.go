package main

import (
	"downloader/src/parser"
	"fmt"
)

func main() {
	testUrl := "https://www.youtube.com/watch?v=ADlGkXAz1D0"
	metaUrl := parser.GetMetaUrl(testUrl)

	fmt.Println(metaUrl)
	parser.DownloadMetaData(metaUrl)
}
