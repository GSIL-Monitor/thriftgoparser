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
package proto

import (
	"strings"
	"fmt"
	"os"
)

const(
	INCLUDES = 1
	PROGRAM = 2
)

var(
	G_program_doctext_candidate string

	/**
 	* Parsing mode, two passes up in this gin rummy!
 	*/
	G_parse_mode = PROGRAM

	/**
	  * Global program tree
      */
	G_program *TProgram

	/**
	 * Search path for inclusions
	 */
	G_incl_searchpath []string

	/**
	 * Current directory of file being parsed
	 */
	G_curdir string

	/**
	 * The First doctext comment
	 */
	//G_program_doctext_candidate string

	G_doctext string

	//G_doctext_lineno = 0

	// Initialize global types
	G_type_void *TBaseType
	G_type_string *TBaseType
	G_type_binary *TBaseType
	G_type_slist  *TBaseType
	G_type_bool *TBaseType
	G_type_i8 *TBaseType
	G_type_i16 *TBaseType
	G_type_i32 *TBaseType
	G_type_i64 *TBaseType
	G_type_double *TBaseType

	/**
 	* Global scope
 	*/
	G_scope *TScope

	/**
	 * Parent scope to also parse types
	 */
	G_parent_scope *TScope

	G_parent_prefix string
)

/**
 * Top level class representing an entire thrift program. A program consists
 * fundamentally of the following:
 *
 *   Typedefs
 *   Enumerations
 *   Constants
 *   Structs
 *   Exceptions
 *   Services
 *
 * The program module also contains the definitions of the base types.
 *
 */
type TProgram struct{
	TDoc              // 从DOC类派生
	Path_ string      // File path
	Name_ string // Name
	Namespace_ string // Namespace
	Language_ string // 语言
	Scope_ *TScope
	PSM_ string // 所在的psm

	// Annotations for dynamic namespaces
	NamespaceAnnotations_ map[string] map[string]string

	// C++ extra includes
	CppIncludes_ []string

	// C extra includes
	CIncludes_ []string


	// Included programs
	Includes_ []*TProgram

	// Include prefix for this program, if any
	IncludePrefix_ string

	// Components to generate code for
	Typedefs_ []*TTypedef
	Enums_ []*TEnum
	Consts_ []*TConst
	Objects_ []*TStruct
	Structs_ []*TStruct
	Xceptions_ []*TStruct
	Services_ []*TService
}

/**
 * Converts a string filename into a thrift program name
 */
func ProgramName(filename string) string{
	last_slash_pos := strings.LastIndexByte(filename, '/')
	if (last_slash_pos  != -1) {
		filename = filename[last_slash_pos+1:];
	}

	last_dot_pos := strings.LastIndexByte(filename, '.')
	if -1 != last_dot_pos {
		filename = filename[: last_dot_pos]
	}

	return filename;
}

func NewProgram(path, name, psm string) *TProgram{
	fmt.Println("NewProgram")
	ret := new(TProgram)
	ret.Path_ = path
	ret.Name_ = name
	ret.Scope_ = NewScope()
	ret.PSM_ = psm
	return ret
}

// Path accessor
func (p *TProgram)GetPath() string{
	return p.Path_
}

// Name accessor
func (p *TProgram)GetName() string{
	return p.Name_
}

// Namespace
func (p *TProgram)GetNamespace() string{
	return p.Namespace_
}

// Include prefix accessor
func (p *TProgram)GetIncludePrefix() string{
	return p.IncludePrefix_
}

// Accessors for program elements
func (p *TProgram) GetTypedefs() []*TTypedef{
	return p.Typedefs_
}

func (p *TProgram) GetEnums() []*TEnum{
	return p.Enums_
}

func (p *TProgram) GetConsts() []*TConst{
	return p.Consts_
}

func (p *TProgram) GetStructs() []*TStruct{
	return p.Structs_
}

func (p *TProgram) GetXceptions() []*TStruct{
	return p.Xceptions_
}

func (p *TProgram) GetObjects() []*TStruct{
	return p.Objects_
}

func (p *TProgram) GetServices() []*TService{
	return p.Services_
}

// Program elements
func (p *TProgram)AddTypedef  (td *TTypedef) {
	fmt.Printf("[AddTypedef]%v\n", td)
	p.Typedefs_ = append(p.Typedefs_, td)
}

func (p *TProgram)AddEnum     (te *TEnum) {
	fmt.Printf("[AddEnum]%v\n", te)
	p.Enums_ = append(p.Enums_, te)
}

