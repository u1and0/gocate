Parallel search for local files and directories.

# Usage


## Installation

```
$ go get github.com/u1and0/gocate
```


## Generate database

```
# $ gocate -init /var/lib/mlocate <- Not yet implemented

$ sudo updatedb
# Generate /var/lib/mlocate/mlocate.db

$ sudo updatedb -U /mnt --output /var/lib/mlocate/another.db
# Generate /var/lib/mlocate/another.db
```

Generate some database files to /var/lib/mlocate


## Register `LOCATE_PATH`

```
$ export LOCATE_PATH="/var/lib/mlocate/mlocate.db:/var/lib/mlocate/another.db"
```

## Searching

```
$ gocate keyword
```

Search the path included "keyword"

```
$ LOCATE_PATH=$(sudo find /var/lib/mlocate -name '*.db' | paste -sd:) gocate keyword
$ gocate keyword -d $(sudo find /var/lib/mlocate -name '*.db' | paste -sd:)
```

Those of two commands are same command.

## `gocate` option

```
$ gocate -h
parallel find files by name

Usage of gocate
	gocate [OPTION]... PATTERN...

-v, -version
	Show version
-d, -database string
	Path of locate database file (ex: /path/something.db:/path/another.db)
-- [OPTION]...
	locate command option
```

Recomend to separate db files for multi thread search.

