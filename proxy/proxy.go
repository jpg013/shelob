package proxy

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Proxy represents a structure that holds relevant ip info to proxy
// an http request
type Proxy struct {
	ID          primitive.ObjectID `bson:"_id"`
	IPAddress   string             `bson:"ip_address"`
	Port        int                `bson:"port"`
	Protocol    string             `bson:"protocol"`
	Location    string             `bson:"location"`
	CreatedAt   time.Time          `bson:"created_at"`
	IsAnonymous bool               `bson:"is_anonymous"`
	IsActive    bool               `bson:"is_deactivated"`
	LastChecked time.Time          `bson:"last_checked"`
	Source      string             `bson:"source"`
}

// NewProxy factor returns a pointer to a new proxy instance
func NewProxy(data map[string]interface{}) (*Proxy, error) {
	port := data["port"].(int)

	proxy := &Proxy{
		ID:          primitive.NewObjectID(),
		IPAddress:   data["ipAddress"].(string),
		Port:        port,
		Protocol:    data["protocol"].(string),
		Location:    data["location"].(string),
		CreatedAt:   time.Now(),
		IsAnonymous: false,
		IsActive:    false,
		LastChecked: time.Time{},
		Source:      data["source"].(string),
	}

	return proxy, nil
}
