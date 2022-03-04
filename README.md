# gstellar

Easy PG snapshot tool for development environment. Inspired by stellar in python.

## Description

Creates and restores snapshots on a postgresql db without pg_restore and pg_dump

### Installing

```bash
go install github.com/schaerli/gstellar@v0.1.1
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
web
```

* Snapshots commands are
```bash
create
restore
list
drop
```

## Version History

* 0.0.1
  * Initial Release
* 0.0.2
  * Web gui
* 0.1.0
  * drop database
* 0.1.1
  * set default values in init for host and port
* 0.1.2
  * fix restore function

## License

THE BEERWARE LICENSE" (Revision 42):
schaerli wrote this code. As long as you retain this
notice, you can do whatever you want with this stuff. If we
meet someday, and you think this stuff is worth it, you can
buy me a beer in return.

