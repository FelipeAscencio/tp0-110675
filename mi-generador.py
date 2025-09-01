# Import de las librerías necesarias.
import sys

# Constantes para índices de argumentos.
INDICE_NOMBRE_ARCHIVO = 1
INDICE_CANT_CLIENTES = 2

# Metadatos del stack.
PROJECT_NAME            = "tp0"

def generar_archivo(nombre_archivo, clientes):
    compose = f"""name: {PROJECT_NAME}
services:
  server:
    container_name: server
    image: server:latest
    entrypoint: python3 /main.py
    environment:
    - PYTHONUNBUFFERED=1
    - CLIENTS={clientes}
    networks:
    - testing_net
    volumes:
      - ./server/config.ini:/config.ini
"""
    
    for i in range(1, clientes + 1):
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
    
    with open(nombre_archivo, "w") as f:
        f.write(compose)
    
    print(f"Archivo {nombre_archivo} generado con {clientes} clientes.")

if __name__ == "__main__":
    nombre_archivo = sys.argv[INDICE_NOMBRE_ARCHIVO]
    clientes = int(sys.argv[INDICE_CANT_CLIENTES])
    generar_archivo(nombre_archivo, clientes)
