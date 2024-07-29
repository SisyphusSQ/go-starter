package utils

import "github.com/speps/go-hashids"

type Hash struct {
	secret string
	length int
}

func New(secret string, length int) *Hash {
	return &Hash{
		secret: secret,
		length: length,
	}
}

func (h *Hash) HashidsEncode(params []int) (string, error) {
	hd := hashids.NewData()
	hd.Salt = h.secret
	hd.MinLength = h.length

	hashStr, err := hashids.NewWithData(hd).Encode(params)
	if err != nil {
		return "", err
	}

	return hashStr, nil
}

func (h *Hash) HashidsDecode(hash string) ([]int, error) {
	hd := hashids.NewData()
	hd.Salt = h.secret
	hd.MinLength = h.length

	ids, err := hashids.NewWithData(hd).DecodeWithError(hash)
	if err != nil {
		return nil, err
	}

	return ids, nil
}
