# gstellar

Easy PG snapshot tool for development environment. Inspired by stellar in python.

## Description

Creates and restores snapshots on a postgresql db without pg_restore and pg_dump

### Installing

```bash
go install github.com/schaerli/gstellar@latest
```

### Executing program

* show commands
```bash
gstellar
```

* commands are
```bash
init
snapshot
```

* Snapshots commands are
```bash
create
restore
list
```

## Version History

* 0.1
    * Initial Release

## License

THE BEERWARE LICENSE" (Revision 42):
<dschaerli@gmail.com> wrote this code. As long as you retain this
notice, you can do whatever you want with this stuff. If we
meet someday, and you think this stuff is worth it, you can
buy me a beer in return.

