# Importación de los paquetes necesarios.
import socket
import logging
import signal
from common import utils

# Constantes.
MSJ_APUESTA = "BET"
MSJ_RESULTADOS = "RESULTS"


# Clase que implementa el servidor.
class Server:
    # Constructor del servidor.
    def __init__(self, port, listen_backlog, clients):
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(("", port))
        self._server_socket.listen(listen_backlog)
        self.running = False
        self.clients = clients  # Modificación de código para implementar ejercicio 7.
        self.client_sockets = []
        self.finished_clients = (
            0  # Modificación de código para implementar la lógica del ejercicio 7.
        )
        self.winners = (
            []
        )  # Modificación de código para implementar la lógica del ejercicio 7.

    # Loop principal del servidor (Modificación de código para manejar la señal pedida).
    def run(self):
        """
        Bucle de servidor ficticio.
        El servidor acepta nuevas conexiones y establece comunicación con un cliente.
        Una vez finalizada la comunicación con el cliente, el servidor vuelve
        a aceptar nuevas conexiones.
        """

        self.running = True
        signal.signal(signal.SIGTERM, self.shutdown)
        while self.running:
            try:
                client_sock = self.__accept_new_connection()
                if client_sock:
                    self.client_sockets.append(client_sock)
                    self.__handle_client_connection(client_sock)
            except:
                if not self.running:
                    break

    # Función que apaga el servidor de manera ordenada.
    def shutdown(self, signum=None, frame=None):
        self.running = False
        logging.info(f"action: shutdown | result: in_progress")
        if self._server_socket:
            self._server_socket.close()
        for client_sock in self.client_sockets:
            client_sock.close()
        logging.info(f"action: shutdown | result: success")

    # Función que maneja la comunicación con un cliente específico.
    def __handle_client_connection(self, client_sock):
        """
        Lee el mensaje de un socket de cliente específico y lo cierra.
        Si surge un problema en la comunicación con el cliente, el socket del
        cliente también se cerrará.
        """

        # Modificación de código para implementar la lógica del ejercicio 7.
        try:
            mensaje = utils.receive_message(client_sock)
            if mensaje == MSJ_APUESTA:
                bet_count = 0
                while True:
                    bets, termino = utils.decode_bets(client_sock, bet_count)
                    if termino:
                        utils.acknowledge_bets(client_sock, bet_count)
                        logging.info(
                            f"action: apuesta_recibida | result: success | "
                            f"cantidad: {bet_count}"
                        )
                        
                        self.finished_clients += 1
                        break

                    utils.store_bets(bets)
                    bet_count += len(bets)

            if mensaje == MSJ_RESULTADOS:
                agency_id = utils.receive_message(client_sock)
                if int(self.finished_clients) == int(self.clients):
                    if not self.winners:
                        bets = utils.load_bets()
                        ganadores = [bet for bet in bets if utils.has_won(bet)]
                        self.winners = ganadores
                        logging.info(f"action: sorteo | result: success")

                    agencia_ganadores = [
                        bet for bet in self.winners if bet.agency == int(agency_id)
                    ]
                    utils.send_results(client_sock, agencia_ganadores)
        except OSError as e:
            logging.error("action: receive_message | result: fail | error: {e}")
        finally:
            client_sock.close()
            self.client_sockets.remove(client_sock)

    # Función que acepta nuevas conexiones de clientes.
    def __accept_new_connection(self):
        """
        Aceptar nuevas conexiones
        La función se bloquea hasta que se establece la conexión con un cliente.
        Después, se imprime y devuelve la conexión creada.
        """

        logging.info("action: accept_connections | result: in_progress")
        c, addr = self._server_socket.accept()
        logging.info(f"action: accept_connections | result: success | ip: {addr[0]}")
        return c
