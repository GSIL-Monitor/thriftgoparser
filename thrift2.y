%{
package main

import(
    "fmt"
    "github.com/AfLnk/thriftgoparser/proto"
)
%}

%{
/**
 * This global variable is used for automatic numbering of field indices etc.
 * when parsing the members of a struct. Field values are automatically
 * assigned starting from -1 and working their way down.
 */
    var y_field_val = -1
/**
 * This global variable is used for automatic numbering of enum values.
 * y_enum_val is the last value assigned; the next auto-assigned value will be
 * y_enum_val+1, and then it continues working upwards.  Explicitly specified
 * enum values reset y_enum_val to that value.
 */
    var y_enum_val int64 = -1
    var G_arglist = 0
    var struct_is_struct int64 = 0
    var struct_is_union int64 = 1
    var g_strict = 255

    var SHRT_MIN int = 0
    var SHRT_MAX int = 2147483647
%}

/**
 * This structure is used by the parser to hold the data types associated with
 * various parse nodes.
 */
%union {
  id string
  iconst int64
  dconst float64
  bconst bool
  tbool bool
  tdoc *proto.TDoc
  caftype proto.ICafType
  tbase *proto.TBaseType
  ttypedef *proto.TTypedef
  tenum *proto.TEnum
  tenumv *proto.TEnumValue
  tconst *proto.TConst
  tconstv *proto.TConstValue
  tstruct *proto.TStruct
  tservice *proto.TService
  tfunction *proto.TFunction
  tcontainer *proto.TContainer
  tlist *proto.TList
  tmap *proto.TMap
  tset *proto.TSet
  tfield *proto.TField
  dtext string
  ereq int32 // TField:EReq
  tannot *proto.TAnnotation
  tfieldid *proto.TFieldId
}

/**
 * Strings identifier
 */
%token<id>     tok_identifier
%token<id>     tok_literal
%token<dtext>  tok_doctext

/**
 * Constant values
 */
%token<iconst> tok_int_constant
%token<dconst> tok_dub_constant
%token<bconst> tok_bool_constant

/**
 * Header keywords
 */
%token tok_include
%token tok_namespace
%token tok_cpp_include
%token tok_cpp_type
%token tok_xsd_all
%token tok_xsd_optional
%token tok_xsd_nillable
%token tok_xsd_attrs

/**
 * Base datatype keywords
 */
%token tok_void
%token tok_bool
%token tok_string
%token tok_binary
%token tok_slist
%token tok_senum
%token tok_i8
%token tok_i16
%token tok_i32
%token tok_i64
%token tok_double

/**
 * Complex type keywords
 */
%token tok_map
%token tok_list
%token tok_set

/**
 * Function modifiers
 */
%token tok_oneway

/**
 * Thrift language keywords
 */
%token tok_typedef
%token tok_struct
%token tok_xception
%token tok_throws
%token tok_extends
%token tok_service
%token tok_enum
%token tok_const
%token tok_required
%token tok_optional
%token tok_union
%token tok_reference

/**
 * Grammar nodes
 */

%type<tbase>     BaseType
%type<tbase>     SimpleBaseType
%type<caftype>     ContainerType
%type<caftype>     SimpleContainerType
%type<tmap>     MapType
%type<tset>     SetType
%type<tlist>     ListType

%type<tdoc>      Definition
%type<caftype>     TypeDefinition

%type<ttypedef>  Typedef

%type<caftype>     TypeAnnotations
%type<caftype>     TypeAnnotationList
%type<tannot>    TypeAnnotation
%type<id>        TypeAnnotationValue

%type<tfield>    Field
%type<tfieldid>  FieldIdentifier
%type<ereq>      FieldRequiredness
%type<caftype>     FieldType
%type<tconstv>   FieldValue
%type<tstruct>   FieldList
%type<tbool>     FieldReference

%type<tenum>     Enum
%type<tenum>     EnumDefList
%type<tenumv>    EnumDef
%type<tenumv>    EnumValue

%type<ttypedef>  Senum
%type<tbase>     SenumDefList
%type<id>        SenumDef

%type<tconst>    Const
%type<tconstv>   ConstValue
%type<tconstv>   ConstList
%type<tconstv>   ConstListContents
%type<tconstv>   ConstMap
%type<tconstv>   ConstMapContents

