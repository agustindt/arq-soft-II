package main

import (
	"arq-soft-II/config/cache"
	"fmt"
	"time"
)

func main() {
	//  Conectar
	c, err := cache.New("localhost:11211")
	if err != nil {
		panic(err)
	}
	fmt.Println(" Conectado a Memcached")

	// Guardar un valor
	err = c.Set("greeting", []byte("Hola José desde Memcached "), 10*time.Second)
	if err != nil {
		panic(err)
	}
	fmt.Println(" Valor guardado en caché")

	// Leerlo
	val, err := c.Get("greeting")
	if err != nil {
		panic(err)
	}
	fmt.Println(" Valor leído desde caché:", string(val))

	// Esperar y volver a intentar
	time.Sleep(11 * time.Second)
	val, err = c.Get("greeting")
	if err != nil {
		fmt.Println(" El valor expiró:", err)
	} else {
		fmt.Println(" Valor aún disponible:", string(val))
	}
}
