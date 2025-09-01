// Declaración del paquete common.
package common

// Importación de los paquetes necesarios.
import (
	"net"
	"os"
	"strconv"
	"strings"
	"time"
	"github.com/op/go-logging"
)

// Constantes para índices de los campos de la respuesta.
const (
	INDEX_DOCUMENTO = 0
	INDEX_NUMERO   = 1
)

// Declaración de una variable global para el logger.
var log = logging.MustGetLogger("log")

// Estructura que contiene la configuración del cliente.
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopAmount    int
	LoopPeriod    time.Duration
}

// Estructura que representa al cliente.
type Client struct {
	config ClientConfig
	conn   net.Conn
	bet    Bet
}

// Función que crea un nuevo cliente con la configuración dada.
func NewClient(config ClientConfig, bet Bet) *Client {
	client := &Client{
		config: config,
		bet:    bet,
	}

	return client
}

// Función que crea el socket del cliente y se conecta al servidor.
func (c *Client) createClientSocket() error {
	conn, err := net.Dial("tcp", c.config.ServerAddress)
	if err != nil {
		log.Criticalf(
			"action: connect | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
	}

	c.conn = conn
	return nil
}

// Función que inicia el loop del cliente, enviando mensajes al servidor.
func (c *Client) StartClientLoop(sigChan chan os.Signal) {
	select {
	case <-sigChan:
		log.Infof("action: shutdown | result: success")
		return
	default:

		// Modificación de código para el ejercicio 5.
		c.createClientSocket()
		err := sendBet(c.conn, c.bet)
		if err != nil {
			return
		}

		mensaje, err := receiveMessage(c.conn)
		if err != nil {
			log.Errorf("action: finalizar_envio | result: fail | error: %v",
				c.config.ID,
				err,
			)

			return
		}

		c.conn.Close()
		informacion_respuesta := strings.Split(mensaje, ",")
		doc_apuesta, _ := strconv.Atoi(strings.TrimSpace(c.bet.Document))
		numero_apuesta, _ := strconv.Atoi(strings.TrimSpace(c.bet.Number))
		doc_respuesta, _ := strconv.Atoi(strings.TrimSpace(informacion_respuesta[INDEX_DOCUMENTO]))
		numero_respuesta, _ := strconv.Atoi(strings.TrimSpace(informacion_respuesta[INDEX_NUMERO]))
		if doc_respuesta == doc_apuesta && numero_respuesta == numero_apuesta {
			log.Infof("action: apuesta_enviada | result: success | dni: %v | numero: %v", c.bet.Document, c.bet.Number)
		} else {
			log.Errorf("action: apuesta_enviada | result: fail | dni: %v | numero: %v", c.bet.Document, c.bet.Number)
		}
	}
}
