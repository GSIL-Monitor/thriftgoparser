package main

import (
	"time"
	"fmt"
	"github.com/AfLnk/thriftgoparser/proto"
)

var(
	G_TypeCache map[string] *IDLType
	G_FuncCache map[string] *IDLFunction
)

func load_types(){
	newTypeCache := make(map[string] *IDLType)
	err := LoadTypesFromDb(newTypeCache)
	if err != nil{
		fmt.Printf("LoadTypes failed%s\n", err.Error())
		return
	}

	if len(newTypeCache) > 0 {
		G_TypeCache = newTypeCache
	}
}

func load_functions() {
	newFuncCache := make(map[string]*IDLFunction)
	err := LoadFunctionsFromDb(newFuncCache)
	if err != nil {
		fmt.Printf("LoadTypes failed:%s\n", err.Error())
		return
	}

	if len(newFuncCache) > 0 {
		G_FuncCache = newFuncCache
	}
}

func InitLoader(db_usr, db_passwd, db_host string, db_port uint16){
	InitMysql(db_usr, db_passwd, db_host, db_port)
	go func(){
		load_types()
		load_functions()
		for {
			select {
			case <- time.After(300 * time.Second): // 5分钟同步一次
				load_types()
				load_functions()
			default:
				fmt.Println("[IDLE]................")
				time.Sleep(1*time.Second)
			}
		}
	}()
}


func GetIDLType(typename string) (*IDLType, error){
	// 先判断是否基础类型
	fmt.Printf("type cache:%v\n", G_TypeCache)
	result, ok := G_TypeCache[typename]
	if ok{
		return result, nil
	}

	return nil, fmt.Errorf("Type not exist")
}

func get_function(psm string, servicename string) (*IDLFunction, error){
	fmt.Printf("nservice cache:%v\n", G_FuncCache)
	real_service_name := psm + servicename
	result, ok := G_FuncCache[real_service_name]
	if ok{
		return result, nil
	}

	return nil, fmt.Errorf("service not exist")
}

// 查询请求包的节点信息
// 节点的Key为name
func GetReqType(psm string, service string) (*IDLType, error) {
	idlService, err := get_function(psm, service)
	if err != nil{
		return nil, fmt.Errorf("GetService fail.")
	}

	idlType, err2 := GetIDLType(idlService.ReqType1)
	if err2 != nil{
		return nil, fmt.Errorf("GetType fail")
	}

	return idlType, nil
}

// 查询请求包的节点信息
// 节点的Key为name
func GetReqTypes(psm string, service string) (map[int16]*IDLType, error) {
	idlService, err := get_function(psm, service)
	if err != nil{
		return nil, fmt.Errorf("GetService fail.")
	}

	ret := make(map[int16]*IDLType)

	idlType1, _ := GetIDLType(idlService.ReqType1)
	if idlType1 != nil{
		ret[1] = idlType1
	}

	idlType2, _ := GetIDLType(idlService.ReqType2)
	if idlType2 != nil{
		ret[2] = idlType2
	}

	idlType3, _ := GetIDLType(idlService.ReqType3)
	if idlType3 != nil{
		ret[3] = idlType3
	}

	idlType4, _ := GetIDLType(idlService.ReqType4)
	if idlType4 != nil{
		ret[4] = idlType4
	}

	idlType5, _ := GetIDLType(idlService.ReqType5)
	if idlType5 != nil{
		ret[5] = idlType5
	}

	idlType6, _ := GetIDLType(idlService.ReqType6)
	if idlType6 != nil{
		ret[6] = idlType6
	}

	idlType7, _ := GetIDLType(idlService.ReqType7)
	if idlType7 != nil{
		ret[7] = idlType7
	}

	idlType8, _ := GetIDLType(idlService.ReqType8)
	if idlType8 != nil{
		ret[8] = idlType8
	}

	idlType9, _ := GetIDLType(idlService.ReqType9)
	if idlType9 != nil{
		ret[9] = idlType9
	}

	return ret, nil
}

// 查询相应包中的节点信息
// 节点的key为tag
func GetRspType(psm string, service string) (*IDLType, error) {
	idlFunc, err := get_function(psm, service)
	if err != nil{
		return nil, fmt.Errorf("GetService fail.")
	}

	idlType, err2 := GetIDLType(idlFunc.RspType)
	if err2 != nil{
		return nil, fmt.Errorf("GetType fail")
	}

	return idlType, nil
}

func GetStructTypeId(ts *proto.TStruct) string{
	return ts.GetProgram().PSM_ + ":" + ts.GetProgram().GetNamespace() + "." + ts.GetName()
}

