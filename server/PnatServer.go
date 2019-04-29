package server

import (
    "log"
    "net"
    "sync"
)

type PNatManager struct {

  Host string
  Port string
  Conn net.Conn
  Addr *net.UDPAddr
  SubNets   map[string]interface{}
  closed chan struct{}

  Buffered  []byte

  MaxPort int

  WorkerPort  int

  mutex sync.Mutex
}

type RequestCommand []byte

var request RequestCommand = []byte("ping request")

func NewPNAT(host string, port string)*PNatManager {
    manager := PNatManager{
        Host: host,
        Port: port,
        SubNets: make(map[string]interface{}),
        closed: make(chan struct{}),
    }

    remoteAddr,err := net.ResolveUDPAddr("udp", net.JoinHostPort(manager.Host, manager.Port))
    if err!=nil{
       panic("init pnat manager error")
    }
    manager.Addr = remoteAddr

   conn, err :=  net.DialUDP("udp", nil,remoteAddr )
   if err!=nil{
       panic("dial udp server error!")
   }
   manager.Conn = conn
   return  &manager
}

func (p *PNatManager)FetchNetworks() {
   if conn, ok := p.Conn.(*net.UDPConn);ok{
      _, err := conn.Write(request)
      if err!= nil{
         log.Println("fetch request connection fail", err.Error())
          return
      }
      n, err := conn.Read(p.Buffered)
      if err!=nil{
          log.Println("fetch response connection fail", err.Error())
          return
      }
      p.handleFetch(p.Buffered[:n])
      p.Buffered = p.Buffered[:0]
   }
}

func (p *PNatManager)handleFetch(msg []byte) {
    _ := NewWorker(p, "", "")

}

func (p *PNatManager)work(worker *Worker){

    go worker.Start()
}
