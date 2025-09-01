#!/usr/bin/env python3

# Importación de los paquetes necesarios.
import sys
import signal
from configparser import ConfigParser
from common.server import Server
import logging
import os

# Variable global que contiene la instancia del servidor.
global server
server = None

# Función que maneja la señal de terminación para apagar el servidor de manera ordenada.
def graceful_shutdown(signum=None, frame=None):
    global server
    server.shutdown()
    sys.exit(0)

# Función que inicializa la configuración del servidor.
def initialize_config():
    """ Carga la configuración del programa.
    Busca primero en variables de entorno y luego en un archivo de configuración.
    - Si falta algún parámetro → lanza KeyError.
    - Si un valor no se puede parsear → lanza ValueError.
    - Si todo está bien → devuelve un objeto ConfigParser con los parámetros.
    """

    config = ConfigParser(os.environ)

    # Si config.ini no existe, el objeto de configuración original no se modifica.
    config.read("config.ini")
    config_params = {}
    try:
        config_params["port"] = int(os.getenv('SERVER_PORT', config["DEFAULT"]["SERVER_PORT"]))
        config_params["listen_backlog"] = int(os.getenv('SERVER_LISTEN_BACKLOG', config["DEFAULT"]["SERVER_LISTEN_BACKLOG"]))
        config_params["logging_level"] = os.getenv('LOGGING_LEVEL', config["DEFAULT"]["LOGGING_LEVEL"])
    except KeyError as e:
        raise KeyError("Key was not found. Error: {} .Aborting server".format(e))
    except ValueError as e:
        raise ValueError("Key could not be parsed. Error: {}. Aborting server".format(e))

    return config_params

# Principal función que inicializa el servidor y comienza su ejecución.
def main():
    global server
    config_params = initialize_config()
    logging_level = config_params["logging_level"]
    port = config_params["port"]
    listen_backlog = config_params["listen_backlog"]
    initialize_log(logging_level)

    # Registrar los parámetros de configuración al inicio del programa para verificar la configuración del componente.
    logging.debug(f"action: config | result: success | port: {port} | "
                  f"listen_backlog: {listen_backlog} | logging_level: {logging_level}")

    # Inicializar el servidor e iniciar el bucle del servidor.
    server = Server(port, listen_backlog)
    signal.signal(signal.SIGTERM, graceful_shutdown)
    server.run()

# Inicialización del logging de Python.
def initialize_log(logging_level):
    """Inicializa el sistema de logs.
    Agrega la marca de tiempo para identificar en los logs de Docker
    la fecha y hora en que llegó cada mensaje.
    """
    
    logging.basicConfig(
        format='%(asctime)s %(levelname)-8s %(message)s',
        level=logging_level,
        datefmt='%Y-%m-%d %H:%M:%S',
    )

# Punto de entrada del programa.
if __name__ == "__main__":
    main()
