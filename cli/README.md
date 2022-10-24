# Looks CLI

# Usage
Run `looks --help` to get a list of all commands or `looks <command> --help` for help with a specific command

# Config
Looks supports configuration via a `config.json` in the directory it is being called from or from the root directory of a project when using the API.

The configuration schema is as follow (WIP)

```json
{
  "input": {
    "local": {
      "filename": "%s-%s.png",
      "pathname": "pieces"
    }
  },
  "output": {
    "local": {
      "directory": "output"
    },
    "image-count": 1
  },
  "settings": {
    "max-workers": 3,
    "piece-order": [
      "background",
      "color"
    ],
    "rarity": {
      "order": [
        "common"
      ],
      "chances": {
        "common": 1
      }
    }
  },
  "attributes": {
    "background": {
      "pieces": {
        "dark": {
          "rarity": "common"
        },
        "green": {
          "rarity": "common"
        }
      }
    },
    "color": {
      "pieces": {
        "golden": {
          "rarity": "common"
        },
        "blue": {
          "rarity": "common"
        }
      }
    }
  }
}
```