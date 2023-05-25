package main

import (
	"log"
	"net/http"
	"os"

	wk "github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

var (
	XToken = os.Getenv("TOKEN")
	Port   = os.Getenv("BIND_ADDR")
)

func main() {
	if XToken == "" {
		log.Fatal("TOKEN is not set")
	}
	http.HandleFunc("/input", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "X-Token, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		w.Header().Set("content-security-policy", "default-src 'none'")

		if r.Method == "OPTIONS" {
			return
		}
		// write cors headers

		if r.Header.Get("X-Token") != XToken {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// wk.SetPath("/bin/wkhtmltopdf")

		pdfg, err := wk.NewPDFGenerator()
		if err != nil {
			log.Fatalf("NewPDFGenerator %v", err)
		}

		pdfg.Dpi.Set(300)
		pdfg.Orientation.Set(wk.OrientationPortrait)
		// pdfg.Grayscale.Set(true)

		page := wk.NewPageReader(r.Body)

		page.FooterRight.Set("[page]")
		page.FooterFontSize.Set(10)
		page.Zoom.Set(0.95)

		pdfg.AddPage(page)

		err = pdfg.Create()
		if err != nil {
			log.Fatalf("pdfg.Create: %v", err)
		}

		_, err = w.Write(pdfg.Buffer().Bytes())
		if err != nil {
			log.Fatal("Write: %v", err)
		}
	})
	log.Printf("Starting server on port :" + Port)

	log.Fatal(http.ListenAndServe(":"+Port, nil))
}
