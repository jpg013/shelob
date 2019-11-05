package proxy

import (
	"crypto/sha1"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UnknownProxyPort struct {
	ID        primitive.ObjectID `bson:"_id"`
	OCRText   string             `bson:"ocr_text"`
	HashState string             `bson:"hash_state"`
}

// NewUnknownProxyPort factory func returns a pointer to a new UnknownProxySource instance structure
func NewUnknownProxyPort(src string, txt string) *UnknownProxyPort {
	h := sha1.New()

	h.Write([]byte(src))

	hash := fmt.Sprintf("%x", h.Sum(nil))

	return &UnknownProxyPort{
		ID:        primitive.NewObjectID(),
		OCRText:   txt,
		HashState: hash,
	}
}
