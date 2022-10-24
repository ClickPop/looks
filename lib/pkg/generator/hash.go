package generator

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
)

func checkHashes(assets []GeneratedAsset) ([]string, error) {
	hashes := make(map[string]bool)
	collisions := make([]string, 0)
	for _, asset := range assets {
		data := asset.Image.Bytes()
		reader := bytes.NewReader(data)
		hasher := md5.New()
		_, err := io.Copy(hasher, reader)
		if err != nil {
			return nil, err
		}
		hash := hex.EncodeToString(hasher.Sum(nil))
		if _, ok := hashes[hash]; ok {
			collisions = append(collisions, hash)
		} else {
			hashes[hash] = true
		}
	}
	if len(collisions) > 0 {
		return collisions, errors.New("collisions found")
	}
	return nil, nil
}
