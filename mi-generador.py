# Import de las librerías necesarias.
import sys

# Constantes para índices de argumentos.
INDICE_NOMBRE_ARCHIVO = 1
INDICE_CANT_CLIENTES = 2

# Metadatos del stack.
PROJECT_NAME            = "tp0"

# Servicio Servidor.
SERVER_NAME             = "server"
SERVER_IMAGE            = "server:latest"
SERVER_ENTRYPOINT       = "python3 /main.py"
ENV_PYTHON_UNBUFFERED   = "PYTHONUNBUFFERED=1"
CLIENT_PREFIX           = "client"
SERVER_CONFIG           = "./server/config.ini:/config.ini"

# Servicio Cliente.
CLIENT_IMAGE            = "client:latest"
CLIENT_ENTRYPOINT       = "/client"
ENV_CLI_ID              = "CLI_ID"
CLIENT_NOMBRE           = "Santiago Lionel"
CLIENT_APELLIDO         = "Lorca"
CLIENT_DOCUMENTO        = "30904465"
CLIENT_NACIMIENTO       = "1999-03-17"
CLIENT_NUMERO           = "7574"

# Red.
NETWORK_NAME            = "testing_net"
NETWORK_DRIVER          = "default"
NETWORK_SUBNET          = "172.25.125.0/24"
CLIENT_CONFIG           = "./client/config.yaml:/config.yaml"

# Función para generar el archivo docker-compose.
def generate_file(filename, clients):
    compose = f"""name: {PROJECT_NAME}
services:
  {SERVER_NAME}:
    container_name: {SERVER_NAME}
    image: {SERVER_IMAGE}
    entrypoint: {SERVER_ENTRYPOINT}
    environment:
    - {ENV_PYTHON_UNBUFFERED}
    networks:
    - {NETWORK_NAME}
    volumes:
      - {SERVER_CONFIG}
"""
    
    for i in range(1, clients + 1):
        compose += f"""
  {CLIENT_PREFIX}{i}:
    container_name: {CLIENT_PREFIX}{i}
    image: {CLIENT_IMAGE}
    entrypoint: {CLIENT_ENTRYPOINT}
    environment:
      - {ENV_CLI_ID}={i}
      - NOMBRE={CLIENT_NOMBRE}
      - APELLIDO={CLIENT_APELLIDO}
      - DOCUMENTO={CLIENT_DOCUMENTO}
      - NACIMIENTO={CLIENT_NACIMIENTO}
      - NUMERO={CLIENT_NUMERO}
    networks:
      - {NETWORK_NAME}
    volumes:
      - {CLIENT_CONFIG}
    depends_on:
      - {SERVER_NAME}
"""
    
    compose += f"""
networks:
  {NETWORK_NAME}:
    ipam:
      driver: {NETWORK_DRIVER}
      config:
        - subnet: {NETWORK_SUBNET}
"""
    
    with open(filename, "w") as f:
        f.write(compose)

    print(f"Archivo {filename} generado con {clients} clientes.")

# Punto de entrada del script.
if __name__ == "__main__":
    filename = sys.argv[INDICE_NOMBRE_ARCHIVO]
    clients = int(sys.argv[INDICE_CANT_CLIENTES])
    generate_file(filename, clients)
