// Todo este archivo es generado para el ejercicio 5.

// Declaración del paquete common.
package common

// Importación de los paquetes necesarios.
import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"strings"
)

// Constantes de configuración del protocolo
const (
	MAX_TAMANIO_MENSAJE = 8192 // Máximo tamaño del mensaje (8 KB),
	TAMANIO_BYTES       = 2    // Cantidad de bytes reservados para representar el tamaño,
)

// Estructura que representa una apuesta.
type Bet struct {
	AgencyId  string
	Name      string
	LastName  string
	Document  string
	BirthDate string
	Number    string
}

// Función que envía un mensaje al servidor a través del socket.
// Controlando el largo del mensaje.
func sendMessage(conn net.Conn, message string) error {
	bytesMensaje := []byte(message)
	largoMensaje := len(message)
	if largoMensaje > MAX_TAMANIO_MENSAJE {
		log.Error("action: send_message | result: fail | error: message exceeds 8kb")
		return errors.New("message exceeds 8kb")
	}

	tamanioMensaje := uint16(largoMensaje)
	sizeBuffer := make([]byte, TAMANIO_BYTES)
	binary.BigEndian.PutUint16(sizeBuffer, tamanioMensaje)
	_, err := conn.Write(sizeBuffer)
	if err != nil {
		log.Errorf("action: send_message | result: fail | error: %v", err)
		return err
	}

	_, err = conn.Write(bytesMensaje)
	if err != nil {
		log.Errorf("action: send_message | result: fail | error: %v", err)
		return err
	}

	return nil
}

// Función que recibe un mensaje del servidor a través del socket.+
func receiveMessage(conn net.Conn) (string, error) {
	mensaje, err := bufio.NewReader(conn).ReadString('\n')
	mensaje = strings.TrimSpace(mensaje)
	return mensaje, err
}

// Función que envía una apuesta al servidor a través del socket.
func sendBet(conn net.Conn, bet Bet) error {
	mensaje := fmt.Sprintf(
		"%s,%s,%s,%s,%s,%s\n",
		bet.AgencyId,
		bet.Name,
		bet.LastName,
		bet.Document,
		bet.BirthDate,
		bet.Number,
	)

	return sendMessage(conn, mensaje)
}
