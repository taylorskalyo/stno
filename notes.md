Potential usage:

```sh
stno add       path [--infile path][--auto-resolve-conflicts]
stno list|ls   [path...][--where cond][--select toml-path...][--order toml-path][--format type]
stno edit      [path...][--where cond][--select toml-path...][--order toml-path]
stno remove|rm [path...][--where cond][--select toml-path...]
stno move|mv   src dest
```

Use path as UID instead of datetime + title + sequence
Add bash completion for paths (like password-store)

Datastore changes:
- No need to initialize a datastore with a path
- No need to create a new unique entry ID
- New interface is:
  - List(string) ([]string, error)
  - NewWriteCloser(string) (io.WriteCloser, error)
  - NewReadCloser(string) (io.ReadCloser, error)
  - Remove(string) error

Add
- Touch file at path
- Write out template to tmpfile
- Loop
  - Open tmpfile in editor
  - Open tmpfile and lint; if errors prompt to fix (y/n)
    - If y, go back to loop
    - Else exit
- Copy tree to file

List
- If dir, glob path/\*
- For each path
  - Filter based on where condition; return on channel
  - Select subtree based on select path; return on channel
  - Order based on order path; return array
- Switch on format type
  - stdout (default): output to terminal
  - toml: single toml file
  - etc.
