## LOG-RIPPER

Log-Ripper is a command line tool for stripping out large log files in a directory or a standalone file.

While it might be really easy to just delete the file and make a new file. The file descriptor is lost with that method and 
other processes communicating through the logs files redirects stdin && stderr elsewhere. This approach makes sure the the FD is presevered while trimming the file.


##### sample usage
```
run 'go build -o logripper'

logripper -n 100 -path /path/  --> trim all files in the folder to the last 100 lines

logripper -n 100 -ext .log -path /path/
```


#TODO
Write TEST
