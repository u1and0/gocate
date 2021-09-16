Parallel search for local files and directories.


# Installation

```
$ go install github.com/u1and0/gocate@latest
```

# Requirement

* mlocate
* updatedb


# Usage

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

Indexing for directories, neither files nor symbolic links.


## Search

```
$ gocate keyword
```

Search the path included "keyword"


## `gocate` option

```
$ gocate -h
parallel find files by name

Usage of gocate
gocate [OPTION]... PATTERN...

-v, -version
Show version
  -d, -database DIRECTORY
Path of locate database directory (default: /var/lib/mlocate)
  -init
  updatedb mode
  -U, -database-root DIRECTORY
  Store only results of scanning the file system subtree rooted at PATH  to  the  generated  database.
  -o, -output DIRECTORY
Write the database to DIRECTORY instead of using the default database directory. (default: /var/lib/mlocate)
  -dryrun
  Just print command, do NOT run updatedb command.
  -- [OPTION]...
  locate or updatedb command option
```

## Dryrun
Print only command if `-dryrun` option


```
$ gocate -dryrun keyword
/usr/sbin/locate -d /var/lib/mlocate/var.db keyword
/usr/sbin/locate -d /var/lib/mlocate/boot.db keyword
/usr/sbin/locate -d /var/lib/mlocate/dev.db keyword
/usr/sbin/locate -d /var/lib/mlocate/etc.db keyword
/usr/sbin/locate -d /var/lib/mlocate/run.db keyword
/usr/sbin/locate -d /var/lib/mlocate/mnt.db keyword
/usr/sbin/locate -d /var/lib/mlocate/srv.db keyword
/usr/sbin/locate -d /var/lib/mlocate/home.db keyword
/usr/sbin/locate -d /var/lib/mlocate/sys.db keyword
/usr/sbin/locate -d /var/lib/mlocate/tmp.db keyword
/usr/sbin/locate -d /var/lib/mlocate/usr.db keyword
/usr/sbin/locate -d /var/lib/mlocate/opt.db keyword
/usr/sbin/locate -d /var/lib/mlocate/proc.db keyword
/usr/sbin/locate -d /var/lib/mlocate/root.db keyword
```

```
$ gocate -init -dryrun
warning: /.dockerenv is not directory, it will be ignored for indexing.
warning: /bin is not directory, it will be ignored for indexing.
warning: /lib is not directory, it will be ignored for indexing.
warning: /lib64 is not directory, it will be ignored for indexing.
warning: /sbin is not directory, it will be ignored for indexing.
/usr/sbin/updatedb -U /var --output /var/lib/mlocate/var.db
/usr/sbin/updatedb -U /boot --output /var/lib/mlocate/boot.db
/usr/sbin/updatedb -U /dev --output /var/lib/mlocate/dev.db
/usr/sbin/updatedb -U /etc --output /var/lib/mlocate/etc.db
/usr/sbin/updatedb -U /tmp --output /var/lib/mlocate/tmp.db
/usr/sbin/updatedb -U /sys --output /var/lib/mlocate/sys.db
/usr/sbin/updatedb -U /usr --output /var/lib/mlocate/usr.db
/usr/sbin/updatedb -U /root --output /var/lib/mlocate/root.db
/usr/sbin/updatedb -U /run --output /var/lib/mlocate/run.db
/usr/sbin/updatedb -U /home --output /var/lib/mlocate/home.db
/usr/sbin/updatedb -U /srv --output /var/lib/mlocate/srv.db
/usr/sbin/updatedb -U /mnt --output /var/lib/mlocate/mnt.db
/usr/sbin/updatedb -U /opt --output /var/lib/mlocate/opt.db
/usr/sbin/updatedb -U /proc --output /var/lib/mlocate/proc.db
```



## Available whole locate & updatedb option

```
$ gocate -- -i --regex ".*ls.*cpp$"                       # locate mode
$ gocate -init -- --debug-pruning --add-prunenames NAMES  # updatedb mode
```

Use command option after double dash `--`.
See `man locate`, `man updatedb`.


## Parallel updatedb

Recomend to separate db files for multi thread searching.
Use `gocate -init`

