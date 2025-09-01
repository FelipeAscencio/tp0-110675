# Import de las librerías necesarias.
import sys

# Constantes para índices de argumentos.
INDICE_NOMBRE_ARCHIVO = 1
INDICE_CANT_CLIENTES = 2

# Metadatos del stack.
PROJECT_NAME         = "tp0"
SERVICE_SERVER_NAME  = "server"
SERVICE_CLIENT_PREF  = "client"
NETWORK_NAME         = "testing_net"

# Imágenes / entrypoints.
SERVER_IMAGE         = "server:latest"
CLIENT_IMAGE         = "client:latest"
SERVER_ENTRYPOINT    = "python3 /main.py"
CLIENT_ENTRYPOINT    = "/client"

# Variables de entorno.
ENV_PYTHON_UNBUF_KV  = "PYTHONUNBUFFERED=1"
ENV_CLIENTS_KEY      = "CLIENTS"
ENV_CLI_ID_KEY       = "CLI_ID"

# Volúmenes.
SERVER_CONFIG_VOL    = "./server/config.ini:/config.ini"
CLIENT_CONFIG_VOL    = "./client/config.yaml:/config.yaml"
CLIENT_AGENCY_VOL_TMPL = "./.data/agency-{i}.csv:/agency.csv"

# Red.
NETWORK_DRIVER       = "default"
NETWORK_SUBNET       = "172.25.125.0/24"

# Función para generar el archivo docker-compose.
def generar_archivo(nombre_archivo, clientes):
    compose = (
        f"name: {PROJECT_NAME}\n"
        f"services:\n"
        f"  {SERVICE_SERVER_NAME}:\n"
        f"    container_name: {SERVICE_SERVER_NAME}\n"
        f"    image: {SERVER_IMAGE}\n"
        f"    entrypoint: {SERVER_ENTRYPOINT}\n"
        f"    environment:\n"
        f"    - {ENV_PYTHON_UNBUF_KV}\n"
        f"    - {ENV_CLIENTS_KEY}={clientes}\n"
        f"    networks:\n"
        f"    - {NETWORK_NAME}\n"
        f"    volumes:\n"
        f"      - {SERVER_CONFIG_VOL}\n"
    )

    for i in range(1, clientes + 1):
        compose += (
            f"\n"
            f"  {SERVICE_CLIENT_PREF}{i}:\n"
            f"    container_name: {SERVICE_CLIENT_PREF}{i}\n"
            f"    image: {CLIENT_IMAGE}\n"
            f"    entrypoint: {CLIENT_ENTRYPOINT}\n"
            f"    environment:\n"
            f"      - {ENV_CLI_ID_KEY}={i}\n"
            f"    networks:\n"
            f"      - {NETWORK_NAME}\n"
            f"    volumes:\n"
            f"      - {CLIENT_CONFIG_VOL}\n"
            f"      - {CLIENT_AGENCY_VOL_TMPL.format(i=i)}\n"
            f"    depends_on:\n"
            f"      - {SERVICE_SERVER_NAME}\n"
        )

    compose += (
        f"\nnetworks:\n"
        f"  {NETWORK_NAME}:\n"
        f"    ipam:\n"
        f"      driver: {NETWORK_DRIVER}\n"
        f"      config:\n"
        f"        - subnet: {NETWORK_SUBNET}\n"
    )

    with open(nombre_archivo, "w") as f:
        f.write(compose)

    print(f"Archivo {nombre_archivo} generado con {clientes} clientes.")

# Punto de entrada del script.
if __name__ == "__main__":
    nombre_archivo = sys.argv[INDICE_NOMBRE_ARCHIVO]
    clientes = int(sys.argv[INDICE_CANT_CLIENTES])
    generar_archivo(nombre_archivo, clientes)
