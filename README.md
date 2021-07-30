# Cron-parser
CLI tool to parse the cron expression passed as an input string.

## Table of Contents

- [Setup](#setup)
- [Usage](#usage)
- [Other](#other)

### Setup

`go get github.com/akroshchenko/cron-parser`

### Usage

---
**NOTE**

IMPORTANT: currently only one expression type(("n-n", "n/n", "n,n")) to generate time sequence in scope of one field (minute, hours etc) can be used.

IMPORTANT: Program may not correcly work with days of the month. It does not check the correct number of days in the month.

---

Example:
```
cron-parser ＂*/15 0 1,15 * 1-5 /usr/bin/find"

Output:
minute        0 15 30 45
hour          0
day of month  1 15
month         1 2 3 4 5 6 7 8 9 10 11 12
day of week   1 2 3 4 5
command       /usr/bin/find
```

### Other

https://cronjob.xyz/ - usefull online service to check the correctnes of working the program