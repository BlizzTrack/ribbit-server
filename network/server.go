/*
 * Copyright (c) 2020. BlizzTrack
 */

package network

import (
	"bufio"
	"net"
	"strings"

	log "github.com/sirupsen/logrus"
)

type handleData func(Command) (bool, []byte, string, error)

type Server struct {
	addr   string
	client *BlizzardClient

	running bool

	onData handleData
}

func NewServer(addr string, client *BlizzardClient, onData handleData) *Server {
	server := new(Server)
	server.addr = addr
	server.client = client
	server.onData = onData

	server.running = false

	return server
}

func (server *Server) Run() error {
	server.running = true

	l, err := net.Listen("tcp", server.addr)
	if err != nil {
		return err
	}
	defer l.Close()

	for server.running {
		conn, err := l.Accept()
		if err != nil {
			continue
		}

		log.Infof("[%s] has connected to the cluster", conn.RemoteAddr().String())
		go server.handle(conn)
	}

	return nil
}

func (server *Server) handle(conn net.Conn) {
	defer conn.Close()

	in, _, err := bufio.NewReader(conn).ReadLine()
	if err != nil {
		log.Errorf("[%s] failed to read inbound data", conn.RemoteAddr().String())
		return
	}

	command, err := ParseCommand(strings.ToLower(string(in)))
	if err != nil {
		log.Errorf("[%s] command error %s", conn.RemoteAddr().String(), err)
		return
	}

	remote, out, seqn, err := server.onData(command)
	if err != nil {
		log.Errorf("[%s] callback error %s", conn.RemoteAddr().String(), err)
		return
	}

	mode := "local"
	if remote {
		mode = "remote"
	}
	log.Infof("[%s] %s -> %s %s", conn.RemoteAddr().String(), command, mode, seqn)
	if _, err := conn.Write([]byte(out)); err != nil {
		log.Errorf("[%s] failed to write %s", conn.RemoteAddr().String(), err)
	}

}
