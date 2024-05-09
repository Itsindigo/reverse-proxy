package crypto

import "encoding/base64"

func B64Encode(s []byte) string {
	return base64.StdEncoding.EncodeToString(s)
}

func B64Decode(s string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
