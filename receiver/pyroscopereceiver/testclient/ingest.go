package testclient

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
)

func Ingest(addr string, urlParams map[string]string, jfr string) error {
	data, err := os.ReadFile(jfr)
	if err != nil {
		return err
	}

	body := new(bytes.Buffer)

	mw := multipart.NewWriter(body)
	part, err := mw.CreateFormFile("jfr", "jfr")
	if err != nil {
		return fmt.Errorf("failed to create form file: %w", err)
	}
	gw := gzip.NewWriter(part)
	if _, err := gw.Write(data); err != nil {
		return err
	}
	gw.Close()
	mw.Close()

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/ingest", addr), body)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", mw.FormDataContentType())

	q := req.URL.Query()
	for k, v := range urlParams {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("failed to upload profile; http status code: %d", resp.StatusCode)
	}
	return nil
}
