package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"arq-soft-II/backend/search-api/controllers"
	"arq-soft-II/backend/search-api/services"
	"arq-soft-II/config/cache"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}

	// Conectarse a Memcached
	memc, err := cache.New("localhost:11211")
	if err != nil {
		log.Fatal(" Error conectando con Memcached:", err)
	}
	defer log.Println("🧹 Conexión Memcached cerrada.")

	// Crear service y controller
	service := services.NewSearchService(memc)
	controller := controllers.NewSearchController(service)

	//  Endpoints HTTP
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Search API is running")
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Search API - Sports Activities Platform")
	})

	// Nuevo endpoint de búsqueda
	http.HandleFunc("/search", controller.HandleSearch)

	log.Printf("Search API corriendo en puerto %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
