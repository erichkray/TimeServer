Author: Erich Ray
TimeServer Version 3

Compile src\github.com\TimeServer\TimeServer\TimeServer.go

Yes, there are multiple layers of TimeServer, because I didn't feel like renaming the project once I had the extra packages added to the directory.

usage:
TimeServer -log [path to log file]
	no log file will cause all logging to output to the screen.
Other optional flags:
	-v: print version number and terminate
	-debug: extra spew to the console
	-port: override for web port