func (p *TProgram)AddConst    (tc *TConst) {
	fmt.Printf("[AddConst]%v\n", tc)
	p.Consts_ = append(p.Consts_, tc)
}

func (p *TProgram)AddStruct   (ts *TStruct) {
	fmt.Printf("[AddStruct]%v\n", ts)
	p.Objects_ = append(p.Objects_, ts)
	p.Structs_ = append(p.Structs_, ts)
}

func (p *TProgram)AddXception (tx *TStruct) {
	fmt.Printf("[AddXCeption]%v\n", tx)
	p.Objects_ = append(p.Objects_, tx)
	p.Xceptions_ = append(p.Xceptions_, tx)
}

func (p *TProgram)AddService (ts *TService) {
	fmt.Printf("[AddService]%v\n", ts)
	p.Services_ = append(p.Services_, ts)
}

// Programs to include
func (p *TProgram) GetIncludes() []*TProgram{
	return p.Includes_
}

// Typename collision detection
/**
 * Search for typename collisions
 * @param t    the type to test for collisions
 * @return     true if a certain collision was found, otherwise false
 */
func (p *TProgram) IsUniqueTypename(caftype ICafType) bool{
	ret := p.ProgramTypenameCount(p, caftype)
	for _,v := range p.Includes_ {
		ret += p.ProgramTypenameCount(v, caftype);
	}
	return 0 == ret;
}

/**
 * Search all type collections for duplicate typenames
 * @param prog the program to search
 * @param t    the type to test for collisions
 * @return     the number of certain typename collisions
 */
func (p *TProgram)ProgramTypenameCount(prog *TProgram, caftype ICafType) int{
	ret := 0
	var tmp ICafType
	for _, v := range prog.Typedefs_{
		tmp = v
		if caftype != tmp && caftype.GetName() == v.GetName() && p.IsCommonNamespace(prog, caftype) {
			ret ++
		}
	}

	for _, v := range prog.Enums_{
		tmp = v
		if caftype != tmp && caftype.GetName() == v.GetName() && p.IsCommonNamespace(prog, caftype) {
			ret ++
		}
	}

	for _, v := range prog.Objects_{
		tmp = v
		if caftype != tmp && caftype.GetName() == v.GetName() && p.IsCommonNamespace(prog, caftype) {
			ret ++
		}
	}

	for _, v := range prog.Services_{
		tmp = v
		if caftype != tmp && caftype.GetName() == v.GetName() && p.IsCommonNamespace(prog, caftype) {
			ret ++
		}
	}

	return ret
}

/**
 * Determine whether identical typenames will collide based on namespaces.
 *
 * Because we do not know which languages the user will generate code for,
 * collisions within programs (IDL files) having namespace declarations can be
 * difficult to determine. Only guaranteed collisions return true (cause an error).
 * Possible collisions involving explicit namespace declarations produce a warning.
 * Other possible collisions go unreported.
 * @param prog the program containing the preexisting typename
 * @param t    the type containing the typename match
 * @return     true if a collision within namespaces is found, otherwise false
 */
func (p *TProgram) IsCommonNamespace(prog *TProgram, caftype ICafType) bool{
	return prog.GetNamespace() == caftype.GetProgram().GetNamespace()
}

// Scoping and namespacing
func (p *TProgram)SetNamespace(name string) {
	p.Namespace_ = name;
}

// Scope accessor
func (p *TProgram)Scope() *TScope{
	return p.Scope_
}

// Includes

func (p *TProgram)AddInclude(path, include_site string) {


	// include prefix for this program is the site at which it was included
	// (minus the filename)
	include_prefix := include_site
	last_slash_pos := strings.LastIndexByte(include_site, '/')
	if last_slash_pos != -1 {
		include_prefix = include_site[:last_slash_pos]
	}
	program := NewProgram(path, include_prefix, p.PSM_)
	program.SetIncludePrefix(include_prefix);
	p.Includes_ = append(p.Includes_, program)

}

func (p *TProgram) SetIncludePrefix(include_prefix string) {
	p.IncludePrefix_ = include_prefix

	// this is intended to be a directory; add a trailing slash if necessary
	if (len(p.IncludePrefix_) > 0 && p.IncludePrefix_[len(p.IncludePrefix_) - 1] != '/') {
		p.IncludePrefix_ += "/"
	}
}

// Language neutral namespace / packaging
func (p *TProgram) SetNamespaceWithLang(language string, name_space string) {
	p.Namespace_ = name_space
	p.Language_ = language
}

func (p *TProgram)GetNamespaceByLang(language string) string{
	if (language == p.Language_){
		return p.Namespace_
	}

	return ""
}

