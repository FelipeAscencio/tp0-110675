# Trabajo Práctico 0 - Sistemas Distribuidos 1 (Roca)

## Ejercicio 2 - Introducción a Docker (2/4)

### Descripción

El objetivo de este ejercicio es modificar el funcionamiento del 'Cliente' y el 'Servidor' para persistir los archivos de configuración fuera de la imagen.

Para conseguir este objetivo se agregaron volumenes dentro de los archivos 'docker-compose-dev.yaml' y 'mi-generador.py'.

### En caso de detectar errores de permisos a la hora de ejecutar el script

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

#### Archivos modificados

- mi-generador.py.
- docker-compose-dev.yaml.
- client/Dockerfile.
