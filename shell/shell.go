// Shell
// -----

// Author:

// * Rony Novianto (rony@novianto.tech)

// Copyright © Novianto.tech

package shell

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
)

func CopyFile(sourcePath, destinationPath string, uid, gid int, fileMode, directoryMode *os.FileMode) error {
	var err error
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()
	directoryPath := path.Dir(destinationPath)
	if directoryMode == nil {
		err = os.MkdirAll(directoryPath, os.ModePerm)
	} else {
		err = os.MkdirAll(directoryPath, *directoryMode)
	}
	if err != nil {
		return err
	}
	destinationFile, err := os.Create(destinationPath)
	if err != nil {
		return err
	}
	defer destinationFile.Close()
	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}
	if uid >= 0 || gid >= 0 {
		err = destinationFile.Chown(uid, gid)
		if err != nil {
			return err
		}
	}
	if fileMode != nil {
		err = destinationFile.Chmod(*fileMode)
		if err != nil {
			return err
		}
	}
	return destinationFile.Sync()
}

func Copy(source, destination string, uid, gid int, fileMode, directoryMode *os.FileMode) error {
	info, err := os.Stat(source)
	originalMode := info.Mode()
	if err != nil {
		return err
	}
	if !info.IsDir() {
		if fileMode == nil {
			err = CopyFile(source, destination, uid, gid, &originalMode, directoryMode)
		} else {
			err = CopyFile(source, destination, uid, gid, fileMode, directoryMode)
		}
		if err != nil {
			return err
		}
		return nil
	}
	if directoryMode == nil {
		err = os.MkdirAll(destination, originalMode)
	} else {
		err = os.MkdirAll(destination, *directoryMode)
	}
	if err != nil {
		return err
	}
	infos, err := ioutil.ReadDir(source)
	if err != nil {
		return err
	}
	for _, info := range infos {
		name := info.Name()
		path1 := path.Join(source, name)
		path2 := path.Join(destination, name)
		err = Copy(path1, path2, uid, gid, fileMode, directoryMode)
		if err != nil {
			return err
		}
	}
	return nil
}

func RunScript(content string, environment []string) error {
	command := exec.Command("sh")
	command.Env = environment
	command.Stdin = bytes.NewBufferString(content)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	return command.Run()
}