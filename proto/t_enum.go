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
 * An enumerated type. A list of constant objects with a name for the type.
 *
 */

type TEnum struct{
	TTypeDesc
	Constants_ []*TEnumValue
}

func NewEnum(program *TProgram) *TEnum{
	var ret = new(TEnum)
	ret.Program_ = program
	return ret
}

func (em *TEnum)GetConstants() []*TEnumValue {
	return em.Constants_
}

func (em *TEnum)Append(value *TEnumValue) {
	em.Constants_ = append(em.Constants_, value)
}

func (em *TEnum)GetConstantByName(name string) *TEnumValue {
	for _, v := range em.Constants_{
		if (v.Name_ == name) {
			return v
		}
	}
	return nil
}

func (em *TEnum)GetConstantByValue(value int64) *TEnumValue{
	for _, v := range em.Constants_ {
		if (v.Value_  == value) {
			return v
		}
	}
	return nil
}

func (em *TEnum)GetMinValue() *TEnumValue{
	if len(em.Constants_) == 0 {
		return nil
	}

	var ret *TEnumValue
	for _, v := range em.Constants_{
		if nil == ret || ret.Value_ > v.Value_ {
			ret = v
		}
	}

	return ret
}

func (em* TEnum)GetMaxValue() *TEnumValue {
	if len(em.Constants_) == 0 {
		return nil
	}


	var ret *TEnumValue
	for _, v := range em.Constants_{
		if nil == ret || ret.Value_ < v.Value_ {
			ret = v
		}
	}

	return ret
}

func (em *TEnum) IsEnum() bool{
	return true
}

//virtual std::string get_fingerprint_material() const {
func (e *TEnum)GetFingerprintMaterial() string{
	return "enum";
}

