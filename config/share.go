package config

import (
	"crypto/rand"
	"crypto/sha3"
	"encoding/hex"
)

// A host can share its ping times on u00.io, a public page of live values.
// Each host has its own API key: 16 random bytes in hex, secret, it allows
// to set the values. The page is found by the public address: the first 16
// bytes of SHA3-256 of the key, in hex. See https://u00.io/docs

// ShareSite is where the values are sent and shown
const ShareSite = "https://u00.io"

// NewShareKey returns a new random API key
func NewShareKey() string {
	bs := make([]byte, 16)
	rand.Read(bs)
	return hex.EncodeToString(bs)
}

// ShareAddress returns the public address of the key, "" for an invalid key
func ShareAddress(key string) string {
	bs, err := hex.DecodeString(key)
	if err != nil || len(bs) != 16 {
		return ""
	}
	sum := sha3.Sum256(bs)
	return hex.EncodeToString(sum[:16])
}

// ShareURL returns the public page of the host's values, "" without a valid key
func (c ConfigHost) ShareURL() string {
	address := ShareAddress(c.ShareKey)
	if address == "" {
		return ""
	}
	return ShareSite + "/ch/" + address
}

// ensureShareKeys gives a key to the hosts that have none (added by an older
// version or by hand); returns whether one was added
func (c *Config) ensureShareKeys() bool {
	added := false
	for _, host := range c.Hosts {
		if ShareAddress(host.ShareKey) == "" {
			host.ShareKey = NewShareKey()
			added = true
		}
	}
	return added
}
