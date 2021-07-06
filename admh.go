/*
 * Copyright (c) 2021 Juby210
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 * http://www.apache.org/licenses/LICENSE-2.0
 */

package admh

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type DownloadUrlResolver func(packageName string, version int) (string, error)

var (
	DebugLogs = true

	logger = log.New(os.Stdout, "[admh] ", log.LstdFlags)
)

func DownloadAndExtractAPK(packageName string, version int, workDir, output string, resolveDownloadUrl DownloadUrlResolver) error {
	versionStr := strconv.Itoa(version)
	apkDir := filepath.Join(workDir, "apk")
	apkFile := filepath.Join(apkDir, packageName+"-"+versionStr+".apk")

	if _, err := os.Stat(apkFile); err != nil {
		//goland:noinspection GoBoolExpressions
		if DebugLogs {
			logger.Printf("Downloading apk (%s %s)\n", packageName, versionStr)
		}
		url, err := resolveDownloadUrl(packageName, version)
		if err != nil {
			return err
		}
		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		os.MkdirAll(apkDir, 0777)
		file, err := os.Create(apkFile)
		if err != nil {
			return err
		}
		io.Copy(file, resp.Body)
		file.Close()
	}

	return ExtractAPK(apkFile, workDir, output)
}

//goland:noinspection GoBoolExpressions
func ExtractAPK(apkFile, workDir, output string) error {
	jadxPath := filepath.Join(workDir, "jadx", "bin", "jadx")
	DownloadJADX(workDir, jadxPath)

	if DebugLogs {
		logger.Println("Deleting old app dir")
	}
	os.RemoveAll(filepath.Join(output, "app"))
	if DebugLogs {
		logger.Println("Running jadx")
	}
	cmd := exec.Command(
		jadxPath,
		"-e",
		"--show-bad-code",
		"--no-debug-info",
		"--no-inline-anonymous",
		"--no-inline-methods",
		"--no-generate-kotlin-metadata",
		"--no-replace-consts",
		"--respect-bytecode-access-modifiers",
		"--fs-case-sensitive",
		"-d", output,
		apkFile,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func Push(dir, versionName, versionCode string) (err error) {
	if err = execCmd(dir, "git", "add", "."); err != nil {
		return
	}
	if err = execCmd(dir, "git", "commit", "-m", fmt.Sprintf("%s (%s)", versionName, versionCode)); err != nil {
		return
	}
	return execCmd(dir, "git", "push")
}

func DownloadJADX(workDir, jadxPath string) (err error) {
	if _, err = os.Stat(jadxPath); err == nil {
		return
	}

	//goland:noinspection GoBoolExpressions
	if DebugLogs {
		logger.Println("Downloading jadx")
	}
	resp, err := http.Get("https://github.com/Juby210/jadx/releases/download/v1.2.0.82-fork1/jadx-1.2.0.82-fork1.zip")
	if err != nil {
		return
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	reader := bytes.NewReader(b)
	zipReader, err := zip.NewReader(reader, reader.Size())
	if err != nil {
		return
	}

	jadxDir := filepath.Join(workDir, "jadx")
	os.MkdirAll(jadxDir, 0777)
	for _, f := range zipReader.File {
		fPath := filepath.Join(jadxDir, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fPath, f.Mode())
		} else {
			fDir := ""
			if lastIndex := strings.LastIndex(fPath, string(os.PathSeparator)); lastIndex > -1 {
				fDir = fPath[:lastIndex]
			}

			os.MkdirAll(fDir, f.Mode())
			file, err := os.OpenFile(fPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
			if err != nil {
				return err
			}

			fReader, err := f.Open()
			if err != nil {
				return err
			}
			io.Copy(file, fReader)
			file.Close()
			fReader.Close()
		}
	}
	return
}

func execCmd(dir, command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
