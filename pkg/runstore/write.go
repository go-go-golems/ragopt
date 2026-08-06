package runstore

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func digestBytes(data []byte) string {
	sum := sha256.Sum256(data)
	return "sha256:" + hex.EncodeToString(sum[:])
}

func canonicalJSON(value any) ([]byte, error) {
	raw, err := json.Marshal(value)
	if err != nil {
		return nil, errors.Wrap(err, "marshal JSON")
	}
	return canonicalizeJSON(raw)
}

func canonicalizeJSON(raw []byte) ([]byte, error) {
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()
	var value any
	if err := decoder.Decode(&value); err != nil {
		return nil, errors.Wrap(err, "decode JSON")
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			return nil, errors.New("JSON contains multiple values")
		}
		return nil, errors.Wrap(err, "check JSON trailing data")
	}
	canonical, err := json.Marshal(value)
	if err != nil {
		return nil, errors.Wrap(err, "marshal canonical JSON")
	}
	return canonical, nil
}

func (run *Run) writeJSON(ctx context.Context, relative string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return errors.Wrapf(err, "marshal %s", relative)
	}
	return run.writeBytes(ctx, relative, append(data, '\n'))
}

func (run *Run) writeBytes(ctx context.Context, relative string, data []byte) error {
	path, err := joinWithin(run.dir, relative)
	if err != nil {
		return err
	}
	return atomicWrite(ctx, path, data)
}

func atomicWrite(ctx context.Context, path string, data []byte) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	directory := filepath.Dir(path)
	if err := os.MkdirAll(directory, 0o700); err != nil {
		return errors.Wrap(err, "create artifact parent directory")
	}
	temporary, err := os.CreateTemp(directory, ".artifact-*")
	if err != nil {
		return errors.Wrap(err, "create temporary artifact")
	}
	temporaryName := temporary.Name()
	defer func() {
		_ = temporary.Close()
		_ = os.Remove(temporaryName)
	}()
	if err := temporary.Chmod(0o600); err != nil {
		return errors.Wrap(err, "set temporary artifact permissions")
	}
	if _, err := temporary.Write(data); err != nil {
		return errors.Wrap(err, "write temporary artifact")
	}
	if err := temporary.Sync(); err != nil {
		return errors.Wrap(err, "sync temporary artifact")
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := temporary.Close(); err != nil {
		return errors.Wrap(err, "close temporary artifact")
	}
	if err := os.Rename(temporaryName, path); err != nil {
		return errors.Wrap(err, "publish artifact")
	}
	return syncDirectory(directory)
}

func syncDirectory(path string) error {
	directory, err := os.Open(path)
	if err != nil {
		return errors.Wrap(err, "open artifact directory for sync")
	}
	if err := directory.Sync(); err != nil {
		_ = directory.Close()
		return errors.Wrap(err, "sync artifact directory")
	}
	return errors.Wrap(directory.Close(), "close artifact directory")
}
