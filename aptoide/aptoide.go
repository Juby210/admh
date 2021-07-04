/*
 * Copyright (c) 2021 Juby210
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 * http://www.apache.org/licenses/LICENSE-2.0
 */

package aptoide

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type (
	AppMeta struct {
		Data   *Data
		Info   *Info
		Errors []*Error
	}

	Data struct {
		File *File
		Icon string
		Name string
	}

	Error struct {
		Description string
	}

	File struct {
		Path    string
		VerCode int
		VerName string
	}

	Info struct {
		Status string
	}
)

func DownloadUrlResolver(packageName string, version int) (string, error) {
	res, err := GetAppMeta(packageName, strconv.Itoa(version))
	if err != nil {
		return "", err
	}
	return res.Data.File.Path, nil
}

func GetAppMeta(packageName, version string) (res *AppMeta, err error) {
	if version != "" {
		version = "&vercode=" + version
	}
	resp, err := http.Get(fmt.Sprintf("https://ws75.aptoide.com/api/7/app/getMeta?package_name=%s%s", packageName, version))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return
	}
	if resp.StatusCode > 202 {
		errMsg := ""
		if len(res.Errors) > 0 {
			errMsg = res.Errors[0].Description
		}
		err = fmt.Errorf("aptoide api returned %d %s", resp.StatusCode, errMsg)
	}
	return
}
