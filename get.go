package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
)

func DecompressString(compressedString string) string {

	// Decode base64
	data, _ := base64.StdEncoding.DecodeString(compressedString)

	// Decompress gzip
	rdata := bytes.NewReader(data)
	r, _ := gzip.NewReader(rdata)
	s, _ := ioutil.ReadAll(r)

	return string(s)
}
func main() {
	url := "http://127.0.0.1:8080/snippet/create"

	var jsonStr = []byte(`snippets`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println(DecompressString(resp.Header.Get("Set-Cookie")))
}