package encoding

import (
	"time"

	"github.com/speps/go-hashids/v2"
)

type Encoder struct {
	salt      string
	minLength int
}

func NewEncoder(salt string, minLength int) *Encoder {
	return &Encoder{
		salt:      salt,
		minLength: minLength,
	}
}

func (e *Encoder) EncodeID(id int64) (string, error) {
	h, err := e.createNewHashID()
	if err != nil {
		return "", err
	}
	encodedID, err := h.EncodeInt64([]int64{id})
	if err != nil {
		return "", err
	}

	return encodedID, err
}

func (e *Encoder) createNewHashID() (*hashids.HashID, error) {
	data := hashids.NewData()
	data.Salt = e.salt
	data.MinLength = e.minLength

	return hashids.NewWithData(data)
}

func (e *Encoder) GenerateNewId() (string, error) {
	return e.EncodeID(time.Now().UnixMilli())
}
