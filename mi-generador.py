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
SERVER_CONFIG           = "./server/config.ini:/config.ini"

# Servicio Cliente.
CLIENT_IMAGE            = "client:latest"
CLIENT_PREFIX           = "client"
CLIENT_ENTRYPOINT       = "/client"
ENV_CLI_ID              = "CLI_ID"
CLIENT_NOMBRE           = "Santiago Lionel"
CLIENT_APELLIDO         = "Lorca"
CLIENT_DOCUMENTO        = "30904465"
CLIENT_NACIMIENTO       = "1999-03-17"
CLIENT_NUMERO           = "7574"
CLIENT_CONFIG           = "./client/config.yaml:/config.yaml"

# Red.
NETWORK_NAME            = "testing_net"
NETWORK_DRIVER          = "default"
NETWORK_SUBNET          = "172.25.125.0/24"

def generate_file(filename, clients):
    compose = f"""name: tp0
services:
  server:
    container_name: server
    image: server:latest
    entrypoint: python3 /main.py
    environment:
    - PYTHONUNBUFFERED=1
    - CLIENTS={clients}
    networks:
    - testing_net
    volumes:
      - ./server/config.ini:/config.ini
"""
    
    for i in range(1, clients + 1):
        compose += f"""
  client{i}:
    container_name: client{i}
    image: client:latest
    entrypoint: /client
    environment:
      - CLI_ID={i}
    networks:
      - testing_net
    volumes:
      - ./client/config.yaml:/config.yaml
      - ./.data/agency-{i}.csv:/agency.csv
    depends_on:
      - server
"""
    
    compose += """
networks:
  testing_net:
    ipam:
      driver: default
      config:
        - subnet: 172.25.125.0/24
"""
    
    with open(filename, "w") as f:
        f.write(compose)
    
    print(f"Archivo {filename} generado con {clients} clientes.")

if __name__ == "__main__":
    filename = sys.argv[1]
    clients = int(sys.argv[2])
    
    generate_file(filename, clients)