%type<iconst>    StructHead
%type<tstruct>   Struct
%type<tstruct>   Xception
%type<tservice>  Service

%type<tfunction> Function
%type<caftype>   FunctionType
%type<tservice>  FunctionList

%type<tstruct>   Throws
%type<tservice>  Extends
%type<tbool>     Oneway
%type<tbool>     XsdAll
%type<tbool>     XsdOptional
%type<tbool>     XsdNillable
%type<tstruct>   XsdAttributes
%type<id>        CppType

%type<dtext>     CaptureDocText

%%

/**
 * Thrift Grammar Implementation.
 *
 * For the most part this source file works its way top down from what you
 * might expect to find in a typical .thrift file, i.e. type definitions and
 * namespaces up top followed by service definitions using those types.
 */

Program:
  HeaderList DefinitionList
    {
      fmt.Println("Program -> Headers DefinitionList");
      if proto.G_program_doctext_candidate != ""  && proto.G_program_doctext_status != proto.ALREADY_PROCESSED{
        proto.G_program.SetDoc(proto.G_program_doctext_candidate)
        proto.G_program_doctext_status = proto.ALREADY_PROCESSED
      }
      proto.Clear_doctext()
    }

CaptureDocText:
    {
      if (proto.G_parse_mode == proto.PROGRAM) {
        $$ = proto.G_doctext
        proto.G_doctext = ""
      } else {
        $$ = ""
      }
    }

/* TODO(dreiss): Try to DestroyDocText in all sorts or random places. */
DestroyDocText:
    {
      if (proto.G_parse_mode == proto.PROGRAM) {
        proto.Clear_doctext()
      }
    }

/* We have to DestroyDocText here, otherwise it catches the doctext
   on the first real element. */
HeaderList:
  HeaderList DestroyDocText Header
    {
      fmt.Println("HeaderList -> HeaderList Header")
    }
|
    {
      fmt.Println("HeaderList -> ")
    }

Header:
  Include
    {
      fmt.Println("Header -> Include");
    }
| tok_namespace tok_identifier tok_identifier TypeAnnotations
    {
      fmt.Printf("Header -> tok_namespace tok_identifier tok_identifier");
      proto.Declare_valid_program_doctext()
      if (proto.G_parse_mode == proto.PROGRAM) {
        proto.G_program.SetNamespaceWithLang($2, $3);
      }
      if ($4 != nil) {
        proto.G_program.SetNamespaceAnnotations($2, $4.GetAnnotations())
      }
    }
| tok_namespace '*' tok_identifier
    {
      fmt.Println("Header -> tok_namespace * tok_identifier")
      proto.Declare_valid_program_doctext()
      if (proto.G_parse_mode == proto.PROGRAM) {
        proto.G_program.SetNamespaceWithLang("*", $3);
      }
    }
| tok_cpp_include tok_literal
    {
      fmt.Println("Header -> tok_cpp_include tok_literal");
      proto.Declare_valid_program_doctext()
      if (proto.G_parse_mode == proto.PROGRAM) {
        proto.G_program.AddCppInclude($2)
      }
    }

Include:
  tok_include tok_literal
    {
      fmt.Println("Include -> tok_include tok_literal");
      proto.Declare_valid_program_doctext()
      if (proto.G_parse_mode == proto.INCLUDES) {
        path := proto.IncludeFile($2);
        if path != "" {
          proto.G_program.AddInclude(path, $2)
        }
      }
    }

DefinitionList:
  DefinitionList CaptureDocText Definition
    {
      fmt.Println("DefintionList -> DefinitionList Definition")
      if $2 != "" && $3 != nil {
        $3.SetDoc($2)
      }
    }
|
    {
      fmt.Println("DefinitionList -> ");
    }

Definition:
  Const
    {
      fmt.Println("Definition -> Const")
      if (proto.G_parse_mode == proto.PROGRAM) {
        proto.G_program.AddConst($1);
      }
      $$ = proto.NewDoc($1.GetDoc())
    }
