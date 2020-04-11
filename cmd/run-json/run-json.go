package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

type JsonCmd struct {
	Title    string     `json:"title"`
	Command  string     `json:"command"`
	PreArgs  []string   `json:"pre_args"`
	Envs     []string   `json:"envs"`
	WorkDir  string     `json:"work_dir"`
	OptFiles [][]string `json:"opt_files"`
	Files    []string   `json:"files"`
	OptDirs  [][]string `json:"opt_dirs"`
	Dirs     []string   `json:"dirs"`
}

func runJsonCmd(jcmd *JsonCmd) error {
	args := []string{}
	args = append(args, jcmd.PreArgs...)

	for _, v := range jcmd.OptFiles {
		args = append(args, v...)
	}
	for _, v := range jcmd.Files {
		args = append(args, v)
	}
	for _, v := range jcmd.OptDirs {
		args = append(args, v...)
	}
	for _, v := range jcmd.Dirs {
		args = append(args, v)
	}
	os.Chdir(jcmd.WorkDir)
	fmt.Printf("%v\n", args)
	cmd := exec.Command(jcmd.Command, args...)
	pIn, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	defer pIn.Close()
	pErr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	defer pErr.Close()

	jcmd.Envs = append(jcmd.Envs, os.Environ()...)
	for _, v := range jcmd.Envs {
		cmd.Env = append(cmd.Env, v)
	}

	err = cmd.Start()
	if err != nil {
		return err
	}

	mrd := io.MultiReader(pIn, pErr)
	brd := bufio.NewReader(mrd)
	for {
		line, _, err := brd.ReadLine()
		if err != nil {
			return err
		}
		fmt.Println(string(line))
	}
}

func readJsonCmd(fn string) (*JsonCmd, error) {
	res := new(JsonCmd)
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func main() {
	if len(os.Args) < 2 {
		return
	}
	jc, err := readJsonCmd(os.Args[1])
	if err != nil {
		panic(err)
	}
	err = runJsonCmd(jc)
	fmt.Println("Exit:", err.Error())
}
