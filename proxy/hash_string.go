package proxy

import (
	"crypto/sha1"
	"fmt"
)

func hashString(s string) string {
	h := sha1.New()
	h.Write([]byte(s))

	return fmt.Sprintf("%x", h.Sum(nil)) // fmt hash as string
}
