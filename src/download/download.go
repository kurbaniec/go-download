package download

import (
	"bufio"
	"downloader/src/parser"
	. "downloader/src/utils"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
)

// Struct used for downloading parts of the file
type downloadPart struct {
	fileBuffer []byte
	url        string
	chunkSize  int
	startRange int
	endRange   int
}

// Downloads the audio track of a YouTube video
// Returns the filename
func DownloadStream(stream parser.AudioStream, fileInfo FileInfo) {
	fmt.Println("Downloading...")
	fo, err := os.OpenFile(fileInfo.FileName(), os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	// Close output file on exit
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()
	// Setup variables
	fileSize := stream.ContentLength
	rangeSize := 9437184 / 10
	chunkSize := 4096
	startRange := 0
	downloadParts := make([]*downloadPart, 0)
	var wg sync.WaitGroup
	// Build download information for each part of the video
	for i := 0; startRange < fileSize; i++ {
		var endRange int
		if (startRange + rangeSize) < fileSize {
			endRange = startRange + rangeSize - 1
		} else {
			endRange = fileSize - 1
		}
		downloadParts = append(downloadParts, &downloadPart{
			fileBuffer: make([]byte, 0),
			url:        stream.Url,
			chunkSize:  chunkSize,
			startRange: startRange,
			endRange:   endRange,
		})
		wg.Add(1)
		startRange += rangeSize
	}
	// Download audio in parts
	for _, part := range downloadParts {
		go part.download(&wg)
	}
	wg.Wait()
	// Save audio buffer to file
	for _, part := range downloadParts {
		if _, err := fo.Write(part.fileBuffer); err != nil {
			panic(err)
		}
	}
}

func (part *downloadPart) download(wg *sync.WaitGroup) {
	defer wg.Done()
	client := &http.Client{}
	req, _ := http.NewRequest("GET", part.url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Range", "bytes="+strconv.Itoa(part.startRange)+"-"+strconv.Itoa(part.endRange))

	res, err := client.Do(req)
	if err != nil {
		panic(err.Error())
	}
	defer res.Body.Close()

	if strconv.Itoa(res.StatusCode)[0] == '2' {
		for {
			reader := bufio.NewReader(res.Body)
			buffer := make([]byte, part.chunkSize)
			n, err := reader.Read(buffer)
			if err == nil || err == io.EOF {
				part.fileBuffer = append(part.fileBuffer, buffer[:n]...)
				if err == io.EOF {
					break
				}
			} else {
				fmt.Println("Something went wrong downloading file: " + err.Error())
				os.Exit(1)
			}
		}
	} else {
		fmt.Println("Could not retrieve file from server.")
		fmt.Println("Server response: " + strconv.Itoa(res.StatusCode))
		os.Exit(1)
	}
}
