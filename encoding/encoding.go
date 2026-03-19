package encoding

import "encoding/base64"

func Encode(value []byte) string {
	return base64.StdEncoding.EncodeToString(value)
}

func Decode(value string) ([]byte, error) {
	res, err := base64.StdEncoding.DecodeString(value)

	if err != nil {
		return nil, err
	}

	return res, nil
}
