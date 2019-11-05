package proxy

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Proxy represents a structure that holds relevant ip info to proxy
// an http request
type Proxy struct {
	ID            primitive.ObjectID `bson:"_id"`
	IPAddress     string             `bson:"ip_address"`
	Port          int                `bson:"port"`
	Protocol      string             `bson:"protocol"`
	Location      string             `bson:"location"`
	CreatedAt     time.Time          `bson:"created_at"`
	IsAnonymous   bool               `bson:"is_anonymous"`
	IsDeactivated bool               `bson:"is_deactivated"`
	Socket        string             `bson:"socket"`
}

// NewProxy factor returns a pointer to a new proxy instance
func NewProxy(addr string, port int, protocol string, loc string) *Proxy {
	return &Proxy{
		ID:            primitive.NewObjectID(),
		IPAddress:     addr,
		Port:          port,
		Protocol:      protocol,
		Location:      loc,
		CreatedAt:     time.Now(),
		IsAnonymous:   false,
		IsDeactivated: false,
		Socket:        fmt.Sprintf("%v://%v:%v", protocol, addr, string(port)),
	}
}
