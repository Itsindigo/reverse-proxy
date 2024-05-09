package crypto

import "encoding/base64"

type StringOrByteSlice interface {
	string | []byte
}

func B64Encode[T StringOrByteSlice](s T) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func B64Decode(s string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
