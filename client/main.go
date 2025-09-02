// Declaración del paquete main.
package main

// Importación de los paquetes necesarios.
import (
	"encoding/csv"
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
	_ = v.BindEnv("batch", "maxAmount")
	v.SetConfigFile("./config.yaml")
	if err := v.ReadInConfig(); err != nil {
		fmt.Printf("Configuration could not be read from config file. Using env variables instead")
	}

	if _, err := time.ParseDuration(v.GetString("loop.period")); err != nil {
		return nil, errors.Wrapf(err, "Could not parse CLI_LOOP_PERIOD env var as time.Duration.")
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

// Añadido para la resolución del ejercicio 6.
// Lee el archivo agency.csv y devuelve un slice de apuestas (common.Bet).
func GetAgencyData(agencyId string) ([]common.Bet, error) {
	archivo, err := os.Open("agency.csv")
	if err != nil {
		log.Errorf("action: read_file | result: fail | error: %v",
			err,
		)

		return nil, err
	}

	lector := csv.NewReader(archivo)
	lineas, err := lector.ReadAll()
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil, err
	}

	var bets []common.Bet
	for _, linea := range lineas {
		bet := common.Bet{
			AgencyId:  agencyId,
			Name:      linea[0],
			LastName:  linea[1],
			Document:  linea[2],
			BirthDate: linea[3],
			Number:    linea[4],
		}

		bets = append(bets, bet)
	}

	return bets, nil
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

	// Agregado para el ejercicio 6.
	agencyId := v.GetString("id")
	bets, err := GetAgencyData(agencyId)
	if err != nil {
		log.Criticalf("%s", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM)
	PrintConfig(v)
	clientConfig := common.ClientConfig{
		ServerAddress:  v.GetString("server.address"),
		ID:             v.GetString("id"),
		LoopAmount:     v.GetInt("loop.amount"),
		LoopPeriod:     v.GetDuration("loop.period"),
		BatchMaxAmount: v.GetInt("batch.maxAmount"),
	}

	client := common.NewClient(clientConfig, bets)
	client.StartClientLoop(sigChan)
}
