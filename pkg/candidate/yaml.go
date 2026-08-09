package candidate

import (
	"bytes"
	"io"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

func readStrictYAMLWithDigest(path string, destination any) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)
	if err := decoder.Decode(destination); err != nil {
		return "", errors.Wrap(err, "decode strict YAML")
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			return "", errors.New("YAML contains multiple documents")
		}
		return "", errors.Wrap(err, "check YAML trailing document")
	}
	return digestBytes(data), nil
}
