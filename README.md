[locate](https://en.wikipedia.org/wiki/Locate_(Unix)) command by goroutine.
Fast search for local files and directories.

# Usage

```
$ gocate -init /path/to/index/root
$ gocate -i "ls"
```

# Installation

```
$ go get github.com/u1and0/gocate
```



# MEMO
600万ファイル超があるNSにて実験

```
$ export LOCATE_PATH=$(sudo find /var/lib/mlocate -name '*.db' | paste -sd:)
```

として、gocateとlocateの実行速度比較

```
$ time ./gocate grep-server | wc -l
68
./gocate grep-server  0.00s user 0.03s system 0% cpu 13.332 total
wc -l  0.00s user 0.00s system 0% cpu 13.331 total

$ time locate grep-server | wc -l
locate: stat () `'/var/lib/mlocate/mlocate.db' できません: そのようなファイルやディレクトリはありません
68
locate grep-server  23.64s user 0.16s system 99% cpu 23.843 total
wc -l  0.00s user 0.00s system 0% cpu 23.842 total
```

Parallel locate: 13.3sec
Normal locate  : 23.4sec

並列化すると実行時間を40%程削減できた
