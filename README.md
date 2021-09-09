Parallel search for local files and directories.

# Usage


## Installation

```
$ go install github.com/u1and0/gocate@latest
```


## Generate database

```
$ sudo gocate -init
```

Generate some database files to /var/lib/mlocate

`sudo gocate -init` is same as below

```
$ sudo updatedb -U /var --output /var/lib/mlocate/var.db
$ sudo updatedb -U /root --output /var/lib/mlocate/root.db
$ sudo updatedb -U /run --output /var/lib/mlocate/run.db
$ sudo updatedb -U /srv --output /var/lib/mlocate/srv.db
$ sudo updatedb -U /proc --output /var/lib/mlocate/proc.db
$ sudo updatedb -U /sys --output /var/lib/mlocate/sys.db
$ sudo updatedb -U /etc --output /var/lib/mlocate/etc.db
$ sudo updatedb -U /home --output /var/lib/mlocate/home.db
$ sudo updatedb -U /tmp --output /var/lib/mlocate/tmp.db
$ sudo updatedb -U /boot --output /var/lib/mlocate/boot.db
$ sudo updatedb -U /mnt --output /var/lib/mlocate/mnt.db
$ sudo updatedb -U /usr --output /var/lib/mlocate/usr.db
$ sudo updatedb -U /dev --output /var/lib/mlocate/dev.db
$ sudo updatedb -U /opt --output /var/lib/mlocate/opt.db
```

But indexing only directory NOT files, symbolic links


## Register `LOCATE_PATH`

```
$ export LOCATE_PATH="/var/lib/mlocate/var.db:/var/lib/mlocate/root.db"
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
-init
	updatedb mode
-U, -database-root
	Store only results of scanning the file system subtree rooted at PATH  to  the  generated  database.
-dryrun
	Just print command, do NOT run updatedb command.
-- [OPTION]...
	locate command option
```

Recomend to separate db files for multi thread searching.
