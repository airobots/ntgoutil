// Shell
// -----

// Author:

// * Rony Novianto (rony@novianto.tech)

// Copyright © Novianto.tech

package shell

import {
	"io"
	"io/ioutil"
	"os"
	"path"
}

func CopyFile(src, dest string, specifiedMode ...os.FileMode) error {
	var err error
	var copyMode os.FileMode

	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()
	
	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}
	
	if len(specifiedMode) > 0 {
		copyMode = specifiedMode[0]
	} else {
		info, err := os.Stat(src)
		if err != nil {
			return err
		}
		copyMode = info.Mode()
	}
	
	os.Chmod(dest, copyMode)
	
	return destFile.Sync()
}

func CopyDirectory(sourceDirectory, destinationDirectory string, specifiedMode ...os.FileMode) error {
	var copyMode os.FileMode
	
	infos, err := ioutil.ReadDir(sourceDirectory)
	if err != nil {
		return err
	}
	
	if len(specifiedMode) > 0 {
		copyMode = specifiedMode[0]
	} else {
		info, err := os.Stat(sourceDirectory)
		if err != nil {
			return err
		}
		copyMode = info.Mode()
	}
	
	// Create any missing directory
	err = os.MkdirAll(destinationDirectory, copyMode)
	if err != nil {
		return err
	}
	
	for _, info := range infos {
		name := info.Name()
		sourceFile := path.Join(sourceDirectory, name)
		destinationFile := path.Join(destinationDirectory, name)
		
		if len(specifiedMode) > 0 {
			copyMode = specifiedMode[0]
		} else {
			copyMode = info.Mode()
		}
		
		if info.IsDir() {
			return CopyDirectory(sourceFile, destinationFile, copyMode)
		} else {
			err = CopyFile(sourceFile, destinationFile, copyMode)
			if err != nil {
				return err
			}
		}
	}
	
	return nil
}