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


/**
 * A map is a lightweight container type that just wraps another two data
 * types.
 *
 */

type TMap struct{
	TContainer
	KeyType_ ICafType
	ValType_ ICafType
}

func NewMap(key_type, val_type ICafType) *TMap{
	ret := new(TMap)
	ret.ValType_ = val_type
	ret.KeyType_ = key_type
	return ret
}

func (m *TMap)GetKeyType() ICafType{
	return m.KeyType_
}

func (m *TMap)GetValType() ICafType{
	return m.ValType_;
}

func (m *TMap)IsMap() bool{
	return true
}

func (m *TMap)GetFingerprintMaterial() string{
	return fmt.Sprintf("map<:%s,%s>", m.KeyType_.GetFingerprintMaterial(), m.ValType_.GetFingerprintMaterial())
}

func (m *TMap)GenerateFingerprint() {
	GenerateBaseFingerprint(m)
	m.KeyType_.GenerateFingerprint()
	m.ValType_.GenerateFingerprint()
}