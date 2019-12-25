# gig

Gig is an attempt at implementing a git-compatible command in pure Go (using [go-git](https://github.com/src-d/go-git)).

## Why?

This is intended for Plan 9, where git doesn’t work natively.
Gig requires these two PRs on Plan 9:
* https://github.com/src-d/go-billy/pull/78
* https://github.com/src-d/go-git/pull/1269

## See also
* https://github.com/oridb/git9
* https://github.com/driusan/dgit
* [Wrapping Git in rc shell](https://blog.gopheracademy.com/advent-2014/wrapping-git/)
* [Git port to 9legacy by Kyohei Kadota](https://9fans.topicbox.com/groups/9fans/Te3752ec266e3a002-M7286f7236d8aab10096f7946/9fans-git-client)
