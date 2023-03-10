package main

import (
	"bufio"
	"bytes"
	"errors"
	"os/exec"
)

type VaultClient struct {
	rootPath string
}

func NewVaultClient(rootPath string) (*VaultClient, error) {
	return &VaultClient{rootPath}, nil
}

func (c *VaultClient) List(path Path) ([]Path, error) {
	cmd := exec.Command("vault", "kv", "list", path.Path())
	stdout, err := cmd.Output()
	exitErr, isExitError := err.(*exec.ExitError)
	if isExitError {
		return nil, errors.New(string(exitErr.Stderr))
	}

	list := make([]Path, 0)
	err = parseList(stdout, func(localPath string) {
		list = append(list, NewPath(path, path.Path()+localPath, localPath, isFinal(localPath)))
	})
	if err != nil {
		return nil, err
	}
	return list, nil
}

func parseList(stdout []byte, append func(localPath string)) error {
	scanner := bufio.NewScanner(bytes.NewReader(stdout))
	// Check the output header
	if !scanner.Scan() || scanner.Text() != "Keys" || !scanner.Scan() || scanner.Text() != "----" {
		return errors.New("Invalid vault list header:" + string(stdout))
	}
	for scanner.Scan() {
		localPath := scanner.Text()
		append(localPath)
	}
	return nil
}

func (c *VaultClient) Get(path Path) (Content, error) {
	cmd := exec.Command("vault", "kv", "get", path.Path())
	stdout, err := cmd.Output()
	exitErr, isExitError := err.(*exec.ExitError)
	if isExitError {
		return nil, errors.New(string(exitErr.Stderr))
	}
	return NewContent(path, string(stdout)), nil
}

func isFinal(path string) bool {
	return path[len(path)-1] != '/'
}

func (c *VaultClient) RootPath() Path {
	return NewPath(nil, c.rootPath, c.rootPath, false)
}
