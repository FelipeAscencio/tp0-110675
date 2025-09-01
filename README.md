# Trabajo Práctico 0 - Sistemas Distribuidos 1 (Roca)

## Ejercicio 3 - Introducción a Docker (3/4)

### Descripción

Para este ejercicio se creó un script de shell llamado 'validar-echo-server.sh' que verifica la comunicación con el Echo Server utilizando netcat dentro de un contenedor temporal (no instala nada en la máquina host y no expone puertos).

Funcionamiento del script:
- Levanta un contenedor Alpine efímero en la misma red Docker que el servidor (tp0_testing_net).
- Envía el mensaje "hello" al servicio server en el puerto '12345'.
- printf 'hello\n' envía el mensaje con salto de línea.
- '-w 2' agrega timeout para evitar bloqueos si no hay respuesta.
- Normaliza la respuesta (elimina \r\n) y la compara con el mensaje enviado.
- Imprime:
    - 'action: test_echo_server | result: success' si coincide.
    - 'action: test_echo_server | result: fail' si no coincide.

Nota: Todo ocurre dentro de contenedores y en la red de Docker (tp0_testing_net), cumpliendo el requisito de no instalar netcat en el host y no exponer puertos.

### En caso de detectar errores de permisos a la hora de ejecutar el script

```bash

chmod +x validar-echo-server.sh

```

### Uso

#### Generar la definición de Docker Compose:

```bash

./generar-compose.sh docker-compose-dev.yaml 5

```

#### Levantar los servicios:

```bash

make docker-compose-up

```

#### Ejecutar el validador:

```bash

./validar-echo-server.sh

```

#### Salida esperada:

action: test_echo_server | result: success

### Archivos modificados

#### Archivos creados

- validar-echo-server.sh
