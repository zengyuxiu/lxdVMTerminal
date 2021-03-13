package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	lxd "github.com/lxc/lxd/client"
	"github.com/lxc/lxd/shared/api"
	"io"
	"net"
)

func vga(d lxd.InstanceServer, name string) error {
	var err error

	// We currently use the control websocket just to abort in case of errors.
	controlDone := make(chan struct{}, 1)
	handler := func(control *websocket.Conn) {
		<-controlDone
		closeMsg := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "")
		control.WriteMessage(websocket.CloseMessage, closeMsg)
	}

	// Prepare the remote console.
	req := api.InstanceConsolePost{
		Type: "vga",
	}

	consoleDisconnect := make(chan bool)
	sendDisconnect := make(chan bool)
	defer close(sendDisconnect)

	consoleArgs := lxd.InstanceConsoleArgs{
		Control:           handler,
		ConsoleDisconnect: consoleDisconnect,
	}

	go func() {
		<-sendDisconnect
		close(consoleDisconnect)
	}()

	var socket string
	var listener net.Listener

	listener, err = net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return err
	}
	defer listener.Close()

	addr := listener.Addr().(*net.TCPAddr)
	socket = fmt.Sprintf("spice://127.0.0.1:%d", addr.Port)
	fmt.Printf(socket)
	op, connect, err := d.ConsoleInstanceDynamic(name, req, &consoleArgs)
	if err != nil {
		return err
	}

	// Handle connections to the socket.
	go func() {
		count := 0

		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			count++

			go func(conn io.ReadWriteCloser) {
				err = connect(conn)
				if err != nil {
					sendDisconnect <- true
				}
				count--
				if count == 0 {
					sendDisconnect <- true
				}
			}(conn)
		}
	}()

	// Wait for the operation to complete.
	err = op.Wait()
	if err != nil {
		return err
	}

	return nil
}
