# fetc

fetc is the command line tool of the golang template engine fet.

## Usage

```bash
# install fetc
GO111MODULE=on go get -u -v github.com/fefit/fetc

# in the parent directory of your templates directory, initialize the fet configs
fetc init

# watch the file changes and recompile
fetc watch

# compile the files, if no file given, compile all the template files.
fetc compile x.html xx.html
```
