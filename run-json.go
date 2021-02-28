package main

import (
	"encoding/json"
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
	OptTexts [][]string `json:"opt_texts"`
	Texts    []string   `json:"texts"`
	TailArgs []string   `json:"tail_args"`
	Help     string     `json:"help"`
}

func getArgs(jcmd *JsonCmd) []string {
	args := []string{}
	args = append(args, jcmd.PreArgs...)

	for _, v := range jcmd.OptFiles {
		args = append(args, v...)
	}
	for _, v := range jcmd.OptDirs {
		args = append(args, v...)
	}
	for _, v := range jcmd.Files {
		args = append(args, v)
	}
	for _, v := range jcmd.Dirs {
		args = append(args, v)
	}
	for _, v := range jcmd.TailArgs {
		args = append(args, v)
	}
	return args
}

func runJsonCmd(jcmd *JsonCmd) (pIn io.ReadCloser, pOut io.WriteCloser, pErr io.ReadCloser, err error) {
	args := getArgs(jcmd)
	//os.Chdir(jcmd.WorkDir)
	//fmt.Printf("%v\n", args)
	cmd := exec.Command(jcmd.Command, args...)
	cmd.Dir = jcmd.WorkDir
	pIn, err = cmd.StdoutPipe()
	if err != nil {
		return nil, nil, nil, err
	}
	//defer pIn.Close()
	pOut, err = cmd.StdinPipe()
	if err != nil {
		return nil, nil, nil, err
	}

	pErr, err = cmd.StderrPipe()
	if err != nil {
		return nil, nil, nil, err
	}
	//defer pErr.Close()

	jcmd.Envs = append(jcmd.Envs, os.Environ()...)
	for _, v := range jcmd.Envs {
		cmd.Env = append(cmd.Env, v)
	}

	err = cmd.Start()
	if err != nil {
		return nil, nil, nil, err
	}

	if err != nil {
		return nil, nil, nil, err
	}
	return
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

func runJsonCfg(fn string) error {
	jc, err := readJsonCmd(fn)
	if err != nil {
		return err
	}
	pIn, pOut, pErr, err := runJsonCmd(jc)
	if err != nil {
		return err
	}
	defer pIn.Close()
	defer pOut.Close()
	defer pErr.Close()
	go io.Copy(os.Stdout, pIn)
	_, err = io.Copy(os.Stderr, pErr)
	return err
}
