package cli

import (
	b64 "encoding/base64"
)

func B64Encode(byteValue []byte) string {
	return b64.StdEncoding.EncodeToString(byteValue)
}

func B64Decode(value string) ([]byte, error) {
	scValueByte, err := b64.StdEncoding.DecodeString(value)
	if err != nil {
		return nil, err
	}
	return scValueByte, nil
}



//read file or value  -- type!!!
