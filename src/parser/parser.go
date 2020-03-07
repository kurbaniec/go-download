package parser

import (
	"fmt"
	"os"
	"strings"
)

func GetMetaUrl(videoUrl string) string {
	index := strings.Index(videoUrl, "watch?v=") + 8
	return "https://www.youtube.com/get_video_info?video_id=" + videoUrl[index:]
}

func DownloadMetaData(metaUrl string) {
	err := downloadFile("meta.txt", metaUrl)
	if err != nil {
		fmt.Println("Something went wrong")
		os.Exit(0)
	}
}
