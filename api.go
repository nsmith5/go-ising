package main

import (
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/nsmith5/mjpeg"
)

func (m *Model) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		handler := mjpeg.Handler{
			Next: func() (image.Image, error) {
				time.Sleep(100 * time.Millisecond)
				return m.Image(), nil
			},
			Options: &jpeg.Options{Quality: 100},
		}
		handler.ServeHTTP(w, r)
	case "POST":
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(400)
			io.WriteString(w, "Failed to read body")
			return
		}

		beta, err := strconv.ParseFloat(string(b), 64)
		if err != nil {
			w.WriteHeader(400)
			io.WriteString(w, "Failed to parse float")
			return
		}

		m.Beta = beta
		log.Println("Beta modified to", m.Beta)
	}

	w.Header().Set("Access-Control-Allow-Origin", "https://www.nfsmith.ca")
}
