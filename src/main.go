package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func serveWebsite() {
	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/p/", downloadHandler)

	port := ":8000"
	if os.Getenv("PORT") != "" {
		port = ":" + os.Getenv("PORT")
	}

	http.ListenAndServe(port, "cert.pem", "key.pem", nil)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(fmt.Sprintf("Method %s is not allowed on this endpoint.", r.Method)))
		return
	}

	body := r.Body
	id := GenerateID(6)
	file, err := os.Create("./files/" + id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	defer body.Close()
	defer file.Close()

	_, err = io.Copy(file, body)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(id))

	// Waits 5 minutes before deleting the file.
	go func() {
		time.Sleep(time.Minute * 5)

		if _, err := os.Stat("./files/" + id); err == nil {
			os.Remove("./files/" + id)
		}
	}()
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.Path, "/")[1:]
	id := path[1]

	if _, err := os.Stat("./files/" + id); os.IsNotExist(err) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`error("The requested item could not be found.", 0)`))
		return
	}

	file, err := os.Open("./files/" + id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	buf := make([]byte, 32*1000) // Reads in 32KB chunks (which is the default for io.Copy)
	w.Write([]byte("<html><body>"))

	for {
		_, err := file.Read(buf)

		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		w.Write(buf)
	}
	w.Write([]byte("</body></html>"))

	file.Close()
	os.Remove("./files/" + id)
}

func main() {
	// Clear files of the previous session
	if _, err := os.Stat("./files"); err == nil {
		os.RemoveAll("./files")
	}
	if _, err := os.Stat("./files"); os.IsNotExist(err) {
		os.Mkdir("./files", 0644)
	}

	// Serves the website
	serveWebsite()
}