| TypeDefinition
    {
      fmt.Println("Definition -> TypeDefinition");
      if (proto.G_parse_mode == proto.PROGRAM) {
        proto.G_scope.AddType($1.GetName(), $1)
        if (proto.G_parent_scope != nil) {
          proto.G_parent_scope.AddType(proto.G_parent_prefix + $1.GetName(), $1)
        }
        if !proto.G_program.IsUniqueTypename($1){
          panic(fmt.Sprintf("Type \"%s\" is already defined.", $1.GetName()))
        }
      }
      $$ = proto.NewDoc($1.GetDoc())
    }
| Service
    {
      fmt.Println("Definition -> Service")
      if (proto.G_parse_mode == proto.PROGRAM) {
        proto.G_scope.AddService($1.GetName(), $1)
        if (proto.G_parent_scope != nil) {
          proto.G_parent_scope.AddService(proto.G_parent_prefix + $1.GetName(), $1);
        }
        proto.G_program.AddService($1)
        if !proto.G_program.IsUniqueTypename($1) {
          panic(fmt.Sprintf("Type \"%s\" is already defined.", $1.GetName()))
        }
      }
      $$ = proto.NewDoc($1.GetDoc())
    }

TypeDefinition:
  Typedef
    {
      fmt.Println("TypeDefinition -> Typedef");
      if (proto.G_parse_mode == proto.PROGRAM) {
        proto.G_program.AddTypedef($1)
        $$ = $1
      }
    }
| Enum
    {
      fmt.Println("TypeDefinition -> Enum");
      if (proto.G_parse_mode == proto.PROGRAM) {
        proto.G_program.AddEnum($1)
        $$ = $1
      }
    }
| Senum
    {
      fmt.Println("TypeDefinition -> Senum");
      if proto.G_parse_mode == proto.PROGRAM {
        proto.G_program.AddTypedef($1)
        $$ = $1
      }
    }
| Struct
    {
      fmt.Println("TypeDefinition -> Struct");
      if (proto.G_parse_mode == proto.PROGRAM) {
        proto.G_program.AddStruct($1);
        $$ = $1
      }
    }
| Xception
    {
      fmt.Println("TypeDefinition -> Xception");
      if (proto.G_parse_mode == proto.PROGRAM) {
        proto.G_program.AddXception($1)
        $$ = $1
      }
    }

CommaOrSemicolonOptional:
  ','
    {}
| ';'
    {}
|
    {}

Typedef:
  tok_typedef FieldType tok_identifier TypeAnnotations CommaOrSemicolonOptional
    {
      fmt.Println("TypeDef -> tok_typedef FieldType tok_identifier");
      //proto.Validate_simple_identifier($3);
      td := proto.NewTypedef(proto.G_program, $2, $3, false);
      $$ = td
      if $4 != nil {
        $$.SetAnnotations($4.GetAnnotations())
      }
    }

Enum:
  tok_enum tok_identifier '{' EnumDefList '}' TypeAnnotations
    {
      fmt.Println("Enum -> tok_enum tok_identifier { EnumDefList }")
      $$ = $4
      //proto.Validate_simple_identifier( $2);
      $$.SetName($2)
      if $6 != nil {
        $$.SetAnnotations($6.GetAnnotations())
      }

      // make constants for all the enum values
      if (proto.G_parse_mode == proto.PROGRAM) {
        for _, v := range $$.GetConstants(){
          const_name := $$.GetName() + "." + v.GetName()
          const_val := proto.NewConstValue(v.GetValue())
          const_val.SetEnum($$)
          proto.G_scope.AddConstant(const_name, proto.NewConst(proto.G_type_i32, v.GetName(), const_val))
          if (proto.G_parent_scope != nil) {
            proto.G_parent_scope.AddConstant(proto.G_parent_prefix + const_name, proto.NewConst(proto.G_type_i32, v.GetName(), const_val))
          }
        }
      }
    }

EnumDefList:
  EnumDefList EnumDef
    {
      fmt.Println("EnumDefList -> EnumDefList EnumDef")
      $$ = $1
      $$.Append($2)
    }
|
    {
      fmt.Println("EnumDefList -> ")
      $$ = proto.NewEnum(proto.G_program)
      y_enum_val = -1
    }

