package generator

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"os"

	conf "github.com/clickpop/looks/pkg/config"
)

func storeFile(config *conf.Config, img image.Image, jsonData []byte, genMeta bool, i int) error {
	out, err := os.Create(fmt.Sprintf("%s/%d.png", config.Output.Local.Directory, i))
	if err != nil {
		return err
	}
	png.Encode(out, img)
	out.Close()
	log.Printf("Image #%d.png created\n", i)
	if genMeta {
		err = os.WriteFile(fmt.Sprintf("./%s/%d.json", config.Output.Local.Directory, i), jsonData, 0666)
		if err != nil {
			return err
		}
		log.Printf("Metadata #%d.json created\n", i)
	}
	return nil
}