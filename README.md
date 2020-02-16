

## Threader

### Info:
Threader executes a given command (-run) in x parallel threads. It can be used to
just execute the Command a defined number of times (-runs) or to pass input given
by STDIN split by a delimiter and provide each result part as \\$INPUTSTR param to
your -run command. 

Threader is shipped with /bin/bash as only dependency. 

### Install.
#### All
```sh
$ git clone https://github.com/voodooEntity/threader
$ cd threader 
```
#### If you want to build it yourself with go installed
```sh
$ make && make build && make install 
```
#### If you want to use the shipped binary simply download it and execute (currently only debian64bit build)
```sh
$ sudo cp bin/deb64/threader /usr/bin/threader 
```

### Args: 
* -run \"yourcommand\" 
  * Can include \\$INPUTSTR \\$INPUTID \\$THREADID
* -runs INT 
  * Amount of run executions to be done if no input is given
* -delimiter \"delimiterstring\" 
  * String to split stdin given input up to single command inputstr , default delimiter=\"\\n\"
* -verbose on 
  * Sets threaders core output to verbose mode for debugging purposes
* -threads INT 
  * Define a number of threads to be used for parallel execution, default threads=amount of cpus
### Vars: 
The following vars can be used in your execution command.
* \\$INPUTSTR    
  * This variable will include a single input part provided by the result of splitting the STDIN input by -delimiter string.
* \\$INPUTID     
  * This variable will include a the id of the given \\$INPUTSTR. This variable is only unique for each thread, not in total.
* \\$THREADID    
  * This variable will include a the id of the thread executing the current command. It can be used to create unique identifiers combined with \\$INPUTID



### Usage:

##### Running a command without input 100 times in 4 threads
 ```sh
$ threader -run "curl http://some.domain.com > /dev/null" -runs 100 -threads 4
```
##### Running a command without input 100 times auto thread amount detection (by cpu count)
 ```sh
$ threader -run "curl http://some.domain.com > /dev/null" -runs 100
```
##### Running a command with input from a file (f.e. 1 url per line)
 ```sh
$ cat urllist.txt | threader -run "curl \$INPUTSTR > /dev/null"
```
##### Running a command with input from cli and custom delimiter
 ```sh
$ echo "/etc,/home,/srv" | threader -run "stat \$INPUTSTR" -delimiter "," 
```
##### Running a command with input from a cli command
 ```sh
$ ls -1 / | threader -run "stat /\$INPUTSTR"
```
##### Get the total size of each directory in / (using function output, dynamic amount of threads, default delimiter \n )
 ```sh
$ ls -D1 / | threader -run "du -h /\$INPUTSTR | tail -1"
```




