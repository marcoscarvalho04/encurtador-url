package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/marcos.carvalho/encurtador/url"
)

var (
	porta   int
	urlBase string
)

func init() {
	porta = 8888
	urlBase = fmt.Sprintf("http://localhost:%d", porta)
}

type Headers map[string]string

func main() {
	log.Println("[Encurtador de URL] Iniciando servidor!")
	http.HandleFunc("/api/encurtar", Encurtador)
	http.HandleFunc("/r/", Redirecionador)
	log.Printf("\n[Encurtador de URL] servidor iniciado em %d", porta)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", porta), nil))

}

func Encurtador(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		responderCom(w, http.StatusMethodNotAllowed, Headers{"Allow": "POST"})
	}
	urlExtraida := url.ExtrairUrl(r)

	url, nova, err := url.BuscarOuCriarNovaUrl(urlExtraida)
	if err != nil {
		log.Printf("\n[Encurtador de URL] Erro ao encurtar a URL: %s", urlExtraida)
		responderCom(w, http.StatusBadRequest, nil)
		return
	}

	var status int
	if nova {
		status = http.StatusCreated
	} else {
		status = http.StatusOK
	}

	urlCurta := fmt.Sprintf("%s/r/%s", urlBase, url.Id)
	log.Printf("\n[Encurtador de URL] URL encurtada: %s, passou a ser: %s ", urlExtraida, urlCurta)
	responderCom(w, status, Headers{"Location": urlCurta})
	return
}

func responderCom(w http.ResponseWriter, status int, headers Headers) {
	for k, v := range headers {
		w.Header().Set(k, v)
	}
	w.WriteHeader(status)
}

func Redirecionador(w http.ResponseWriter, r *http.Request) {
	caminho := strings.Split(r.URL.Path, "/")
	id := caminho[len(caminho)-1]

	if url := url.Buscar(id); url != nil {
		http.Redirect(w, r, url.Destino, http.StatusMovedPermanently)
	} else {
		http.NotFound(w, r)
	}

}
