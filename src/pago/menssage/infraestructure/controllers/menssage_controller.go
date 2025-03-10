package controllers

import (
	"Receive/src/pago/menssage/application/useCases"
	"Receive/src/pago/menssage/domain/entities.go"
	"Receive/src/pago/menssage/domain/repository"
	"Receive/src/pago/menssage/infraestructure/hub"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Utilizamos un DTO para adaptarnos a distintos formatos que pueda enviar el front o consumidor.
type mensajeDTO struct {
	Message   string `json:"message"`
	Contenido string `json:"contenido"`
}

type MensajeController struct {
	useCase *useCases.ProcesarMensajeUseCase
	wsHub   *hub.Hub
}

func NewMensajeController(rabbitRepo repository.RabbitMQRepository, wsHub *hub.Hub) *MensajeController {
	return &MensajeController{
		useCase: useCases.NewProcesarMensajeUseCase(rabbitRepo),
		wsHub:   wsHub,
	}
}

func (c *MensajeController) RecibirMensaje(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error al leer el cuerpo de la solicitud", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	fmt.Println("JSON recibido:", string(body))

	// Deserializamos en un DTO para admitir tanto "message" como "contenido"
	var dto mensajeDTO
	if err := json.Unmarshal(body, &dto); err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	// Mapear el DTO al objeto de dominio. Si "contenido" está vacío, usamos "message".
	contenido := dto.Contenido
	if contenido == "" {
		contenido = dto.Message
	}
	mensaje := entities.Mensaje{
		ID:        "default-id", // Se asigna un ID por defecto (puedes generar un UUID si prefieres)
		Contenido: contenido,
	}

	fmt.Println("Mensaje después de mapeo:", mensaje)

	// Ejecutar el use case para enviar el mensaje a RabbitMQ.
	if err := c.useCase.Execute(mensaje); err != nil {
		http.Error(w, "Error al procesar el mensaje", http.StatusInternalServerError)
		return
	}

	// Convertir el mensaje a JSON para emitirlo vía WebSocket.
	messageBytes, err := json.Marshal(mensaje)
	if err != nil {
		http.Error(w, "Error al serializar el mensaje", http.StatusInternalServerError)
		return
	}

	// Emitir el mensaje a todos los clientes conectados.
	c.wsHub.Broadcast(messageBytes)
	fmt.Println("Mensaje procesado y emitido vía WebSocket")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Mensaje recibido y emitido correctamente"))
}