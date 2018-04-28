package proto

import (
	"fmt"
)

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

// Forward declare for xsd_attrs

/**
 * A simple struct for the parser to use to store a field ID, and whether or
 * not it was specified by the user or automatically chosen.
 */
type TFieldId struct{
	Value int
	AutoAssigned bool
}

func NewFieldId(val int, auto_assigned bool) *TFieldId{
	ret := new(TFieldId)
	ret.Value = val
	ret.AutoAssigned = auto_assigned
	return ret
}

const( //枚举 EReq
	T_REQUIRED = iota
	T_OPTIONAL
	T_OPT_IN_REQ_OUT
)

/**
 * Class to represent a field in a thrift structure. A field has a data type,
 * a symbolic name, and a numeric identifier.
 *
 */
type TField struct{
	TDoc
	Type_ ICafType
	Name_ string
	Key_ int32
	Req_ int32 // EReq
	Value_ *TConstValue
	XsdOptional_ bool
	XsdNillable_ bool
	XsdAttrs_ *TStruct
	Reference_ bool
	Annotations_ map[string] string
}

func NewField(t ICafType, name string, key int32) *TField{
	var ret = new(TField)
	ret.Type_ = t
	ret.Name_ = name
	ret.Key_ = key
	ret.Value_ = nil
	ret.XsdOptional_ = false
	ret.XsdNillable_ = false
	ret.XsdAttrs_ = nil
	ret.Reference_ = false
	return ret
}

func (f *TField)GetType() ICafType{
	return f.Type_
}

func (f *TField)GetName() string{
	return f.Name_
}

func (f *TField)GetKey() int32{
	return f.Key_
}

func (f *TField)SetReq(req int32) {
	f.Req_ = req;
}

func (f *TField)GetReq() int32{
	return f.Req_;
}

func (f *TField)SetValue(value *TConstValue){
	f.Value_ = value;
}

func (f *TField)GetValue() *TConstValue {
	return f.Value_;
}

func (f *TField)SetXsdOptional(v bool){
	f.XsdOptional_ = v
}

func (f *TField)GetXsdOptional() bool{
	return f.XsdOptional_
}

func (f *TField)SetXsdNillable(v bool) {
	f.XsdNillable_ = v
}

func (f *TField)GetXsdNillable() bool{
	return f.XsdNillable_
}

func (f* TField)SetXsdAttrs(attrs *TStruct) {
	f.XsdAttrs_ = attrs
}

func (f *TField)GetXsdAttrs() *TStruct{
	return f.XsdAttrs_
}

// This is not the same function as t_type::get_fingerprint_material,
// but it does the same thing.
func (f *TField)GetFingerprintMaterial() string{
	if f.Req_ == T_OPTIONAL {
		return fmt.Sprintf("%d:opt-%s", f.Key_, f.Type_.GetFingerprintMaterial())
	}else{
		return fmt.Sprintf("%d:%s", f.Key_, f.Type_.GetFingerprintMaterial())
	}
}

func (f *TField)GetReference() bool{
	return f.Reference_
}

func (f *TField)SetReference(ref bool) {
	f.Reference_ = ref
}

/**
 * Comparator to sort fields in ascending order by key.
 * Make this a functor instead of a function to help GCC inline it.
 * The arguments are (const) references to const pointers to const t_fields.
 */
func Operator(left, right *TField) bool{
	return left.GetKey() < right.GetKey()
}


func (f *TField)GetAnnotations() map[string]string{
	return f.GetAnnotations()
}

func (f *TField)SetAnnotations(annotation map[string]string) {
	f.Annotations_ = annotation
}