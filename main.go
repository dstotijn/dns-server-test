package main

import (
	"fmt"
	"log"
	"os"

	"github.com/miekg/dns"
)

func main() {
	udpAddr := ":53"
	tcpAddr := ":53"

	if v := os.Getenv("TCP_ADDR"); v != "" {
		tcpAddr = v
	}

	if v := os.Getenv("UDP_ADDR"); v != "" {
		udpAddr = v
	}

	udpServer := &dns.Server{
		Addr:    udpAddr,
		Net:     "udp",
		Handler: dns.HandlerFunc(handler),
	}

	tcpServer := &dns.Server{
		Addr:    tcpAddr,
		Net:     "tcp",
		Handler: dns.HandlerFunc(handler),
	}

	go func() {
		if err := udpServer.ListenAndServe(); err != nil {
			log.Panic(err)
		}
	}()
	if err := tcpServer.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}

func handler(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)

	log.Printf("Received request: %v", r.Question[0])

	m.Answer = append(m.Answer, &dns.TXT{
		Hdr: dns.RR_Header{
			Name:   r.Question[0].Name,
			Rrtype: dns.TypeTXT,
			Class:  dns.ClassINET,
			Ttl:    3600,
		},
		Txt: []string{fmt.Sprintf("Remote addr: %v", w.RemoteAddr().String())},
	})

	log.Println(m.Answer)

	if err := w.WriteMsg(m); err != nil {
		log.Printf("Error writing message: %v", err)
	}
}
