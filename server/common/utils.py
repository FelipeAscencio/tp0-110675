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

# Constantes de tamaños, indices y límites.
TAMANIO_HEADER = 2
CAMPOS_APUESTA_ESPERADA = 6
INDEX_AGENCIA = 0
INDEX_NOMBRE = 1
INDEX_APELLIDO = 2
INDEX_DOCUMENTO = 3
INDEX_FECHA_NACIMIENTO = 4
INDEX_NUMERO = 5
FINISH_MSJ = "FINISH"

# Clase que representa una apuesta.
""" A lottery bet registry. """
class Bet:
    def __init__(self, agency: str, first_name: str, last_name: str, document: str, birthdate: str, number: str):
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
    with open(STORAGE_FILEPATH, 'a+') as file:
        writer = csv.writer(file, quoting=csv.QUOTE_MINIMAL)
        for bet in bets:
            writer.writerow([bet.agency, bet.first_name, bet.last_name,
                             bet.document, bet.birthdate, bet.number])

"""
Loads the information all the bets in the STORAGE_FILEPATH file.
Not thread-safe/process-safe.
"""
def load_bets() -> list[Bet]:
    with open(STORAGE_FILEPATH, 'r') as file:
        reader = csv.reader(file, quoting=csv.QUOTE_MINIMAL)
        for row in reader:
            yield Bet(row[0], row[1], row[2], row[3], row[4], row[5])

# Implementación de las funciones de comunicación (Para el ejercicio 5).
# Funcion para recibir mensajes a través de sockets.
def receive_message(client_sock):
    tamanio = int.from_bytes(client_sock.recv(TAMANIO_HEADER), byteorder='big')
    informacion = b""
    while len(informacion) < tamanio: # Garantiza que no tengamos un 'Short-Read'.
        paquete = client_sock.recv(tamanio - len(informacion))
        if not paquete:
            raise ConnectionError("Connection closed unexpectedly")
        informacion += paquete
    
    mensaje = informacion.decode('utf-8').strip()
    return mensaje

# Función para enviar mensajes a través de sockets.
def send_message(client_sock, message): #CAMBIAR POR SENDALL SI O SI PARA QUE GARANTICE QUE NO TENEMOS 'SHORT-WRITE'.
    client_sock.send("{}\n".format(message).encode('utf-8')) # Garantiza que no tenemos 'Short-write'.

# Función para decodificar una apuesta recibida a través de un socket.
def decode_bet(client_sock):
    mensaje = receive_message(client_sock)
    if mensaje == FINISH_MSJ:
        return None, True
    
    bets_decodificadas = []
    bets = mensaje.split(';')
    for bet in bets:
        bet_informacion = bet.split(',')
        if len(bet_informacion) != CAMPOS_APUESTA_ESPERADA:
            logging.error(f"action: apuesta_recibida | result: fail | cantidad: {bet_count} | msg: {mensaje} | error: Invalid bet data")
            raise ValueError("Invalid bet data")
    
        bet_decodificada = Bet(bet_informacion[INDEX_AGENCIA], bet_informacion[INDEX_NOMBRE], bet_informacion[INDEX_APELLIDO], bet_informacion[INDEX_DOCUMENTO], bet_informacion[INDEX_FECHA_NACIMIENTO], bet_informacion[INDEX_NUMERO])
        bets_decodificadas.append(bet_decodificada)

    return bets_decodificadas, False

# Función para enviar el acuse de recibo de una apuesta a través de un socket.
def acknowledge_bet(client_sock, document, number):
    send_message(client_sock, bet_count)
