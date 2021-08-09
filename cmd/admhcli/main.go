/*
 * Copyright (c) 2021 Juby210
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 * http://www.apache.org/licenses/LICENSE-2.0
 */

package main

import (
	"flag"
	"log"
	"path/filepath"
	"strconv"

	"github.com/Juby210/admh"
	"github.com/Juby210/admh/aptoide"
)

var (
	apkFile     = flag.String("apk", "", "Apk file")
	packageName = flag.String("packageName", "", "Package name")
	version     = flag.Int("version", 0, "App version")
	output      = flag.String("output", "", "Output dir")
	push        = flag.Bool("push", false, "(git) Commit and push")
	workDir     = flag.String("workDir", "work", "Work dir")
)

func main() {
	flag.StringVar(packageName, "pn", *packageName, "Alias for packageName")
	flag.IntVar(version, "v", *version, "Alias for version")
	flag.StringVar(output, "o", *output, "Alias for output")
	flag.Parse()

	apkFileMode := *apkFile != ""
	if (*packageName == "" || *version == 0) && !apkFileMode {
		log.Fatal("Missing required flags [packageName, version] or [apk]")
	}

	out := *output
	if out == "" {
		out = filepath.Join(*workDir, "out")
	}
	jadxFlags := admh.GetDefaultJadxFlags()
	var err error
	if apkFileMode {
		err = admh.ExtractAPK(*apkFile, *workDir, out, jadxFlags)
	} else {
		err = admh.DownloadAndExtractAPK(*packageName, *version, *workDir, out, jadxFlags, aptoide.DownloadUrlResolver)
	}
	if err != nil {
		log.Fatal(err)
	}
	if *push {
		versionStr := strconv.Itoa(*version)
		res, err := aptoide.GetAppMeta(*packageName, versionStr)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("[admh] Pushing")
		err = admh.Push(out, res.Data.File.VerName, versionStr)
		if err != nil {
			log.Fatal(err)
		}
	}
}
