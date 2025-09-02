# Importación de los paquetes necesarios.
import socket
import logging
import signal
import threading  # Biblioteca utilizada para manejar múltiples hilos de ejecución.
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
        self.clients = clients
        self.clientes_terminados = 0
        self.winners = []
        self.candado_mutex_1 = threading.Lock()
        self.candado_mutex_2 = threading.Lock()
        self.hilos = []

    # Loop principal del servidor (Modificación de código para manejar la señal pedida).
    def run(self):
        """
        Bucle de servidor ficticio.
        El servidor acepta nuevas conexiones y establece comunicación con un cliente.
        Una vez finalizada la comunicación con el cliente, el servidor vuelve a aceptar nuevas conexiones.
        """

        self.running = True
        signal.signal(signal.SIGTERM, self.shutdown)
        while self.running:
            try:
                client_sock = self.__accept_new_connection()
                if client_sock:
                    hilo_cliente = threading.Thread(
                        target=self.__handle_client_connection, args=(client_sock,)
                    )
                    hilo_cliente.daemon = True
                    hilo_cliente.start()
                    self.hilos.append(hilo_cliente)
            except:
                if not self.running:
                    break

        self.__esperar_otros_hilos()

    # Función que espera a que terminen todos los hilos del servidor.
    def __esperar_otros_hilos(self):
        for thread in self.hilos:
            thread.join()

    # Función que apaga el servidor de manera ordenada.
    def shutdown(self, signum=None, frame=None):
        logging.info("action: shutdown | result: in_progress")
        self.running = False
        if self._server_socket:
            self._server_socket.close()

        self.__esperar_otros_hilos()
        logging.info("action: shutdown | result: success")

    # Función que maneja la comunicación con un cliente específico.
    def __handle_client_connection(self, client_sock):
        try:
            mensaje = utils.receive_message(client_sock)
            if mensaje == MSJ_APUESTA:
                contador_apuesta = 0
                while True:
                    bets, terminada = utils.decode_bets(client_sock, contador_apuesta)
                    if terminada:
                        utils.acknowledge_bets(client_sock, contador_apuesta)
                        logging.info(
                            f"action: apuesta_recibida | result: success | cantidad: {contador_apuesta}"
                        )
                        with self.candado_mutex_1:
                            self.clientes_terminados += 1
                        break

                    with self.candado_mutex_2:
                        utils.store_bets(bets)
                        contador_apuesta += len(bets)

            if mensaje == MSJ_RESULTADOS:
                agency_id = utils.receive_message(client_sock)
                with self.candado_mutex_1:
                    if int(self.clientes_terminados) == int(self.clients):
                        if not self.winners:
                            bets = utils.load_bets()
                            winners = [bet for bet in bets if utils.has_won(bet)]
                            self.winners = winners
                            logging.info(f"action: sorteo | result: success")

                        agencia_ganadores = [
                            bet for bet in self.winners if bet.agency == int(agency_id)
                        ]
                        utils.send_results(client_sock, agencia_ganadores)
        except OSError as e:
            logging.error("action: receive_message | result: fail | error: {e}")
            utils.send_message(client_sock, "ERROR")
        finally:
            client_sock.close()

    # Función que acepta nuevas conexiones de clientes.
    def __accept_new_connection(self):
        """
        Aceptar nuevas conexiones
        La función se bloquea hasta que se establece la conexión con un cliente.
        Después, se imprime y devuelve la conexión creada.
        """

        logging.info("action: accept_connections | result: in_progress")
        c, direccion = self._server_socket.accept()
        logging.info(
            f"action: accept_connections | result: success | ip: {direccion[0]}"
        )
        return c
