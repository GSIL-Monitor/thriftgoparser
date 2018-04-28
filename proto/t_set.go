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
 * A set is a lightweight container type that just wraps another data type.
 *
 */

type TSet struct{
	TContainer
	ElemType_ ICafType
}

func NewSet(elem_type ICafType) *TSet{
	ret := new(TSet)
	ret.ElemType_ = elem_type
	return ret
}

func (s *TSet)GetElemType() ICafType{
	return s.ElemType_
}

func (s *TSet)IsSet() bool{
	return true
}

func (s *TSet)GetFingerprintMaterial() string{
	return "set<" + s.ElemType_.GetFingerprintMaterial() + ">";
}

func (s *TSet)GenerateFingerprint() {
	GenerateBaseFingerprint(s)
	s.ElemType_.GenerateFingerprint()
}