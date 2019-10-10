package proxy

import (
	"log"
	"shelob/db"
)

// Proxy represents a structure that holds relevant ip info to proxy
// an http request
type Proxy struct {
	IPAddress string
	Port      int
	Protocol  string
	Location  string
}

// NewProxy factor returns a pointer to a new proxy instance
func NewProxy(addr string, port int, protocol string, loc string) *Proxy {
	return &Proxy{
		IPAddress: addr,
		Port:      port,
		Protocol:  protocol,
		Location:  loc,
	}
}

// Insert creates a new record in the database
func (p *Proxy) Insert() int64 {
	sql := "INSERT INTO proxy(ip_address, port, protocol, location) VALUES(?,?,?,?)"
	stmt, err := db.Conn.Prepare(sql)

	if err != nil {
		log.Fatal(err)
	}

	res, err := stmt.Exec(
		p.IPAddress,
		p.Port,
		p.Protocol,
		p.Location,
	)

	if err != nil {
		log.Fatal(err)
	}

	id, err := res.LastInsertId()

	if err != nil {
		log.Fatal(err)
	}

	return id
}
