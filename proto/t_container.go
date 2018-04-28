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

type TContainer struct{
	TTypeDesc
	CppName_ string
	HasCppName_ bool
}


func NewContainer() *TContainer{
	ret := new(TContainer)
	ret.HasCppName_ = false
	return ret
}

func (t* TContainer)SetCppName(cpp_name string) {
	t.CppName_ = cpp_name
	t.HasCppName_ = true
}

func (t* TContainer)has_cpp_name() bool{
	return t.HasCppName_
}

func (t* TContainer)get_cpp_name() string{
	return t.CppName_
}

func (t* TContainer) is_container() bool {
	return true
}