EnumDef:
  CaptureDocText EnumValue TypeAnnotations CommaOrSemicolonOptional
    {
      fmt.Println("EnumDef -> EnumValue")
      $$ = $2
      if $1 != "" {
        $$.SetDoc($1)
      }
	  if $3 != nil {
        $$.SetAnnotations($3.GetAnnotations())
      }
    }

EnumValue:
  tok_identifier '=' tok_int_constant
    {
      fmt.Println("EnumValue -> tok_identifier = tok_int_constant")
      y_enum_val = $3
      $$ = proto.NewEnumValue($1, y_enum_val);
    }
 |
  tok_identifier
    {
      fmt.Println("EnumValue -> tok_identifier")
      //proto.Validate_simple_identifier( $1)
      y_enum_val ++
      $$ = proto.NewEnumValue($1, y_enum_val)
    }

Senum:
  tok_senum tok_identifier '{' SenumDefList '}' TypeAnnotations
    {
      fmt.Println("Senum -> tok_senum tok_identifier { SenumDefList }")
      //proto.Validate_simple_identifier( $2);
      $$ = proto.NewTypedef(proto.G_program, $4, $2, false)
      if $6 != nil {
        $$.SetAnnotations($6.GetAnnotations())
      }
    }

SenumDefList:
  SenumDefList SenumDef
    {
      fmt.Println("SenumDefList -> SenumDefList SenumDef")
      $$ = $1
      $$.AddStringEnumVal($2)
    }
|
    {
      fmt.Println("SenumDefList -> ")
      $$ = proto.NewBaseType("string", proto.TYPE_STRING)
      $$.SetStringEnum(true)
    }

SenumDef:
  tok_literal CommaOrSemicolonOptional
    {
      fmt.Println("SenumDef -> tok_literal")
      $$ = $1
    }

Const:
  tok_const FieldType tok_identifier '=' ConstValue CommaOrSemicolonOptional
    {
      fmt.Println("Const -> tok_const FieldType tok_identifier = ConstValue")
      if (proto.G_parse_mode == proto.PROGRAM) {
        //proto.Validate_simple_identifier($3)
        proto.G_scope.ResolveConstValue($5, $2)
        $$ = proto.NewConst($2, $3, $5)
        //proto.Validate_const_type($$)

        proto.G_scope.AddConstant($3, $$)
        if (proto.G_parent_scope != nil) {
          proto.G_parent_scope.AddConstant(proto.G_parent_prefix + $3, $$)
        }
      } else {
        $$ = nil
      }
    }

ConstValue:
  tok_int_constant
    {
      fmt.Println("ConstValue => tok_int_constant");
      $$ = proto.NewConstValue(nil)
      $$.SetInteger($1)
    }
| tok_dub_constant
    {
      fmt.Println("ConstValue => tok_dub_constant");
      $$ = proto.NewConstValue(nil)
      $$.SetDouble($1)
    }
| tok_literal
    {
      fmt.Println("ConstValue => tok_literal");
      $$ = proto.NewConstValue($1);
    }
| tok_identifier
    {
      fmt.Println("ConstValue => tok_identifier");
      $$ = proto.NewConstValue(nil)
      $$.SetIdentifier($1)
    }
| ConstList
    {
      fmt.Println("ConstValue => ConstList")
      $$ = $1;
    }
| ConstMap
    {
      fmt.Println("ConstValue => ConstMap")
      $$ = $1;
    }

ConstList:
  '[' ConstListContents ']'
    {
      fmt.Println("ConstList => [ ConstListContents ]")
      $$ = $2;
    }

ConstListContents:
  ConstListContents ConstValue CommaOrSemicolonOptional
    {
      fmt.Println("ConstListContents => ConstListContents ConstValue CommaOrSemicolonOptional");
      $$ = $1;
      $$.AddList($2);
    }
|
    {
      fmt.Println("ConstListContents =>");
      $$ = proto.NewConstValue(nil)
      $$.SetList()
    }

ConstMap:
  '{' ConstMapContents '}'
    {
      fmt.Println("ConstMap => { ConstMapContents }");
      $$ = $2;
    }

ConstMapContents:
  ConstMapContents ConstValue ':' ConstValue CommaOrSemicolonOptional
    {
      fmt.Println("ConstMapContents => ConstMapContents ConstValue CommaOrSemicolonOptional")
      $$ = $1
      $$.AddMap($2, $4);
    }
