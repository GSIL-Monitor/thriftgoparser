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
 * A thrift base type, which must be one of the defined enumerated types inside
 * this definition.
 *
 */
const(//Tbase
	TYPE_VOID = iota
	TYPE_STRING
	TYPE_BOOL
	TYPE_I8
	TYPE_I16
	TYPE_I32
	TYPE_I64
	TYPE_DOUBLE
)

type TBaseType struct{
	TTypeDesc
	Base_ int32 // 参考TBase
	StringList_ bool
	Binary_ bool
	StringEnum_ bool
	StringEnumVals_ []string
}

func NewBaseType(name string, base int32) *TBaseType{
	ret := new(TBaseType)
	ret.Base_ = base
	ret.Name_ = name
	ret.StringList_ = false
	ret.Binary_ = false
	ret.StringEnum_ = false
	return ret
}

func (t *TBaseType)GetBase() int32{
	return t.Base_;
}

func (t *TBaseType)IsVoid() bool{
	return t.Base_ == TYPE_VOID
}

func (t *TBaseType)IsString() bool{
	return t.Base_ == TYPE_STRING
}

func (t *TBaseType)IsBool() bool{
	return t.Base_ == TYPE_BOOL
}

func (t *TBaseType)SetStringList(val bool) {
	t.StringList_ = val;
}

func (t *TBaseType) IsStringList() bool {
	return t.Base_ == TYPE_STRING && t.StringList_
}

func (t *TBaseType)SetBinary(val bool){
	t.Binary_ = val
}

func (t *TBaseType)IsBinary() bool{
	return (t.Base_ == TYPE_STRING) && t.Binary_
}

func (t *TBaseType)SetStringEnum(val bool) {
	t.StringEnum_ = val
}

func (t *TBaseType)IsStringEnum() bool{
	return t.Base_ == TYPE_STRING && t.StringEnum_
}

func (t *TBaseType)AddStringEnumVal(val string) {
	t.StringEnumVals_ = append(t.StringEnumVals_, val)
}

func (t *TBaseType)GetStringEnumVals() []string{
	return t.StringEnumVals_
}

func (t *TBaseType)IsBaseType() bool{
	return true
}

func BaseName(base int32) string{
	switch base {
	case TYPE_VOID:
		return "void"
	case TYPE_STRING:
		return "string"
	case TYPE_BOOL:
		return "bool"
	case TYPE_I8:
		return "i8"
	case TYPE_I16:
		return "i16"
	case TYPE_I32:
		return "i32"
	case TYPE_I64:
		return "i64"
	case TYPE_DOUBLE:
		return "double"
	default:
		return "(unknown)"
	}
}

func (t *TBaseType)GetFingerprintMaterial() string{
	ret := BaseName(t.Base_)
	if (ret == "(unknown)") {
		panic("BUG: Can't get fingerprint material for this base type.")
	}
	return ret
}
