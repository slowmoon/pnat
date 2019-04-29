package  main

import (
   "flag"
   "os"
   "os/signal"
   "pnat/server"
   "syscall"
   "time"
)

var (
   host string
   port string
)


func main()  {

   flag.StringVar(&host, "host", "", "special the host")
   flag.StringVar(&port, "port", "", "special the port")
   flag.Parse()

   natserver := server.NewPNAT(host, port)
   ticker := time.NewTicker(time.Minute * 5)

   ch := make(chan os.Signal, 1)
   signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)


   for {
      select {
      case <-ticker.C:
         natserver.FetchNetworks()
      case <-ch:
      	break
      }
   }

}
