// Declaración del paquete main e imports.
package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/op/go-logging"
	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/common"
)

// Declaración de una variable global para el logger.
var log = logging.MustGetLogger("log")

// Carga la configuración de la app usando Viper.
// Lee primero variables de entorno (con prefijo CLI) y después el archivo config.yaml.
// Las variables de entorno tienen prioridad. 
// Valida que loop.period sea un tiempo válido.
// Devuelve la configuración lista para usar o un error.
func InitConfig() (*viper.Viper, error) {
	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvPrefix("cli")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	_ = v.BindEnv("id")
	_ = v.BindEnv("server", "address")
	_ = v.BindEnv("loop", "period")
	_ = v.BindEnv("loop", "amount")
	_ = v.BindEnv("log", "level")
	v.SetConfigFile("./config.yaml")
	if err := v.ReadInConfig(); err != nil {
		fmt.Printf("No se pudo leer la configuración del archivo de configuración. Se utilizan variables de entorno.")
	}

	if _, err := time.ParseDuration(v.GetString("loop.period")); err != nil {
		return nil, errors.Wrapf(err, "No se pudo analizar la variable de entorno CLI_LOOP_PERIOD como time.Duration.")
	}

	return v, nil
}

// Configura el sistema de logs según el nivel recibido (DEBUG, INFO, ERROR, etc).
// Si el nivel es válido lo aplica, si no devuelve un error.
func InitLogger(logLevel string) error {
	baseBackend := logging.NewLogBackend(os.Stdout, "", 0)
	format := logging.MustStringFormatter(
		`%{time:2006-01-02 15:04:05} %{level:.5s}     %{message}`,
	)

	backendFormatter := logging.NewBackendFormatter(baseBackend, format)
	backendLeveled := logging.AddModuleLevel(backendFormatter)
	logLevelCode, err := logging.LogLevel(logLevel)
	if err != nil {
		return err
	}

	backendLeveled.SetLevel(logLevelCode, "")
	logging.SetBackend(backendLeveled)
	return nil
}

// Muestra por log los valores actuales de configuración.
// Útil para depuración y ver qué parámetros se están usando.
func PrintConfig(v *viper.Viper) {
	log.Infof("action: config | result: success | client_id: %s | server_address: %s | loop_amount: %v | loop_period: %v | log_level: %s",
		v.GetString("id"),
		v.GetString("server.address"),
		v.GetInt("loop.amount"),
		v.GetDuration("loop.period"),
		v.GetString("log.level"),
	)
}

// Función main del programa.
func main() {
	v, err := InitConfig()
	if err != nil {
		log.Criticalf("%s", err)
	}

	if err := InitLogger(v.GetString("log.level")); err != nil {
		log.Criticalf("%s", err)
	}

	sigChan := make(chan os.Signal, 1)   // Modificación de código para manejar la señal pedida.
	signal.Notify(sigChan, syscall.SIGTERM)   // Modificación de código para manejar la señal pedida.

	PrintConfig(v)
	clientConfig := common.ClientConfig{
		ServerAddress: v.GetString("server.address"),
		ID:            v.GetString("id"),
		LoopAmount:    v.GetInt("loop.amount"),
		LoopPeriod:    v.GetDuration("loop.period"),
	}

	client := common.NewClient(clientConfig)

	client.StartClientLoop(sigChan)   // Modificación de código para manejar la señal pedida.
}
