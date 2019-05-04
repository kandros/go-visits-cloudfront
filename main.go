package main

import (
	"log"
	"net/http"
	"os"

	"github.com/kandros/visits/internal/visit"

	apexlog "github.com/apex/log"
	"github.com/apex/log/handlers/json"
)

func main() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// logRequestHeaders(r)
		v := visit.NewFromRequest(r)

		if err := v.Persist(); err != nil {
			w.Write([]byte(err.Error()))
			return
		}

		w.Write([]byte("ok"))

	})

	addr := ":" + os.Getenv("PORT")

	log.Fatal(http.ListenAndServe(addr, handler))
}

func logRequestHeaders(r *http.Request) {
	apexlog.SetHandler(json.New(os.Stdout))
	apexlog.WithFields(apexlog.Fields{
		"headers": r.Header,
	}).Info("xxx")
}
