package generator

import (
	"fmt"
	"os"

	conf "github.com/clickpop/looks/pkg/config"
)

func StoreFile(config *conf.Config, asset *GeneratedAsset) error {
	var err error
	if asset.Image != nil {
		err = os.WriteFile(fmt.Sprintf("./%s/%s.json", config.Output.Directory, asset.Name), asset.Image.Bytes(), 0666)
		if err != nil {
			return err
		}
	}
	if config.Output.IncludeMeta && config.Output.MetaFormat == conf.JSON && asset.Meta != nil {
		err = os.WriteFile(fmt.Sprintf("./%s/%s.json", config.Output.Directory, asset.Name), asset.Meta.Bytes(), 0666)
		if err != nil {
			return err
		}
	}
	return nil
}
