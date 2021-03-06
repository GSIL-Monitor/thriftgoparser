package main

import (
	"fmt"
	"strings"
	"strconv"
	"github.com/AfLnk/thriftgoparser/proto"
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
 * Base class for a thrift code generator. This class defines the basic
 * routines for code generation and contains the top level method that
 * dispatches code generation across various components.
 *
 */
type TGenerator struct{
	/**
	 * Current code indentation level
	 */
	Indent_ int

	/**
	 * Temporary variable counter, for making unique variable names
	 */
	Tmp_ int

	/**
	 * The program being generated
	 */
	Program_ *proto.TProgram

	/**
	 * Quick accessor for formatted program name that is currently being
	 * generated.
	 */
	ProgramName_ string

	/**
	 * Quick accessor for formatted service name that is currently being
	 * generated.
	 */
	ServiceName_ string

	/**
	 * Output type-specifc directory name ("gen-*")
	 */
	OutDirBase_ string

	/**
	 * Map of characters to escape in string literals.
	 */
	Escape_ map[byte]string
}

func NewGenerator(program *proto.TProgram) *TGenerator{
	ret := new(TGenerator)
	ret.Tmp_ = 0
	ret.Indent_ = 0
	ret.Program_ = program
	ret.ProgramName_ = program.GetName()
	ret.Escape_ = make(map[byte]string)
	ret.Escape_['\n'] = "\\n"
	ret.Escape_['\r'] = "\\r"
	ret.Escape_['\t'] = "\\t"
	ret.Escape_['\t'] = "\\\""
	ret.Escape_['\\'] = "\\\\"
	return ret
}

func (g *TGenerator)GetProgram() *proto.TProgram{
	return g.Program_
}

/**
 * check whether sub-namespace declaraction is used by generator.
 * e.g. allow
 * namespace py.twisted bar
 * to specify namespace to use when -gen py:twisted is specified.
 * Will be called with subnamespace, i.e. is_valid_namespace("twisted")
 * will be called for the above example.
 */
func (g *TGenerator)IsValidNamespace(sub_namespace string) bool{ // 没看懂
	return false
}

func (g *TGenerator)get_escaped_string(constval *proto.TConstValue) string{
	return constval.GetString()
}

/**
 * Optional methods that may be implemented by subclasses to take necessary
 * steps at the beginning or end of code generation.
 */

func (g *TGenerator)init_generator() {

}

func (g *TGenerator)close_generator() {

}
/**
 * Pure virtual methods implemented by the generator subclasses.
 */
func (g *TGenerator) generate_xception(txception *proto.TStruct) {
	// By default exceptions are the same as structs
	g.generate_struct(txception, true)
}

/**
 * Method to get the program name, may be overridden
 */
func (g *TGenerator)get_program_name(program *proto.TProgram) string {
	return program.GetName()
}

/**
 * Method to get the service name, may be overridden
 */
func (g *TGenerator)get_service_name(tservice *proto.TService)  string{
	return tservice.GetName();
}

/**
 * Creates a unique temporary variable name, which is just "name" with a
 * number appended to it (i.e. name35)
 */
func (g *TGenerator)Tmp(name string) string{
	g.Tmp_ ++
	return fmt.Sprintf("%s%d", name, g.Tmp_)
}

/**
 * Generates a comment about this code being autogenerated, using C++ style
 * comments, which are also fair game in Java / PHP, yay!
 *
 * @return C-style comment mentioning that this file is autogenerated.
 */
func (g *TGenerator)autogen_comment() string {
	return "/**\n* Autogenerated by Thrift Compiler(Lobster)\n *\n * DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING\n *  @generated\n */\n)"
}

/**
 * Indentation level modifiers
 */

func (g *TGenerator)indent_up() {
	g.Indent_ ++
}

func (g *TGenerator)indent_down() {
	g.Indent_ --
}

/**
* Indentation validation helper
*/
func (g *TGenerator) indent_count() int{
	return g.Indent_
}

func (g *TGenerator) Indent_validate(expected int, func_name string) {
	if (g.Indent_ != expected) {
		fmt.Printf("Wrong indent count in %s: difference = %i \n", func_name, (expected - g.Indent_))
	}
}

/**
 * Indentation print function
 */

func (g *TGenerator) Indent() string{
	ind := ""
	for i := 0; i < g.Indent_; i ++ {
		ind += "  "
	}
	return ind;
}

/**
 * Capitalization helpers
 */
func (g *TGenerator)capitalize(in string) string{
	s := strings.ToUpper(string(in[0]))
	s += in[1:]
	return s
}

func (g *TGenerator)decapitalize(in string) string{
	s := strings.ToLower(string(in[0]))
	s += in[1:]
	return s
}

func (g *TGenerator)lowercase(in string) string {
	return strings.ToLower(in)
}

func (g *TGenerator)uppercase(in string) string {
	return strings.ToUpper(in)
}
/**
 * Transforms a camel case string to an equivalent one separated by underscores
 * e.g. aMultiWord -> a_multi_word
 *      someName   -> some_name
 *      CamelCase  -> camel_case
 *      name       -> name
 *      Name       -> name
 */
func (g *TGenerator)underscore(in string) string{
	// TODO
	return in
}
/**
  * Transforms a string with words separated by underscores to a camel case equivalent
  * e.g. a_multi_word -> aMultiWord
  *      some_name    ->  someName
  *      name         ->  name
  */
func (g *TGenerator)camelcase(in string) string{
	return in
}

func (g *TGenerator)emit_double_as_string(val float64) string{
	return strconv.FormatFloat(val, 'f', 2,10)
}

/**
 * Generates code for an enumerated type. Done using a class to scope
 * the values.
 *
 * @param tenum The enumeration
 */
func (g *TGenerator)generate_enum(em *proto.TEnum) {
	fmt.Printf("[Em]%v\n", em)
}

func (g *TGenerator)Init_generator(){
	fmt.Printf("Init_generator\n")
}


func (g *TGenerator)Generate_typedef(td *proto.TTypedef) {
	fmt.Printf("[TypeDef]%s\n", td.GetName())
}

func (g *TGenerator)generate_forward_declaration(ts *proto.TStruct) {
	fmt.Printf("[Object]%s\n", ts.GetName())
}
/**
 * Generate a constant value
 */
func (g *TGenerator)generate_const(con *proto.TConst) {
	fmt.Printf("[Const]%v\n", con.GetName())
}

/**
 * Generates a struct definition for a thrift data type.
 *
 * @param tstruct The struct definition
 */
func (g *TGenerator)generate_struct(ts *proto.TStruct, is_exception bool){
	InsertOrUpdateStruct(ts)
}

//
// 生成五个元素
//   ServiceName
//     FuncName
//
//
func (g *TGenerator)generate_service(ts *proto.TService) {
	fmt.Printf("[Service]PSM:%s, NameSpace:%s, Name:%s\n", ts.GetProgram().PSM_, ts.GetProgram().GetNamespace(), ts.GetName())

}

func (g *TGenerator)Generate_program() {
	// Initialize the generator
	g.Init_generator()

	// Generate enums
	for _,v := range g.Program_.GetEnums() {
		g.generate_enum(v)
	}

	// Generate typedefs
	for _,v := range g.Program_.GetTypedefs(){
		g.Generate_typedef(v)
	}

	// Generate structs, exceptions, and unions in declared order
	for _, v := range g.Program_.GetObjects(){
		g.generate_forward_declaration(v)
	}

	for _, v := range g.Program_.GetObjects(){
		if (v.IsXception()) {
			g.generate_xception(v)
		} else {
			g.generate_struct(v, false)
		}
	}

	// Generate constants
	for _, v := range g.Program_.GetConsts() {
		g.generate_const(v)
	}

	// Generate services
	for _, v := range g.Program_.GetServices() {
		g.generate_service(v)
	}
}