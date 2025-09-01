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
	MAX_TAMANIO_MENSAJE = 8192
	TAMANIO_BYTES       = 2
	FINISH_MSJ          = "FINISH"
	APUESTA_MSJ         = "BET"
	RESULTADOS_MSJ      = "RESULTS"
	PERDEDORES_MSJ      = "NOWINNERS"
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
		log.Error("action: send_message | result: fail | error: message exceeds 8kb")
		return errors.New("message exceeds 8kb")
	}

	tamanio_mensaje := uint16(largo_mensaje)
	sizeBuffer := make([]byte, TAMANIO_BYTES)
	binary.BigEndian.PutUint16(sizeBuffer, tamanio_mensaje)
	_, err := conn.Write(sizeBuffer)
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
	apuestas_string := make([]string, 0, len(batch))
	for _, bet := range batch {
		apuesta_string := fmt.Sprintf(
			"%s,%s,%s,%s,%s,%s",
			bet.AgencyId,
			bet.Name,
			bet.LastName,
			bet.Document,
			bet.BirthDate,
			bet.Number,
		)

		apuestas_string = append(apuestas_string, apuesta_string)
	}

	mensaje := strings.Join(apuestas_string, ";")
	return sendMessage(conn, mensaje)
}

// Función que envía el mensaje de final de transmisión.
func sendFinishMessage(conn net.Conn) error {
	return sendMessage(conn, FINISH_MSJ)
}

// Función que envía el mensaje de inicio de apuestas.
func sendStartBets(conn net.Conn) error {
	return sendMessage(conn, APUESTA_MSJ)
}

// Función que envía el mensaje de solicitud de resultados y recibe la respuesta.
func sendAskForResults(conn net.Conn, agencyId string) ([]string, error) {
	mensaje1 := RESULTADOS_MSJ
	mensaje2 := agencyId
	err := sendMessage(conn, mensaje1)
	if err != nil {
		return nil, err
	}

	err = sendMessage(conn, mensaje2)
	if err != nil {
		return nil, err
	}

	msj, err := receiveMessage(conn)
	if err != nil {
		return nil, err
	}

	if msj == PERDEDORES_MSJ {
		return []string{}, nil
	}

	return strings.Split(msj, ","), nil
}
