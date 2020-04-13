package parser

import (
	. "downloader/src/opts"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// Return audio stream information of the given video and specified quality
func GetAudioStream(
	videoUrl string,
	cipherStore map[string]*CipherOperations,
	opts *Opts,
) AudioStream {
	audioStreams := make([]AudioStream, 0, 10)
	getStreams(videoUrl, cipherStore, &audioStreams)
	audioStreams = filterAndSortAudioStreams(audioStreams)
	if opts.Quality == "high" {
		return audioStreams[len(audioStreams)-1]
	} else if opts.Quality == "medium" {
		return audioStreams[len(audioStreams)/2]
	} else {
		return audioStreams[0]
	}
}

// Filter audio streams to only feature audio containers with "webm" and sort them according
// to their quality (= bitrate).
func filterAndSortAudioStreams(audioStreams []AudioStream) []AudioStream {
	filtered := make([]AudioStream, 0)
	for _, stream := range audioStreams {
		if stream.Container == Webm {
			filtered = append(filtered, stream)
		}
	}
	sort.SliceStable(filtered, func(i, j int) bool {
		return filtered[i].Bitrate < filtered[j].Bitrate
	})
	return filtered
}

// Returns all audio streams of the given video
func getStreams(
	videoUrl string,
	cipherStore map[string]*CipherOperations,
	audioStreams *[]AudioStream,
) {
	// Get all video information
	var infoFile, embedFile, assetUrl string
	var wg sync.WaitGroup
	wg.Add(2)
	go GetVideoInfo(videoUrl, &infoFile, &wg)
	go GetVideoEmbedPage(videoUrl, &embedFile, &wg)
	wg.Wait()
	assetUrl = getAssetUrl(embedFile)
	// Check if cipher to decrypt url was already scrambled
	if _, ok := cipherStore[assetUrl]; !ok {
		// If not, get cipher
		assetFile := getAssetFile(assetUrl)
		cipherStore[assetUrl] = getCipherSrc(assetFile)
	}
	cipher := cipherStore[assetUrl]
	// Parse meta information
	metaEncoded := infoFile
	metaDecoded, err := url.QueryUnescape(metaEncoded)
	if err != nil {
		fmt.Println("Something went wrong reading meta-data")
		os.Exit(0)
	}
	metaArray := strings.Split(metaDecoded, "&")
	var info string
	for _, entry := range metaArray {
		if strings.HasPrefix(entry, "player_response") {
			info = entry[16:]
		}
	}
	var dat map[string]interface{}
	if err := json.Unmarshal([]byte(info), &dat); err != nil {
		fmt.Println("MARSHALL ERROr")
	}
	title := dat["videoDetails"].(map[string]interface{})["title"].(string)
	formats := dat["streamingData"].(map[string]interface{})["adaptiveFormats"].([]interface{})

	wg.Add(len(formats))
	// Parse meta information of all audio streams
	for _, entry := range formats {
		stream := entry.(map[string]interface{})
		mime := stream["mimeType"].(string)
		if strings.HasPrefix(mime, "audio") {
			go addAudioStream(title, audioStreams, stream, cipher, &wg)
		} else {
			// Video streams are not supported at the moment
			/**url, ok := stream["url"].(string)
			if !ok {
				url = buildStreamUrl(stream["cipher"].(string), cipher)
			}
			//fmt.Println(url)*/
			wg.Done()
		}
	}
	wg.Wait()
}

func getCipherSrc(assetFile string) *CipherOperations {
	regexFunc, err := regexp.Compile(`(\w+)=function\(\w+\){(\w+)=\w+\.split\(.*\);.*return.*\.join\(.*\)}`)
	errorHandler(err)
	cipherFunc := regexFunc.FindString(assetFile)
	if cipherFunc != "" {
		//fmt.Println(cipherFunc)
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
			} else if check, _ := regexp.MatchString(statementName+`:\bfunction\b\(\w+\).\b\w+\.reverse\(\)`, cipherAlgorithmBody);
			// Check if reverse operation
			check {
				operations.addOperation(newCipherReverse())
			}
		}
		// cipherName := cipherFunc[:strings.Index(cipherFunc, "=")]
		//fmt.Println(cipherName)
		//fmt.Println(cipherBody)
		//fmt.Println(cipherAlgorithmName)
		//fmt.Println(cipherAlgorithmBody)
		return operations
	} else {
		fmt.Println("Seems like the YouTube API has changed...")
		fmt.Println("Please contact the repo owner about this breaking change.")
		os.Exit(1)
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
		"&el=embedded&eurl=" + eurl + "&hl=en_US"
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
		metaUrl := "https://youtube.com/embed/" + videoId + "?hl=en"
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
