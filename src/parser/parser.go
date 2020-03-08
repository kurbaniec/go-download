package parser

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func GetConfig(videoUrl string) {
	infoFile := GetVideoInfo(videoUrl)
	embedFile := GetVideoEmbedPage(videoUrl)
	assetFile := getAssetFile(embedFile)
	//watchFile := GetVideoWatchPage(videoUrl)
	getCipherSrc(assetFile)

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
	fmt.Println(dat["playabilityStatus"])
	fmt.Println(dat["streamingData"].(map[string]interface{})["formats"])
	formats := dat["streamingData"].(map[string]interface{})["adaptiveFormats"]
	fmt.Println(formats)
	for _, entry := range formats.([]interface{}) {
		set := entry.(map[string]interface{})
		fmt.Println(set["cipher"])
		/**
		streamUrl, err := url.QueryUnescape(set["ciphe`].(string))
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

func getCipherSrc(assetFile string) {
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
	} else {
		// TODO breaking youtube api change error
	}

	/**
	regex, err := regexp2.Compile(
		`(\w+)=function\(\w+\){(\w+)=\2\.split\(\x22{2}\);.*?return\s+\2\.join\(\x22{2}\)}`, 0)
	errorHandler(err)
	isMatch, err2 := regex.MatchString(assetFile)
	errorHandler(err2)
	if isMatch {
		funcName, _ := regex.FindStringMatch(assetFile)
		fmt.Println(funcName)
		r2, _ := regexp2.Compile(`(?!h\.)` + regexp2.Escape(funcName.String()) + `=function\(\w+\)\{(.*?)\}`, 0)
		isMatch, _ = r2.MatchString(assetFile)
		if isMatch {
			fmt.Println("OK!")
		}
	}*/

	/**
	patterns := [12]string{
		`\b[cs]\s*&&\s*[adf]\.set\([^,]+\s*,\s*encodeURIComponent\s*\(\s*(?P<sig>[a-zA-Z0-9$]+)\(`,
		`\b[a-zA-Z0-9]+\s*&&\s*[a-zA-Z0-9]+\.set\([^,]+\s*,\s*encodeURIComponent\s*\(\s*(?P<sig>[a-zA-Z0-9$]+)\(`,
		`\b(?P<sig>[a-zA-Z0-9$]{2})\s*=\s*function\(\s*a\s*\)\s*{\s*a\s*=\s*a\.split\(\s*""\s*\)`,
		`(?P<sig>[a-zA-Z0-9$]+)\s*=\s*function\(\s*a\s*\)\s*{\s*a\s*=\s*a\.split\(\s*""\s*\)`,
		`(["\'])signature\1\s*,\s*(?P<sig>[a-zA-Z0-9$]+)\(`,
		`\.sig\|\|(?P<sig>[a-zA-Z0-9$]+)\(`,
		`yt\.akamaized\.net/\)\s*\|\|\s*.*?\s*[cs]\s*&&\s*[adf]\.set\([^,]+\s*,\s*(?:encodeURIComponent\s*\()?\s*(?P<sig>[a-zA-Z0-9$]+)\(`,
		`\b[cs]\s*&&\s*[adf]\.set\([^,]+\s*,\s*(?P<sig>[a-zA-Z0-9$]+)\(`,
		`\b[a-zA-Z0-9]+\s*&&\s*[a-zA-Z0-9]+\.set\([^,]+\s*,\s*(?P<sig>[a-zA-Z0-9$]+)\(`,
		`\bc\s*&&\s*a\.set\([^,]+\s*,\s*\([^)]*\)\s*\(\s*(?P<sig>[a-zA-Z0-9$]+)\(`,
		`\bc\s*&&\s*[a-zA-Z0-9]+\.set\([^,]+\s*,\s*\([^)]*\)\s*\(\s*(?P<sig>[a-zA-Z0-9$]+)\(`,
		`\bc\s*&&\s*[a-zA-Z0-9]+\.set\([^,]+\s*,\s*\([^)]*\)\s*\(\s*(?P<sig>[a-zA-Z0-9$]+)\(`,
	}

	for _, pattern := range patterns {
		rr  := regexp.MustCompile(pattern)
		match := rr.MatchString(assetFile)
		if match {
			kek := rr.FindString(assetFile)
			fmt.Println(kek)
		}
	}*/
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
	return result
}

func getAssetFile(embedFile string) string {
	search := "\"assets\":{\"js\":\""
	searchLen := len(search)
	index := strings.Index(embedFile, search)
	assetBegin := embedFile[index+searchLen:]
	assetEncoded := assetBegin[:strings.Index(assetBegin, "\"")]
	assetDecoded := strings.ReplaceAll(assetEncoded, "\\/", "/")
	asset := "https://youtube.com" + assetDecoded
	fmt.Println(asset)
	assetFile, err := downloadAsString(asset)
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
