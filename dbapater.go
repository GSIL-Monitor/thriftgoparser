/*
	本文件负责把所有类型从数据库中拉出来
*/
package main

import (
	//_ "code.byted.org/gopkg/mysql-driver"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	//"code.byted.org/gopkg/gorm"
)

var (
	db_ *sql.DB
)

func InitMysql(usr, passwd, host string, port uint16) {
	s := fmt.Sprintf("%s:%s@tcp(%s:%d)/apigate?charset=utf8&parseTime=True", usr, passwd, host, port)
	fmt.Println(s)
	var err error
	db_, err = sql.Open("mysql", s)
	if err != nil {
		fmt.Printf("InitMysql: connection db failed, error:%v\n", err)
		panic("InitMysql: connection db failed!")
	}

	db_.SetMaxOpenConns(200)
	db_.SetMaxIdleConns(200)
}

// 定义一个协议文件
type IDLFile struct{
	PSM string
	Name string
	Content string
}

// 定义一个接口
type IDLFunction struct{
	Module   string
	FuncName string     // ServiceName, 带namespace和psm psm:namespace.name
	ReqType1 string // 带Namespace
	ReqType2 string // 带Namespace
	ReqType3 string // 带Namespace
	ReqType4 string // 带Namespace
	ReqType5 string // 带Namespace
	ReqType6 string // 带Namespace
	ReqType7 string // 带Namespace
	ReqType8 string // 带Namespace
	ReqType9 string // 带Namespace
	RspType string // 带Namespace
	UpdateTime uint32
}

type IDLMember struct{
	TypeName string   // 带PSM-Namespace
	FieldTag int32   // 字段的Tag
	FieldName string  // 字段的默认值
	FieldType string    // 基本类型，MAP, Struct, List等
	FieldTypeId string    // 带PSM-Namespace
	FieldDefault string // 字段的默认值
	FieldKeyType string // 如果是Map
	FieldElemType string // 如果是Map,Set, List使用该字段做ValueType
}

type IDLType struct{
	Name string
	TagMembers map[int32] *IDLMember  `gorm:"-"`
	NameMembers map[string] *IDLMember `gorm:"-"`
}

func (it *IDLType)GetName() string{
	return it.Name
}

func (it *IDLType)SetName(name string){
	it.Name = name
}

func (it *IDLType)IsBase() bool{
	return true
}

func (it *IDLType)IsMap() bool{
	return it.Name == "MAP"
}

func (it *IDLType)IsList() bool{
	return it.Name == "LIST"
}

func (it *IDLType)IsSet() bool{
	return it.Name == "SET"
}

func (it *IDLType)IsStruct() bool{
	return it.Name == "STRUCT"
}

func (it *IDLType)GetKey() *IDLMember{
	if it.IsMap() {
		return it.TagMembers[1]
	}

	return nil
}

func (it *IDLType)GetValue() *IDLMember{
	if it.IsMap() || it.IsSet() || it.IsList()  {
		return it.TagMembers[0]
	}

	return nil
}

func LoadFunctionsFromDb(idlfuncs map[string]*IDLFunction) error{
	fmt.Println("LoadFunctions")
	rows, err := db_.Query("select module_name, func_name, reqtype1,reqtype2,reqtype3,reqtype4,reqtype5,reqtype6,reqtype7,reqtype8,reqtype9, rsptype from idl_functions")
	if err != nil{
		return fmt.Errorf("db query error:%v", err)
	}
	defer rows.Close()

	for rows.Next() {
		fmt.Println("LoadServices rows.Next()")
		var idlFunc IDLFunction
		rows.Scan(&idlFunc.Module, &idlFunc.FuncName,
			      &idlFunc.ReqType1, &idlFunc.ReqType2, &idlFunc.ReqType3,
				  &idlFunc.ReqType4, &idlFunc.ReqType5, &idlFunc.ReqType6,
				  &idlFunc.ReqType7, &idlFunc.ReqType8, &idlFunc.ReqType9,
				  &idlFunc.RspType)

		idlfuncs[idlFunc.FuncName] = &idlFunc
		fmt.Printf("[AddIdlFunc]:%v\n", idlFunc)
	}

	return nil
}

