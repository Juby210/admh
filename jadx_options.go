/*
 * Copyright (c) 2021 Juby210
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 * http://www.apache.org/licenses/LICENSE-2.0
 */

package admh

type JadxOptions struct {
	GradleProject,
	ShowBadCode,
	DebugInfo,
	InlineAnonymous,
	InlineMethods,
	GenerateKotlinMetadata,
	ReplaceConsts,
	RespectBytecodeAccessModifiers bool

	JadxRelease string
}

func GetDefaultJadxOptions() *JadxOptions {
	return &JadxOptions{
		GradleProject:                  true,
		ShowBadCode:                    true,
		DebugInfo:                      false,
		InlineAnonymous:                false,
		InlineMethods:                  false,
		GenerateKotlinMetadata:         false,
		ReplaceConsts:                  true,
		RespectBytecodeAccessModifiers: true,

		JadxRelease: "https://github.com/Aliucord/jadx/releases/download/v1.3.2.311-fork1/jadx-1.3.2.311-fork1.zip",
	}
}

func (options *JadxOptions) GetRawFlags() (raw []string) {
	if options.GradleProject {
		raw = append(raw, "-e")
	}
	if options.ShowBadCode {
		raw = append(raw, "--show-bad-code")
	}
	if !options.DebugInfo {
		raw = append(raw, "--no-debug-info")
	}
	if !options.InlineAnonymous {
		raw = append(raw, "--no-inline-anonymous")
	}
	if !options.InlineMethods {
		raw = append(raw, "--no-inline-methods")
	}
	if !options.GenerateKotlinMetadata {
		raw = append(raw, "--no-generate-kotlin-metadata")
	}
	if !options.ReplaceConsts {
		raw = append(raw, "--no-replace-consts")
	}
	if options.RespectBytecodeAccessModifiers {
		raw = append(raw, "--respect-bytecode-access-modifiers")
	}
	return
}