```
$ sudo gocate -init
/usr/sbin/updatedb -U /var --output /var/lib/mlocate/var.db
/usr/sbin/updatedb -U /root --output /var/lib/mlocate/root.db
/usr/sbin/updatedb -U /run --output /var/lib/mlocate/run.db
/usr/sbin/updatedb -U /srv --output /var/lib/mlocate/srv.db
/usr/sbin/updatedb -U /proc --output /var/lib/mlocate/proc.db
/usr/sbin/updatedb -U /sys --output /var/lib/mlocate/sys.db
/usr/sbin/updatedb -U /etc --output /var/lib/mlocate/etc.db
/usr/sbin/updatedb -U /home --output /var/lib/mlocate/home.db
/usr/sbin/updatedb -U /tmp --output /var/lib/mlocate/tmp.db
/usr/sbin/updatedb -U /boot --output /var/lib/mlocate/boot.db
/usr/sbin/updatedb -U /mnt --output /var/lib/mlocate/mnt.db
/usr/sbin/updatedb -U /usr --output /var/lib/mlocate/usr.db
/usr/sbin/updatedb -U /dev --output /var/lib/mlocate/dev.db
/usr/sbin/updatedb -U /opt --output /var/lib/mlocate/opt.db
```

By default, `updatedb` path are all directory under root directory.


## Select directory for database
Use `-U DIRECTORY` or `-database-root DIRECTORY` option.

```
$ sudo gocate -init -U /mnt -U /usr
warning: /usr/lib64 is not directory, it will be ignored for indexing.
warning: /usr/sbin is not directory, it will be ignored for indexing.
warning: /var/lock is not directory, it will be ignored for indexing.
warning: /var/mail is not directory, it will be ignored for indexing.
warning: /var/run is not directory, it will be ignored for indexing.
/usr/sbin/updatedb -U /usr/lib32 --output /var/lib/mlocate/usr_lib32.db
/usr/sbin/updatedb -U /usr/include --output /var/lib/mlocate/usr_include.db
/usr/sbin/updatedb -U /var/tmp --output /var/lib/mlocate/var_tmp.db
/usr/sbin/updatedb -U /var/empty --output /var/lib/mlocate/var_empty.db
/usr/sbin/updatedb -U /var/games --output /var/lib/mlocate/var_games.db
/usr/sbin/updatedb -U /usr/local --output /var/lib/mlocate/usr_local.db
/usr/sbin/updatedb -U /var/lib --output /var/lib/mlocate/var_lib.db
/usr/sbin/updatedb -U /usr/bin --output /var/lib/mlocate/usr_bin.db
/usr/sbin/updatedb -U /usr/share --output /var/lib/mlocate/usr_share.db
/usr/sbin/updatedb -U /var/local --output /var/lib/mlocate/var_local.db
/usr/sbin/updatedb -U /var/cache --output /var/lib/mlocate/var_cache.db
/usr/sbin/updatedb -U /var/db --output /var/lib/mlocate/var_db.db
/usr/sbin/updatedb -U /var/log --output /var/lib/mlocate/var_log.db
/usr/sbin/updatedb -U /usr/src --output /var/lib/mlocate/usr_src.db
/usr/sbin/updatedb -U /var/opt --output /var/lib/mlocate/var_opt.db
/usr/sbin/updatedb -U /var/spool --output /var/lib/mlocate/var_spool.db
/usr/sbin/updatedb -U /usr/lib --output /var/lib/mlocate/usr_lib.db
```

You can make database for the directories under /mnt and /usr.
**Note that, files and symbolic links directly under the selected directory will be ignored by the database.**
That will display as warning.

## Output directory for database

You'd like to change the directory where db files are placed, use `-o`, `-output` option.

By default, it will be /var/lib/mlocate (Same as updatedb)

```
$ sudo gocate -init -o /home/myhome/.locate
```


If you select the output directory except /var/lib/mlocate, you must use `-d DIRECTORY` option for search.

By default, it will be /var/lib/mlocate (Same as updatedb)

```
$ gocate -d /home/myhome/.locate keyword
```


[日本語README](https://qiita.com/u1and0/items/964be5817da800b82603)
