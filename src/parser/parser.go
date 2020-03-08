package parser

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"
)

func GetConfig(videoUrl string) {
	infoFile := GetVideoInfo(videoUrl)
	embedFile := GetVideoEmbedPage(videoUrl)
	//watchFile := GetVideoWatchPage(videoUrl)
	getCipherSrc(embedFile)

	metaEncoded := infoFile
	//fmt.Println(metaEncoded)
	metaDecoded, err := url.QueryUnescape(metaEncoded)
	if err != nil {
		fmt.Println("Something went wrong reading meta-data")
		os.Exit(0)
	}
	//fmt.Println(metaDecoded)
	metaArray := strings.Split(metaDecoded, "&")
	var info string
	for _, entry := range metaArray {
		//fmt.Println(entry)
		if strings.HasPrefix(entry, "player_response") {
			info = entry[16:]
		}
	}

	fmt.Println(info)

	var dat map[string]interface{}
	if err := json.Unmarshal([]byte(info), &dat); err != nil {
		fmt.Println("MARSHALL ERROR")
	}
	fmt.Println(dat["playabilityStatus"])
	fmt.Println(dat["streamingData"].(map[string]interface{})["formats"])
	formats := dat["streamingData"].(map[string]interface{})["adaptiveFormats"]
	fmt.Println(formats)
	for _, entry := range formats.([]interface{}) {
		set := entry.(map[string]interface{})
		fmt.Println(set["cipher"])
		/**
		streamUrl, err := url.QueryUnescape(set["cipher"].(string))
		if err != nil {
			fmt.Println("STREAMURL AHHHHH!")
			os.Exit(0)
		}
		fmt.Println(streamUrl)
		fmt.Println(entry)
		fmt.Println(entry["url"])
		fmt.Println("---")*/
	}

	//var raw map[string]interface{}
	//json.Unmarshal(meta)
}

func getCipherSrc(embedFile string) {
	search := "\"assets\":{\"js\":\""
	searchLen := len(search)
	index := strings.Index(embedFile, search)
	assetBegin := embedFile[index+searchLen:]
	assetEncoded := assetBegin[:strings.Index(assetBegin, "\"")]
	assetDecoded, _ := url.QueryUnescape(assetEncoded)
	asset := "https://youtube.com" + assetDecoded
	fmt.Println(asset)
}

func getVideoId(videoUrl string) string {
	return videoUrl[strings.Index(videoUrl, "watch?v=")+8:]
}

func GetVideoInfo(videoUrl string) string {
	videoId := getVideoId(videoUrl)
	eurl := url.QueryEscape("https://youtube.googleapis.com/v/" + videoId)
	metaUrl := "https://youtube.com/get_video_info?video_id=" + videoId +
		"&el=embedded&eurl=" + eurl + "&hl=en"
	data, err := downloadAsString(metaUrl)
	errorHandler(err)
	return data
}

func GetVideoEmbedPage(videoUrl string) string {
	videoId := getVideoId(videoUrl)
	metaUrl := "https://youtube.com/embed/" + videoId + "?disable_polymer=true&hl=en"
	data, err := downloadAsString(metaUrl)
	errorHandler(err)
	return data
}

func GetVideoWatchPage(videoUrl string) string {
	videoId := getVideoId(videoUrl)
	metaUrl := "https://youtube.com/watch?v=" + videoId + "&disable_polymer=true&bpctr=9999999999&hl=en"
	data, err := downloadAsString(metaUrl)
	errorHandler(err)
	return data
}

func download(fileName string, metaUrl string) {
	err := downloadFile(fileName, metaUrl)
	if err != nil {
		fmt.Println("Something went wrong, cant get meta-data")
		os.Exit(0)
	}
}
