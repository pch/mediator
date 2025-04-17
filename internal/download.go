package internal

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"time"
)

type DownloadedFile struct {
	OriginalURL   string
	FinalURL      string
	ContentType   string
	ContentLength int
	StatusCode    int
	Buffer        bytes.Buffer
}

func (f *DownloadedFile) Size() int {
	return f.Buffer.Len()
}

const UserAgent = "mediator"

type requestHandler func(req *http.Request)
type responseHandler func(resp *http.Response)

func DownloadFile(url string, maxSize int, requestTimeout time.Duration) (file *DownloadedFile, err error) {
	var out bytes.Buffer

	slog.Debug("Downloading file", "url", url)

	downloadedFile, err := doRequest(url, maxSize, requestTimeout,
		func(req *http.Request) {},
		func(resp *http.Response) {},
		&out,
	)

	if err != nil {
		return nil, err
	}

	if downloadedFile.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-200 response code: %d, ", downloadedFile.StatusCode)
	}

	downloadedFile.Buffer = out

	size := downloadedFile.Size()
	contentLength := downloadedFile.ContentLength

	// make sure we downloaded the whole thing
	// content-length is sometimes 0 (server doesn't return it...)
	if (contentLength > 0 && contentLength != size) || size == 0 {
		return nil, fmt.Errorf("incomplete download: size: %d, content-length: %d", size, contentLength)
	}

	return downloadedFile, nil
}

var copiedResponseHeaders = []string{
	"Content-Type",
	"Content-Length",
	"Content-Encoding",
	"Transfer-Encoding",
}

// Download file and stream it back as response
func ProxyFile(url string, maxSize int, requestTimeout time.Duration, r *http.Request, w http.ResponseWriter) (*DownloadedFile, error) {
	return doRequest(url, maxSize, requestTimeout,
		func(req *http.Request) {
			// forward request headers
			for k, vv := range r.Header {
				for _, v := range vv {
					req.Header.Set(k, v)
				}
			}
		},
		func(resp *http.Response) {
			// copy selected response headers, preserving any existing ones
			for _, h := range copiedResponseHeaders {
				if w.Header().Get(h) == "" { // only set if not already set
					w.Header().Set(h, resp.Header.Get(h))
				}
			}

			w.WriteHeader(resp.StatusCode)
		},
		w,
	)
}

func checkFileSize(maxSize int, contentLength string) (int, error) {
	size, _ := strconv.Atoi(contentLength)

	if size > maxSize {
		return size, fmt.Errorf("file too big: %d (max: %d)", size, maxSize)
	}

	return size, nil
}

func doRequest(url string, maxSize int, requestTimeout time.Duration, reqHandler requestHandler, respHandler responseHandler, out io.Writer) (*DownloadedFile, error) {
	client := newHttpClient(requestTimeout)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("request error. %+v", err)
	}

	req.Header.Set("User-Agent", UserAgent)
	reqHandler(req)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request error. %+v", err)
	}
	defer resp.Body.Close()

	size, err := checkFileSize(maxSize, resp.Header.Get("Content-Length"))
	if err != nil {
		return nil, err
	}

	respHandler(resp)

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return nil, err
	}

	return &DownloadedFile{
		OriginalURL:   url,
		FinalURL:      resp.Request.URL.String(), // the last URL client tried to access
		ContentType:   resp.Header.Get("Content-Type"),
		ContentLength: size,
		StatusCode:    resp.StatusCode,
	}, nil
}

func newHttpClient(timout time.Duration) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: timout,
			}).Dial,
			TLSHandshakeTimeout: timout,
		},
	}
}
