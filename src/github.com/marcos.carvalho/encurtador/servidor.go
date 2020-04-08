package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/marcos.carvalho/encurtador/url"
)

type Redirecionador struct {
	stats chan string
}
type Confs struct {
	port    int
	urlBase string
}

func (red *Redirecionador) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	caminho := strings.Split(r.URL.Path, "/")
	id := caminho[len(caminho)-1]

	if url := url.Buscar(id); url != nil {
		http.Redirect(w, r, url.Destino, http.StatusMovedPermanently)
		log.Printf("\n[Encurtador de URL] Redirecionando... ")
		red.stats <- id
	} else {
		http.NotFound(w, r)
	}

}

var (
	urlBase string
)

type Headers map[string]string

func main() {
	porta := 8888

	urlBase = fmt.Sprintf("http://localhost:%d", porta)
	stats := make(chan string)
	log.Println("[Encurtador de URL] Iniciando servidor!")
	http.HandleFunc("/api/encurtar", Encurtador)
	http.Handle("/r/", &Redirecionador{stats})
	http.HandleFunc("/api/stats/", Visualizador)
	log.Printf("\n[Encurtador de URL] servidor iniciado em %d", porta)

	go registrarEstatisticas(stats)
	defer close(stats)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", porta), nil))

}
func Visualizador(w http.ResponseWriter, r *http.Request) {
	caminho := strings.Split(r.URL.Path, "/")
	id := caminho[len(caminho)-1]

	if url := url.Buscar(id); url != nil {
		log.Println("[Encurtador de URL] URL encontrada! ", url)
		json, err := json.Marshal(url)
		log.Printf("\n[Encurtador de URL] JSON gerado: %s", json)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		responderComJSON(w, string(json))

	} else {
		http.NotFound(w, r)
	}

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
	responderCom(w, status, Headers{"Location": urlCurta, "Link": fmt.Sprintf("<%s/api/stats/%s>", urlBase, url.Id)})
	return
}

func responderCom(w http.ResponseWriter, status int, headers Headers) {
	for k, v := range headers {
		w.Header().Set(k, v)
	}
	w.WriteHeader(status)
}

func registrarEstatisticas(stats <-chan string) {
	fmt.Printf("[Encurtador de URL] Iniciando serviço de estatisticas.")
	for id := range stats {
		url.RegistrarClick(id)
		fmt.Printf("[Encurtador de URL] Click registrado com sucesso para %s\n", id)
	}
	fmt.Printf("[Encurtador de URL] Serviço de estatisticas finalizado!")
}

func responderComJSON(w http.ResponseWriter, json string) {
	responderCom(w, http.StatusOK, Headers{"Content-Type": "application/json"})
	fmt.Fprintf(w, json)

}
