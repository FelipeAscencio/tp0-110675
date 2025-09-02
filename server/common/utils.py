# Importación de los paquetes necesarios.
import csv
import datetime
import time
import logging

# Constantes de configuración.
""" Bets storage location. """
STORAGE_FILEPATH = "./bets.csv"
""" Simulated winner number in the lottery contest. """
LOTTERY_WINNER_NUMBER = 7574

# Constantes de tamaños y límites.
TAMANIO_HEADER = 2
CAMPOS_APUESTA_ESPERADA = 6
FINISH_MSJ = "FINISH"
PERDEDORES_MSJ = "NOWINNERS"


# Clase que representa una apuesta.
class Bet:
    def __init__(
        self,
        agency: str,
        first_name: str,
        last_name: str,
        document: str,
        birthdate: str,
        number: str,
    ):
        """
        agency must be passed with integer format.
        birthdate must be passed with format: 'YYYY-MM-DD'.
        number must be passed with integer format.
        """
        self.agency = int(agency)
        self.first_name = first_name
        self.last_name = last_name
        self.document = document
        self.birthdate = datetime.date.fromisoformat(birthdate)
        self.number = int(number)


""" Checks whether a bet won the prize or not. """


def has_won(bet: Bet) -> bool:
    return bet.number == LOTTERY_WINNER_NUMBER


"""
Persist the information of each bet in the STORAGE_FILEPATH file.
Not thread-safe/process-safe.
"""


def store_bets(bets: list[Bet]) -> None:
    with open(STORAGE_FILEPATH, "a+") as file:
        writer = csv.writer(file, quoting=csv.QUOTE_MINIMAL)
        for bet in bets:
            writer.writerow(
                [
                    bet.agency,
                    bet.first_name,
                    bet.last_name,
                    bet.document,
                    bet.birthdate,
                    bet.number,
                ]
            )


"""
Loads the information all the bets in the STORAGE_FILEPATH file.
Not thread-safe/process-safe.
"""


def load_bets() -> list[Bet]:
    with open(STORAGE_FILEPATH, "r") as file:
        reader = csv.reader(file, quoting=csv.QUOTE_MINIMAL)
        for row in reader:
            yield Bet(row[0], row[1], row[2], row[3], row[4], row[5])


# Funcion para recibir mensajes a través de sockets.
def receive_message(client_sock):
    tamanio = int.from_bytes(client_sock.recv(TAMANIO_HEADER), byteorder="big")
    data = b""
    while len(data) < tamanio:
        paquete = client_sock.recv(tamanio - len(data))
        if not paquete:
            raise ConnectionError("Connection closed unexpectedly")
        data += paquete

    mensaje = data.decode("utf-8").strip()
    return mensaje


# Función para enviar mensajes a través de sockets.
def send_message(client_sock, message):
    client_sock.sendall(
        "{}\n".format(message).encode("utf-8")
    )  # Garantizo que no haya 'Short-Write'.


# Función para decodificar una apuesta recibida a través de un socket.
def decode_bets(client_sock, bet_count):
    msj = receive_message(client_sock)
    if msj == FINISH_MSJ:
        return None, True

    apuestas_decodificadas = []
    bets = msj.split(";")
    for bet in bets:
        bet_data = bet.split(",")
        if len(bet_data) != CAMPOS_APUESTA_ESPERADA:
            logging.error(
                f"action: apuesta_recibida | result: fail | cantidad: {bet_count} | "
                f"msg: {msj} | error: Invalid bet data"
            )
            
            raise ValueError("Invalid bet data")

        apuesta_decodificada = Bet(
            bet_data[0], bet_data[1], bet_data[2], bet_data[3], bet_data[4], bet_data[5]
        )
        apuestas_decodificadas.append(apuesta_decodificada)

    return apuestas_decodificadas, False


# Función para enviar el acuse de recibo de una apuesta a través de un socket.
def acknowledge_bets(client_sock, bet_count):
    send_message(client_sock, bet_count)


# Función para enviar los resultados de las apuestas a través de un socket.
def send_results(client_sock, winners):
    if len(winners) == 0:
        mensaje = PERDEDORES_MSJ
    else:
        mensaje = ",".join([f"{winner.document}" for winner in winners])

    send_message(client_sock, mensaje)
