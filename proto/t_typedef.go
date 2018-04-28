package proto

/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements. See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership. The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License. You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

/**
 * A typedef is a mapping from a symbolic name to another type. In dymanically
 * typed languages (i.e. php/python) the code generator can actually usually
 * ignore typedefs and just use the underlying type directly, though in C++
 * the symbolic naming can be quite useful for code clarity.
 *
 */

type TTypedef struct{
	TTypeDesc
	Type_ ICafType
	Symbolic_ string
	Forward_ bool
	Seen_ bool
}

func NewTypedef(program *TProgram, caftype ICafType, symbolic string, forward bool) *TTypedef{
	ret := new(TTypedef)
	ret.Type_ = caftype
	ret.Symbolic_ = symbolic
	ret.Forward_ = false
	ret.Seen_ = false
	ret.Forward_ = forward
	return ret
}

func (td *TTypedef) GetType() ICafType{
	return td.Type_
}

func (td *TTypedef)GetSymbolic() string{
	return td.Symbolic_
}

func (td *TTypedef)IsforwardTypedef() bool{
	return td.Forward_
}

func (td *TTypedef)IsTypedef() bool{
	return true
}

func (td *TTypedef)GetFingerprintMaterial() string{
	if !td.Seen_ {
		td.Seen_ = true
		ret := td.GetType().GetFingerprintMaterial()
		td.Seen_ = false
		return ret;
	}
	return "";
}

func (td *TTypedef)GenerateFingerprint() {
	GenerateBaseFingerprint(ICafType(td))
	if td.GetType().HasFingerprint() {
		td.GetType().GenerateFingerprint()
	}
}

