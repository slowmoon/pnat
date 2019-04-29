package  main

import (
	"log"
	"net"
)

func main()  {

	addr ,err := net.ResolveUDPAddr("udp", "localhost:12345")
	if err !=nil{
		panic(err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err!= nil{
		panic(err)
	}
	buf := make([]byte, 1024)

	for {
		n, addr ,err := conn.ReadFromUDP(buf)
		if err!=nil{
			log.Println(err)
			continue
		}
		data := buf[:n]
		log.Printf("receive data %s", data)
		data[8] = 4
		n, err = conn.WriteToUDP(data, addr)
		if err!=nil{
		  log.Println(err)
		}
	}
}
