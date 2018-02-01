# Angela

![angela merkle tree](https://i.imgur.com/m26QGtS.jpg)

raft based distributed merkel tree for model caching

our use case will be our api instances writing model updates to a raft log and constructing a merkle tree of object hashes, allowing quick distributed validation of model consistency, and smaller delta updates for follower nodes.

raft implementation taken almost verbatim from github.com/otoolep/hraftd

to run example:

start master


```./cmd -id node0 ~/node0```


start followers


`./cmd -id node1 -haddr :11001 -raddr :12001 -join :11000 ~/node1`


`/cmd -id node2 -haddr :11002 -raddr :12002 -join :11000 ~/node2`

#### TODO

* implement http handlers for checking the merkle tree and performing delta updates based on that
* enable writing from hosts other than the leader?
