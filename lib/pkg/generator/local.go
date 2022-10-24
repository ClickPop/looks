package generator

import (
	"fmt"
	"image"
	"image/png"
	"os"

	conf "github.com/clickpop/looks/pkg/config"
)

func storeFile(config *conf.Config, img image.Image, jsonData []byte, i int) error {
	out, err := os.Create(fmt.Sprintf("%s/%d.png", config.Output.Local.Directory, i))
	if err != nil {
		return err
	}
	png.Encode(out, img)
	out.Close()
	if config.Output.IncludeMeta && config.Output.MetaFormat == conf.JSON && jsonData != nil {
		err = os.WriteFile(fmt.Sprintf("./%s/%d.json", config.Output.Local.Directory, i), jsonData, 0666)
		if err != nil {
			return err
		}
	}
	return nil
}
