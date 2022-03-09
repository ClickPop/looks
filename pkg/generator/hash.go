package generator

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func checkHashes(outputDir string) ([]string, error) {
	hashes := make(map[string]string)
	collisions := make([]string, 0)
	log.Println("Checking hashes for collisions...")
	_, err := os.Stat(outputDir)
	if os.IsNotExist(err) {
		os.Mkdir(outputDir, 0777)
	} else if err != nil {
		return []string{}, err
	}
	err = filepath.WalkDir(outputDir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		if strings.Contains(d.Name(), ".png") {
			data, _ := os.ReadFile(path)
			reader := bytes.NewReader(data)
			hasher := md5.New()
			_, err := io.Copy(hasher, reader)
			if err != nil {
				return err
			}
			hash := hex.EncodeToString(hasher.Sum(nil))
			if hashes[hash] != "" {
				log.Println("COLLISION", hashes[hash], path)
				collisions = append(collisions, path)
			} else {
				hashes[hash] = path
			}
		}
		return err
	})
	if err != nil {
		return []string{}, err
	}
	if len(collisions) > 0 {
		return collisions, err
	}
	log.Println("All hashes unique")
	return []string{}, nil
}