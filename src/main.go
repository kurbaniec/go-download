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

	cipherStore := map[string]*parser.CipherOperations{}

	parser.GetConfig(testUrl, cipherStore)
}
