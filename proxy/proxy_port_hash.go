package proxy

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProxyPortHash struct {
	ID        primitive.ObjectID `bson:"_id"`
	Port      int                `bson:"port"`
	HashState string             `bson:"hash_state"`
}

// NewProxyPortHash factory func returns a pointer to a new UnknownProxySource instance structure
func NewProxyPortHash(src string, port int) (p *ProxyPortHash) {
	p = &ProxyPortHash{
		ID:        primitive.NewObjectID(),
		HashState: hashString(src),
	}

	if port != 0 {
		p.Port = port
	}

	return p
}
