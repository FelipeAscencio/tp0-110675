// Declaración del paquete main.
package main

// Importación de los paquetes necesarios.
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
func InitConfig() (*viper.Viper, *common.Bet, error) {
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
		fmt.Printf("Configuration could not be read from config file. Using env variables instead")
	}

	if _, err := time.ParseDuration(v.GetString("loop.period")); err != nil {
		return nil, nil, errors.Wrapf(err, "Could not parse CLI_LOOP_PERIOD env var as time.Duration.")
	}

	// Leer las variables de entorno para la apuesta (Modificación de código para el ejercicio 5).
	nombre_apuesta := os.Getenv("NOMBRE")
	apellido_apuesta := os.Getenv("APELLIDO")
	documento_apuesta := os.Getenv("DOCUMENTO")
	nacimiento_apuesta := os.Getenv("NACIMIENTO")
	numero_apuesta := os.Getenv("NUMERO")
	bet := &common.Bet{
		AgencyId:  v.GetString("id"),
		Name:      nombre_apuesta,
		LastName:  apellido_apuesta,
		Document:  documento_apuesta,
		BirthDate: nacimiento_apuesta,
		Number:    numero_apuesta,
	}

	return v, bet, nil
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
	v, bet, err := InitConfig()
	if err != nil {
		log.Criticalf("%s", err)
	}

	if err := InitLogger(v.GetString("log.level")); err != nil {
		log.Criticalf("%s", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM)
	PrintConfig(v)
	clientConfig := common.ClientConfig{
		ServerAddress: v.GetString("server.address"),
		ID:            v.GetString("id"),
		LoopAmount:    v.GetInt("loop.amount"),
		LoopPeriod:    v.GetDuration("loop.period"),
	}

	client := common.NewClient(clientConfig, *bet)
	client.StartClientLoop(sigChan)
}
