# Trabajo Práctico 0 - Sistemas Distribuidos 1 (Roca)

## Ejercicio 6 - Repaso de comunicaciones (2/3)

### Descripción

Se modificó la lógica de los 'Clientes' y el 'Servidor' para admitir varias apuestas a la vez (Con la limitación de '8kB', configuración 'Default') procesados por 'Batchs'.

Estas multiapuestas se leen a partir de los '.csv' brindados por la cátedra, siguiendo el contrado de funcionamiento establecido por el enunciado.

Mediante a el valor 'MaxAmount', que se encuentra dentro del archivo de configuración, se puede modificar la limitación de cantidad de apuestas.

El 'Servidor' se comporta igual que en el ejercicio anterior pero con las adaptaciones necesarias para procesar múltiples apuestas por cada 'Cliente'.

### Modificación del Protocolo

#### Cliente → Servidor (Go ➜ Python)

Se mantiene el framing con cabecera de 2 bytes (big-endian) que indica el tamaño del payload.

El payload ahora es un batch: Varias apuestas en una sola “unidad de mensaje”, con este formato:

Campos de cada apuesta (6): 'agency', 'first_name', 'last_name', 'document', 'birthdate' y 'number'.

Separador de campos: ','.

Separador entre apuestas del batch: ';'.

Ejemplo de payload: 'A1,John,Doe,12345678,1990-01-01,1111;A1,Jane,Doe,87654321,1992-02-02,2222'.

Cuando el cliente termina de enviar lotes, manda el mensaje de control FINISH (también con framing de 2 bytes).

#### Servidor → Cliente (Python ➜ Go)

Tras recibir 'FINISH', el servidor responde una sola vez con la cantidad total de apuestas recibidas en la sesión (ACK final), como línea terminada en '\n' (delimitado por newline).

#### Lectura de archivo

Al trabajar con archivos de tamaño pequeño, que no implican una problemática a la hora de cargarlos en memoria, se tomó la decisión de diseño de leerlos directamente, y luego trabajar los batchs con los límites definidos para facilitar la implementación, dando la verdadera restricción en el punto clave que es en la fase de envío de mensajes.

Para futuras versiones, si se decide trabajar con archivos de gran tamaño, debe modificarse esa implementación para que soporte el nuevo requerimiento.

### Uso

#### Generar la definición de Docker Compose:

```bash

./generar-compose.sh docker-compose-dev.yaml 5

```

#### Levantar los servicios:

```bash

make docker-compose-up  

```

### Archivos modificados

#### Archivos editados

- .gitignore
- docker-compose-dev.yaml
- mi-generador.py
- client/main.go
- client/common/client.go
- client/common/utils.go
- server/common/server.py
- server/common/utils.py
