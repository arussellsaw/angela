# Angela

raft based distributed merkel tree for model caching

our use case will be our api instances writing model updates to a raft log and constructing a merkel tree of object hashes, allowing quick distributed validation of model consistency, and smaller delta updates for follower nodes.
