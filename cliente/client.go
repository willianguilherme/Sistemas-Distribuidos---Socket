package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

var (
	StopCharacter = "\r\n\r\n"
	loop          = true
	jumpReq       = false
)

type Employee struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Cpf     string `json:"cpf"`
	Method  string `json:"method"`
}

func list() []byte {
	fmt.Println("\nLISTANDO")
	userList := ""
	fmt.Println("\nDigite o nome")
	fmt.Scanln(&userList)
	client := &Employee{Method: "list", Name: userList}
	listClient, err := json.Marshal(client)
	if err != nil {
		fmt.Println(err)
	}
	return listClient
}
func create() []byte {
	fmt.Println("\nCRIANDO")
	client := Employee{}
	fmt.Println("Digite o nome")
	fmt.Scanln(&client.Name)
	fmt.Println("Digite o cpf")
	fmt.Scanln(&client.Cpf)
	fmt.Println("Digite o endereco")
	fmt.Scanln(&client.Address)
	client.Method = "create"
	newClient, err := json.Marshal(client)
	if err != nil {
		fmt.Println(err)
	}
	return newClient
}
func delete() []byte {
	fmt.Println("\nDELETANDO")

	userDelete := ""
	fmt.Println("Digite o nome")
	fmt.Scanln(&userDelete)
	client := &Employee{Name: userDelete, Method: "delete"}
	deleteClient, err := json.Marshal(client)
	if err != nil {
		fmt.Println(err)
	}
	return deleteClient
}
func exit() {
	fmt.Println("\nSAINDO")
	loop = false
}

func SocketClient(ip string, port int) {
	addr := strings.Join([]string{ip, strconv.Itoa(port)}, ":")
	conn, err := net.Dial("tcp", addr)

	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	defer conn.Close()

	for loop {
		option := "-1"
		jumpReq = false
		fmt.Println("\n\nSelecione uma opcao: \n1 - criar\n2 - listar\n3 - deletar\n9 - sair")
		fmt.Scanln(&option)
		fmt.Println("Voce selecionou - ", option)

		var sendReq []byte

		switch option {
		case "1":
			sendReq = create()
		case "2":
			sendReq = list()
		case "3":
			sendReq = delete()
		case "9":
			exit()
			jumpReq = true
		default:
			fmt.Println("Opcao invalida -- Tente novamente")
			jumpReq = true
		}
		if !jumpReq {
			conn.Write(sendReq)
			conn.Write([]byte(StopCharacter))
			log.Printf("Send: %s", sendReq)

			buff := make([]byte, 1024)
			n, _ := conn.Read(buff)
			log.Printf("Receive: %s", string(buff[:n]))
		}
		fmt.Println("")
	}
}

func main() {

	var (
		ip   = "127.0.0.1"
		port = 3333
	)

	SocketClient(ip, port)
}
