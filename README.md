## LOG-RIPPER

Log-Ripper is a command line tool for striiping out large log files in a directory or a standalone file


##### sample usage
```
run go build.

logripper -n 100 -path /path/  --> trim all files in the folder to the last 100 lines

logripper -n 100 -ext .log -path /path/
```


#TODO
Write TEST