|
    {
      fmt.Println("ConstMapContents =>");
      $$ = proto.NewConstValue(nil)
      $$.SetMap();
    }

StructHead:
  tok_struct
    {
      $$ = struct_is_struct
    }
| tok_union
    {
      $$ = struct_is_union
    }

Struct:
  StructHead tok_identifier XsdAll '{' FieldList '}' TypeAnnotations
    {
      fmt.Println("Struct -> tok_struct tok_identifier { FieldList }");
      //proto.validate_simple_identifier( $2);
      $5.SetXsdAll($3)
      $5.SetUnion($1 == struct_is_union)
      $$ = $5
      $$.SetName($2);
      if ($7 != nil) {
        $$.SetAnnotations($7.GetAnnotations())
      }
    }

XsdAll:
  tok_xsd_all
    {
      $$ = true;
    }
|
    {
      $$ = false;
    }

XsdOptional:
  tok_xsd_optional
    {
      $$ = true
    }
|
    {
      $$ = false
    }

XsdNillable:
  tok_xsd_nillable
    {
      $$ = true
    }
|
    {
      $$ = false
    }

XsdAttributes:
  tok_xsd_attrs '{' FieldList '}'
    {
      $$ = $3
    }
|
    {
      $$ = nil
    }

Xception:
  tok_xception tok_identifier '{' FieldList '}' TypeAnnotations
    {
      fmt.Println("Xception -> tok_xception tok_identifier { FieldList }")
      //proto.Validate_simple_identifier( $2)
      $4.SetName($2)
      $4.SetXception(true)
      $$ = $4
      if $6 != nil {
        $$.SetAnnotations($6.GetAnnotations())
      }
    }

Service:
  tok_service tok_identifier Extends '{' FlagArgs FunctionList UnflagArgs '}' TypeAnnotations
    {
      fmt.Println("Service -> tok_service tok_identifier { FunctionList }")
      //proto.Validate_simple_identifier( $2);
      $$ = $6;
      $$.SetName($2);
      $$.SetExtends($3);
      if ($9 != nil) {
        $$.SetAnnotations($9.GetAnnotations())
      }
    }

FlagArgs:
    {
       G_arglist = 1;
    }

UnflagArgs:
    {
       G_arglist = 0;
    }

Extends:
  tok_extends tok_identifier
    {
      fmt.Println("Extends -> tok_extends tok_identifier")
      $$ = nil
      if (proto.G_parse_mode == proto.PROGRAM) {
        $$ = proto.G_scope.GetService($2)
        if ($$ == nil) {
          panic(fmt.Sprintf("Service \"%s\" has not been defined.", $2))
        }
      }
    }
|
    {
      $$ = nil
    }

FunctionList:
  FunctionList Function
    {
      fmt.Println("FunctionList -> FunctionList Function")
      $$ = $1;
      $1.AddFunction($2)
    }
|
    {
      fmt.Println("FunctionList -> ");
      $$ = proto.NewService(proto.G_program);
    }

Function:
  CaptureDocText Oneway FunctionType tok_identifier '(' FieldList ')' Throws TypeAnnotations CommaOrSemicolonOptional
    {
      //proto.Validate_simple_identifier( $4);
      fmt.Printf("Function -> FunctionType")
      $6.SetName($4 + "_args")
      $$ = proto.NewFunctionWithXception($3, $4, $6, $8, $2)
      if ($1 != "") {
        $$.SetDoc($1)
      }
      if ($9 != nil) {
        $$.SetAnnotations($9.GetAnnotations())
      }
    }

Oneway:
  tok_oneway
    {
      $$ = true
    }
|
    {
      $$ = false
    }

Throws:
  tok_throws '(' FieldList ')'
    {
      fmt.Println("Throws -> tok_throws ( FieldList )");
      $$ = $3;
      if proto.G_parse_mode == proto.PROGRAM {
        panic("Throws clause may not contain non-exception types")
      }
    }
|
    {
      $$ = proto.NewStruct(proto.G_program)
    }

