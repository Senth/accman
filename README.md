# accman
Accounting Manager

## Usage

### `parser.json`
Contains parsers that automatically books the verification correctly.

```
{
    "identifier": "Tax account",
    // optional, by default will select layout
    "options": {
      "pdfLayout": "layout" // [layout, raw]
    },
    // containing the groups date, name, and amount
    "regexp": " (?P<date>\\d{6})\\s+(?P<name>[\\w\\d- åäöÅÄÖ]+?)\\s{2}(?P<amount>-?\\d+ \\d+|-?\\d+)\\s+[\\d ]*\\n",
    // for date formatting, see https://stackoverflow.com/questions/20234104/how-to-format-current-time-using-a-yyyymmddhhmmss-format
    "dateFormat": "060102",
    "verificationParsers": [
      {
        // searches for the matched name to see if this should create a verification.
        "identifier": "Preliminary tax",
        // optional name, will use the 'name' from regexp group if omitted.
        "name": "Preliminary taxes",
        "type": 32,
        "accountFrom": 1630,
        "accountTo": 2518
      }
    ]
  }
```

## Installation
TODO

### Dependencies
For parsing files additional dependencies are required on linux

```bash
$ sudo apt-get install poppler-utils wv unrtf tidy
```

