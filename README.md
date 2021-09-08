# Looks CLI

# Usage
Run `looks --help` to get a list of all commands or `looks <command> --help` for help with a specific command

# Config
Looks supports configuration via a `config.json` in the directory it is being called from or from the root directory of a project when using the API.

The configuration schema is as follow (WIP)

```json
{
  "piece-order": ["background", "foreground"] /** Array of strings representing the piece names for layering order. */,
  "filename": "%s-%s.png" /** Format string to use for loading the individual assets. You will need a %s for each sub-category of assets */,
  "pathname": "pieces" /** The pathname where your input files can be found */,
  "output-directory": "rats" /** Output directory where generated assets will be stored */,
  "output-image-count": 10 /** How many images to generate */,
  "max-workers": 3 /** The number of workers to run. Generally you want to keep this below like 6 or 7 unless you have a really beefy machine, as generating can use a lot of memory depending on the number of layers */,
  "attributes": /** Object of key value pairs that signify your different pieces and their attributes */ {
    "background": {
      "dark": {
        "rarity": 6
      },
      "light": {
        "rarity": 4
      }
    },
    "foreground": {
      "pink": {
        "rarity": 3,
      },
      "green": {
        "rarity": 4,
        "cuteness": 1
      }
    }
  }
}
```