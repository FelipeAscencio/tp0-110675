# Import de las librerías necesarias.
import sys

# Constantes para índices de argumentos.
INDICE_NOMBRE_ARCHIVO = 1
INDICE_CANT_CLIENTES = 2

# Metadatos del stack.
NOMBRE_STACK = "tp0"

# Servicio servidor.
SERVIDOR_NOMBRE_SERVICIO = "server"
SERVIDOR_NOMBRE_CONTENEDOR = "server"
SERVIDOR_IMAGEN = "server:latest"
SERVIDOR_ENTRYPOINT = "python3 /main.py"
SERVIDOR_ENV_PYTHONUNBUFFERED = "1"
VOLUMEN_SERVIDOR = "./server/config.ini:/config.ini"

# Servicio cliente.
CLIENTE_IMAGEN = "client:latest"
CLIENTE_ENTRYPOINT = "/client"
VOLUMEN_CLIENTE  = "./client/config.yaml:/config.yaml"

# Constantes generales de cliente
CLIENTE_NOMBRE = "Santiago Lionel"
CLIENTE_APELLIDO = "Lorca"
CLIENTE_DOCUMENTO = "30904465"
CLIENTE_FECHA_NACIMIENTO = "1999-03-17"
CLIENTE_NUMERO = "7574"

# Red.
NOMBRE_RED = "testing_net"
SUBRED_CIDR = "172.25.125.0/24"

# Función para generar el archivo docker-compose.
def generar_archivo(nombre_archivo: str, cantidad_clientes: int) -> None:
    """
    Genera un archivo docker-compose con un servicio servidor y N clientes.
    """

    # Bloque inicial (name + servicio server).
    compose = f"""name: {NOMBRE_STACK}
services:
  {SERVIDOR_NOMBRE_SERVICIO}:
    container_name: {SERVIDOR_NOMBRE_CONTENEDOR}
    image: {SERVIDOR_IMAGEN}
    entrypoint: {SERVIDOR_ENTRYPOINT}
    environment:
    - PYTHONUNBUFFERED={SERVIDOR_ENV_PYTHONUNBUFFERED}
    networks:
    - {NOMBRE_RED}
    volumes:
      - {VOLUMEN_SERVIDOR}
"""

    # Bloques de cada cliente.
    for i in range(1, cantidad_clientes + 1):
        compose += f"""
  client{i}:
    container_name: client{i}
    image: {CLIENTE_IMAGEN}
    entrypoint: {CLIENTE_ENTRYPOINT}
    environment:
      - CLI_ID={i}
      - NOMBRE={CLIENTE_NOMBRE}
      - APELLIDO={CLIENTE_APELLIDO}
      - DOCUMENTO={CLIENTE_DOCUMENTO}
      - NACIMIENTO={CLIENTE_FECHA_NACIMIENTO}
      - NUMERO={CLIENTE_NUMERO}
    networks:
      - {NOMBRE_RED}
    volumes:
      - {VOLUMEN_CLIENTE}
    depends_on:
      - {SERVIDOR_NOMBRE_SERVICIO}
"""

    # Bloque de redes.
    compose += f"""
networks:
  {NOMBRE_RED}:
    ipam:
      driver: default
      config:
        - subnet: {SUBRED_CIDR}
"""

    # Escritura del archivo.
    with open(nombre_archivo, "w") as f:
        f.write(compose)

    # Mensaje final.
    print(f"Archivo {nombre_archivo} generado con {cantidad_clientes} clientes.")

# Punto de entrada del programa.
if __name__ == "__main__":
    nombre_archivo = sys.argv[INDICE_NOMBRE_ARCHIVO]
    cant_clientes = int(sys.argv[INDICE_CANT_CLIENTES])
    generar_archivo(nombre_archivo, cant_clientes)
