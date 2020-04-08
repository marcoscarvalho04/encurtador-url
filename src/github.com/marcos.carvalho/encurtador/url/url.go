package url

import (
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

type Url struct {
	Id      string `json: "id"`
	Criacao time.Time `json: "criacao"`
	Destino string `json: "destino"`
	Clicks  int    `json: "clicks"`
}

type Repositorio interface {
	IdExiste(id string) bool
	BuscarPorId(id string) *Url
	BuscarPorUrl(url string) *Url
	Salvar(url Url) error
	RegistrarClick(id string)
	BuscarClicks(id string) int
}

var repo Repositorio

const (
	tamanho  = 5
	simbolos = "abcdefghijklmnopqrstuvwyzABCDEFGIJKLMNOPQRSTUVWYZ123456789+-_"
)

func init() {
	rand.Seed(time.Now().UnixNano())
	ConfigurarRepositorio(NovoRepositorioMemoria())
}

func ConfigurarRepositorio(r Repositorio) {
	repo = r
}
func ExtrairUrl(r *http.Request) string {
	log.Printf("\n[Encurtador de URL] Iniciando extração da URL para ser encurtada")
	url := make([]byte, r.ContentLength, r.ContentLength)
	r.Body.Read(url)
	log.Printf("\n[Encurtador de URL] Extração de URL %s finalizada", url)
	return string(url)
}

func BuscarOuCriarNovaUrl(destino string) (u *Url, nova bool, err error) {
	log.Printf("\n[Encurtador de URL] Iniciando criação de nova URL (Ou busca de já existente)")
	if u = repo.BuscarPorUrl(destino); u != nil {
		log.Printf("\n[Encurtador de URL] URL encontrada! retornando: %s ", u.Destino)
		return u, false, nil
	}

	if _, err = url.ParseRequestURI(destino); err != nil {
		return nil, false, err
	}
	log.Printf("\n[Encurtador de URL] Iniciando criação de nova URL (Ou busca de já existente)")
	url := Url{gerarId(), time.Now(), destino, 0}
	repo.Salvar(url)
	return &url, true, nil
}

func gerarId() string {
	log.Printf("\n[Encurtador de URL] Iniciar geração de ID")
	novoID := func() string {
		id := make([]byte, tamanho, tamanho)
		for i := range id {
			id[i] = simbolos[rand.Intn(len(simbolos))]
		}
		return string(id)
	}
	for {
		if id := novoID(); !repo.IdExiste(id) {
			log.Printf("\n[Encurtador de URL] ID gerado: %s", id)
			return id
		}
	}

}

func Buscar(id string) *Url {
	url := repo.BuscarPorId(id)
	return url
}
func RegistrarClick(id string) {
	repo.RegistrarClick(id)
}