func CreateNewList(tlist *proto.TList) (string, error){
	val_type, val_id, err := GenerateTypeAndId(tlist.GetElemType())
	if val_type == "" || err != nil {
		return "", fmt.Errorf("fail to generate type&id for:%s, error:%v", tlist.GetElemType().GetName(), err)
	}

	newName := fmt.Sprintf("LIST:<%s>", val_id)

	newList := new(IDLType)
	newList.TagMembers = make(map[int32]*IDLMember)

	elem := new(IDLMember)
	elem.TypeName = newName
	elem.FieldTag = 0
	elem.FieldName = "VALUE"
	elem.FieldType = val_type
	elem.FieldTypeId = val_id
	newList.TagMembers[0] = elem

	InsertIdlType(newList)

	return newName, nil
}

func CreateNewMap(tmap *proto.TMap) (string, error){
	if !tmap.GetKeyType().IsBaseType(){
		return "", fmt.Errorf("Key type must be base type, Actual:%", tmap.GetKeyType())
	}
	keyType := tmap.GetKeyType().GetName()

	var err error
	val_type, val_id, err := GenerateTypeAndId(tmap.GetValType())
	if val_type == "" || err != nil {
		return "", fmt.Errorf("fail to generate type&id for:%s, error:%v", tmap.GetValType().GetName(), err)
	}

	newName := fmt.Sprintf("MAP:<%s, %s>", keyType, val_id)
	newMap := new(IDLType)
	newMap.TagMembers = make(map[int32]*IDLMember)

	newMap.SetName(newName)

	key := new(IDLMember)
	key.TypeName = newName
	key.FieldTag = 1
	key.FieldName = "KEY"
	key.FieldType = keyType
	key.FieldTypeId = keyType
	newMap.TagMembers[key.FieldTag] = key

	val := new(IDLMember)
	val.TypeName = newName
	val.FieldTag = 0
	val.FieldName = "VALUE"
	val.FieldType = val_type
	val.FieldTypeId = val_id
	newMap.TagMembers[val.FieldTag] = val

	InsertIdlType(newMap)

	return newName, nil
}

func CreateNewSet(tset *proto.TSet) (string, error){
	val_type, val_id, err := GenerateTypeAndId(tset.GetElemType())
	if val_type == "" || err != nil {
		return "", fmt.Errorf("fail to generate type&id for:%s, error:%v", tset.GetElemType().GetName(), err)
	}

	newName := fmt.Sprintf("SET:<%s>", val_id)

	newSet := new(IDLType)
	newSet.TagMembers = make(map[int32]*IDLMember)

	elem := new(IDLMember)
	elem.TypeName = newName
	elem.FieldTag = 0
	elem.FieldName = "VALUE"
	elem.FieldType = val_type
	elem.FieldTypeId = val_id
	newSet.TagMembers[0] = elem

	InsertIdlType(newSet)

	return newName, nil
}

func GenerateTypeAndId(caftype proto.ICafType) (string, string, error){
	ret_type := ""
	ret_id  := ""
	var err error
	if caftype.IsBaseType(){ // 基本类型
		ret_type = caftype.GetName()
		ret_id = caftype.GetName()
	}else if caftype.IsStruct() {
		ret_type = "STRUCT"
		ret_id = GetStructTypeId(caftype.(*proto.TStruct))
	}else if caftype.IsMap(){
		ret_id, err = CreateNewMap(caftype.(*proto.TMap))
		ret_type = "MAP"
	} else if caftype.IsSet(){
		ret_type = "SET"
		ret_id, err = CreateNewSet(caftype.(*proto.TSet))
	}else if caftype.IsList(){
		ret_id, err = CreateNewList(caftype.(*proto.TList))
		ret_type = "LIST"
	}

	fmt.Printf("[GenerateTypeAndId]type:%s, id:%s, err:%v\n", ret_type, ret_id, err)

	if err == nil {
		return ret_type, ret_id, nil
	}

	return "","", fmt.Errorf("fail to get Id:%v", err)
}

func InsertOrUpdateStruct(ts *proto.TStruct){
	fmt.Printf("[Struct]name:%s, number_count:%d\n", ts.GetName(), len(ts.GetMembers()))
	dbIDLType := new(IDLType)
	dbIDLType.TagMembers = make(map[int32]*IDLMember)
	dbIDLType.SetName(GetStructTypeId(ts))

	for _, v := range ts.GetMembers() {
		fmt.Printf("[Member]%d %s %s\n", v.GetKey(), v.GetName(), v.GetType().GetName())
		idlMember :=  new(IDLMember)
		idlMember.TypeName = dbIDLType.GetName()
		idlMember.FieldTag = v.GetKey()
		idlMember.FieldName = v.GetName()
		var err error
		idlMember.FieldType, idlMember.FieldTypeId, err = GenerateTypeAndId(v.GetType())
		if nil != err{
			fmt.Printf("[InvalidMember]struct:%s, member:%s(%d), error:%v", ts.GetName(), v.GetName(), v.GetKey(), err.Error())
			continue
		}

		dbIDLType.TagMembers[idlMember.FieldTag] = idlMember
	}

	err := InsertIdlType(dbIDLType)
	fmt.Println(err)
}