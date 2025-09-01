// Declaración del paquete common.
package common

// Importación de los paquetes necesarios.
import (
	"bufio"
	"fmt"
	"net"
	"os"   // Modificación de código para manejar la señal pedida.
	"time"
	"github.com/op/go-logging"
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
}

// Función que crea un nuevo cliente con la configuración dada.
func NewClient(config ClientConfig) *Client {
	client := &Client{
		config: config,
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
func (c *Client) StartClientLoop(sigChan chan os.Signal) {   // Modificación de código para manejar la señal pedida.
	for msgID := 1; msgID <= c.config.LoopAmount; msgID++ {
		select {   // Modificación de código para manejar la señal pedida.
		case <-sigChan:   // Modificación de código para manejar la señal pedida.
			log.Infof("action: shutdown | result: success")   // Modificación de código para manejar la señal pedida.
			return   // Modificación de código para manejar la señal pedida.
		default:   // Modificación de código para manejar la señal pedida.
			c.createClientSocket()
			fmt.Fprintf(
				c.conn,
				"[CLIENT %v] Message N°%v\n",
				c.config.ID,
				msgID,
			)

			msg, err := bufio.NewReader(c.conn).ReadString('\n')
			c.conn.Close()
			if err != nil {
				log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
					c.config.ID,
					err,
				)

				return
			}

			log.Infof("action: receive_message | result: success | client_id: %v | msg: %v",
				c.config.ID,
				msg,
			)

			time.Sleep(c.config.LoopPeriod)
		}
	}

	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}
