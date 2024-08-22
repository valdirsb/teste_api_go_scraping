package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	http.HandleFunc("/api/images", imagesHandler)
	port := "8080"
	fmt.Printf("Servidor rodando na porta %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// Handler para o endpoint /api/images
func imagesHandler(w http.ResponseWriter, r *http.Request) {
	// Habilita o CORS
	enableCORS(w, r)

	// Verifica se o método é OPTIONS para suporte CORS
	if r.Method == http.MethodOptions {
		return
	}

	// Recupera os parâmetros de consulta
	urlOrigem := r.URL.Query().Get("url_origem")
	urlDestino := r.URL.Query().Get("url_destino")

	if urlOrigem == "" || urlDestino == "" {
		http.Error(w, "Os parâmetros 'url_origem' e 'url_destino' são obrigatórios", http.StatusBadRequest)
		return
	}

	// Faz o scraping e gera a lista de URLs
	urls, err := scrapeImages(urlOrigem, urlDestino)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Define o header como JSON
	w.Header().Set("Content-Type", "application/json")

	// Codifica a lista de URLs em JSON e escreve na resposta
	if err := json.NewEncoder(w).Encode(urls); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Função para habilitar CORS
func enableCORS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

func scrapeImages(urlOrigem string, urlDestino string) ([]string, error) {
	resp, err := http.Get(urlOrigem)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status da resposta: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var urls []string
	doc.Find("div.thumb a img").Each(func(index int, item *goquery.Selection) {
		imgSrc, exists := item.Attr("data-src")
		if exists {
			fileName := path.Base(imgSrc)
			nameWithoutPrefixAndExt := removePrefixAndExt(fileName)
			newURL := createNewURL(nameWithoutPrefixAndExt, urlDestino)
			urls = append(urls, newURL)
		}
	})

	return urls, nil
}

// Função para remover a parte antes do primeiro ponto e a extensão do arquivo
func removePrefixAndExt(fileName string) string {
	dotIndex := strings.Index(fileName, ".")
	if dotIndex == -1 {
		return "" // Retorna vazio se não houver ponto
	}
	baseName := fileName[0:dotIndex]
	ext := path.Ext(baseName)
	nameWithoutExt := strings.TrimSuffix(baseName, ext)
	return nameWithoutExt
}

func createNewURL(fileName string, urlDestino string) string {
	if len(fileName) < 4 {
		return ""
	}
	part1 := fileName[:2]
	part2 := fileName[2:4]
	part3 := fileName[4:6]
	part4 := fileName
	return fmt.Sprintf("%s/%s/%s/%s/%s_169.mp4", urlDestino, part1, part2, part3, part4)
}
