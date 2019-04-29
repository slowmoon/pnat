package server

import (
	"bytes"
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
		closed: make(chan struct{}, 1),
		ReplyHost: host,
		ReplyPort: port,
		Create: time.Now(),
		Buffered: make([]byte, 1024),
	}

	addr ,err := net.ResolveUDPAddr("udp", net.JoinHostPort("0.0.0.0", worker.localPort))

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
		select {
		case <- w.closed:
			break
		default:
		n, addr ,err := w.conn.ReadFromUDP(w.Buffered)
		if err!=nil{
			log.Print(err)
			continue
		}
		data := w.Buffered[:n]

		if n <=16 {
			continue
		}

		dataType := data[8]
		host := data[9:13]
		port := data[13:15]

		log.Printf("receive data type %v , host %s, port %s \n", dataType, host, port)

		if bytes.Equal(host, []byte{0, 0, 0, 0}){
			// this is from the client request
			data[8] = 3
			sip := addr.IP.To4()
			sport := addr.Port
			log.Printf("receive remote ip %s port %d\n", sip, sport)

			copy(data[9:13], sip)
			copy(data[13:15], []byte{byte(sport >>8), byte(sport & 0xff)})

			//redirect the request
			log.Println("begin to redirect data:", data)
			_, err :=  w.conn.WriteToUDP(data, w.ReplyAddr)
			if err!=nil{
				log.Println("send data error ", err.Error() )
				continue
			}
		}else {
			// this is from the pnat request
			//clear the message
			port0 := int(data[13])<<8 + int(data[14])
			ip := net.IP(host)

			udpaddr ,err := net.ResolveUDPAddr("udp", net.JoinHostPort(ip.String(), strconv.Itoa(port0)))
			if err!=nil {
				log.Println("receive data fail", err)
				continue
			}

			copy(data[8:15], []byte{0, 0, 0, 0,0,0,0})
			_, err = w.conn.WriteToUDP(data, udpaddr)
			if err!= nil{
				log.Println("redirect data fail",err)
			}
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
