package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"arq-soft-II/config/httpx"
	"arq-soft-II/config/rabbitmq"
)

// Activity representa una actividad simulada
type Activity struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Conexión global a RabbitMQ
var mq *rabbitmq.Rabbit

func main() {
	// Configurar puerto HTTP
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	// Conectarse a RabbitMQ
	var err error
	mq, err = rabbitmq.New("amqp://admin:admin@localhost:5672/")
	if err != nil {
		log.Fatal("Error conectando a RabbitMQ:", err)
	}
	defer mq.Close()

	// Declarar exchange y queue
	err = mq.DeclareSetup("entity.events", "search-sync", "activities.*")
	if err != nil {
		log.Fatal("Error declarando exchange y queue:", err)
	}
	log.Println("RabbitMQ listo para publicar mensajes")

	// Crear router multiplexer
	mux := http.NewServeMux()

	// Endpoint público de salud
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Activities API is running")
	})

	// Endpoint público raíz
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Activities API - Sports Activities Platform")
	})

	// Endpoint protegido: requiere JWT válido
	mux.Handle("/create",
		httpx.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			a := Activity{
				ID:          "A1",
				Name:        "Clase de Muay Thai",
				Description: "Profesor José",
			}

			// Publicar evento en RabbitMQ
			body, _ := json.Marshal(a)
			err := mq.Publish("entity.events", "activities.created", body)
			if err != nil {
				log.Println("Error publicando mensaje en RabbitMQ:", err)
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, "Error publicando mensaje en RabbitMQ")
				return
			}

			// Llamar internamente a search-api con X-Service-Token
			payload := []byte(`{"activity_id":"A1"}`)
			req, err := http.NewRequest(http.MethodPost, "http://localhost:8083/internal/reindex", bytes.NewReader(payload))
			if err != nil {
				log.Println("Error creando request a search-api:", err)
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, "Error creando request a search-api")
				return
			}
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-Service-Token", os.Getenv("SERVICE_TOKEN"))

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				log.Println("Error llamando a search-api:", err)
				w.WriteHeader(http.StatusBadGateway)
				fmt.Fprint(w, "Error llamando a search-api")
				return
			}
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				log.Println("search-api devolvió estado:", resp.Status)
			}

			log.Println("Mensaje publicado para actividad:", a.Name)
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "Actividad creada y mensaje enviado: %s", a.Name)
		})),
	)

	// Endpoint solo para administradores (JWT + is_admin=true)
	mux.Handle("/admin/stats",
		httpx.RequireAuth(httpx.RequireAdmin(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "Panel de estadísticas de administración (solo admins)")
		}))),
	)

	// Iniciar servidor HTTP
	log.Printf("Activities API iniciando en puerto %s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
