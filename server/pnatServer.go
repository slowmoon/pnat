package server

import (
    "fmt"
    "log"
    "net"
    "sync"
    "time"
)

const (
    LIFECYCLE = time.Minute * 5
    HEALTHYCHECK = time.Minute *1
    TIMEOUT = time.Minute*5
    PORT_MAX = 65535
    PORT_LEASE = 20000
)
type RequestCommand []byte
var request RequestCommand = []byte("ping request")


type PNatManager struct {

  Host string

  Port string

  Conn net.Conn

  Addr *net.UDPAddr

  SubNets   map[string]*Worker

  closed chan struct{}

  Buffered  []byte


  WorkerPort  int

  mutex sync.Mutex
}



func NewPNAT(host string, port string)*PNatManager {
    manager := PNatManager{
        Host: host,
        Port: port,
        SubNets: make(map[string]*Worker),
        closed: make(chan struct{}),
        WorkerPort: PORT_LEASE,
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
   log.Printf("manager start ok @ host %s port %s\n", host, port)
   return  &manager
}

func (p *PNatManager)RandPort()int {
    if p.WorkerPort >= PORT_MAX {
        p.WorkerPort = PORT_LEASE
    }
    p.WorkerPort++
    return  p.WorkerPort
}

func (p *PNatManager)FetchNetworks() {
	fmt.Println("fetch networks begin:")
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

    var networkId  string

    if  wk, ok := p.SubNets[networkId];!ok {

        worker, err := NewWorker(p, "", "")

        if err!=nil{
            log.Printf("create worker fail error %s\n",err.Error())
        }
        go worker.Start()

        p.SubNets["networkId"] = worker
    }else {
        wk.Create = time.Now()
    }
}

func (p *PNatManager)LifecycleChecks() {
	p.mutex.Lock()
	for _, v := range p.SubNets {
	   v.Ping()
    }
	p.mutex.Unlock()
}

func (p *PNatManager)healthyChecks() {

   p.mutex.Lock()
   subnets := make(map[string]*Worker)
   for i, k :=range p.SubNets {
       subnets[i] = k
   }
   p.mutex.Unlock()

   for n, v := range subnets {
       if time.Since(v.Create) >= TIMEOUT{
       	   go func() {
       	       v.Stop()
           }()
           delete(subnets,n )
       }
   }
   p.SubNets = subnets

}