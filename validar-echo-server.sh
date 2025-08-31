#!/usr/bin/env sh

# Mensaje esperado (Es el que se enviará al servidor).
MESSAGE="hello"

# Ejecuta contenedor Alpine temporal en la red tp0_testing_net.
# Envía el mensaje vía netcat (-w 2 = timeout de 2 segundos).
# Remueve retornos de carro y saltos de línea (\r\n) de la respuesta.
RESPONSE=$(
  docker run --rm --network=tp0_testing_net alpine sh -c "
    printf 'hello\n' | nc -w 2 server 12345
  " | tr -d '\r\n'
)

# Verifica si la respuesta coincide con el mensaje enviado.
if [ "$RESPONSE" = "$MESSAGE" ]; then
  echo "action: test_echo_server | result: success"
else
  echo "action: test_echo_server | result: fail"
fi
