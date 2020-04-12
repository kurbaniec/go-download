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

func DownloadStream(stream parser.AudioStream) {
	fmt.Println("Donwloading")
	fmt.Println(stream)

	// open output file
	fo, err := os.Create(string("output." + stream.Container))
	if err != nil {
		panic(err)
	}
	// close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()

	var wg sync.WaitGroup

	rangeSize := 9437184
	partsCount := stream.ContentLength / rangeSize
	endSize := stream.ContentLength % rangeSize
	hasEndPart := false
	if endSize != 0 {
		partsCount += 1
		hasEndPart = true
	}

	fileBuffer := make([][]byte, partsCount)
	for i := 0; i < partsCount; i++ {
		startRange := i * rangeSize
		var endRange int
		if hasEndPart && i == partsCount-1 {
			endRange = endSize - 1
		} else {
			endRange = ((i + 1) * rangeSize) - 1
		}
		fileBuffer[i] = make([]byte, 0)
		wg.Add(1)
		go donwloadPart(i, &fileBuffer[i], stream.Url, startRange, endRange, &wg)
	}

	wg.Wait()

	for i := 0; i < partsCount; i++ {
		if _, err := fo.Write(fileBuffer[i]); err != nil {
			panic(err)
		}
	}

}

func donwloadPart(index int, fileBuffer *[]byte, url string, startRange int, endRange int, wg *sync.WaitGroup) {
	defer wg.Done()
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Range", "bytes="+strconv.Itoa(startRange)+"-"+strconv.Itoa(endRange))

	res, err := client.Do(req)
	if err != nil {
		panic(err.Error())
	}
	defer res.Body.Close()

	if strconv.Itoa(res.StatusCode)[0] == '2' {
		reader := bufio.NewReader(res.Body)
		buffer := make([]byte, 81920)
		for {
			n, err := reader.Read(buffer)
			if err == nil {
				*fileBuffer = append(*fileBuffer, buffer[:n]...)
			} else {
				if err == io.EOF {
					break
				} else {
					fmt.Println("Something went wrong downloading file: " + err.Error())
					os.Exit(1)
				}
			}
		}
	} else {
		fmt.Println("Could not retrieve file from server.")
		fmt.Println("Server response: " + strconv.Itoa(res.StatusCode))
		os.Exit(1)
	}
}
