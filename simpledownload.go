package simpledownload

import (
	"io"
	"net/http"
	"os"
)

const NEW = 0
const IN_PROGRESS = 1
const DONE = 2
const ERROR = 3

type SimpleDownloader struct {
	requests []SimpleDownloaderRequest
}

type SimpleDownloaderRequest struct {
	url             string
	fileDestination string
	status          int
	err             error
}

type SimpleDownloaderReponse struct {
	Url             string
	FileDestination string
	Err             error
}

func (sd *SimpleDownloader) AddDownload(url, fileDestination string) {
	sd.requests = append(sd.requests, SimpleDownloaderRequest{url: url, fileDestination: fileDestination, status: NEW})
}

func (sd *SimpleDownloader) GetStatus() []SimpleDownloaderReponse {
	responses := []SimpleDownloaderReponse{}

	for _, r := range (*sd).requests {
		responses = append(responses, SimpleDownloaderReponse{Url: r.url, FileDestination: r.fileDestination, Err: r.err})
	}

	return responses
}

func (sd *SimpleDownloader) Download() {
	c := make(chan SimpleDownloaderReponse)

	for k, r := range (*sd).requests {
		if r.status != NEW && r.status != ERROR {
			continue
		}
		sd.requests[k].status = IN_PROGRESS
		go sd.downloadUrl(r.url, r.fileDestination, c)
	}

	for k := range sd.requests {
		if sd.requests[k].status != IN_PROGRESS {
			continue
		}
		response := <-c
		if response.Url == sd.requests[k].url {
			if response.Err == nil {
				sd.requests[k].status = DONE
			} else {
				sd.requests[k].status = ERROR
				sd.requests[k].err = response.Err
			}
		}
	}
}

func (SimpleDownloader) downloadUrl(url, fileDestination string, c chan SimpleDownloaderReponse) {
	out, err := os.Create(fileDestination)
	if err != nil {
		c <- SimpleDownloaderReponse{Url: url, Err: err}

		return
	}

	defer out.Close()
	resp, err := http.Get(url)
	if err != nil {
		c <- SimpleDownloaderReponse{Url: url, Err: err}

		return
	}

	defer resp.Body.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		c <- SimpleDownloaderReponse{Url: url, Err: err}

		return
	}
	c <- SimpleDownloaderReponse{Url: url, Err: nil}
}
