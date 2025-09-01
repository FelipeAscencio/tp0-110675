# Trabajo Práctico 0 - Sistemas Distribuidos 1 (Roca)

## Ejercicio 5 - Repaso de comunicaciones (1/3)

### Descripción

Se realizó la implementación necesaria para modelar el nuevo caso de uso planteado.

El 'Cliente' tiene los datos de la apuesta definidos en variables de entorno (Tal como se pide en el enunciado).

El 'Servidor' recibe la apuesta y la almacena en un archivo '.csv' usando la función correspondiente.

---

## Análisis de la completitud de los requisitos pedidos

Esta sección describe cómo la implementación actual cumple con los siguientes requisitos:

- **Definición del protocolo para el envío de mensajes**
- **Serialización y deserialización de los datos**
- **Separación correcta de responsabilidades (Modelo vs. Comunicación)**
- **Manejo de errores de sockets (short-read / short-write)**

### Protocolo de comunicación

#### Cliente → Servidor (Go ➜ Python)
**Esquema:** _length-prefixed_ (cabecera de 2 bytes, big-endian) + **payload UTF-8**

- **Header (2 bytes, big-endian):** indica el largo en bytes del payload.
  - Constantes: `TAMANIO_BYTES = 2` (Go) / `TAMANIO_HEADER = 2` (Python).
- **Payload:** cadena CSV de 6 campos (apuesta) **terminada en `\n`**.
- **Límite de tamaño:** `MAX_TAMANIO_MENSAJE = 8192` bytes (Go) para evitar mensajes excesivos.

**Implementación:**
- **Go (envío):** `utils.go → sendMessage`  
  - Serializa la apuesta (`sendBet`), calcula el tamaño, escribe **header** + **payload** usando `writeAll` (envío garantizado).
- **Python (recepción):** `utils.py → receive_message`  
  - Lee exactamente 2 bytes para el tamaño (big-endian) y luego **acumula** bytes hasta completar el tamaño declarado (evita short-read).
  - Decodifica UTF-8 y `strip()` (elimina `\r\n` si los hubiera).

#### Servidor → Cliente (Python ➜ Go)
**Esquema:** **línea terminada en `\n`** (ACK de la apuesta)

- **Payload:** `"{document},{number}\n"` en UTF-8.
- **Delimitación:** el fin de mensaje se identifica con `\n`.

**Implementación:**
- **Python (envío ACK):** `utils.py → acknowledge_bet`  
  - Usa `sendall` para garantizar el envío completo.
- **Go (recepción ACK):** `utils.go → receiveMessage`  
  - Lee con `ReadString('\n')` y luego `TrimSpace()`; si la línea llega incompleta, se recibe `err` y el caller lo maneja.

> Nota: El protocolo es **asimétrico a propósito**: **Cliente→Servidor** usa framing por longitud; **Servidor→Cliente** usa terminación por nueva línea. Ambos son estables y están implementados de forma segura.

### Serialización y Deserialización

#### Apuesta (Cliente→Servidor)
- **Formato:** CSV con 6 campos en este orden:  
  `agency,first_name,last_name,document,birthdate,number\n`
- **Serialización (Go):** `utils.go → sendBet`  
  - Arma el CSV (`fmt.Sprintf`) y lo envía mediante `sendMessage` (aplica header + validación de tamaño).
- **Deserialización (Python):** `utils.py → decode_bet`  
  - Recibe el mensaje, `split(',')`, **valida** que haya 6 campos (`CAMPOS_APUESTA_ESPERADA = 6`), y construye `Bet(...)`.
  - La clase `Bet` convierte tipos: `agency` y `number` a `int`, `birthdate` con `fromisoformat`.

#### ACK (Servidor→Cliente)
- **Formato:** `"{document},{number}\n"` (2 campos)
- **Serialización (Python):** `utils.py → acknowledge_bet` (con `sendall`)
- **Deserialización (Go):** `client.go → StartClientLoop`  
  - `receiveMessage` ➜ `strings.Split(",")` ➜ `strconv.Atoi` y comparación con lo enviado.

### Separación de responsabilidades

- **Modelo de dominio:**
  - **Python:** `utils.py → class Bet`
  - **Go:** `utils.go → type Bet`
  - Encapsulan los datos y (en Python) validan tipos al construir.

- **Capa de comunicación (I/O + protocolo):**
  - **Python:** `receive_message`, `decode_bet`, `acknowledge_bet` (en `utils.py`)
  - **Go:** `sendMessage`, `receiveMessage`, `sendBet` (en `utils.go`)

- **Orquestación / Aplicación:**
  - **Servidor (Python):** `server.py → Server.__handle_client_connection`
    - Lee apuesta ➜ loguea ➜ persiste (`store_bets`) ➜ responde ACK.
  - **Cliente (Go):** `client.go → StartClientLoop`
    - Conecta ➜ envía apuesta ➜ recibe ACK ➜ valida respuesta ➜ loguea resultado.

### Manejo robusto de errores ('Short-Read' y 'Short-Write')

#### Short-Read (lectura parcial)
- **Python (servidor):** `utils.py → receive_message`
  - Bucle `while len(informacion) < tamanio:` acumula datos hasta completar el payload.
  - Si `recv()` devuelve vacío (cierre inesperado), se lanza `ConnectionError("Connection closed unexpectedly")`.

- **Go (cliente):** `utils.go → receiveMessage`
  - `ReadString('\n')` bloquea hasta encontrar `\n` o `err`.  
  - Si el servidor cerrara antes de enviar `\n`, el error se propaga y se loguea en `client.go`.

#### Short-Write (escritura parcial)
- **Go (cliente):** `utils.go`  
  - **Se agrega `writeAll`** (bucle que reintenta hasta enviar todo el buffer).  
  - `sendMessage` usa `writeAll` para escribir **header** y **payload** (garantiza envío completo).

- **Python (servidor):** `utils.py → acknowledge_bet`  
  - **Se cambia a `sendall`** para garantizar que el ACK se envíe por completo.

#### Validaciones adicionales
- **Límite de tamaño:** `MAX_TAMANIO_MENSAJE = 8192` (Go) — Protege de payloads inesperadamente grandes (Requisito del ejercicio 6, pero se aprovecho para implementarlo prematuramente).
- **Validación de campos:** `decode_bet` comprueba `len == 6` antes de construir la `Bet`.
- **Logging consistente:** se registran `success/fail` con contexto (IP, mensaje, DNI, número), facilitando diagnóstico.

---

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

- mi-generador.py
- docker-compose-dev.yaml
- client/config.yaml
- client/main.go
- client/common/client.go
- server/config.ini
- server/common/server.py
- server/common/utils.py

#### Archivos creados

- client/common/utils.go
