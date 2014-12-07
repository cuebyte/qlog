qlog
====

Usage:
```Go
log, err := InitLogger(TRAC)
if err != nil {
    log.Fatal(err)
}
log.Trace("I'd logged!")
log.Close()
```
