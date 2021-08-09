/*
 * Copyright (c) 2021 Juby210
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 * http://www.apache.org/licenses/LICENSE-2.0
 */

package admh

type JadxFlags struct {
	GradleProject,
	ShowBadCode,
	DebugInfo,
	InlineAnonymous,
	InlineMethods,
	GenerateKotlinMetadata,
	ReplaceConsts,
	RespectBytecodeAccessModifiers bool
}

func GetDefaultJadxFlags() *JadxFlags {
	return &JadxFlags{
		GradleProject:                  true,
		ShowBadCode:                    true,
		DebugInfo:                      false,
		InlineAnonymous:                false,
		InlineMethods:                  false,
		GenerateKotlinMetadata:         false,
		ReplaceConsts:                  true,
		RespectBytecodeAccessModifiers: true,
	}
}

func (flags *JadxFlags) GetRawFlags() (raw []string) {
	if flags.GradleProject {
		raw = append(raw, "-e")
	}
	if flags.ShowBadCode {
		raw = append(raw, "--show-bad-code")
	}
	if !flags.DebugInfo {
		raw = append(raw, "--no-debug-info")
	}
	if !flags.InlineAnonymous {
		raw = append(raw, "--no-inline-anonymous")
	}
	if !flags.InlineMethods {
		raw = append(raw, "--no-inline-methods")
	}
	if !flags.GenerateKotlinMetadata {
		raw = append(raw, "--no-generate-kotlin-metadata")
	}
	if !flags.ReplaceConsts {
		raw = append(raw, "--no-replace-consts")
	}
	if flags.RespectBytecodeAccessModifiers {
		raw = append(raw, "--respect-bytecode-access-modifiers")
	}
	return
}
