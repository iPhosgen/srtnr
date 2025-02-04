package shortener

import (
	"encoding/hex"
	"hash"
	"strconv"
	"time"

	"github.com/spaolacci/murmur3"
)

const emptyString = ""

type Shortener interface {
	Shorten(url string, userId string) (string, error)
}

type UrlShortener struct {
	hasher hash.Hash64
}

func NewUrlShortener() Shortener {
	return &UrlShortener{hasher: murmur3.New64()}
}

func (us *UrlShortener) Shorten(url string, userId string) (string, error) {
	us.hasher.Reset()

	ts := time.Now().UnixMilli()

	if _, err := us.hasher.Write([]byte(url)); err != nil {
		return emptyString, err
	}

	if _, err := us.hasher.Write([]byte(strconv.FormatInt(ts, 10))); err != nil {
		return emptyString, err
	}

	if _, err := us.hasher.Write([]byte(userId)); err != nil {
		return emptyString, err
	}

	return hex.EncodeToString(us.hasher.Sum(nil))[:8], nil
}