// Language specific namespace / packaging
func (p *TProgram)AddCppInclude(path string) {
	p.CppIncludes_ = append(p.CppIncludes_, path)
}

func (p *TProgram)GetCppIncludes() []string{
	return p.CppIncludes_
}

func (p *TProgram)AddCInclude(path string) {
	p.CIncludes_ = append(p.CIncludes_, path)
}

func (p *TProgram)GetCIncludes() []string{
	return p.CIncludes_;
}


func (p *TProgram) SetNamespaceAnnotations(lang string, annotations map[string]string) {
	p.NamespaceAnnotations_[lang] = annotations
}

/**
 * Gets the directory path of a filename
 */
func DirName(filename string) string {
	last_slash_pos := strings.LastIndexByte(filename, '/')
	if last_slash_pos == -1 {
		return "."
	}

	return filename[:last_slash_pos]
}

/**
 * We are sure the program doctext candidate is really the program doctext.
 */
func Declare_valid_program_doctext() {
	if G_program_doctext_candidate != "" && G_program_doctext_status == STILL_CANDIDATE {
		G_program_doctext_status = ABSOLUTELY_SURE
		fmt.Printf("%s\n", "program doctext set to ABSOLUTELY_SURE")
	} else {
		G_program_doctext_status = NO_PROGRAM_DOCTEXT
		fmt.Printf("%s\n", "program doctext set to NO_PROGRAM_DOCTEXT");
	}
}

/**
 * Finds the appropriate file path for the given filename
 */
func IncludeFile(filename string) string{
	// Absolute path? Just try that
	if filename[0] == '/'{ //真实地址，判断文件是不是打得开
		// Realpath!
		// cppcheck-suppress uninitvar
		_, err := os.Stat(filename)
		if err != nil && os.IsNotExist(err) {
			fmt.Printf("Cannot open include file %s\n", filename)
			return ""
		}

		return filename
	} else { // relative path, start searching
		// new search path with current dir global
		sp := []string{G_curdir}
		sp = append(sp, G_incl_searchpath...)
		// iterate through paths
		for  _, v := range sp {
			sfilename := v + "/" + filename
			_, err := os.Stat(sfilename)
			if err != nil && os.IsNotExist(err) {
				fmt.Printf("Cannot open include file %s\n", sfilename)
				continue
			}
		}
	}

	// Uh oh
	fmt.Printf("Could not find include file %s\n", filename)
	return ""
}

/**
 * Clears any previously stored doctext string.
 * Also prints a warning if we are discarding information.
 */
func Clear_doctext() {
	if G_doctext != "" {
		fmt.Printf("Uncaptured doctext at on line %d\n", G_doctext_lineno)
	}
	G_doctext = "";
}

/**
 * Emits a warning on list<byte>, binary type is typically a much better choice.
 */
func Check_for_list_of_bytes(list_elem_type ICafType) {
	if G_parse_mode == PROGRAM && list_elem_type == nil && list_elem_type.IsBaseType() {
		if list_elem_type.(*TBaseType).GetBase() == TYPE_I8{
			fmt.Println("Consider using the more efficient \"binary\" type instead of \"list<byte>\".");
		}
	}
}

func InitGlobals(){
	// Initialize global types
	G_type_void   = NewBaseType("void",   TYPE_VOID)
	G_type_string = NewBaseType("string", TYPE_STRING)
	G_type_binary = NewBaseType("string", TYPE_STRING)
	G_type_binary.SetBinary(true)
	G_type_slist  = NewBaseType("string", TYPE_STRING);
	G_type_slist.SetStringList(true);
	G_type_bool   = NewBaseType("bool",   TYPE_BOOL);
	G_type_i8   = NewBaseType("i8",   TYPE_I8);
	G_type_i16    = NewBaseType("i16",    TYPE_I16);
	G_type_i32    = NewBaseType("i32",    TYPE_I32);
	G_type_i64    = NewBaseType("i64",    TYPE_I64);
	G_type_double = NewBaseType("double", TYPE_DOUBLE);
}

/**
 * You know, when I started working on Thrift I really thought it wasn't going
 * to become a programming language because it was just a generator and it
 * wouldn't need runtime type information and all that jazz. But then we
 * decided to add constants, and all of a sudden that means runtime type
 * validation and inference, except the "runtime" is the code generator
 * runtime.
 */
