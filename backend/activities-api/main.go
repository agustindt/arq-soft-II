package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"arq-soft-II/config/rabbitmq"
)

// Estructura de una actividad (simulada)
type Activity struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Estructura para manejar la conexión a RabbitMQ (global)
var mq *rabbitmq.Rabbit

func main() {
	// 1️ Configurar puerto HTTP
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	// 2️ Conectarse a RabbitMQ
	var err error
	mq, err = rabbitmq.New("amqp://admin:admin@localhost:5672/")
	if err != nil {
		log.Fatal(" Error conectando a RabbitMQ:", err)
	}
	defer mq.Close()

	// 3️ Declarar exchange y queue
	err = mq.DeclareSetup("entity.events", "search-sync", "activities.*")
	if err != nil {
		log.Fatal(" Error declarando exchange y queue:", err)
	}
	log.Println(" RabbitMQ listo para publicar mensajes")

	// 4️ Endpoint de salud
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Activities API is running")
	})

	// 5️ Endpoint raíz
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Activities API - Sports Activities Platform")
	})

	// 6️ Endpoint para crear una actividad (simulada)
	http.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		a := Activity{
			ID:          "A1",
			Name:        "Clase de Muay Thai",
			Description: "Profesor José ",
		}

		body, _ := json.Marshal(a)
		err := mq.Publish("entity.events", "activities.created", body)
		if err != nil {
			log.Println(" Error publicando mensaje en RabbitMQ:", err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "Error publicando mensaje en RabbitMQ")
			return
		}

		log.Println(" Mensaje publicado para actividad:", a.Name)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Actividad creada y mensaje enviado: %s", a.Name)
	})

	// 7️ Iniciar servidor HTTP
	log.Printf(" Activities API iniciando en puerto %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
