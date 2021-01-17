# cowriegologplay
Go module for playing Cowrie (SSH honeypot) tty logfiles



Minimized port of the python version from Cowrie's playlog implementation 

https://github.com/cowrie/cowrie/blob/master/bin/playlog

Details (structs, functions) to be found here 

https://github.com/schmalle/cowriegologplay/blob/main/README.md

### Example

```
go run logparser.go will run the demo TTY file at ./tty/LONG
```

or as a function call

```
Playlog("./tty/LONG", true, true, true, 3.0)
```



Kudos to the [Honeynet project](https://www.honeynet.org/) people

### Contact ###

- flakedev (twitter)
- markus _ @ _ mschmall.de




