package server

import (
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

func NewWorker(manager *PNatManager, host string, port string)*Worker  {
   worker := Worker{
     manager: manager,
     localPort: strconv.Itoa(manager.WorkerPort),
     closed: make(chan struct{}),
     ReplyHost: host,
     ReplyPort: port,
     Create: time.Now(),
   }

   addr ,err := net.ResolveUDPAddr("udp", net.JoinHostPort("", worker.localPort))

	if err!=nil{
		panic("create worker fail")
	}

	worker.ReplyAddr,err = net.ResolveUDPAddr("udp", net.JoinHostPort(host, port))

	if err!=nil{
		panic("create worker fail")
	}

	conn, err := net.ListenUDP("udp", addr)

	if err!=nil{
		panic(err)
	}

	worker.conn = conn
   return  &worker
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

	go func() {
		ticker := time.NewTicker(time.Second * 5)
        defer  ticker.Stop()

		for  range ticker.C {
		   w.Ping()
        }
    }()


	for {

		n, addr ,err := w.conn.ReadFromUDP(w.Buffered)
		if err!=nil{
			log.Print(err)
			continue
		}

		switch  w.Buffered[:n] {
		case []byte("") :
			// 客户端转发
		  w.Redirect(w.Buffered[:n], addr)
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