func Validate_const_rec(name string, cafType ICafType, val *TConstValue) {
	if (cafType.IsVoid()) {
		panic("type error: cannot declare a void const: " + name)
	}

	if cafType.IsBaseType() {
		switch(cafType.(*TBaseType).GetBase()) {
		case TYPE_STRING:
			if val.GetType() != CV_STRING {
				panic("type error: const \"" + name + "\" was declared as string")
			}
		case TYPE_BOOL:
			if (val.GetType() != CV_INTEGER) {
				panic("type error: const \"" + name + "\" was declared as bool")
			}
		case TYPE_I8:
			if (val.GetType() != CV_INTEGER) {
				panic("type error: const \"" + name + "\" was declared as byte")
			}
		case TYPE_I16:
			if (val.GetType() != CV_INTEGER) {
				panic("type error: const \"" + name + "\" was declared as i16")
			}
		case TYPE_I32:
			if (val.GetType() != CV_INTEGER) {
				panic("type error: const \"" + name + "\" was declared as i32")
			}

		case TYPE_I64:
			if (val.GetType() != CV_INTEGER) {
				panic("type error: const \"" + name + "\" was declared as i64")
			}
		case TYPE_DOUBLE:
			if val.GetType() != CV_INTEGER && val.GetType() != CV_DOUBLE {
				panic("type error: const \"" + name + "\" was declared as double")
			}
		default:
			panic("compiler error: no const of base type " + name)
		}
	} else if (cafType.IsEnum()) {
		if (val.GetType() != CV_IDENTIFIER) {
			panic("type error: const \"" + name + "\" was declared as enum")
		}
		/*
// see if there's a dot in the identifier
std::string name_portion = value->get_identifier_name();

const vector<t_enum_value*>& enum_values = ((t_enum*)type)->get_constants();
vector<t_enum_value*>::const_iterator c_iter;
bool found = false;

for (c_iter = enum_values.begin(); c_iter != enum_values.end(); ++c_iter) {
if ((*c_iter)->get_name() == name_portion) {
found = true;
break;
}
}
if (!found) {
	panic("type error: const " + name + " was declared as type " + cafType.GetName() + " which is an enum, but " + val.GetIdentifier() + " is not a valid value for that enum")
}
} else if (cafType.IsStruct() || cafType.IsXception()) {
if (val.GetType() != CV_MAP) {
panic("type error: const \"" + name + "\" was declared as struct/xception")
}
const vector<t_field*>& fields = ((t_struct*)type)->get_members();
vector<t_field*>::const_iterator f_iter;

const map<t_const_value*, t_const_value*, t_const_value::value_compare>& val = value->get_map();
map<t_const_value*, t_const_value*, t_const_value::value_compare>::const_iterator v_iter;
for (v_iter = val.begin(); v_iter != val.end(); ++v_iter) {
if (v_iter->first->get_type() != t_const_value::CV_STRING) {
throw "type error: " + name + " struct key must be string";
}
t_type* field_type = NULL;
for (f_iter = fields.begin(); f_iter != fields.end(); ++f_iter) {
if ((*f_iter)->get_name() == v_iter->first->get_string()) {
field_type = (*f_iter)->get_type();
}
}
if (field_type == NULL) {
throw "type error: " + type->get_name() + " has no field " + v_iter->first->get_string();
}

validate_const_rec(name + "." + v_iter->first->get_string(), field_type, v_iter->second);
}
} else if (type->is_map()) {
t_type* k_type = ((t_map*)type)->get_key_type();
t_type* v_type = ((t_map*)type)->get_val_type();
const map<t_const_value*, t_const_value*, t_const_value::value_compare>& val = value->get_map();
map<t_const_value*, t_const_value*, t_const_value::value_compare>::const_iterator v_iter;
for (v_iter = val.begin(); v_iter != val.end(); ++v_iter) {
validate_const_rec(name + "<key>", k_type, v_iter->first);
validate_const_rec(name + "<val>", v_type, v_iter->second);
}
} else if (type->is_list() || type->is_set()) {
t_type* e_type;
if (type->is_list()) {
e_type = ((t_list*)type)->get_elem_type();
} else {
e_type = ((t_set*)type)->get_elem_type();
}
const vector<t_const_value*>& val = value->get_list();
vector<t_const_value*>::const_iterator v_iter;
for (v_iter = val.begin(); v_iter != val.end(); ++v_iter) {
validate_const_rec(name + "<elem>", e_type, *v_iter);
}
*/
	}
}

/**
 * Check the type of a default value assigned to a field.
 */
func Validate_field_value(field *TField, cv *TConstValue) {
	Validate_const_rec(field.GetName(), field.GetType(), cv);
}

func Get_true_type(caftype ICafType) ICafType{
	//return caftype.GetTrueType()
	return nil
}