package download

import (
	"bufio"
	"downloader/src/parser"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
)

type donwloadPart struct {
	fileBuffer []byte
	url        string
	chunkSize  int
	startRange int
	endRange   int
}

func DownloadStream(stream parser.AudioStream) {
	fmt.Println("Donwloading")
	fmt.Println(stream)

	// Create and open output file
	_, err := os.Create(string("output." + stream.Container))
	if err != nil {
		panic(err)
	}
	fo, err := os.OpenFile(string("output."+stream.Container), os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	// Close output file on exit
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()

	var wg sync.WaitGroup

	fileSize := stream.ContentLength
	rangeSize := 9437184 / 3
	chunkSize := 4096
	startRange := 0
	downloadParts := make([]*donwloadPart, 0)
	for i := 0; startRange < fileSize; i++ {
		var endRange int
		if (startRange + rangeSize) < fileSize {
			endRange = startRange + rangeSize - 1
		} else {
			endRange = fileSize - 1
		}
		downloadParts = append(downloadParts, &donwloadPart{
			fileBuffer: make([]byte, 0),
			url:        stream.Url,
			chunkSize:  chunkSize,
			startRange: startRange,
			endRange:   endRange,
		})
		wg.Add(1)
		startRange += rangeSize
	}

	for _, part := range downloadParts {
		go downloadPart(part, &wg)
	}
	wg.Wait()

	for _, part := range downloadParts {
		if _, err := fo.Write(part.fileBuffer); err != nil {
			panic(err)
		}
	}
}

func downloadPart(part *donwloadPart, wg *sync.WaitGroup) {
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

	chunkPosition := 0

	if strconv.Itoa(res.StatusCode)[0] == '2' {
		for {
			reader := bufio.NewReader(res.Body)
			var buffer []byte
			if chunkPosition+part.chunkSize < part.endRange {
				buffer = make([]byte, part.chunkSize)
			} else {
				buffer = make([]byte, (chunkPosition+part.chunkSize)%part.endRange)
			}
			n, err := reader.Read(buffer)
			if err == nil {
				part.fileBuffer = append(part.fileBuffer, buffer[:n]...)
			} else {
				if err == io.EOF {
					break
				} else {
					fmt.Println("Something went wrong downloading file: " + err.Error())
					os.Exit(1)
				}
			}
			chunkPosition += part.chunkSize
		}
	} else {
		fmt.Println("Could not retrieve file from server.")
		fmt.Println("Server response: " + strconv.Itoa(res.StatusCode))
		os.Exit(1)
	}
}
