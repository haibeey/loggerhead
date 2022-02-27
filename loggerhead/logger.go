package loggerhead

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"time"
)

var (
	offset = int64(0) // TODO store on redis
	bytes  = make([]byte, 1000)
)

func GetLogs(program, workDir, programArgs, stdout string) (int, *os.File) {
	if err := os.Chdir(workDir); err != nil {
		log.Fatal(err)
	}

	if program == "" {
		log.Fatal("You must provide a program to run in this mode")
	}
	if programArgs == "" {
		fmt.Fprintln(os.Stderr, "Runing program without arguments")
	}

	var (
		std *os.File
		err error
	)

	if stdout == "" {
		stdout = "loggerhead"
	}

	std, err = os.OpenFile(getHomeDir()+"/"+stdout, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Not able to create output file ", err.Error())
	}

	programArgs = strings.Map(reWithSpace, programArgs)
	argsList := strings.Split(programArgs, " ")

	//strip out space
	for i := 0; i < len(argsList); i++ {
		argsList[i] = strings.TrimSpace(argsList[i])
	}

	cmd := exec.Command(program, argsList...)

	if std != nil {
		cmd.Stderr = std
		cmd.Stdout = std
	}

	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	return cmd.Process.Pid, std
}

//reWithSpace subtitutes all white space characters to space
func reWithSpace(r rune) rune {
	switch {
	case r == '\t':
		return ' '
	case r == '\n':
		return ' '
	case r == '\r':
		return ' '
	}
	return r
}

func getHomeDir() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	if len(usr.HomeDir) <= 0 {
		log.Fatal(errors.New("current user have no home directory"))
	}

	return usr.HomeDir
}

func Watch(program, workDir, programArgs, stdout string) {
	_, f := GetLogs(program, workDir, programArgs, stdout)
	for {
		lastLog := getLastLog(f, offset)
		time.Sleep(1 * time.Second)
		socketServer.BroadcastToNamespace("/", "message", lastLog)
	}
}

func getLastLog(f *os.File, off int64) string {

	n, _ := f.ReadAt(bytes, offset)
	offset += int64(n)
	if n == 0 {
		return ""
	}
	return string(bytes)
}
