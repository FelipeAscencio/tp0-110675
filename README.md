# Trabajo Práctico 0 - Sistemas Distribuidos 1 (Roca)

## Ejercicio 8 - Repaso de Concurrencia (1/1)

### Descripción

Modificar el servidor para **aceptar conexiones y procesar mensajes en paralelo**.

La implementación utiliza **multithreading** en Python, respetando las limitaciones del lenguaje (GIL) y asegurando **seguridad de concurrencia** mediante **locks**.

### Explicación de la implementación

Se actualizó el servidor (`server/common/server.py`) para que:
- **Acepte múltiples clientes en paralelo**: Cada conexión se atiende en un **thread dedicado**.
- **Sincronice el acceso** a recursos compartidos con **candados (locks)**:
  - Un lock protege el contador `clientes_terminados`.
  - Otro lock protege las operaciones críticas de **persistencia** (`store_bets`).
- **Calcule resultados** una única vez y los **cachee** en memoria (`self.winners`), evitando recomputar.
- **Apague ordenadamente** ante `SIGTERM` (graceful shutdown): Cierra el socket de escucha, espera la finalización de todos los hilos y registra el resultado.


### Cambios Clave en `server.py`

- **Multithreading**: 
  - Se importa `threading` y se crea un **thread por conexión** con `target=self.__handle_client_connection`.
  - Los hilos se marcan como `daemon = True` para no bloquear la salida del proceso y se **joinean** al final para orden.
- **Locks**:
  - `self.candado_mutex_1` protege `self.clientes_terminados` y la sección donde se decide si ya finalizaron todos los clientes (antes de enviar resultados).
  - `self.candado_mutex_2` protege la **persistencia** de apuestas (`utils.store_bets`) y el contador local por conexión.
- **Caché de resultados**:
  - `self.winners` se calcula **una sola vez** cuando todos los clientes terminaron, y luego se reutiliza para responder `RESULTS`.
- **Manejo de señales**:
  - Se instala un handler para `SIGTERM` con `signal.signal(signal.SIGTERM, self.shutdown)`.
  - `shutdown()` cambia `self.running`, cierra el socket y espera a los hilos (`join`), dejando el sistema en estado coherente.
- **Manejo de errores**:
  - Ante errores de E/S, se loguea y se envía un mensaje `ERROR` al cliente, cerrando la conexión en `finally`.

### Flujo de Ejecución

1. **Arranque**:
   - `Server.run()` crea el socket, lo `bind`ea y `listen`.
   - En un loop mientras `self.running`, acepta conexiones (`__accept_new_connection`).

2. **Aceptación de conexiones**:
   - Por cada conexión entrante se crea un **hilo** (`threading.Thread`) que ejecuta `__handle_client_connection`.

3. **Protocolo de mensajes**:
   - Primer mensaje define el modo:
     - `BET`: El cliente enviará **batches** de apuestas.
       - Por cada batch (`utils.decode_bets`), se protege con `candado_mutex_2` la llamada a `utils.store_bets`.
       - Al recibir `FINISH`, se hace `ack` con `utils.acknowledge_bets` y, bajo `candado_mutex_1`, se incrementa `self.clientes_terminados`.
     - `RESULTS`: El cliente solicita resultados por `agency_id`.
       - Bajo `candado_mutex_1`, si `clientes_terminados == self.clients`:
         - Si `self.winners` está vacío, se cargan apuestas (`utils.load_bets`), se filtran ganadores (`utils.has_won`) y se cachean.
         - Se filtra por `agency_id` y se responde con `utils.send_results`.

4. **Apagado ordenado**:
   - `shutdown()` cierra el socket de escucha y **espera** a que todos los hilos terminen (`__esperar_otros_hilos`).

### Concurrencia y Sincronización

- **Modelo**: *Thread-per-connection*.
- **Secciones críticas**:
  - **Contador global** `clientes_terminados`: Protegido por `candado_mutex_1`.
  - **Persistencia / estado compartido** (`store_bets` y el contador/ensamblado de batches): Protegido por `candado_mutex_2`.
- **Resultados cacheados**:
  - La primera vez que se solicitan y **todos los clientes terminaron**, se computan y guardan en `self.winners`.
  - Este patrón evita condiciones de carrera y computaciones repetidas.
- **GIL y Python**:
  - Aunque el GIL impide paralelismo puro de CPU-bound, este servidor es **I/O-bound** (sockets), por lo que **threading** es adecuado: Mientras un hilo bloquea en I/O, otros pueden avanzar.
  - Por otro lado, a la hora de la implementación se asumió el tope de 5 clientes indicado en el enunciado. Si en un futuro se quiere trabajar con una cantidad considerablemente mas grande, se deberá cambiar la resolución para evitar vulnerabilidades ante ataques del tipo 'DDoS' (Una variante recomendada es la utilización de un 'Threadpool').
- **Invariantes**:
  - `0 ≤ clientes_terminados ≤ clients`.
  - `winners` solo se define cuando `clientes_terminados == clients`.
  - No hay escritura concurrente a la persistencia de apuestas sin lock.

### Interfaz y Mensajes

- **Mensajes soportados**:
  - `BET`: Recepción de apuestas en lotes, termina con `FINISH`.
  - `RESULTS`: Solicitud de resultados para una `agency_id`.
  - `ERROR`: Si ocurre un fallo interno, el servidor responde con `ERROR` y cierra.

### Manejo de Errores

- **Lectura/Escritura**:
  - Excepciones `OSError` son capturadas, se loguea y se responde `ERROR`.
- **Persistencia**:
  - `store_bets` se protege con lock para evitar corrupción/condiciones de carrera.
- **Cierre de conexiones**:
  - Los sockets de cliente se **cierran en `finally`**.

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

- server/common/server.py
