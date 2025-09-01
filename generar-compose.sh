#!/bin/bash

# Habilita opciones para mejorar la seguridad y manejo de errores en el script.
set -euo pipefail
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Imprime en pantalla el primer argumento recibido.
echo "Nombre del archivo de salida: $1"

# Imprime en pantalla el segundo argumento recibido.
echo "Cantidad de clientes: $2"

# Ejecuta el script de Python generado pasando como parámetros el archivo y la cantidad de clientes.
python3 mi-generador.py $1 $2
