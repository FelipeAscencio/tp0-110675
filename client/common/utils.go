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
    // Máximo tamaño del mensaje (8 KB).
	MAX_TAMANIO_MENSAJE = 8192
	// Cantidad de bytes reservados para representar el tamaño.
	TAMANIO_BYTES       = 2
	// Mensaje de finalización.
	FINISH_MSJ          = "FINISH"
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
// Controlando el largo del mensaje y el correcto uso del 'Socket'.
func sendMessage(conn net.Conn, message string) error {
	bytes_mensaje := []byte(message)
	largo_mensaje := len(message)
	if largo_mensaje > MAX_TAMANIO_MENSAJE {
		// Limitación de tamaño pedida por default.
		log.Error("action: send_message | result: fail | error: message exceeds 8kb")
		return errors.New("message exceeds 8kb")
	}

	tamanio_mensaje := uint16(largo_mensaje)
	tamanio_buffer := make([]byte, TAMANIO_BYTES)
	binary.BigEndian.PutUint16(tamanio_buffer, tamanio_mensaje)
	_, err := conn.Write(tamanio_buffer)
	if err != nil {
		log.Errorf("action: send_message | result: fail | error: %v", err)
		return err
	}

	_, err = conn.Write(bytes_mensaje)
	if err != nil {
		log.Errorf("action: send_message | result: fail | error: %v", err)
		return err
	}

	return nil
}

// Función que recibe un mensaje del servidor a través del socket.
// Controlando el largo del mensaje para garantizar que no ocurre un 'Short-Read'.
func receiveMessage(conn net.Conn) (string, error) {
	mensaje, err := bufio.NewReader(conn).ReadString('\n')
	mensaje = strings.TrimSpace(mensaje)
	return mensaje, err
}

// Función que envía un batch de apuestas al servidor a través del socket.
func sendBetBatch(conn net.Conn, batch []Bet, betCount int) error {
	bets_str := make([]string, 0, len(batch))
	for _, bet := range batch {
		bet_str := fmt.Sprintf(
			"%s,%s,%s,%s,%s,%s",
			bet.AgencyId,
			bet.Name,
			bet.LastName,
			bet.Document,
			bet.BirthDate,
			bet.Number,
		)

		bets_str = append(bets_str, bet_str)
	}

	mensaje := strings.Join(bets_str, ";")
	return sendMessage(conn, mensaje)
}

//  Función que envía el mensaje de final de transmisión.
func sendFinishMessage(conn net.Conn) error {
	return sendMessage(conn, FINISH_MSJ)
}
