# Data

All data was parsed from https://github.com/BradyStephenson/bible-data/tree/main


## Example nodes

Nodes that have been scraped is directly from [commandments](https://raw.githubusercontent.com/BradyStephenson/bible-data/refs/heads/main/BibleData-Commandments.csv)

```bash
$ python3 convert.py https://raw.githubusercontent.com/BradyStephenson/bible-data/refs/heads/main/BibleData-Commandments.csv commandments.csv
Done. Wrote transformed CSV to: commandments.csv
```


```bash
$ sqlite3 ./example/bible.sqlite
.mode csv
.import --skip 1 ./example/commandments.csv nodes
```