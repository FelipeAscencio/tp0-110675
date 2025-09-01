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

// Función que escribe todos los bytes en el socket (Garantiza que no tenemos 'Short-write').
func writeAll(conn net.Conn, buf []byte) error {
	total := 0
	for total < len(buf) {
		n, err := conn.Write(buf[total:])
		if err != nil {
			return err
		}

		total += n
	}

	return nil
}

// Función que envía un mensaje al servidor a través del socket.
// Controlando el largo del mensaje.
func sendMessage(conn net.Conn, message string) error {
	bytes_mensaje := []byte(message)
	largo_mensaje := len(message)
	if largo_mensaje > MAX_TAMANIO_MENSAJE {
		log.Error("action: send_message | result: fail | error: message exceeds 8kb")
		return errors.New("message exceeds 8kb")
	}

	tamanio_mensaje := uint16(largo_mensaje)
	sizeBuffer := make([]byte, TAMANIO_BYTES)
	binary.BigEndian.PutUint16(sizeBuffer, tamanio_mensaje)
	if err := writeAll(conn, sizeBuffer); err != nil {
		log.Errorf("action: send_message | result: fail | error: %v", err)
		return err
	}

	if err := writeAll(conn, bytes_mensaje); err != nil {
		log.Errorf("action: send_message | result: fail | error: %v", err)
		return err
	}

	return nil
}

// Función que recibe un mensaje del servidor a través del socket.
// Controlando el largo del mensaje para garantizar que no ocurre un 'Short-Read'.
func receiveMessage(conn net.Conn) (string, error) {
	msg, err := bufio.NewReader(conn).ReadString('\n')
	msg = strings.TrimSpace(msg)
	return msg, err
}

// Función que envía una apuesta al servidor a través del socket.
func sendBet(conn net.Conn, bet Bet) error {
	message := fmt.Sprintf(
		"%s,%s,%s,%s,%s,%s\n",
		bet.AgencyId,
		bet.Name,
		bet.LastName,
		bet.Document,
		bet.BirthDate,
		bet.Number,
	)

	return sendMessage(conn, message)
}
