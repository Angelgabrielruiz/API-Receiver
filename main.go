package main

import (
	"Receive/src/pago/menssage/infraestructure/controllers"
	"Receive/src/pago/menssage/infraestructure/database"
	"Receive/src/pago/menssage/infraestructure/routes"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	ws "Receive/src/pago/menssage/infraestructure/hub"
)

func main() {
	fmt.Println("Iniciando API...")

	// Iniciar el Hub de WebSocket
	hub := ws.NewHub()
	go hub.Run()

	// Conectar a RabbitMQ (implementación de RabbitMQRepository)
	rmq, err := database.NewRabbitMQ()
	if err != nil {
		log.Fatal("No se pudo conectar a RabbitMQ:", err)
	}
	defer rmq.Close()

	// Inyectar la implementación de RabbitMQ y el hub en el controlador.
	mensajeController := controllers.NewMensajeController(rmq, hub)

	// Configurar rutas usando mux.
	router := mux.NewRouter()
	routes.SetupRoutes(router, mensajeController)

	// Ruta para las conexiones WebSocket.
	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWs(hub, w, r)
	})

	port := ":8081"
	fmt.Println("API escuchando en", port)
	log.Fatal(http.ListenAndServe(port, router))
}