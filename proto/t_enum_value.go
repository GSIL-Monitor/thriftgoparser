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
 * A constant. These are used inside of enum definitions. Constants are just
 * symbol identifiers that may or may not have an explicit value associated
 * with them.
 *
 */

type TEnumValue struct{
	TDoc
	Name_ string
	Value_ int64
	Annotations_ map[string]string
}

func NewEnumValue(name string, val int64) *TEnumValue{
	ret := new(TEnumValue)
	ret.Name_ = name
	ret.Value_ = val
	return ret
}

func (ev *TEnumValue)GetName() string{
	return ev.Name_
}

func (ev *TEnumValue)GetValue() int64{
	return ev.Value_
}


func (ev *TEnumValue)GetAnnotations() map[string]string{
	return ev.GetAnnotations()
}

func (ev *TEnumValue)SetAnnotations(annotation map[string]string) {
	ev.Annotations_ = annotation
}