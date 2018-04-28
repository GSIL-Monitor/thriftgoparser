%{
// Copyright (c) 2011 CZ.NIC z.s.p.o. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// blame: jnml, labs.nic.cz

// +build ignore

package main

import (
	"fmt"
	"strconv"
	"io/ioutil"
	"os"
	//"unicode/utf8"
	"runtime/debug"
)

func thrift_reserved_keyword(keyword []byte) {
    panic(fmt.Sprintf("Cannot use reserved language keyword: %q\n", keyword))
}

func integer_overflow(text string) {
    panic(fmt.Sprintf("This integer is too big: \"%s\"\n", text))
}

func unexpected_token(text string) {
    panic(fmt.Sprintf("Unexpected token in input: \"%s\"\n", text))
}

func error_unsupported_namespace_decl(name string){
    panic(fmt.Sprintf("unsupported namespace:%s\n", name))
}

func emit_byte_type_warning() {
	fmt.Printf("The \"byte\" type is a compatibility alias for \"i8\". Use \"i8\" to emphasize the signedness of this type.\n");
}

type CafLexer struct {
	FilePath  string
	HasErrors bool
	ErrMsg    string

	Src      []byte   //当前文档的指针
	CurChar    byte   // 当前读的字符
	In_stack   []byte // 当前已经在栈内
	Un_stack   []byte // 已读的吐回
	Mark int
	CurLineNo    int    // 当前所在行数
}

func NewFileLexer(file_path string) *CafLexer{
	cl := new(CafLexer)
	cl.FilePath = file_path
	cl.Init()
	return cl
}

func (cl *CafLexer) Init(){
	file, err := os.Open(cl.FilePath)
	if err != nil {
		panic(err)
	}

	defer file.Close()
	defer func() {
		if x := recover(); x != nil {
			debug.PrintStack()
			fmt.Println(x)
		}
	}()

	cl.Src, err = ioutil.ReadAll(file)
	cl.Next()
}

func (cl *CafLexer)Next() {
	if len(cl.Un_stack) != 0 {
		cl.CurChar = cl.Un_stack[len(cl.Un_stack)-1]
		cl.Un_stack = cl.Un_stack[:len(cl.Un_stack)-1]
		fmt.Printf("Get content from un stack:%q\n", cl.Un_stack)
		return
	}

	cl.In_stack = append(cl.In_stack, cl.CurChar)
	if len(cl.Src) == 0 {
		cl.CurChar = 0 // EOF
		return
	}

	cl.CurChar = cl.Src[0]
	//fmt.Printf("----next: %q, in:%q, un:%q sc:%d\n", c, in_stack, un, sc)
	cl.Src = cl.Src[1:]
}

func (cl *CafLexer)Unget(b byte) {
	cl.Un_stack = append(cl.Un_stack, b)
}

func (cl *CafLexer) Error(s string) {
	cl.ErrMsg = s
	cl.HasErrors = true
	fmt.Println("[ERROR]Parser:", s)
}

func (cl *CafLexer)Reduced(rule, state int, lval *yySymType) bool{
	//fmt.Printf("[Debug]rule:%d, state:%d, token:\n%v\n", lval)
	//fmt.Printf("[Debug]in_stack:%q cur:%q\n", cl.In_stack, cl.CurChar)
	return false
}

func (cl *CafLexer)Match(){
    //fmt.Printf("\tstate accepts: %q\n", cl.In_stack)
    cl.Mark = len(cl.In_stack)
}

func (cl *CafLexer)Lex(yylval *yySymType) int{
%}

%yyn cl.Next()
%yyc cl.CurChar
%yym cl.Match()
%yyt sc

intconstant   ([+-]?[0-9]+)
hexconstant   ([+-]?"0x"[0-9A-Fa-f]+)
dubconstant   ([+-]?[0-9]*(\.[0-9]+)?([eE][+-]?[0-9]+)?)
identifier    ([a-zA-Z_](\.[a-zA-Z_0-9]|[a-zA-Z_0-9])*)
whitespace    ([ \t\r\n]*)
sillycomm     ("/*""*"*"*/")
multicm_begin ("/*")
doctext_begin ("/**")
comment       ("//"[^\n]*)
unixcomment   ("#"[^\n]*)
symbol        ([:;\,\{\}\(\)\=<>\[\]])
literal_begin (['\"])

%%

    cl.In_stack = cl.In_stack[:0]
	cl.Mark = -1

\0
    return -1

{whitespace} /* do nothing */
{sillycomm} /* do nothing */

{doctext_begin}
  state := 0  // 0 = normal, 1 = "*" seen, "*/" seen
  for(state < 2){
    cl.Next()
    switch cl.CurChar {
      case 0:
        panic(fmt.Sprintf("Unexpected end of file in doc-comment at %d\n", cl.CurLineNo))
      case '*':
        state = 1
      case '/':
        {
          if (1 == state){
            state = 2
          }else{
            state = 0
          }
        }
      default:
        state = 0
    }
  }
  cl.Next()
  fmt.Printf("[doctext]%q\n", cl.In_stack)

{multicm_begin}
  state := 0 // 0 = normal, 1 = "*" seen, "*/" seen
  for state < 2 {
    cl.Next()
    switch cl.CurChar {
      case 0:
        panic(fmt.Sprintf("Unexpected end of file in multiline comment at %d\n", cl.CurLineNo))
      case '*':
        state = 1
      case '/':
        if state == 1{
            state = 2
        }else{
            state = 0
        }
      default:
        state = 0
    }
  }
  cl.Next()
  fmt.Printf("[multi_comm]%q\n",cl.In_stack)

{comment}
    fmt.Printf("comment\n")

{unixcomment}
    fmt.Printf("unix comment\n")

{symbol}
    return int(cl.In_stack[0])

"*"
    return int(cl.In_stack[0])

"false"              return tok_bool_constant
"true"               return tok_bool_constant

"namespace"          return tok_namespace
"cpp_namespace"      error_unsupported_namespace_decl("cpp"); /* do nothing */
"cpp_include"        return tok_cpp_include
"cpp_type"           return tok_cpp_type
"java_package"       error_unsupported_namespace_decl("java_package") /* do nothing */
"cocoa_prefix"       error_unsupported_namespace_decl("cocoa_prefix") /* do nothing */
"csharp_namespace"   error_unsupported_namespace_decl("csharp") /* do nothing */
"delphi_namespace"   error_unsupported_namespace_decl("delphi") /* do nothing */
"php_namespace"      error_unsupported_namespace_decl("php") /* do nothing */
"py_module"          error_unsupported_namespace_decl("py_module") /* do nothing */
"perl_package"       error_unsupported_namespace_decl("perl_package") /* do nothing */
"ruby_namespace"     error_unsupported_namespace_decl("ruby") /* do nothing */
"smalltalk_category" error_unsupported_namespace_decl("smalltalk_category") /* do nothing */
"smalltalk_prefix"   error_unsupported_namespace_decl("smalltalk_category") /* do nothing */
"xsd_all"            return tok_xsd_all
"xsd_optional"       return tok_xsd_optional
"xsd_nillable"       return tok_xsd_nillable
"xsd_namespace"      error_unsupported_namespace_decl("xsd") /* do nothing */
"xsd_attrs"          return tok_xsd_attrs
"include"            return tok_include
"void"               return tok_void
"bool"               return tok_bool

"byte"
  emit_byte_type_warning()
  return tok_i8

"i8"                 return tok_i8
"i16"                return tok_i16
"i32"                return tok_i32
"i64"                return tok_i64
"double"             return tok_double
"string"             return tok_string
"binary"             return tok_binary
"slist"
  fmt.Printf("\"slist\" is deprecated and will be removed in a future compiler version.  This type should be replaced with \"string\".\n")
  return tok_slist;

"senum"
  fmt.Printf("\"senum\" is deprecated and will be removed in a future compiler version.  This type should be replaced with \"string\".\n")
  return tok_senum;

"map"                return tok_map
"list"               return tok_list
"set"                return tok_set
"oneway"             return tok_oneway
"typedef"            return tok_typedef
"struct"             return tok_struct
"union"              return tok_union
"exception"          return tok_xception
"extends"            return tok_extends
"throws"             return tok_throws
"service"            return tok_service
"enum"               return tok_enum
"const"              return tok_const
"required"           return tok_required
"optional"           return tok_optional
"async"
  fmt.Printf("\"async\" is deprecated.  It is called \"oneway\" now.\n")
  return tok_oneway

"&"                  return tok_reference


"BEGIN"              thrift_reserved_keyword(cl.In_stack)
"END"                thrift_reserved_keyword(cl.In_stack)
"__CLASS__"          thrift_reserved_keyword(cl.In_stack)
"__DIR__"            thrift_reserved_keyword(cl.In_stack)
"__FILE__"           thrift_reserved_keyword(cl.In_stack)
"__FUNCTION__"       thrift_reserved_keyword(cl.In_stack)
"__LINE__"           thrift_reserved_keyword(cl.In_stack)
"__METHOD__"         thrift_reserved_keyword(cl.In_stack)
"__NAMESPACE__"      thrift_reserved_keyword(cl.In_stack)
"abstract"           thrift_reserved_keyword(cl.In_stack)
"alias"              thrift_reserved_keyword(cl.In_stack)
"and"                thrift_reserved_keyword(cl.In_stack)
"args"               thrift_reserved_keyword(cl.In_stack)
"as"                 thrift_reserved_keyword(cl.In_stack)
"assert"             thrift_reserved_keyword(cl.In_stack)
"begin"              thrift_reserved_keyword(cl.In_stack)
"break"              thrift_reserved_keyword(cl.In_stack)
"case"               thrift_reserved_keyword(cl.In_stack)
"catch"              thrift_reserved_keyword(cl.In_stack)
"class"              thrift_reserved_keyword(cl.In_stack)
"clone"              thrift_reserved_keyword(cl.In_stack)
"continue"           thrift_reserved_keyword(cl.In_stack)
"declare"            thrift_reserved_keyword(cl.In_stack)
"def"                thrift_reserved_keyword(cl.In_stack)
"default"            thrift_reserved_keyword(cl.In_stack)
"del"                thrift_reserved_keyword(cl.In_stack)
"delete"             thrift_reserved_keyword(cl.In_stack)
"do"                 thrift_reserved_keyword(cl.In_stack)
"dynamic"            thrift_reserved_keyword(cl.In_stack)
"elif"               thrift_reserved_keyword(cl.In_stack)
"else"               thrift_reserved_keyword(cl.In_stack)
"elseif"             thrift_reserved_keyword(cl.In_stack)
"elsif"              thrift_reserved_keyword(cl.In_stack)
"end"                thrift_reserved_keyword(cl.In_stack)
"enddeclare"         thrift_reserved_keyword(cl.In_stack)
"endfor"             thrift_reserved_keyword(cl.In_stack)
"endforeach"         thrift_reserved_keyword(cl.In_stack)
"endif"              thrift_reserved_keyword(cl.In_stack)
"endswitch"          thrift_reserved_keyword(cl.In_stack)
"endwhile"           thrift_reserved_keyword(cl.In_stack)
"ensure"             thrift_reserved_keyword(cl.In_stack)
"except"             thrift_reserved_keyword(cl.In_stack)
"exec"               thrift_reserved_keyword(cl.In_stack)
"finally"            thrift_reserved_keyword(cl.In_stack)
"float"              thrift_reserved_keyword(cl.In_stack)
"for"                thrift_reserved_keyword(cl.In_stack)
"foreach"            thrift_reserved_keyword(cl.In_stack)
"from"               thrift_reserved_keyword(cl.In_stack)
"function"           thrift_reserved_keyword(cl.In_stack)
"global"             thrift_reserved_keyword(cl.In_stack)
"goto"               thrift_reserved_keyword(cl.In_stack)
"if"                 thrift_reserved_keyword(cl.In_stack)
"implements"         thrift_reserved_keyword(cl.In_stack)
"import"             thrift_reserved_keyword(cl.In_stack)
"in"                 thrift_reserved_keyword(cl.In_stack)
"inline"             thrift_reserved_keyword(cl.In_stack)
"instanceof"         thrift_reserved_keyword(cl.In_stack)
"interface"          thrift_reserved_keyword(cl.In_stack)
"is"                 thrift_reserved_keyword(cl.In_stack)
"lambda"             thrift_reserved_keyword(cl.In_stack)
"module"             thrift_reserved_keyword(cl.In_stack)
"native"             thrift_reserved_keyword(cl.In_stack)
"new"                thrift_reserved_keyword(cl.In_stack)
"next"               thrift_reserved_keyword(cl.In_stack)
"nil"                thrift_reserved_keyword(cl.In_stack)
"not"                thrift_reserved_keyword(cl.In_stack)
"or"                 thrift_reserved_keyword(cl.In_stack)
"package"            thrift_reserved_keyword(cl.In_stack)
"pass"               thrift_reserved_keyword(cl.In_stack)
"public"             thrift_reserved_keyword(cl.In_stack)
"print"              thrift_reserved_keyword(cl.In_stack)
"private"            thrift_reserved_keyword(cl.In_stack)
"protected"          thrift_reserved_keyword(cl.In_stack)
"raise"              thrift_reserved_keyword(cl.In_stack)
"redo"               thrift_reserved_keyword(cl.In_stack)
"rescue"             thrift_reserved_keyword(cl.In_stack)
"retry"              thrift_reserved_keyword(cl.In_stack)
"register"           thrift_reserved_keyword(cl.In_stack)
"return"             thrift_reserved_keyword(cl.In_stack)
"self"               thrift_reserved_keyword(cl.In_stack)
"sizeof"             thrift_reserved_keyword(cl.In_stack)
"static"             thrift_reserved_keyword(cl.In_stack)
"super"              thrift_reserved_keyword(cl.In_stack)
"switch"             thrift_reserved_keyword(cl.In_stack)
"synchronized"       thrift_reserved_keyword(cl.In_stack)
"then"               thrift_reserved_keyword(cl.In_stack)
"this"               thrift_reserved_keyword(cl.In_stack)
"throw"              thrift_reserved_keyword(cl.In_stack)
"transient"          thrift_reserved_keyword(cl.In_stack)
"try"                thrift_reserved_keyword(cl.In_stack)
"undef"              thrift_reserved_keyword(cl.In_stack)
"unless"             thrift_reserved_keyword(cl.In_stack)
"unsigned"           thrift_reserved_keyword(cl.In_stack)
"until"              thrift_reserved_keyword(cl.In_stack)
"use"                thrift_reserved_keyword(cl.In_stack)
"var"                thrift_reserved_keyword(cl.In_stack)
"virtual"            thrift_reserved_keyword(cl.In_stack)
"volatile"           thrift_reserved_keyword(cl.In_stack)
"when"               thrift_reserved_keyword(cl.In_stack)
"while"              thrift_reserved_keyword(cl.In_stack)
"with"               thrift_reserved_keyword(cl.In_stack)
"xor"                thrift_reserved_keyword(cl.In_stack)
"yield"              thrift_reserved_keyword(cl.In_stack)


{intconstant}
  yylval.iconst, _ = strconv.ParseInt(string(cl.In_stack), 10, 64)
  return tok_int_constant


{hexconstant}
  yylval.iconst, _ = strconv.ParseInt(string(cl.In_stack), 10, 64)
  return tok_int_constant

{identifier}
  yylval.id = string(cl.In_stack)
  return tok_identifier

{dubconstant}
  yylval.dconst, _ = strconv.ParseFloat(string(cl.In_stack), 64)
  return tok_dub_constant

{literal_begin}
  tmp_mark := cl.In_stack[0]
  var result string
  for {
    cl.Next()
    switch cl.CurChar {
      case 0:
        panic(fmt.Sprintf("End of file while read string at %d\n", cl.CurLineNo))
      case '\n':
        panic(fmt.Sprintf("End of line while read string at %d\n", cl.CurLineNo - 1))
      case '\\':
        cl.Next()
        switch cl.CurChar {
          case 'r':
            result += "\r"
          case 'n':
            result += "\n"
          case 't':
            result += "\t"
          case '"':
            result += "\""
          case '\'':
            result += "'"
          case '\\':
            result += "\""
          default:
            panic("Bad escape character\n")
        }
        break;
      default:
        if cl.CurChar == tmp_mark {
           yylval.id = result
           return tok_literal
        } else {
           result += string(cl.CurChar)
        }
    }
  }

. unexpected_token(string(cl.In_stack))

%%
    fmt.Printf("Goes here\n")
	if cl.Mark >= 0 {
        if len(cl.In_stack) > cl.Mark {
            cl.Unget(cl.CurChar)
            for i := len(cl.In_stack)-1; i >= cl.Mark; i-- {
                fmt.Printf("cl.Unget:%s\n", string(cl.In_stack[i]))
                cl.Unget(cl.In_stack[i])
            }
            cl.Next()
        }
        cl.In_stack = cl.In_stack[:cl.Mark]
        goto yyAction // Hook: Execute the semantic action of the last matched rule.
    }

    switch n := len(cl.In_stack); n {
    case 0: // [] z
        fmt.Printf("zero: %q\n", cl.In_stack)
        cl.Next()
    case 1: // [x] z
        fmt.Printf("one %q\n", cl.In_stack)
    default: // [x, y, ...], z
        fmt.Printf("two %q\n", cl.In_stack)
        cl.Unget(cl.CurChar) // z
        for i := n - 1; i > 1; i-- {
            cl.Unget(cl.In_stack[i]) // ...
        }
        cl.CurChar = cl.In_stack[1] // y
    }

    return 1024
}
