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
 * A service consists of a set of functions.
 *
 */
type TService struct{
	TTypeDesc
	Functions_ []*TFunction
	Extends_ *TService
}

func NewService(program* TProgram) *TService{
	var ret = new(TService)
	ret.Program_ = program
	return ret
}

func (s *TService)IsService() bool{
	return true
}

func (s *TService)SetExtends(service *TService){
	s.Extends_ = service
}

func (s *TService)AddFunction(f *TFunction){
	fmt.Printf("[AddFunction]%v\n", f)
	for _, v := range s.Functions_ {
		if (f.GetName() == v.GetName()) {
			panic("Function " + f.GetName() + " is already defined")
		}
	}
	fmt.Printf("[AddFunction]%v\n", f)
	s.Functions_ = append(s.Functions_, f);
}

func (s *TService)GetFunctions() []*TFunction{
	return s.Functions_
}

func (s *TService)GetExtends() *TService{
	return s.Extends_
}

func (s *TService)GetFingerprintMaterial() string{
	// Services should never be used in fingerprints.
	panic("BUG: Can't get fingerprint material for service.")
	return "112131"
}
