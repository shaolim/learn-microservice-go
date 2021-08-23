package handlers

import (
	"errors"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"github.com/shaolim/learn-microservice-go/data"
)

type Products struct {
	logger *log.Logger
}

func NewProducts(logger *log.Logger) *Products {
	return &Products{logger}
}

func (p *Products) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.getProducts(rw, r)
		return
	}

	if r.Method == http.MethodPost {
		p.addProduct(rw, r)
		return
	}

	if r.Method == http.MethodPut {
		id, err := getIdFromUrl(r)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
		p.updateProduct(id, rw, r)
	}

	rw.WriteHeader(http.StatusMethodNotAllowed)
}

func getIdFromUrl(r *http.Request) (int, error) {
	reg := regexp.MustCompile(`/([0-9]+)`)
	g := reg.FindAllStringSubmatch(r.URL.Path, -1)
	if len(g) != 1 {
		return -1, errors.New("invalid URI")
	}

	if len(g[0]) != 2 {
		return -1, errors.New("invalid URI")
	}

	idString := g[0][1]
	id, err := strconv.Atoi(idString)
	if err != nil {
		return -1, errors.New("invalid URI")
	}

	return id, nil
}

func (p *Products) getProducts(rw http.ResponseWriter, r *http.Request) {
	p.logger.Println("Handle GET Product")

	products := data.GetProducts()
	err := products.ToJson(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (p *Products) addProduct(rw http.ResponseWriter, r *http.Request) {
	p.logger.Println("Handle POST Product")

	prod := &data.Product{}

	err := prod.FromJson(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
		return
	}

	data.AddProduct(prod)
}

func (p *Products) updateProduct(id int, rw http.ResponseWriter, r *http.Request) {
	p.logger.Println("Handle PUT Product")

	prod := &data.Product{}

	err := prod.FromJson(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
		return
	}

	err = data.UpdateProduct(id, prod)
	if err == data.ErrProductNotFound {
		http.Error(rw, "Product not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(rw, "Product not found", http.StatusInternalServerError)
		return
	}
}
