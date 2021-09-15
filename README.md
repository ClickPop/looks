# Looks CLI

# Usage
Run `looks --help` to get a list of all commands or `looks <command> --help` for help with a specific command

# Config
Looks supports configuration via a `config.json` in the directory it is being called from or from the root directory of a project when using the API.

The configuration schema is as follow (WIP)

```json
{
  "input": { /* Object: key-value pairs that represent the input settings */
    "local": {
      "filename": "%s-%s.png", /* string: Format string to use for loading the individual assets. You will need a %s for each sub-category of assets */
      "pathname": "pieces" /* string: Pathname where input files can be found */
    }
  },
  "output": { /* Object: key-value pairs that represent the input settings */
    "local": { "directory": "output" }, /* WIP: Object: Local Settings */
    "internal": true, /* WIP: bool: Use internal output (for sending to NFT Minting tool) */
    "image-count": 1 /* int: Number of images to generate */
  },
  "settings": {
    "max-workers": 3, /* int: Number of workers to run. Generally you want to keep this below like 6 or 7 unless you have a really beefy machine, as generating can use a lot of memory depending on the number of layers */
    "piece-order": ["background", "foreground"], /* []string: Represents the piece names for layering order (bottom to top) */
    "stats": { /* Object: key-value pairs that represent the stat settings */
      "cunning": { /* string: Object key is the stat "id" */
        "minimum": 0, /* int: Prevents the stat from falling below a given value */
        "maximum": 7, /* int: Prevents the stat from exceeding a given value */
        "name": "Cunning" /* string: Friendly name for stat */
      },
      "cuteness": {
        "minimum": 0,
        "maximum": 7,
        "name": "Cuteness"
      },
      "rattitude": {
        "minimum": 0,
        "maximum": 7,
        "name": "Rattitude"
      }
    },
    "rarity": {
      "order": ["common", "uncommon", "rare", "epic", "legendary"],
      "chances": {
        "common": 600,
        "uncommon": 250,
        "rare": 100,
        "epic": 30,
        "legendary": 20
      }
    }
  },
  "attributes": { /** WIP: Object: key-value pairs that signify your different pieces and their attributes */
    "background": {
      "dark": {
        "rarity": "common",
        "stats": {
          "rattitude": 2,
          "cunning": 1,
          "cuteness": -1
        }
      },
      "light": {
        "rarity": "uncommon",
        "stats": {
          "cunning": 2,
          "rattitude": -1
        }
      }
    },
    "foreground": {
      "pink": {
        "rarity": "epic",
        "stats": {
          "cuteness": 2
        }
      },
      "green": {
        "rarity": "common",
        "stats": {
          "cunning": 1
        }
      }
    }
  },
  "descriptions": { /** WIP: Object: key-value pairs to be used by the description generator */
    "template": "This little rat is a %s, that means %s. Their favorite hobbies include %s.", /* string: Format string to use as template for description generator */
    "hobbies-count": 3, /* int: Number of random hobbies to select */
    "fallback-primary-stat": "fallback", /* string: refrences the fallback type in the below "types" object when there is no primary stat */
    "types": { /* Object: key-value pairs establishing "types" based on the dominant stat. Use the stat "id" as the key */
      "cunning": { /* string: id of primary (dominant) stat to utilize for this type */
        "id": "lab-rat", /* string: id for this type */
        "name": "Lab Rat", /* string: Friendly name for this type */
        "descriptors": [ /* []string: Array of descriptors to be randomly selected by the description generator */
          "they’re constantly tinkering with things they probably shouldn’t",
          "they're smarter than the average rat",
          "they're as interested in where cheese is as why cheese is"
        ],
        "hobbies": [ /* []string: Array of hobbies to be randomly selected by the description generator, utilizes  */
          "reading scraps of thrown out books",
          "naming new constellations",
          "pondering"
        ]
      },
      "fallback": { /** See fallback-primary-stat above **/
        "id": "pack-rat",
        "name": "Pack Rat",
        "descriptors": [
          "they have a trash stash that would make you green with envy",
          "their motto is be prepared",
          "they always have a ketchup packet on hand"
        ],
        "hobbies": [
          "hoarding",
          "tossing back cookies",
          "digging through old computer towers for spare parts"
        ]
      }
    }
  }
} 
```