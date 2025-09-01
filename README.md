# Trabajo Práctico 0 - Sistemas Distribuidos 1 (Roca)

## Ejercicio 1 - Introducción a Docker (1/4)

### Descripción

El objetivo de este ejercicio es crear un **script en Bash** llamado `generar-compose.sh`, que automatice la creación de un archivo `.yaml` con la definición de un entorno de **Docker Compose**.  

La definición generada incluirá:

- Un servicio **servidor**.
- Una cantidad **N de clientes**.
- Una **red de comunicación** compartida.

El script Bash actúa como envoltorio y delega la generación del archivo en el programa **`mi-generador.py`**, el cual produce la estructura final del YAML.
  
En los futuros ejercicios, este script podrá actualizarse o extenderse según nuevos requisitos.

### En caso de detectar errores de permisos a la hora de ejecutar el script

Antes de ejecutar el script, asegurarse de otorgar permisos de ejecución:

```bash

chmod +x generar-compose.sh

```

### Uso

La sintaxis general de ejecución es la siguiente:

```bash

./generar-compose.sh <archivo_salida> <cantidad_clientes>

```

Por ejemplo:

```bash

./generar-compose.sh docker-compose-dev.yaml 5

```

Este comando generará un archivo llamado docker-compose-dev.yaml con 1 servidor y 5 clientes conectados a la misma red.

### Levantar Docker luego de haber generado el archivo

Una vez generado el archivo .yaml, se puede desplegar el entorno completo utilizando Make:

```bash

make docker-compose-up

```

### Archivos modificados

#### Archivos creados

- generar-compose.sh
- mi-generador.py
