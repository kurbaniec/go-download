package parser

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
)

func ReadMetaData() {
	content, err := ioutil.ReadFile("meta1.txt")
	if err != nil {
		// TODO make ONE error exit function
		fmt.Println("Something went wrong reading meta-data")
		os.Exit(0)
	}
	metaEncoded := string(content)
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
		streamUrl, err := url.QueryUnescape(set["cipher"].(string))
		if err != nil {
			fmt.Println("STREAMURL AHHHHH!")
			os.Exit(0)
		}
		fmt.Println(streamUrl)
		fmt.Println(entry)
		//fmt.Println(entry["url"])
		fmt.Println("---")
	}

	//var raw map[string]interface{}
	//json.Unmarshal(meta)
}

func getVideoId(videoUrl string) string {
	return videoUrl[strings.Index(videoUrl, "watch?v=")+8:]
}

func GetVideoInfo(videoUrl string) {
	videoId := getVideoId(videoUrl)
	eurl := url.QueryEscape("https://youtube.googleapis.com/v/" + videoId)
	metaUrl := "https://youtube.com/get_video_info?video_id=" + videoId +
		"&el=embedded&eurl=" + eurl + "&hl=en"
	download("meta1.txt", metaUrl)
}

func GetVideoEmbedPage(videoUrl string) {
	videoId := getVideoId(videoUrl)
	metaUrl := "https://youtube.com/embed/" + videoId + "?disable_polymer=true&hl=en"
	download("meta2.txt", metaUrl)
}

func GetVideoWatchPage(videoUrl string) {
	videoId := getVideoId(videoUrl)
	metaUrl := "https://youtube.com/watch?v=" + videoId + "&disable_polymer=true&bpctr=9999999999&hl=en"
	download("meta3.txt", metaUrl)
}

func download(fileName string, metaUrl string) {
	err := downloadFile(fileName, metaUrl)
	if err != nil {
		fmt.Println("Something went wrong, cant get meta-data")
		os.Exit(0)
	}
}
