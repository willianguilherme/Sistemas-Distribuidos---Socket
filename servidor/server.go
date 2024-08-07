package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net"
	"os"
	"strconv"
)

type Employee struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Cpf     string `json:"cpf"`
	Method  string `json:"method"`
}

type Client struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Cpf     string `json:"cpf"`
}

var clients []Client

const (
	StopCharacter = "\r\n\r\n"
)

func SocketServer(port int) {
	listen, err := net.Listen("tcp4", ":"+strconv.Itoa(port))

	if err != nil {
		log.Fatalf("Socket escutando porta %d falha,%s", port, err)
		os.Exit(1)
	}

	defer listen.Close()

	log.Printf("Iniciando na porta: %d", port)

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatalln(err)
			continue
		}
		if true {
			log.Println("Cliente conectado")
		}
		go handler(conn)
	}
}

func removeClient(slice []Client, s int) []Client {
	return append(slice[:s], slice[s+1:]...)
}

func list(conn net.Conn, obj *Employee) {
	var clientArray []Client
	for i := range clients {
		if clients[i].Name == obj.Name {
			clientArray = append(clientArray, clients[i])
		}
	}
	clientArrayReturn, err := json.Marshal(clientArray)
	if err != nil {
		log.Println(err)
	}
	log.Println(clientArray)
	conn.Write(clientArrayReturn)
}

func create(conn net.Conn, obj *Employee) {
	clients = append(clients, Client{Name: obj.Name, Cpf: obj.Cpf, Address: obj.Address})
	conn.Write([]byte("Cliente criado !!"))
	log.Print("Cliente criado !!")
}

func delete(conn net.Conn, obj *Employee) {
	deleted := false
	message := "Usuario não encontrado"
	for i, item := range clients {
		if item.Name == obj.Name {
			clients = removeClient(clients, i)
			deleted = true
			break
		}
	}
	if deleted {
		conn.Write([]byte("Cliente deletado !!"))
		log.Print("Cliente deletado !!")
	} else {
		conn.Write([]byte(message))
		log.Print(message)
	}
}

func handler(conn net.Conn) {

	for {
		var err error
		var dataRequest []byte

		dataRequest, err = bufio.NewReader(conn).ReadBytes('\n')

		if err != nil {
			log.Printf("Erro ao ler a requisicao: %v\n", err)
			return
		}

		objRequest := &Employee{}

		err = json.Unmarshal(dataRequest, objRequest)

		if err != nil {
			log.Printf("Erro no parse da requisicao: %v\n", err)
			return
		}

		switch objRequest.Method {
		case "list":
			log.Print("\n\n------ENTROU NO LIST-------\n")
			list(conn, objRequest)
		case "create":
			log.Print("\n\n------ENTROU NO CREATE-------\n")
			create(conn, objRequest)
		case "delete":
			log.Print("\n\n------ENTROU NO DELETE-------\n")
			delete(conn, objRequest)
		}

	}

}

func main() {
	port := 3333
	SocketServer(port)
}
