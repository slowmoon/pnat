package  main

import (
    "fmt"
    "net"
)

func main() {

    udpaddr ,err := net.ResolveUDPAddr("udp", net.JoinHostPort("127.0.0.1", "20001"))
    if err!=nil{
        panic(err)
    }
    conn, err := net.DialUDP("udp",nil, udpaddr)
    if err!= nil{
        panic(err)
    }

    buf := make([]byte, 1024)
    for {

         _, err := conn.Write([]byte{1, 2, 3,4, 5, 6, 7, 8,6,0,0,0,0,123,23,1,3,4,6, 97, 98, 99})
         if err!=nil{
             fmt.Println("writdde error ", err)
         }
         n, err := conn.Read(buf)
         if err!=nil{
             panic(err)
         }
         fmt.Printf("receive message %s", buf[:n])
    }

}
