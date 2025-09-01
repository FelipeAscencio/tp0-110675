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
	INDEX_NUMERO = 1
	DELAY = 100
)

// Declaración de una variable global para el logger.
var log = logging.MustGetLogger("log")

// Estructura que contiene la configuración del cliente.
type ClientConfig struct {
	ID             string
	ServerAddress  string
	LoopAmount     int
	LoopPeriod     time.Duration
	BatchMaxAmount int
}

// Estructura que representa al cliente.
type Client struct {
	config ClientConfig
	conn   net.Conn
	bets   []Bet
}

// Función que crea un nuevo cliente con la configuración dada.
func NewClient(config ClientConfig, bets []Bet) *Client {
	client := &Client{
		config: config,
		bets:   bets,
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
	_ = c.createClientSocket()
	bets_totales := len(c.bets)
	contador_bets := 0

// Loop para enviar las apuestas en lotes.
loop:
	for i := 0; i < bets_totales; i += c.config.BatchMaxAmount {
		final := i + c.config.BatchMaxAmount
		if final > bets_totales {
			final = bets_totales
		}

		batch := c.bets[i:final]
		contador_bets += len(batch)
		err := sendBetBatch(c.conn, batch, contador_bets)
		if err != nil {
			return
		}

		select {
		case <-sigChan:
			log.Infof("action: shutdown | result: success")
			break loop
		default:
		}
	}

	log.Infof("action: apuesta_enviada | result: success | cantidad: %v", contador_bets)
	log.Infof("action: finalizar_envio | result: in_progress")
	err := sendFinishMessage(c.conn)
	if err != nil {
		log.Errorf("action: finalizar_envio | result: fail | error: %v",
			c.config.ID,
			err,
		)

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

	contador_respuestas, _ := strconv.Atoi(strings.TrimSpace(mensaje))
	if contador_respuestas != bets_totales {
		log.Errorf("action: finalizar_envio | result: fail | msg: %v | error: unexpected message",
			mensaje,
		)

		return
	} else {
		log.Infof("action: finalizar_envio | result: success")
		_ = c.conn.Close()
		time.Sleep(DELAY * time.Millisecond)
	}
}
