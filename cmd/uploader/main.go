package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/schollz/progressbar/v3"
)

type DiskMetadata struct {
	VirtualMachine string `json:"virtual_machine" xml:"virtual_machine"`
	ByteSize       int64  `json:"byte_size" xml:"byte_size"`
}

func uploadChunk(filePath string, virtualMachine string, index int64, chunkSize int64, totalSize int64, remoteAddress string, httpClient *http.Client, done chan int) {
	var fd *os.File
	var err error
	var byteRead int
	var chunk []byte = make([]byte, chunkSize)
	var resp *http.Response

	fd, err = os.OpenFile(filePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatalf("Unable to open file %s", filePath)
	}
	defer fd.Close()

	_, err = fd.Seek(index, 0)
	if err != nil {
		log.Fatalf("Unable to seek file %s", filePath)
	}
	byteRead, err = fd.Read(chunk)
	if err != nil {
		log.Fatalf("Unable to read from file %s", filePath)
	}

	req, err := http.NewRequest("PUT", remoteAddress, bytes.NewBuffer(chunk[:byteRead]))
	if err != nil {
		done <- 1
		return
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("X-VirtualMachine", virtualMachine)
	req.ContentLength = int64(byteRead)
	req.Header.Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", index, index+int64(byteRead)-1, totalSize))
	resp, err = httpClient.Do(req)
	if err != nil || resp.StatusCode < 200 || resp.StatusCode >= 300 {
		done <- 1
		return
	}

	done <- 0

}

func composeUri(parts ...string) string {
	var res string = ""
	for _, part := range parts {
		if res == "" {
			res = part
			continue
		}
		if res[len(res)-1:] == "/" && part[0] == '/' {
			res += part[1:]
			continue
		}
		if res[len(res)-1:] != "/" && part[0] != '/' {
			res += "/" + part
			continue
		}
		res += part
	}
	return res
}

func main() {
	var filePath string
	var err error
	var fd *os.File
	var fi os.FileInfo
	var byteLength int64
	var chunkSize int64 = 1024 * 1024 // 1 MB
	var transport *http.Transport
	var httpClient *http.Client
	var remoteAddress string
	var filename string
	var resp *http.Response
	var bar *progressbar.ProgressBar
	var virtualMachine string

	flag.StringVar(&filePath, "path", "", "File path to upload")
	flag.StringVar(&virtualMachine, "virtual_machine", "", "Virtual machine id")
	flag.StringVar(&remoteAddress, "host", "127.0.0.1", "url endpoint to upload file to")
	flag.StringVar(&filename, "filename", "", "Name to assign on remote host")

	flag.Parse()

	transport = &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10, // Number of reusable connections per host
		IdleConnTimeout:     90 * time.Second,
	}

	httpClient = &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}

	fd, err = os.OpenFile(filePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatalf("Unable to open file %s", filePath)
	}
	defer fd.Close()

	fi, err = fd.Stat()
	if err != nil {
		log.Fatalf("Unable to get stats for file %s", filePath)
	}

	byteLength = fi.Size()

	var jsonEncoded []byte
	jsonEncoded, err = json.Marshal(DiskMetadata{
		VirtualMachine: virtualMachine,
		ByteSize:       byteLength,
	})
	if err != nil {
		log.Fatal("Unable to create json body")
	}

	fmt.Printf("Init file upload %s\n", filePath)
	resp, err = httpClient.Post(composeUri(remoteAddress, "/upload/", filename, "/begin"), "application/json", bytes.NewBuffer(jsonEncoded))
	if err != nil || resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
		}
		log.Fatalf("Failed to initialize upload")
	}

	var i int64
	var jobs int = 12
	var job int
	var done chan int = make(chan int, jobs)
	bar = progressbar.Default(100, "Upload status")
	var chunks int64 = byteLength / chunkSize

	// This loop tries to keep the pool of jobs full
	for i = 0; i < chunks; i++ {
		job += 1
		go uploadChunk(filePath, virtualMachine, chunkSize*i, chunkSize, byteLength, composeUri(remoteAddress, "/upload/", filename, "/chunk"), httpClient, done)

		// If pool has reached the limit, wait for a job to finish
		if job == jobs {
			<-done
			bar.Set((int)(i * 100 / chunks))
			job -= 1
		}
	}
	for ; job > 0; job-- {
		<-done
		bar.Set((int)(i * 100 / byteLength))
	}
	bar.Set(100)

	fmt.Printf("Finishing file upload %s\n", filePath)
	resp, err = httpClient.Post(composeUri(remoteAddress, "/upload/", filename, "/commit"), "application/json", bytes.NewBuffer(jsonEncoded))
	if err != nil || resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Fatalf("Unable to finish file upload %s", filePath)
	}

	os.Exit(0)
}
