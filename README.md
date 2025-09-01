# Trabajo Práctico 0 - Sistemas Distribuidos 1 (Roca)

## Ejercicio 4 - Introducción a Docker (4/4)

### Descripción

Se modificaron las implementaciones ya brindadas por la cátedra del 'Cliente' y el 'Servidor' para que finalicen de forma correcta al recibir la señal 'SIGTERM'.

Se incoporaron las bibliotecas necesarias para implementar el manejo de las señales y poder gestionar los cierres de forma 'graceful'.

### Uso

#### Generar la definición de Docker Compose:

```bash

./generar-compose.sh docker-compose-dev.yaml 5

```

#### Levantar los servicios:

```bash

make docker-compose-up

```

#### Comando para enviar la señal 'SIGTERM':

```bash

docker compose -f docker-compose-dev.yaml stop -t 10

```

### Archivos modificados

#### Archivos editados

- client/main.go
- client/common/client.go
- server/main.py
- server/common/server.py
