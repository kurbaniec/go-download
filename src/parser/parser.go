package parser

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

func GetStreams(
	videoUrl string,
	cipherStore map[string]*CipherOperations,
	audioStreams []AudioStream,
) {
	var infoFile, embedFile, assetUrl string
	var wg sync.WaitGroup
	wg.Add(2)
	go GetVideoInfo(videoUrl, &infoFile, &wg)
	go GetVideoEmbedPage(videoUrl, &embedFile, &wg)
	wg.Wait()
	assetUrl = getAssetUrl(embedFile)

	if _, ok := cipherStore[assetUrl]; !ok {
		assetFile := getAssetFile(assetUrl)
		cipherStore[assetUrl] = getCipherSrc(assetFile)
	}

	cipher := cipherStore[assetUrl]

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
		fmt.Println("MARSHALL ERROr")
	}
	//fmt.Println(dat["playabilityStatus"])
	//fmt.Println(dat["streamingData"].(map[string]interface{})["formats"])
	formats := dat["streamingData"].(map[string]interface{})["adaptiveFormats"]
	//fmt.Println(formats)

	for _, entry := range formats.([]interface{}) {
		stream := entry.(map[string]interface{})
		mime := stream["mimeType"].(string)
		if strings.HasPrefix(mime, "audio") {
			// fmt.Println(stream["qualityLabel"]) // not available in audio
			itag := int(stream["itag"].(float64))
			audioEncoding, container := getAudioEncodingAndContainer(mime)
			url := buildStreamUrl(stream["cipher"].(string), cipher)
			contentLength, _ := strconv.Atoi(stream["contentLength"].(string))
			bitrate := int(stream["bitrate"].(float64))
			newAudioStream := NewAudioStream(itag, url, contentLength, bitrate, container, audioEncoding)
			audioStreams = append(audioStreams, newAudioStream)
		}
	}
	fmt.Println(audioStreams)
}

func getCipherSrc(assetFile string) *CipherOperations {
	regexFunc, err := regexp.Compile(`(\w+)=function\(\w+\){(\w+)=\w+\.split\(.*\);.*return.*\.join\(.*\)}`)
	errorHandler(err)
	cipherFunc := regexFunc.FindString(assetFile)
	if cipherFunc != "" {
		fmt.Println(cipherFunc)

		cipherName := cipherFunc[:strings.Index(cipherFunc, "=")]
		cipherBody := cipherFunc[strings.Index(cipherFunc, "{")+1 : strings.Index(cipherFunc, "}")]
		cipherStatements := strings.Split(cipherBody, ";")

		regexFunc, err = regexp.Compile(`(\w+).\w+\(\w+,\d+\);`)
		errorHandler(err)
		cipherAlgorithmName := strings.Split(regexFunc.FindString(cipherBody), ".")[0]
		// Important (?s) allows to to interpret newlines and whitespaces for dot character
		regexFunc, err = regexp.Compile(`var\s+` + cipherAlgorithmName + `=\{((?s)\w+:function\(\w+(,\w+)?\)\{(.*?)\}),?\};`)
		errorHandler(err)
		cipherAlgorithmBody := strings.ReplaceAll(regexFunc.FindString(assetFile), "\n", "")

		operations := newCipherOperations()

		for _, statement := range cipherStatements {
			// Get function name of statement
			//regexFunc, err = regexp.Compile(`\w+(?:.|\[)(""?\w+(?:"")?)]?\(`)
			//errorHandler(err)
			//statementName := regexFunc.FindString(statement)

			statementName := statement[strings.Index(statement, ".")+1 : strings.Index(statement, "(")]

			if statementName == "" {
				continue
			}

			if check, _ := regexp.MatchString(statementName+`:\bfunction\b\([a],b\).(\breturn\b)?.?\w+\.`, cipherAlgorithmBody);
			// Check if slice operation
			check {
				regexFunc, err = regexp.Compile(`\d+`)
				errorHandler(err)
				index, _ := strconv.Atoi(regexFunc.FindString(statement))
				operations.addOperation(newCipherSlice(index))
			} else if check, _ := regexp.MatchString(statementName+`:\bfunction\b\(\w+\,\w\).\bvar\b.\bc=a\b`, cipherAlgorithmBody);
			// Check if swap operation
			check {
				regexFunc, err = regexp.Compile(`\d+`)
				errorHandler(err)
				index, _ := strconv.Atoi(regexFunc.FindString(statement))
				operations.addOperation(newCipherSwap(index))
			} else if check, _ := regexp.MatchString(statementName+`:\bfunction\b\(\w+\,\w\).\bvar\b.\bc=a\b`, cipherAlgorithmBody);
			// Check if reverse operation
			check {
				operations.addOperation(newCipherReverse())
			}
		}

		fmt.Println(cipherName)
		fmt.Println(cipherBody)
		fmt.Println(cipherAlgorithmName)
		fmt.Println(cipherAlgorithmBody)

		return operations
	} else {
		// TODO breaking youtube api change error
	}
	return nil
}

func getVideoId(videoUrl string) string {
	return videoUrl[strings.Index(videoUrl, "watch?v=")+8:]
}

func GetVideoInfo(videoUrl string, output *string, wg *sync.WaitGroup) {
	videoId := getVideoId(videoUrl)
	eurl := url.QueryEscape("https://youtube.googleapis.com/v/" + videoId)
	metaUrl := "https://youtube.com/get_video_info?video_id=" + videoId +
		"&el=embedded&eurl=" + eurl + "&hl=en"
	data, err := downloadAsString(metaUrl)
	errorHandler(err)
	*output = data
	wg.Done()
}

func GetVideoEmbedPage(videoUrl string, output *string, wg *sync.WaitGroup) {
	flag := true
	search := "\"assets\":{\"js\":\""
	var result string
	for flag {
		videoId := getVideoId(videoUrl)
		metaUrl := "https://youtube.com/embed/" + videoId + "?disable_polymer=true&hl=en"
		data, err := downloadAsString(metaUrl)
		errorHandler(err)
		if strings.Contains(data, search) {
			flag = false
			result = data
		}
	}
	*output = result
	wg.Done()
}

func getAssetUrl(embedFile string) string {
	search := "\"assets\":{\"js\":\""
	searchLen := len(search)
	index := strings.Index(embedFile, search)
	assetBegin := embedFile[index+searchLen:]
	assetEncoded := assetBegin[:strings.Index(assetBegin, "\"")]
	assetDecoded := strings.ReplaceAll(assetEncoded, "\\/", "/")
	return "https://youtube.com" + assetDecoded
}

func getAssetFile(assetUrl string) string {
	assetFile, err := downloadAsString(assetUrl)
	errorHandler(err)
	return assetFile
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
