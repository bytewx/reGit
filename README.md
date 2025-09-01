# re-git

A minimal Git-like version control system written in Go.

## Features

- Initialize repository
- Add, commit, remove, show files
- Status, log, file history
- Diff, reset, list objects
- Push, pull, fetch, clone, merge (local directory simulation)
- Advanced queries: commit count, tracked files, file versions, etc.

## Usage

### Build

```sh
go build -o regit main.go cli.go
```

### Run

```sh
go run main.go cli.go <command> [args]
```

### Common Commands

- `init`  
  Initialize a new repository.

- `add <file>`  
  Stage a file.

- `commit "<message>"`  
  Commit staged files.

- `status`  
  Show staged files.

- `log`  
  Show commit log.

- `remove <file>`  
  Remove file from staging.

- `show <file>`  
  Show staged file contents.

- `ls-objects`  
  List tracked objects.

- `checkout`  
  Restore files from latest commit.

- `diff`  
  Show differences between staged and working files.

- `list-commits`  
  List all commits.

- `file-history <file>`  
  Show commit history for a file.

- `reset`  
  Clear staging area.

- `istracked <file>`  
  Check if file is tracked.

- `get-file-version <file> <commitIdx>`  
  Show file version from a specific commit.

- `commit-files <commitIdx>`  
  List files in a specific commit.

- `remove-object <oid>`  
  Remove object by ID.

- `commit-count`  
  Show number of commits.

- `find-commit-by-message "<msg>"`  
  Find commits by message substring.

- `find-file-oids <file>`  
  List all object IDs for a file.

- `restore-file-from-commit <file> <commitIdx>`  
  Restore file from a specific commit.

- `purge-unreferenced-objects`  
  Remove objects not referenced by any commit.

- `get-commit-message <commitIdx>`  
  Show commit message.

- `get-commit-date <commitIdx>`  
  Show commit date.

- `get-commit-oid-for-file <file> <commitIdx>`  
  Get object ID for file in commit.

- `list-all-tracked-files`  
  List all files ever tracked.

- `push <remote_path>`  
  Push local repo to remote directory.

- `pull <remote_path>`  
  Pull remote repo into local directory.

- `clone <remote_path> <target_path>`  
  Clone remote repo to target directory.

- `fetch <remote_path>`  
  Fetch objects from remote (no merge).

- `merge <remote_path>`  
  Merge remote log into local log.

- `merge-to-remote <remote_path>`  
  Merge local repo into remote repo.

## Notes

- Remote operations (`push`, `pull`, etc.) work with local directories, not real remote servers.
- All repository data is stored in `.regit` directory.
