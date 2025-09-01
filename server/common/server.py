# Importación de los paquetes necesarios.
import socket
import logging
import signal   # Modificación de código para manejar la señal pedida.

# Clase que implementa el servidor.
class Server:
    # Constructor del servidor.
    def __init__(self, port, listen_backlog):
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self.running = False   # Modificación de código para manejar la señal pedida.
        self.client_sockets = []

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
                    self.client_sockets.append(client_sock)
                    self.__handle_client_connection(client_sock)
            except:
                if not self.running:
                    break
    
    # Función que apaga el servidor de manera ordenada (Modificación de código para manejar la señal pedida).
    def shutdown(self, signum=None, frame=None):
        self.running = False
        logging.info(f'action: shutdown | result: in_progress')
        if self._server_socket:
            self._server_socket.close()

        for client_socket in self.client_sockets:
            client_socket.close()

        logging.info(f'action: shutdown | result: success')

    # Función que maneja la comunicación con un cliente específico.
    def __handle_client_connection(self, client_sock):
        """
        Lee el mensaje de un socket de cliente específico y lo cierra.
        Si surge un problema en la comunicación con el cliente, el socket del
        cliente también se cerrará.
        """

        try:
            msg = client_sock.recv(1024).rstrip().decode('utf-8')
            addr = client_sock.getpeername()
            logging.info(f'action: receive_message | result: success | ip: {addr[0]} | msg: {msg}')
            client_sock.send("{}\n".format(msg).encode('utf-8'))
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

        logging.info('action: accept_connections | result: in_progress')
        c, addr = self._server_socket.accept()
        logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
        return c
