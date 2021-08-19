# thro

`thro` is a cli app that throttles passed command when the same command with the same args is already running.

Main purpose is to use it for throttling apps with long duration that can be called frequently, example:

```shell
#!/bin/sh

while true
do
	inotifywait -e modify -e create -e delete -e attrib -r /myfolder/
	./thro rclone -v sync /myfolder/ myprovider:/myfolder/ &
done
```

In this example several parallel calls will be throttled to one subcommand call and only one more deferred subcommand
call when the first one will be finished