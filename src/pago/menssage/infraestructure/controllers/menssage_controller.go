package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"Receive/src/pago/menssage/application/useCases"
	"Receive/src/pago/menssage/domain/entities.go"
	"Receive/src/pago/menssage/domain/repository"
	"Receive/src/pago/menssage/infraestructure/hub"
)

type mensajeDTO struct {
    Message   string `json:"message"`
    Contenido string `json:"contenido"`
}

type MensajeController struct {
    useCase *useCases.ProcesarMensajeUseCase
    wsHub   *hub.Hub
}

// Recibe la interfaz RabbitMQRepository
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

    var dto mensajeDTO
    if err := json.Unmarshal(body, &dto); err != nil {
        http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
        return
    }

    contenido := dto.Contenido
    if contenido == "" {
        contenido = dto.Message
    }

    mensaje := entities.Mensaje{
        ID:        "default-id",
        Contenido: contenido,
    }

    fmt.Println("Mensaje después de mapeo:", mensaje)

    if err := c.useCase.Execute(mensaje); err != nil {
        http.Error(w, "Error al procesar el mensaje", http.StatusInternalServerError)
        return
    }

    // Emitir mensaje por WebSocket
    messageBytes, err := json.Marshal(mensaje)
    if err != nil {
        http.Error(w, "Error al serializar el mensaje", http.StatusInternalServerError)
        return
    }

    c.wsHub.Broadcast(messageBytes)
    fmt.Println("Mensaje enviado con WebSocket")

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Mensaje recibido y enviado"))
}
