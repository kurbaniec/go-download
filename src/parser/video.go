package parser

import (
	"strconv"
	"strings"
	"sync"
)

type Container string

const (
	Mp4            Container = "mp4"
	Webm           Container = "webm"
	Tgpp           Container = "3gpp"
	OtherContainer Container = "other"
)

func getContainer(container string) Container {
	switch container {
	case "mp4":
		return Mp4
	case "webm":
		return Webm
	case "3gpp":
		return Tgpp
	}
	return OtherContainer
}

type AudioEncoding string

const (
	Aac        AudioEncoding = "mp4a"
	Vorbis     AudioEncoding = "vorbis"
	Opus       AudioEncoding = "opus"
	OtherAudio AudioEncoding = "other"
)

func getAudioEncoding(audioEncoding string) AudioEncoding {
	if strings.HasPrefix(audioEncoding, "mp4a") {
		return Aac
	}
	if strings.HasPrefix(audioEncoding, "vorbis") {
		return Vorbis
	}
	if strings.HasPrefix(audioEncoding, "opus") {
		return Opus
	}
	return OtherAudio
}

type AudioStream struct {
	Title         string
	Itag          int
	Url           string
	ContentLength int
	Bitrate       int
	Container     Container
	AudioEncoding AudioEncoding
}

func NewAudioStream(
	title string,
	itag int,
	url string,
	contentLength int,
	bitrate int,
	container Container,
	audioEncoding AudioEncoding) AudioStream {
	return AudioStream{
		Title:         title,
		Itag:          itag,
		Url:           url,
		ContentLength: contentLength,
		Bitrate:       bitrate,
		Container:     container,
		AudioEncoding: audioEncoding,
	}
}

func addAudioStream(
	title string,
	audioStreams *[]AudioStream,
	stream map[string]interface{},
	cipher *CipherOperations,
	wg *sync.WaitGroup,
) {
	defer wg.Done()
	// fmt.Println(stream["qualityLabel"]) // not available in audio
	mime := stream["mimeType"].(string)
	itag := int(stream["itag"].(float64))
	audioEncoding, container := getAudioEncodingAndContainer(mime)
	url, ok := stream["url"].(string)
	if !ok {
		url = buildStreamUrl(stream["cipher"].(string), cipher)
	}
	contentLength, _ := strconv.Atoi(stream["contentLength"].(string))
	bitrate := int(stream["bitrate"].(float64))
	newAudioStream := NewAudioStream(title, itag, url, contentLength, bitrate, container, audioEncoding)
	*audioStreams = append(*audioStreams, newAudioStream)
}

func getAudioEncodingAndContainer(mimeType string) (AudioEncoding, Container) {
	mimeType = strings.Replace(mimeType, "audio/", "", 1)
	mimeType = strings.ReplaceAll(mimeType, "codecs=\"", "")
	mimeType = mimeType[:len(mimeType)-1]
	values := strings.Split(mimeType, "; ")
	return getAudioEncoding(values[1]), getContainer(values[0])
}

var videoQuality = map[string]string{
	"5":   "144",
	"6":   "240",
	"13":  "144",
	"17":  "144",
	"18":  "360",
	"22":  "720",
	"34":  "360",
	"35":  "480",
	"36":  "240",
	"37":  "1080",
	"38":  "3072",
	"43":  "360",
	"44":  "480",
	"45":  "720",
	"46":  "1080",
	"59":  "480",
	"78":  "480",
	"82":  "360",
	"83":  "480",
	"84":  "720",
	"85":  "1080",
	"91":  "144",
	"92":  "240",
	"93":  "360",
	"94":  "480",
	"95":  "720",
	"96":  "1080",
	"100": "360",
	"101": "480",
	"102": "720",
	"132": "240",
	"151": "144",
	"133": "240",
	"134": "360",
	"135": "480",
	"136": "720",
	"137": "1080",
	"138": "4320",
	"160": "144",
	"212": "480",
	"213": "480",
	"214": "720",
	"215": "720",
	"216": "1080",
	"217": "1080",
	"264": "1440",
	"266": "2160",
	"298": "720",
	"299": "1080",
	"399": "1080",
	"398": "720",
	"397": "480",
	"396": "360",
	"395": "240",
	"394": "140",
	"167": "360",
	"168": "480",
	"169": "720",
	"170": "1080",
	"218": "480",
	"219": "480",
	"242": "240",
	"243": "360",
	"244": "480",
	"245": "480",
	"246": "480",
	"247": "720",
	"248": "1080",
	"271": "1440",
	"272": "2160",
	"278": "144",
	"302": "720",
	"303": "1080",
	"308": "1440",
	"313": "2160",
	"315": "2160",
	"330": "144",
	"331": "240",
	"332": "360",
	"333": "480",
	"334": "720",
	"335": "1080",
	"336": "1440",
	"337": "2160",
}

func getVideoQuality(itag string) string {
	quality, ok := videoQuality[itag]
	if !ok {
		quality = "unknown preset"
	}
	return quality
}
