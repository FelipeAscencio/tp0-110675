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
	INDEX_NUMERO    = 1
	DELAY_CONSULTA  = 100
	DELAY_PAUSA     = 1000
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
	apuestas_totales := len(c.bets)
	contador_apuestas := 0

	// Modificación para la lógica del ejercicio 7.
	log.Infof("action: comenzar_envio | result: in_progress")
	err := sendStartBets(c.conn)
	if err != nil {
		log.Errorf("action: comenzar_envio | result: fail | error: %v",
			c.config.ID,
			err,
		)

		return
	}

	log.Infof("action: comenzar_envio | result: success")

	// Loop para enviar las apuestas en lotes y luego consultar los ganadores
loop1:
	for i := 0; i < apuestas_totales; i += c.config.BatchMaxAmount {
		final := i + c.config.BatchMaxAmount
		if final > apuestas_totales {
			final = apuestas_totales
		}

		batch := c.bets[i:final]
		contador_apuestas += len(batch)
		err := sendBetBatch(c.conn, batch, contador_apuestas)
		if err != nil {
			return
		}

		select {
		case <-sigChan:
			log.Infof("action: shutdown | result: success")
			break loop1
		default:
		}
	}

	log.Infof("action: apuesta_enviada | result: success | cantidad: %v", contador_apuestas)
	log.Infof("action: finalizar_envio | result: in_progress")
	err = sendFinishMessage(c.conn)
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
	if contador_respuestas != apuestas_totales {
		log.Errorf("action: finalizar_envio | result: fail | msg: %v | error: response does not match",
			mensaje,
		)

		return
	}

	log.Infof("action: finalizar_envio | result: success")
	_ = c.conn.Close()
	multiplicador_de_pausa := 0
	for {
		select {
		case <-sigChan:
			log.Infof("action: shutdown | result: success")
			return
		default:
			log.Infof("action: consulta_ganadores | result: in_progress | attempt: %v", multiplicador_de_pausa)
			_ = c.createClientSocket()
			resultados, err := sendAskForResults(c.conn, c.config.ID)
			if err == nil {
				log.Infof("action: consulta_ganadores | result: success | cant_ganadores: %v", len(resultados))
				_ = c.conn.Close()
				time.Sleep(DELAY_CONSULTA * time.Millisecond)
				return
			} else {
				_ = c.conn.Close()
				multiplicador_de_pausa++
				time.Sleep(time.Duration(multiplicador_de_pausa*DELAY_PAUSA) * time.Millisecond)
			}
		}
	}
}
