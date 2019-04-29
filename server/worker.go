package server

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
	"strconv"
	"time"
)

type Worker struct {

  manager *PNatManager

  localPort string

  closed chan struct{}

  conn *net.UDPConn

  Create time.Time

  Buffered []byte

  ReplyHost string

  ReplyPort string

  ReplyAddr *net.UDPAddr
}

func NewWorker(manager *PNatManager, host string, port string)(worker *Worker, err error)  {
	worker = &Worker{
		manager: manager,
		localPort: strconv.Itoa(manager.RandPort()),
		closed: make(chan struct{}),
		ReplyHost: host,
		ReplyPort: port,
		Create: time.Now(),
	}

	addr ,err := net.ResolveUDPAddr("udp", net.JoinHostPort("", worker.localPort))

	if err!=nil{
		return  nil, err
	}

	worker.ReplyAddr,err = net.ResolveUDPAddr("udp", net.JoinHostPort(host, port))

	if err!=nil{
		return  nil, err
	}

	conn, err := net.ListenUDP("udp", addr)

	if err!=nil{
		return  nil, err
	}
	worker.conn = conn

	log.Println("worker start success")
	return
}

func (w *Worker)Ping()(err error) {
       //currently skip verify
      _, err = w.conn.WriteToUDP(request, w.manager.Addr)
   return
}

func (w *Worker)Start()(err error) {
	defer func() {
	    err = w.conn.Close()
    }()

	for {

		n, addr ,err := w.conn.ReadFromUDP(w.Buffered)
		if err!=nil{
			log.Print(err)
			continue
		}
		data := w.Buffered[:n]
	    if n <=15 {
	      log.Println("invalid data")
	      continue
		}
		dataType := data[8]
		host := data[9:13]
		port := data[13:15]

		log.Printf("receive data type %v , host %s, port %s ", dataType, host, port)

		if bytes.Equal(host, []byte{0, 0, 0, 0}){
		   // this is from the client request
		   data[8] = 3
		   sip := addr.IP.To4()
		   sport := addr.Port
		   copy(data[9:13], sip)
		   binary.PutVarint(data[13:15], int64(sport))
		   //redirect the request
		  _, err :=  w.conn.WriteToUDP(data, w.ReplyAddr)
		  if err!=nil{
		  	log.Println("send data error ", err.Error() )
		  	continue
		  }
		}else {
			// this is from the pnat request
			//clear the message
			udpaddr ,err := net.ResolveUDPAddr("udp", net.JoinHostPort(string(host), string(port)))
			if err!=nil {
				//
				log.Println("receive data fail", err)
				continue
			}
		    copy(data[8:15], []byte{0, 0, 0, 0,0,0,0})

			_, err = w.conn.WriteToUDP(data, udpaddr)
			if err!= nil{
				log.Println("redirect data fail")
			}
		}

    }

}

func (w *Worker)Redirect(packet []byte, addr *net.UDPAddr)(err error)  {
	 //modify the packet
	_, err =  w.conn.WriteToUDP(packet, w.ReplyAddr)
    return
}

func (w *Worker)Stop(){
   w.closed <- struct{}{}
}
