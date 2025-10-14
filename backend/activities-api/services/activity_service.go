package services

import (
	"encoding/json"
	"log"

	"arq-soft-II/config/rabbitmq"
)

// Simulamos una estructura de actividad
type Activity struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Estructura del servicio de actividades
type ActivityService struct {
	mq *rabbitmq.Rabbit
}

// Constructor del servicio
func NewActivityService(mq *rabbitmq.Rabbit) *ActivityService {
	return &ActivityService{mq: mq}
}

// Crear una actividad (simulada)
func (s *ActivityService) CreateActivity(a *Activity) error {
	log.Println(" Creando actividad:", a.Name)

	// Publicamos un mensaje en RabbitMQ
	body, _ := json.Marshal(a)
	err := s.mq.Publish("entity.events", "activities.created", body)
	if err != nil {
		return err
	}

	log.Println(" Mensaje enviado a RabbitMQ para la actividad:", a.Name)
	return nil
}
