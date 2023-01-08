package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/sys/unix"
	"io/fs"
	"os"
	"os/exec"
)

type Client interface {
	Name() string
	Init() (Path, error)
	List(Path) ([]Path, error)
	Get(Path) (Content, error)
}

type CmdClient struct {
	name string
	init string
	list string
	get  string
}

func NewCmdClient(name, rootPath string) (*CmdClient, error) {
	return &CmdClient{name, rootPath + "init", rootPath + "list", rootPath + "get"}, nil
}

func (c *CmdClient) Name() string {
	return c.name
}

func (c *CmdClient) Init() (Path, error) {
	if err := validateExecutableFile(c.init); err != nil {
		return nil, err
	}
	if err := validateExecutableFile(c.list); err != nil {
		return nil, err
	}
	if err := validateExecutableFile(c.get); err != nil {
		return nil, err
	}
	cmd := exec.Command(c.init)
	stdout, err := cmd.Output()
	exitErr, isExitError := err.(*exec.ExitError)
	if isExitError {
		return nil, errors.New(string(exitErr.Stderr))
	}
	if len(stdout) == 0 {
		return nil, errors.New("empty init output")
	}
	rootPath := string(stdout)
	return NewPath(nil, rootPath, rootPath, false), nil
}

func (c *CmdClient) List(path Path) ([]Path, error) {
	cmd := exec.Command(c.list, path.Path())
	stdout, err := cmd.Output()
	exitErr, isExitError := err.(*exec.ExitError)
	if isExitError {
		return nil, fmt.Errorf("path %s cannot be listed: %w", path.Path(), exitErr)
	}

	list := make([]Path, 0)
	err = c.parseList(stdout, func(localPath string) {
		list = append(list, NewPath(path, path.Path()+localPath, localPath, isFinal(localPath)))
	})
	if err != nil {
		return nil, fmt.Errorf("path %s cannot be listed: %w", path.Path(), err)
	}
	return list, nil
}

func (c *CmdClient) parseList(stdout []byte, append func(localPath string)) error {
	scanner := bufio.NewScanner(bytes.NewReader(stdout))
	for scanner.Scan() {
		localPath := scanner.Text()
		append(localPath)
	}
	return nil
}

func (c *CmdClient) Get(path Path) (Content, error) {
	cmd := exec.Command(c.get, path.Path())
	stdout, err := cmd.Output()
	exitErr, isExitError := err.(*exec.ExitError)
	if isExitError {
		return nil, errors.New(string(exitErr.Stderr))
	}
	return NewContent(path, string(stdout)), nil
}

func validateExecutableFile(path string) error {
	file, err := os.Stat(path)
	if err != nil {
		return err
	}
	mode := file.Mode()
	if !((mode.IsRegular()) || (mode&fs.ModeSymlink) == 0) {
		return errors.New("file " + path + " should be a normal file or symlink")
	}
	if mode&0111 == 0 {
		return errors.New("file " + path + " should be executable")
	}
	if unix.Access(path, unix.X_OK) != nil {
		return errors.New("file " + path + " should be executed by current user")
	}
	return nil
}
