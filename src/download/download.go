package download

import (
	"bufio"
	"downloader/src/parser"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
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

	client := &http.Client{}
	req, _ := http.NewRequest("GET", stream.Url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0")

	res, err := client.Do(req)
	if err != nil {
		panic(err.Error())
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		reader := bufio.NewReader(res.Body)
		buffer := make([]byte, 81920)
		for {
			n, err := reader.Read(buffer)
			if err == io.EOF {
				break
			}
			// write a chunk
			if _, err := fo.Write(buffer[:n]); err != nil {
				panic(err)
			}
		}
	} else {
		fmt.Println("Error: " + strconv.Itoa(res.StatusCode))
	}

}
