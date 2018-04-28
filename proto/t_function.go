package proto

import "fmt"

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
 * Representation of a function. Key parts are return type, function name,
 * optional modifiers, and an argument list, which is implemented as a thrift
 * struct.
 *
 */

type TFunction struct{
	TDoc
	ReturnType_ ICafType
	Name_ string
	Arglist_ *TStruct
	Xceptions_ *TStruct
	Oneway_ bool
	Annotations_ map[string]string
	OwnXceptions_ bool
}

func NewFunction(returntype ICafType, name string, arglist *TStruct, oneway bool) *TFunction {
	fmt.Printf("[NewFunction]%s\n", name)
	var ret= new(TFunction)
	ret.ReturnType_ = returntype
	ret.Name_ = name
	ret.Arglist_ = arglist
	ret.Oneway_ = oneway
	ret.Xceptions_ = new(TStruct)
	if ret.Oneway_ && (!ret.ReturnType_.IsVoid()) {
		panic("Oneway methods should return void.\n")
	}

	return ret
}

func NewFunctionWithXception(returntype ICafType, name string, xceptions, arglist *TStruct, oneway bool) *TFunction{
		var ret= new(TFunction)
		ret.ReturnType_ = returntype
		ret.Name_ = name
		ret.Arglist_ = arglist
		ret.Oneway_ = oneway
		ret.Xceptions_ = xceptions
		ret.OwnXceptions_ = false

		if ret.Oneway_ && len(ret.Xceptions_.GetMembers()) > 0 {
			panic("Oneway methods can't throw exceptions.")
		}
		if ret.Oneway_ && (!ret.ReturnType_.IsVoid()) {
			panic("Oneway methods should return void.\n")
		}
		return ret
}

func (f *TFunction)GetReturnType() ICafType{
	return f.ReturnType_
}

func (f *TFunction) GetName() string{
	return f.Name_
}

func (f *TFunction)GetArglist() *TStruct{
	return f.Arglist_;
}

func (f *TFunction)GetXceptions() *TStruct{
	return f.Xceptions_;
}

func (f *TFunction)IsOneway() bool{
	return f.Oneway_;
}


func (f *TFunction)GetAnnotations() map[string]string{
	return f.GetAnnotations()
}

func (f *TFunction)SetAnnotations(annotation map[string]string) {
	f.Annotations_ = annotation
}
