package generator

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func checkHashes(outputDir string) error {
	hashes := make(map[string]string)
	log.Println("Checking hashes for collisions...")
	_, err := os.Stat(outputDir)
	if os.IsNotExist(err) {
		os.Mkdir(outputDir, 0777)
	} else if err != nil {
		return err
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
				return fmt.Errorf("COLLISION: %s & %s", hashes[hash], path)
			}
			hashes[hash] = path
		}
		return err
	})
	if err != nil {
		return err
	}
	log.Println("All hashes unique")
	return nil
}