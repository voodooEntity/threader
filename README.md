## Threader 

### Info: 
Threader executes a given command "-run" in x parallel threads. It can be used to just execute the command a defined number of times "-runs" or to pass input given by STDIN split, by a delimiter and provide each result part as \\$INPUTSTR param to
your "-run" command. 

For more information of the current version and changes please check the CHANGELOG.md

### Install: 
##### Use the release shipped binaries 
You can simply download the latest release binaries fitting your os/arch and copy the executable into your PATH to make it executeble. 

##### Build it yourself 
The current mikefile supports systems that include /usr/bin in their PATH. If your system doesn't include that you either add it or compile it manually with "go build -o threader" and copy it into your supported PATH. 
```sh
$ git clone https://github.com/voodooEntity/threader
$ cd threader 
$ make && make build && make install
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

### FAQ 

##### Why don't you  support Windows in your shipped binaries?
The nature of the tool is to build dynamic bash execution strings to execute the command you are threading. Windows CMD works different than the POSIX systems of debian etc. To enable threader for windows i would have to add a lot of special cases for windows. I might do this in future if the feature gets requested, but right now im sticking to linux/mac.

##### Why did you build threader?
I just thought there could be a tool with just simplistic input doing the job, and since golang is very good at simple and stable multithreading i created it. 

##### Why is my os/arch not included in the binaries
Due to the amount of possibilities of executions with this tool its hard to create tests that would cover everything possible or even a big part of it. So for the moment i build binaries im able to test myself. If anyone would suggest a specific os/arch set he/she would like to see in the binaries im absolutly willed to add it, and i would love to get response about this specific binary working .)

##### What are the future plans of threader?
I have no speficic future plans for the tool. The base funcionality is given and i gonne focus on making sure this tool always stays stable and working for a majority of os/arch's. Im still open to suggestions on how to improve the tool, but i can't promise anything .) just feel free to comment.



#### Contributors
Gonne list some people here that helped out with some rubberducking/ideas etc. Later this will include PR creators too .)
* f0o
* luxer