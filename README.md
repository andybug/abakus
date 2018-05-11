# abakus
_abakus_ (andybug backup script) is a git-like repository that tracks files via
snapshots in time. Snapshots are analogous to commits, with the major caveat
that snapshots save the entire state of the tree instead of the delta. This
allows for history to be deleted when snapshots are too old or otherwise not
needed.

Ultimately, _abakus_ will be able to take snapshots locally, push to remote
storage (S3, GCS) with client-side encryption and compression, and pull from
remote storage. A table representing these features in the form of CLI commands
and their respective release milestones is shown below.

| Command       | Version | Completed |
| --- | --- | --- |
| init          |   0.1.0 | X         |
| create        |   0.1.0 | X         |
| lists         |   0.1.0 | X         |
| show          |   0.1.0 | X         |
| status        |   0.1.0 | X         |
| delete        |   0.2.0 |           |
| export        |   0.2.0 |           |
| history       |   0.2.0 |           |
| prune         |   0.2.0 |           |
| validate      |   0.2.0 |           |
| push          |   0.3.0 |           |
| pull          |   0.3.0 |           |
| remote-add    |   0.3.0 |           |
| remote-ls     |   0.3.0 |           |
| remote-show   |   0.3.0 |           |
| remote-delete |   0.3.0 |           |

## Installation

	go get -u "github.com/andybug/abakus/cmd/abakus"

## Usage

	> ls

	> abakus init
	New abakus repository initialized

	> touch a b c
	> abakus status
	No Snapshots
	added:       a
	added:       b
	added:       c

	> abakus create
	Snapshot created

	> abakus list
	ID    TIME              MERKLE      FILES    SIZE
	1     53 seconds ago    bf23c8fd    3        0 B

	> abakus show 1
	PATH    HASH                                                                SIZE    MODE
	a       0e5751c026e543b2e8ab2eb06099daa1d1e5df47778f7787faab45cdf12fe3a8    0 B     644
	b       0e5751c026e543b2e8ab2eb06099daa1d1e5df47778f7787faab45cdf12fe3a8    0 B     644
	c       0e5751c026e543b2e8ab2eb06099daa1d1e5df47778f7787faab45cdf12fe3a8    0 B     644

### Ignoring Files
Abakus looks for a `.abakusignore` file in each directory that contains file
exclusion rules (much like`.gitignore` files).

	version: 1
	excludes:
	  # only exclude .git in this directory
	  - /.git
	  # exclude all files/dirs named temp recursively
	  - temp