FieldList:
  FieldList Field
    {
      fmt.Println("FieldList -> FieldList , Field");
      $$ = $1;
      if (!($$.Append($2))) {
        panic(fmt.Sprintf("\"%d: %s\" - field identifier/name has already been used", $2.GetKey(), $2.GetName()))
      }
    }
|
    {
      fmt.Printf("FieldList -> ");
      y_field_val = -1;
      $$ = proto.NewStruct(proto.G_program);
    }

Field:
  CaptureDocText FieldIdentifier FieldRequiredness FieldType FieldReference tok_identifier FieldValue XsdOptional XsdNillable XsdAttributes TypeAnnotations CommaOrSemicolonOptional
    {
      fmt.Println("tok_int_constant : Field -> FieldType tok_identifier")
      if ($2.AutoAssigned) {
        panic(fmt.Sprintf("No field key specified for %s, resulting protocol may have conflicts or not be backwards compatible!\n", $6))
        if (g_strict >= 192) {
          panic("Implicit field keys are deprecated and not allowed with -strict")
        }
      }
      //proto.Validate_simple_identifier($6);
      $$ = proto.NewField($4, $6, int32($2.Value))
      $$.SetReference($5);
      $$.SetReq($3);
      if ($7 != nil) {
        proto.G_scope.ResolveConstValue($7, $4)
        proto.Validate_field_value($$, $7);
        $$.SetValue($7)
      }
      $$.SetXsdOptional($8)
      $$.SetXsdNillable($9)
      if ($1 != "") {
        $$.SetDoc($1)
      }
      if ($10 != nil) {
        $$.SetXsdAttrs($10)
      }
      if ($11 != nil) {
        $$.SetAnnotations($11.GetAnnotations())
      }
    }

FieldIdentifier:
  tok_int_constant ':'
    {
      if ($1 <= 0) {
          fmt.Printf("Nonpositive value (%d) not allowed as a field key.\n", $1)
          $$ = proto.NewFieldId(65535, true)
      } else {
          $$ = proto.NewFieldId(int($1), false)
      }
      if( (SHRT_MIN > $$.Value) || ($$.Value > SHRT_MAX)) {
          fmt.Println("Field key (%d) exceeds allowed range (%d..%d).\n", $$.Value, SHRT_MIN, SHRT_MAX);
      }
    }
|
    {
      $$ = proto.NewFieldId(65534, true)
      if( (SHRT_MIN > $$.Value) || ($$.Value > SHRT_MAX)) {
        fmt.Printf("Field key (%d) exceeds allowed range (%d..%d).\n", $$.Value, SHRT_MIN, SHRT_MAX)
      }
    }

FieldReference:
  tok_reference
    {
      $$ = true
    }
|
   {
     $$ = false
   }

FieldRequiredness:
  tok_required
    {
      $$ = proto.T_REQUIRED;
    }
| tok_optional
    {
      if G_arglist > 0 {
        if (proto.G_parse_mode == proto.PROGRAM) {
          fmt.Println("optional keyword is ignored in argument lists.\n");
        }
        $$ = proto.T_OPT_IN_REQ_OUT
      } else {
        $$ = proto.T_OPTIONAL
      }
    }
|
    {
      $$ = proto.T_OPT_IN_REQ_OUT
    }

FieldValue:
  '=' ConstValue
    {
      if (proto.G_parse_mode == proto.PROGRAM) {
        $$ = $2
      } else {
        $$ = nil
      }
    }
|
    {
      $$ = nil
    }

FunctionType:
  FieldType
    {
      fmt.Println("FunctionType -> FieldType");
      $$ = $1;
    }
| tok_void
    {
      fmt.Println("FunctionType -> tok_void");
      $$ = proto.G_type_void;
    }

FieldType:
  tok_identifier
    {
      fmt.Println("FieldType -> tok_identifier");
      if (proto.G_parse_mode == proto.INCLUDES) {
        // Ignore identifiers in include mode
        $$ = nil
      } else {
        // Lookup the identifier in the current scope
        $$ = proto.G_scope.GetType($1)
        if ($$ == nil) {
          $$ = proto.NewTypedef(proto.G_program, nil, $1, true)
        }
      }
    }
| BaseType
    {
      fmt.Println("FieldType -> BaseType")
      $$ = $1;
    }
| ContainerType
    {
      fmt.Println("FieldType -> ContainerType")
      $$ = $1
    }

