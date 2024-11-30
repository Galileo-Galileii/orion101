package gz

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
)

func Decompress(out any, input []byte) error {
	gzReader, err := gzip.NewReader(bytes.NewBuffer(input))
	if err != nil {
		return err
	}
	data, err := io.ReadAll(gzReader)
	if err != nil {
		return err
	}

	switch v := out.(type) {
	case *string:
		*v = string(data)
		return nil
	case *[]byte:
		*v = data
		return nil
	default:
		return json.Unmarshal(data, out)
	}
}

func Compress(obj any) ([]byte, error) {
	var (
		input []byte
		err   error
	)
	switch v := obj.(type) {
	case []byte:
		input = v
	case string:
		input = []byte(v)
	default:
		input, err = json.Marshal(obj)
		if err != nil {
			return nil, err
		}
	}

	var (
		out   bytes.Buffer
		gzOut = gzip.NewWriter(&out)
	)

	_, err = gzOut.Write(input)
	if err != nil {
		return nil, err
	}
	if err := gzOut.Close(); err != nil {
		return nil, err
	}
	return out.Bytes(), err
}
