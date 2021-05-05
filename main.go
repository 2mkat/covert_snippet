package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
)

const value = "This you can write you secret message"

func router(w http.ResponseWriter, r *http.Request){
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	files := []string{
		"./html/index.html",
		"./html/layout.html",
		"./html/footer.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

func showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Отображение выбранной заметки с ID %d...", id)
}

func createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Метод запрещен!", 405)
		return
	}

	w.Header().Set("Set-Cookie", CompressString(value))

	w.Write([]byte("Создание новой заметки..."))
}

func SetCookieVal(w http.ResponseWriter, r *http.Request) {
	var C http.Cookie
	C.Value = value
	C.Expires = time.Now().AddDate(0, 1, 0)
	http.SetCookie(w, &C)
}

func CompressString(input string) string {

	//Compress gzip
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)

	zw.Write([]byte(input))
	zw.Close()

	//Encode base64
	gzipEncoded := base64.StdEncoding.EncodeToString(buf.Bytes())

	return gzipEncoded
}


type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}

			return nil, err
		}
	}

	return f, nil
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", router)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./static/")})
	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	
	log.Println("Запуск веб-сервера на http://127.0.0.1:8080")
	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}