func LoadTypesFromDb(idltypes map[string] *IDLType) error {
	fmt.Println(db_)
	rows, err := db_.Query("select m.typename, m.field_tag, m.field_name, m.field_type, m.field_default from idl_types t, idl_type_members m where t.name = m.typename order by t.name")
	if err != nil{
		return fmt.Errorf("query db error:%v", err)
	}
	defer rows.Close()

	for rows.Next() {
		fmt.Println("LoadTypes rows.Next()")

		var idlMember IDLMember
		rows.Scan(&idlMember.TypeName, &idlMember.FieldTag, &idlMember.FieldName,  &idlMember.FieldType, &idlMember.FieldDefault)

		idlType, ok := idltypes[idlMember.TypeName]
		if !ok { // 如果没有，说明是个新类型
			idlType = new(IDLType)
			idlType.SetName(idlMember.TypeName)
			idlType.TagMembers = make(map[int32]* IDLMember)
			idlType.NameMembers = make(map[string]* IDLMember)

			fmt.Printf("[AddNewType]%s\n", idlType.GetName())
			idltypes[idlType.GetName()] = idlType
			idlType, ok = idltypes[idlType.GetName()]
		}

		_, ok2 := idlType.NameMembers[idlMember.FieldName]
		if (ok2){
			return fmt.Errorf("field name:%s already exist", idlMember.FieldName)
		}

		_, ok3 := idlType.TagMembers[idlMember.FieldTag]
		if ok3 {
			return fmt.Errorf("field tag:%s already exist", idlMember.FieldTag)
		}

		idlType.NameMembers[idlMember.FieldName] = &idlMember
		idlType.TagMembers[idlMember.FieldTag] = &idlMember
		fmt.Printf("[AddNewMember]%v\n", idlMember)
	}

	return nil
}

func InsertIDLFile(file *IDLFile) error{
	tx, err := db_.Begin()
	if err != nil {
		return fmt.Errorf("begin transcation fail:%v", err)
	}

	rs, err := tx.Exec("REPLACE INTO idl_files(psm, file, content)values('%s','%s','%s')VALUES ('%s','%s','%s')", file.PSM, file.Name, file.Content)
	if err != nil {
		return fmt.Errorf("exec fail:%v", err)
	}

	rowAffected, err := rs.RowsAffected()
	if err != nil {
		return fmt.Errorf("RowsAffected fail:%v", err)
	}
	fmt.Printf("rows affected:%d\n", rowAffected)

	return nil
}

func InsertIDLFunction(function *IDLFunction) error{
	tx, err := db_.Begin()
	if err != nil {
		return fmt.Errorf("begin transcation fail:%v", err)
	}

	rs, err := tx.Exec("INSERT INTO idl_functions(module_name, func_name, reqtype1, rsptype)values('%s','%s','%s','%s')", function.Module, function.FuncName, function.ReqType1, function.RspType)
	if err != nil {
		return fmt.Errorf("exec fail:%v", err)
	}

	rowAffected, err := rs.RowsAffected()
	if err != nil {
		return fmt.Errorf("RowsAffected fail:%v", err)
	}
	fmt.Printf("rows affected:%d\n", rowAffected)

	return nil
}

func InsertIdlType(idlType *IDLType) error{
	/*tx, err := db_.Begin()
	if err != nil {
		return fmt.Errorf("begin transcation fail:%v", err)
	}*/

	fmt.Printf("[CreateType]:%s", idlType.GetName())

	rs, err := db_.Exec("INSERT IGNORE INTO idl_types(name)VALUES(?)", idlType.GetName())
	if err != nil {
		return fmt.Errorf("exec fail:%v", err)
	}

	fmt.Printf("[CreateType]result:%v\n", rs)

	rowAffected, err := rs.RowsAffected()
	if err != nil {
		return fmt.Errorf("RowsAffected fail:%v", err)
	}
	fmt.Printf("rows affected:%d\n", rowAffected)

	if 0 == rowAffected{
		fmt.Printf("type:[%s] already exist\n", idlType.GetName())
		return nil
	}

	if len(idlType.TagMembers) > 0{
		sql := "INSERT IGNORE INTO idl_type_members(typename, field_tag, field_name, field_type, field_default)VALUES"
		for _, v := range idlType.TagMembers{
			sql += fmt.Sprintf("('%s',%d,'%s','%s','%s'),", v.TypeName, v.FieldTag, v.FieldName, v.FieldType, v.FieldDefault)
		}

		sql = sql[:len(sql)-1]
		fmt.Printf("run sql:%s\n", sql)

		rs, err := db_.Exec(sql)
		if err != nil {
			return fmt.Errorf("exec fail:%v", err)
		}

		rowAffected, err = rs.RowsAffected()
		if err != nil {
			return fmt.Errorf("tx commit fail:%v", err)
		}
		fmt.Println(rowAffected)
	}

	//if err := tx.Commit(); err != nil {
	//	return fmt.Errorf("tx commit fail:%v", err)
	//}

	return nil
}