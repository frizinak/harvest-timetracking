# Harvest timetracking

Harvest API library and a set of commands

## Configuration

~/.timetracking (will be created first time you run timetracking)

```
{
    "account_id": "654321",
    "token": "abc-token-lala",
    "forecast_account_id: "",
    "weekdays_off": [
        "saturday",
        "sunday"
    ],
    "exclude_dates": [
        "2018-11-01",
        "2018-07-04",
    ]
}
```

`exclude_dates` are dates, like weekends, whose tracked hours are moved to the previous workday, if any.

Format YYYY-MM-DD obviously, as it is the only way a date should be formatted.

## Commands

### help
```
Available commands:
  help                 - print list of commands
  off                  - get a list of days off using the forecast api
  tracking             - show tracked hours
  version              - print version
```

### tracking

Retrieves all timetracking entries since `-from | now` and sums them
per day (going back `-days | 20`).

```
  -days int
        Amount of days to retrieve time entries for (default 20)
  -from string
        Custom date to start at [YYYY-MM-DD or end-of-week or next-week]
  -group string
        Group results by day|week|month|year (default "day")
  -hours int
        Amount of hours in a single workweek (default: from harvest api)
  -uid int
        The user id of the user to fetch 
```


#### Examples:

How many hours will I need next week?
```
$> date
Sat Nov 24 16:06:15 UTC 2018

$> timetracking tracking -from next-week -days 10
Running for Bob
ID: 123456
Week: 38h00
Over 10 days: 76h00
From: Sun Dec 02 2018

Fri Nov 30 2018 -  0h00
Thu Nov 29 2018 -  0h00
Wed Nov 28 2018 -  0h00
Tue Nov 27 2018 -  0h00
Mon Nov 26 2018 -  0h00
Fri Nov 23 2018 -  6h31
Thu Nov 22 2018 - 10h36
Wed Nov 21 2018 -  9h10
Tue Nov 20 2018 -  8h12
Mon Nov 19 2018 -  8h01

Total: 42h32 / 76h00 (55.97%)
33h27 remaining...
```

How many hours have I tracked in the last 5 days on which I actually worked?
```
$> timetracking tracking -days 5
Running for Bob
ID: 123456
Week: 38h00
Over 5 days: 38h00
From: Sat Nov 24 2018

Fri Nov 23 2018 -  6h31
Thu Nov 22 2018 - 10h36
Wed Nov 21 2018 -  9h10
Tue Nov 20 2018 -  8h12
Mon Nov 19 2018 -  8h01

Total: 42h32 / 38h00 (111.95%)
target reached!
```

How many hours have I tracked in the last 10 days, not ignoring empty tracking days?
```
$> timetracking tracking -days 10
Running for Bob
ID: 123456
Week: 38h00
Over 10 days: 76h00
From: Sat Nov 24 2018

Fri Nov 23 2018 -  6h31
Thu Nov 22 2018 - 10h36
Wed Nov 21 2018 -  9h10
Tue Nov 20 2018 -  8h12
Mon Nov 19 2018 -  8h01
Fri Nov 16 2018 -  0h00
Thu Nov 15 2018 -  0h00
Wed Nov 14 2018 -  0h00
Tue Nov 13 2018 -  1h56
Mon Nov 12 2018 -  5h53

Total: 50h22 / 76h00 (66.28%)
25h37 remaining...
```

### off

Get forecast days off (only works if a forcast_account_id is stored in ~/.timetracking)

https://forecastapp.com/[forecast_account_id]/schedule

and optionally `-save` them to ~/.timetracking in `excluded_dates`

```
  -hours int
        Amount of hours 'Time Off' should last for it to be an entire day off. (default 7)
  -project string
        Name of the 'Time Off' project (default "Time Off")
  -save
        Save in ~/.timetracking
  -uid int
        The forecast user id of the user to fetch time-off entries for
```
