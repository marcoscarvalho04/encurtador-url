package url

import "log"

type RepositorioMemoria struct {
	urls map[string]*Url
}

func NovoRepositorioMemoria() *RepositorioMemoria {
	log.Printf("\n[Encurtador de URL] Iniciando repositório")
	return &RepositorioMemoria{make(map[string]*Url)}
}

func (r *RepositorioMemoria) IdExiste(id string) bool {
	log.Printf("\n[Encurtador de URL] Procurando por id %s", id)
	_, existe := r.urls[id]
	return existe
}

func (r *RepositorioMemoria) BuscarPorId(id string) *Url {
	log.Printf("\n[Encurtador de URL] Buscando por id: %s", id)
	return r.urls[id]
}

func (r *RepositorioMemoria) BuscarPorUrl(url string) *Url {
	log.Printf("\n[Encurtador de URL] Busca por url %s iniciada!", url)
	for _, value := range r.urls {
		if value.Destino == url {
			log.Printf("\n[Encurtador de URL] URL %s encontrada!", url)
			return value
		}
	}
	log.Printf("\n[Encurtador de URL] URL %s não encontrada!", url)
	return nil
}

func (r *RepositorioMemoria) Salvar(url Url) error {
	log.Printf("\n[Encurtador de URL] URL para salvar: %s", url.Destino)
	r.urls[url.Id] = &url
	return nil
}
