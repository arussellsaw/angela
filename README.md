# Angela

raft based distributed merkel tree for model caching

our use case will be our api instances writing model updates to a raft log and constructing a merkle tree of object hashes, allowing quick distributed validation of model consistency, and smaller delta updates for follower nodes.

raft implementation taken almost verbatim from github.com/otoolep/hraftd

#### TODO

* implement http handlers for checking the merkle tree and performing delta updates based on that
* enable writing from hosts other than the leader?
