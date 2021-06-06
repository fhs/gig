# gig

Gig is a git client written in pure Go (using
[go-git](https://github.com/go-git/go-git)). The main motivation was to
create a git client for platforms that are not supported by the official
git client (mainly Plan 9). Gig tries to be compatible with git CLI,
so anyone familiar with the official git client will already know how
to use gig.

## Install

```
GO111MODULE=on go get github.com/fhs/gig/cmd/git-upload-pack
GO111MODULE=on go get github.com/fhs/gig/cmd/git-receive-pack
GO111MODULE=on go get github.com/fhs/gig/cmd/gig
```

## Testing on Plan 9

To run tests on Plan 9, you need to first run this first:
```
go mod edit -replace 'github.com/rogpeppe/go-internal@v1.6.2=github.com/fhs/go-internal@v1.6.3-0.20201122174144-815d671f4ff9'
```
This is due to https://github.com/rogpeppe/go-internal/pull/115.

## See also
* https://github.com/oridb/git9
* https://github.com/driusan/dgit
* [Wrapping Git in rc shell](https://blog.gopheracademy.com/advent-2014/wrapping-git/)
* [Git port to 9legacy by Kyohei Kadota](https://9fans.topicbox.com/groups/9fans/Te3752ec266e3a002-M7286f7236d8aab10096f7946/9fans-git-client)
