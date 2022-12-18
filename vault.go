package main

import (
	"bufio"
	"bytes"
	"errors"
	"os/exec"
	"strings"
)

type VaultClient struct {
}

func NewVaultClient() *VaultClient {
	return &VaultClient{}
}

func (c *VaultClient) List(path string) ([]ListItem, error) {
	cmd := exec.Command("vault", "kv", "list", path)
	stdout, err := cmd.Output()
	exitErr, isExitError := err.(*exec.ExitError)
	if isExitError {
		return nil, errors.New(string(exitErr.Stderr))
	}

	return parseList(stdout, path)
}

func (c *VaultClient) HasParent(path string) bool {
	return strings.Index(path, "/") > 0
}

func parseList(stdout []byte, parentPath string) ([]ListItem, error) {
	scanner := bufio.NewScanner(bytes.NewReader(stdout))
	// Check the output header
	if !scanner.Scan() || scanner.Text() != "Keys" || !scanner.Scan() || scanner.Text() != "----" {
		return nil, errors.New("Invalid vault list header:" + string(stdout))
	}
	list := make([]ListItem, 0)
	for scanner.Scan() {
		localPath := scanner.Text()
		list = append(list, ListItem{
			Label: localPath,
			Path:  parentPath + localPath,
			Final: isFinal(localPath),
		})
	}
	return list, nil
}

func isFinal(localPath string) bool {
	return !strings.HasSuffix(localPath, "/")
}

func (c *VaultClient) Parent(path string) string {
	return ""
}

func (c *VaultClient) Get(path string) (string, error) {
	return "", nil
}
