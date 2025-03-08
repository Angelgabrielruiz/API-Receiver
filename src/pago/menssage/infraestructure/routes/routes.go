package routes

import (
	"github.com/gorilla/mux"

	"Receive/src/pago/menssage/infraestructure/controllers"
)

func SetupRoutes(router *mux.Router, mensajeController *controllers.MensajeController) {
	router.HandleFunc("/mensaje", mensajeController.RecibirMensaje).Methods("POST")
}