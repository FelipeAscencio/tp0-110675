# Trabajo Práctico 0 - Sistemas Distribuidos 1 (Roca)

## Ejercicio 7 - Repaso de comunicaciones (3/3)

### Descripción

Se adaptó el sistema para que, una vez que cada **Cliente** termina de enviar sus apuestas, pueda consultar directamente los resultados del sorteo.

El **Server**, al recibir la señal de cierre (*FINISH*) de todos los clientes, queda habilitado para ejecutar el sorteo mediante las funciones `load_bets` y `has_won`.

El sorteo se ejecuta una única vez y su resultado queda almacenado en memoria dentro del servidor, evitando recalcular en cada consulta.

De este modo, múltiples clientes pueden solicitar los resultados sin que se repita el proceso.

### Cambios en el protocolo

#### Nuevos mensajes de control

- **BET** → Enviado por el cliente al iniciar el envío de apuestas, para que el servidor se prepare a recibirlas.  
- **RESULTS** → Enviado por el cliente para pedir los resultados. El mensaje incluye el ID del cliente/agencia para que el servidor sepa a qué apuestas filtrar.

#### Flujo de resultados

- Una vez que **todos los clientes** han enviado *FINISH*, el primer mensaje *RESULTS* recibido provoca la ejecución del sorteo.  
- El servidor responde al cliente con la lista de documentos de los ganadores de esa agencia, en el formato: 'DOCUMENTO1,DOCUMENTO2,DOCUMENTO3,DOCUMENTO4'.
- Si la agencia no tuvo ganadores, el servidor devuelve el mensaje **NOWINNERS**.

### Comportamiento antes del cierre total

Si el servidor aún no recibió *FINISH* de todos los clientes, ignora cualquier pedido de *RESULTS*.

En ese caso, los clientes deben volver a intentar la consulta aplicando una estrategia de **reintento con espera incremental** (*sleep* creciente en cada ciclo).

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

- docker-compose-dev.yaml
- mi-generador.py
- client/config.yaml
- client/common/client.go
- client/common/utils.go
- server/main.py
- server/common/server.py
- server/common/utils.py