BaseType: SimpleBaseType TypeAnnotations
    {
      fmt.Println("BaseType -> SimpleBaseType TypeAnnotations");
      if ($2 != nil) {
        $$ = proto.NewBaseType($1.GetName(), $1.GetBase())
        $$.SetAnnotations($2.GetAnnotations())
      } else {
        $$ = $1
      }
    }

SimpleBaseType:
  tok_string
    {
      fmt.Println("BaseType -> tok_string");
      $$ = proto.G_type_string;
    }
| tok_binary
    {
      fmt.Println("BaseType -> tok_binary");
      $$ = proto.G_type_binary;
    }
| tok_slist
    {
      fmt.Println("BaseType -> tok_slist");
      $$ = proto.G_type_slist;
    }
| tok_bool
    {
      fmt.Println("BaseType -> tok_bool");
      $$ = proto.G_type_bool;
    }
| tok_i8
    {
      fmt.Println("BaseType -> tok_i8");
      $$ = proto.G_type_i8;
    }
| tok_i16
    {
      fmt.Println("BaseType -> tok_i16");
      $$ = proto.G_type_i16;
    }
| tok_i32
    {
      fmt.Println("BaseType -> tok_i32");
      $$ = proto.G_type_i32;
    }
| tok_i64
    {
      fmt.Println("BaseType -> tok_i64");
      $$ = proto.G_type_i64;
    }
| tok_double
    {
      fmt.Println("BaseType -> tok_double");
      $$ = proto.G_type_double;
    }

ContainerType: SimpleContainerType TypeAnnotations
    {
      fmt.Println("ContainerType -> SimpleContainerType TypeAnnotations")
      $$ = $1
      if $2 != nil {
        $$.SetAnnotations($2.GetAnnotations())
      }
    }

SimpleContainerType:
  MapType
    {
      fmt.Println("SimpleContainerType -> MapType")
      $$ = $1;
    }
| SetType
    {
      fmt.Println("SimpleContainerType -> SetType")
      $$ = $1;
    }
| ListType
    {
      fmt.Println("SimpleContainerType -> ListType")
      $$ = $1
    }

MapType:
  tok_map CppType '<' FieldType ',' FieldType '>'
    {
      fmt.Println("MapType -> tok_map <FieldType, FieldType>")
      $$ = proto.NewMap($4, $6)
      if ($2 != "") {
        $$.SetCppName($2)
      }
    }

SetType:
  tok_set CppType '<' FieldType '>'
    {
      fmt.Println("SetType -> tok_set<FieldType>")
      $$ = proto.NewSet($4)
      if ($2 != "") {
        $$.SetCppName($2)
      }
    }

ListType:
  tok_list '<' FieldType '>' CppType
    {
      fmt.Println("ListType -> tok_list<FieldType>");
      proto.Check_for_list_of_bytes($3);
      $$ = proto.NewList($3)
      if ($5 != "") {
        $$.SetCppName($5)
      }
    }

CppType:
  tok_cpp_type tok_literal
    {
      $$ = $2
    }
|
    {
      $$ = ""
    }

TypeAnnotations:
  '(' TypeAnnotationList ')'
    {
      fmt.Println("TypeAnnotations -> ( TypeAnnotationList )")
      $$ = $2
    }
|
    {
      $$ = nil
    }

TypeAnnotationList:
  TypeAnnotationList TypeAnnotation
    {
      fmt.Println("TypeAnnotationList -> TypeAnnotationList , TypeAnnotation")
      $$ = $1
      $$.SetAnnotation($2.Key, $2.Value)
    }
|
    {
      $$ = proto.NewStruct(proto.G_program)
    }

TypeAnnotation:
  tok_identifier TypeAnnotationValue CommaOrSemicolonOptional
    {
      fmt.Printf("TypeAnnotation -> TypeAnnotationValue");
      $$ = proto.NewAnnotation($1, $2)
    }

TypeAnnotationValue:
  '=' tok_literal
    {
      fmt.Printf("TypeAnnotationValue -> = tok_literal");
      $$ = $2;
    }
|
    {
      fmt.Printf("TypeAnnotationValue ->");
      $$ = "1";
    }

%%