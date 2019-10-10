package proxy

// Proxy represents a structure that holds relevant ip info to proxy
// an http request
type Proxy struct {
	IPAddress string
	Port      int
	Protocol  string
	Location  string
	Speed     int8
}
