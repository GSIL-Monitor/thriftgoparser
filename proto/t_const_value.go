package proto

import (
	"fmt"
	"strings"
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

/**
 * A const value is something parsed that could be a map, set, list, struct
 * or whatever.
 *
 */

const( // ConstValueType
	CV_INTEGER = iota
	CV_DOUBLE
	CV_STRING
	CV_MAP
	CV_LIST
	CV_IDENTIFIER
)

type TConstValue struct {
	MapVal_ map[*TConstValue]*TConstValue
	ListVal_ []*TConstValue
	StringVal_ string
	IntVal_ int64
	DoubleVal_ float64
	IdentifierVal_ string

	Enum_ *TEnum
	ValType_ int32 // ConstValueType
}

func NewConstValue(value interface{}) *TConstValue{
	var ret = new(TConstValue)
	switch value.(type) {
	case string:
		ret.SetString(value.(string))
	case float64:
		ret.SetDouble(value.(float64))
	case int64:
		ret.SetInteger(value.(int64))
	}

	return ret
}

func (c *TConstValue)SetString(val string){
	c.ValType_ = CV_STRING
	c.StringVal_ = val
}

func (c *TConstValue)GetString() string{
	return c.StringVal_
}

func (c *TConstValue)SetInteger(val int64){
	c.ValType_ = CV_INTEGER
	c.IntVal_ = val
}

func (c *TConstValue)GetInt() int64{
	if c.ValType_ == CV_IDENTIFIER {
		if c.Enum_ == nil {
			panic(fmt.Sprintf("have identifier \"" + c.GetIdentifier() + "\", but unset enum on line!"));
		}


		identifier := c.GetIdentifier();
		pos := strings.IndexByte(identifier, '.')
		if pos != -1 {
			identifier = identifier[pos+1 : len(identifier)-1]
		}

		enum_val := c.Enum_.GetConstantByName("")
		if (enum_val == nil) {
			panic(fmt.Sprintf("Unable to find enum value \"" + identifier + "\" in enum \"" + c.Enum_.GetName() + "\""))
		}

		return enum_val.GetValue()
	} else {
		return c.IntVal_
	}
}

func (c *TConstValue)SetDouble(val float64){
	c.ValType_ = CV_DOUBLE
	c.DoubleVal_ = val
}

func (c *TConstValue)GetDouble() float64{
	return c.DoubleVal_
}

func (c *TConstValue)SetMap(){
	c.ValType_ = CV_MAP
}

func (c *TConstValue)AddMap(key *TConstValue, val *TConstValue){
	c.MapVal_[key] = val
}

func (c *TConstValue)GetMap() (map[*TConstValue] *TConstValue){
	return c.MapVal_
}

func (c *TConstValue)SetList(){
	c.ValType_ = CV_LIST
}

func (c *TConstValue)AddList(val *TConstValue){
	c.ListVal_ = append(c.ListVal_, val)
}

func (c *TConstValue)GetList() ([]*TConstValue){
	return c.ListVal_
}

func (c *TConstValue)SetIdentifier(val string){
	c.ValType_ = CV_IDENTIFIER
	c.IdentifierVal_ = val
}

func (c *TConstValue)GetIdentifier() string{
	return c.IdentifierVal_;
}

func (c *TConstValue)GetIdentifierName() string{
	identifier := c.GetIdentifier();
	pos := strings.IndexByte(identifier, '.')
	if pos == -1 {
		panic(fmt.Sprint("error: identifier " + identifier + " is unqualified!"))
	}

	identifier = identifier[pos+1:]
	pos2 := strings.IndexByte(identifier, '.')
	if pos2 != -1{
		identifier  = identifier[pos2+1:]
	}

	return identifier;
}

func (c *TConstValue)GetIdentifierWithParent() string{
	identifier := c.GetIdentifier();
	pos := strings.IndexByte(identifier, '.')
	if pos == -1 {
		panic(fmt.Sprint("error: identifier " + identifier + " is unqualified!"))
	}

	tmp := identifier[pos+1:]
	pos2 := strings.IndexByte(tmp, '.')
	if pos2 != -1{
		return tmp
	}

	return identifier;
}

func (c *TConstValue)SetEnum(em *TEnum) {
	c.Enum_ = em
}

func (c *TConstValue)GetType() int32{
	return c.ValType_;
}



