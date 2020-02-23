package main

import (
	"bufio"
	"errors"
	"flag"
	"io"
	"log"
	"math"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

type Args struct {
	Delimiter    string
	Input        string
	Verbose      string
	ThreadAmount int
	Run          string
	Runs         int
}

var shell string
var args Args
var loggerOut = log.New(os.Stdout, "", 0)

func init() {
	// presets
}

func main() {
	// detect shell
	executer, exists := os.LookupEnv("SHELL")
	if !exists {
		// Print the value of the environment variable
		loggerOut.Println("No SHELL env given, exiting")
		os.Exit(1)
	}
	shell = executer

	// first we gonne parse the args
	parseArgs()

	// one simple check
	cpuAmount := runtime.NumCPU()
	if args.ThreadAmount > cpuAmount {
		verboseOut("You set '" + strconv.Itoa(args.ThreadAmount) + "' threads but only have '" + strconv.Itoa(cpuAmount) + "' cores. It`s your choice...")
	}

	// prepare the input map
	var input map[int]string
	// is there a command given that provides input?
	if args.Runs != -1 {
		input = make(map[int]string)
	} else {
		args.Input = getPipedInput()
		input = splitInput()
	}

	// decide if its smart minions or stupid minions
	if 0 == len(input) {
		runHeadlessMinions()
	} else {
		runSmartMinions(input)
	}

}

func getPipedInput() string {
	reader := bufio.NewReader(os.Stdin)
	var output []rune

	for {
		input, _, err := reader.ReadRune()
		if err != nil && err == io.EOF {
			break
		}
		output = append(output, input)
	}

	return string(output)
}

func runSmartMinions(input map[int]string) {
	// create an intercom map
	intercom := make(map[int]chan string)

	// some number juggling , probably smarter ways
	// out there but for now its fine
	specialCase := -1
	base := float64(len(input)) / float64(args.ThreadAmount)
	checkBase := math.Floor(float64(base))
	if base != checkBase {
		checkAmount := int(checkBase) * args.ThreadAmount
		specialCase = len(input) - checkAmount
	}

	// start/stop id
	start := 0
	stop := 0

	// run the minions
	for i := 0; i < args.ThreadAmount; i++ {
		// create the channel
		intercom[i] = make(chan string)

		// check how many entries we got to push
		// and calc the stop
		subAmount := int(checkBase)
		if specialCase != -1 && i == 0 {
			subAmount = int(checkBase) + specialCase
		}
		stop = stop + subAmount

		// dont judge me i want to finish this
		subInput := make([]string, subAmount)

		for x := start; x < stop; x++ {
			if val, ok := input[x]; ok {
				subInput[x-start] = val
			}
		}

		// set the next start to the current stop
		start = stop
		// start the minion
		go createSmartMinion(i, subInput, intercom[i])
	}

	// check for the minions to be done
	done := 0
	for {

		// check all thread channels for return
		for i := 0; i < args.ThreadAmount; i++ {
			ret := <-intercom[i]
			if "done" == ret {
				done++
			}
		}

		// all threads are done
		if args.ThreadAmount == done {
			break
		}
	}
	verboseOut("All work is done. Exiting")

}

func createSmartMinion(id int, input []string, intercom chan string) {
	verboseOut("Minion nr " + strconv.Itoa(id) + " spawned. ")
	for iid, singleIn := range input {
		// provide some vars for the command
		// run the command
		execString := prepareExecString(singleIn, id, iid)
		ret, _ := custExec(execString)
		loggerOut.Println(ret)
	}
	verboseOut("Minion nr " + strconv.Itoa(id) + " finished its job.")
	intercom <- "done"
}

func prepareExecString(input string, threadID int, inputID int) string {
	// set the vars as env vars appended with the threadid
	threadIDstr := strconv.Itoa(threadID)
	os.Setenv("TIS"+threadIDstr, input)
	os.Setenv("TID"+threadIDstr, strconv.Itoa(inputID))
	os.Setenv("TTID"+threadIDstr, threadIDstr)

	// build the execution string buy normalising the varnames to standard varnames for erasier command handling
	commandString := "INPUTSTR=$TIS" + threadIDstr + ";THREADID=$TTID" + threadIDstr + ";INPUTID=$TID" + threadIDstr + ";" + args.Run + ";"
	verboseOut(commandString)

	return commandString
}

func runHeadlessMinions() {
	// create an intercom map
	intercom := make(map[int]chan string)

	// check some things
	if args.ThreadAmount > args.Runs {
		verboseOut("You did set '" + strconv.Itoa(args.ThreadAmount) + "' but set only '" + strconv.Itoa(args.Runs) + "' runs. Resetting the Threadamount")
		args.ThreadAmount = args.Runs
	}

	// some number juggling , probably smarter ways
	// out there but for now its fine
	specialCase := -1
	base := float64(args.Runs) / float64(args.ThreadAmount)
	checkBase := math.Floor(float64(base))
	if base != checkBase {
		checkAmount := int(checkBase) * args.ThreadAmount
		specialCase = args.Runs - checkAmount
	}

	// run the minions
	for i := 0; i < args.ThreadAmount; i++ {
		intercom[i] = make(chan string)
		runs := int(checkBase)
		if specialCase != -1 && i == 0 {
			runs = int(checkBase) + specialCase
		}
		go createStupidMinion(i, runs, intercom[i])
	}

	// check for the minions to be done
	done := 0
	for {
		// check all thread channels for return
		for i := 0; i < args.ThreadAmount; i++ {
			ret := <-intercom[i]
			if "done" == ret {
				done++
			}
		}

		// all threads are done
		if args.ThreadAmount == done {
			break
		}
	}
	verboseOut("All work is done. Exiting")
}

func createStupidMinion(id int, runs int, intercom chan string) {
	verboseOut("Minion nr " + strconv.Itoa(id) + " spawned. ")
	for i := 0; i < runs; i++ {
		prepared := prepareExecString("", id, i)
		ret, _ := custExec(prepared)
		if "" != ret {
			loggerOut.Println(ret)
		}
	}
	verboseOut("Minion nr " + strconv.Itoa(id) + " finished its job.")
	intercom <- "done"
}

func splitInput() map[int]string {
	// predefine return map
	data := make(map[int]string)

	// split the input string
	inputSlice := strings.Split(args.Input, args.Delimiter)

	// check if the input amount is less than the given thread amount
	if args.ThreadAmount > len(inputSlice) {
		verboseOut("You set to execute more threads than inputs provided. Resetting threads amount to input amount '" + strconv.Itoa(len(inputSlice)) + "'.")
		args.ThreadAmount = len(inputSlice)
	}

	// transform the input to map
	i := 0
	for _, val := range inputSlice {
		data[i] = val
		i++
	}

	return data
}

func custExec(cmd string) (string, error) {
	output, err := exec.Command(shell, "-c", cmd).Output()
	if nil != err {
		return "", errors.New("Error executin command")
	}
	return strings.TrimRight(string(output), "\n"), nil
}

func verboseOut(out string) {
	if "on" == args.Verbose {
		loggerOut.Println(out)
	}
}

func parseArgs() {
	// first we check for the help flag
	if 1 < len(os.Args) {
		if ok := os.Args[1]; ok == "help" {
			printHelpText()
			os.Exit(1)
		}
	}

	// delimiter to be used for custom generator input
	var delimiter string
	flag.StringVar(&delimiter, "delimiter", "\n", "Cutoff color value")

	// thread amount to be used
	var threadAmount int
	flag.IntVar(&threadAmount, "threads", -1, "Amount of threads to be used")

	// amount of runs in case you dont provide input
	var runs int
	flag.IntVar(&runs, "runs", -1, "Amount of threads to be used")

	// run command by string
	var run string
	flag.StringVar(&run, "run", "", "Job to be executed")

	// input by string
	var verbose string
	flag.StringVar(&verbose, "verbose", "", "Job to be executed")

	// parse the flags
	flag.Parse()

	args = Args{
		Delimiter:    delimiter,
		ThreadAmount: threadAmount,
		Verbose:      verbose,
		Run:          run,
		Input:        "",
		Runs:         runs,
	}

	// check for dynamic thread amount
	if -1 == args.ThreadAmount {
		args.ThreadAmount = runtime.NumCPU()
		verboseOut("Setting output dynamic to " + string(args.ThreadAmount))
	}

}

func printHelpText() {
	helpText := "Threader Help:\n" +
		"> Threader executes a given command (-run) in x parallel threads. It can be used to\n" +
		"  just execute the Command a defined number of times (-runs) or to pass input given\n" +
		"  by stdIn split by a delimiter and provide each result part as \\$INPUTSTR param to\n" +
		"  your command. For examples check https://github.com/voodooEntity/threader readme.\n\n" +
		"  Args: \n" +
		"    -run \"yourcommand\"            | Can include \\$INPUTSTR \\$INPUTID \\$THREADID\n" +
		"    -runs INT                     | Amount of run executions to be done if no input is given\n" +
		"    -delimiter \"delimiterstring\"  | String to split stdin given input up to single command inputstr\n" +
		"                                    default delimiter=\"\\n\"\n" +
		"    -verbose on                   | Sets threaders core output to verbose mode for debugging purposes\n" +
		"    -threads INT                  | Define a number of threads to be used for parallel execution\n" +
		"                                  | default threads=amount of cpus\n" +
		"  Vars: \n" +
		"    The following vars can be used in your execution command. \n" +
		"    - \\$INPUTSTR    This variable will include a single input part provided by the result of \n" +
		"                    splitting the stdIn by -delimiter.\n" +
		"    - \\$INPUTID     This variable will include a the id of the given INPUTSTR. This variable \n" +
		"                    is only unique for each thread, not in total.\n" +
		"    - \\$THREADID    This variable will include a the id of the thread executing the current \n" +
		"                    command. It can be used to create unique identifiers combined with \\INPUTID\n"
	loggerOut.Println(helpText)
}
