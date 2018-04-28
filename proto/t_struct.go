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

// Forward declare that puppy

/**
 * A struct is a container for a set of member fields that has a name. Structs
 * are also used to implement exception types.
 *
 */

type TStruct struct{
	TTypeDesc
	Members_ []*TField
	MembersInIdOrder_ []*TField
	IsXception_ bool
	IsUnion_ bool
	MembersValidated bool
	MembersWithValue int
	XsdAll_ bool
}

func NewStruct(program* TProgram) *TStruct{
	var ret =  new(TStruct)
	ret.Program_ = program
	ret.IsXception_ = false
	ret.IsUnion_ = false
	ret.MembersValidated = false
	ret.MembersWithValue = 0
	ret.XsdAll_ = false
	return ret
}

func (st *TStruct) SetName(name string) {
	st.Name_ = name
	st.ValidateUnionMembers()
}

func (st *TStruct)SetXception(is_xception bool) {
	st.IsXception_ = is_xception
}

func (st *TStruct)ValidateUnionMember(f *TField) {
	if st.IsUnion_ && len(st.Name_)>0 {
		// unions can't have required fields
		if f.GetReq() == T_REQUIRED {
			f.SetReq(T_OPTIONAL)
			panic(fmt.Sprintf("Required field %s of union %s set to optional.\n", f.GetName(), st.Name_))
		}

		// unions may have up to one member defaulted, but not more
		if f.GetValue() != nil {
			st.MembersWithValue ++
			if ( st.MembersWithValue > 1) {
				panic(fmt.Sprintf("Error: Field " + f.GetName() + " provides another default value for union " + st.Name_))
			}
		}
	}
}

func (st *TStruct)ValidateUnionMembers(){
	if st.IsUnion_ && len(st.Name_) > 0 && (!st.MembersValidated) {
		for _, v := range st.MembersInIdOrder_{
			st.ValidateUnionMember(v)
		}
		st.MembersValidated  = true
	}
}

func (st *TStruct)SetUnion(is_union bool) {
	st.IsUnion_ = is_union
	st.ValidateUnionMembers()
}

func (st *TStruct)SetXsdAll(xsd_all bool){
	st.XsdAll_ = xsd_all
}

func (st *TStruct)GetXsdAll() bool{
	return st.XsdAll_
}

func (st *TStruct)Append(f *TField) bool{
	// returns false when there is a conflict of field names
	if nil != st.GetFieldByName(f.GetName()) {
		return false;
	}
	st.Members_ = append(st.Members_, f)
	//members_in_id_order_.insert(bounds.second, elem);
	// 把f按顺序加到membersInOrder里面
	st.ValidateUnionMember(f)
	return true
}

func (st *TStruct)GetMembers() []*TField{
	return st.Members_;
}

func  (st *TStruct)GetSortedMembers() []*TField{
	return st.MembersInIdOrder_;
}

func (st *TStruct)IsStruct() bool{
	return !st.IsXception_
}

func (st *TStruct)IsXception() bool{
	return st.IsXception_;
}

func (st* TStruct) IsUnion() bool{
	return st.IsUnion_;
}

func (st* TStruct)GetFingerprintMaterial() string{
	//doReserve := len(st.MembersInIdOrder_) > 1
	rv := "{"
	for _,v := range st.MembersInIdOrder_ {
		rv += v.GetFingerprintMaterial()
		rv += ";"

		//if len(st.MembersInIdOrder_) > 1 {
		//	estimation := len(st.MembersInIdOrder_) * len(rv) + 16 // 没看懂
		//	append(rv, estinamtion) // 也没看懂
		//	doReserve = false;
		//}
	}
	rv += "}";

	return rv;
}

// **重要**
func (st *TStruct)GenerateFingerprint() {
	GenerateBaseFingerprint(st) // 这里是强制基类的调用
	for _, v := range st.MembersInIdOrder_{
		v.GetType().GenerateFingerprint() // 调子成员的派生类
	}
}

func (st* TStruct)GetFieldByName(field_name string) *TField {
	for _,v := range st.GetMembers() {
		if (v.GetName() == field_name) {
			return v;
		}
	}
	return nil;
}
