package proto

import (
	"fmt"
	"crypto/md5"
)

/**
 * Placeholder struct for returning the key and value of an annotation
 * during parsing.
 */
type TAnnotation struct  {
	Key string
	Value string
};

func NewAnnotation(k, v string) *TAnnotation{
	ret := new(TAnnotation)
	ret.Key = k
	ret.Value = v
	return ret
}

//
// 每一个类型的基类
//
type TTypeDesc struct{
	TDoc
	Program_* TProgram
	Name_ string
	Fingerprint_ []uint8 // 长度固定为16
	Annotations_ map[string] string
}

/**
 * Generic representation of a thrift type. These objects are used by the
 * parser module to build up a tree of object that are all explicitly typed.
 * The generic t_type class exports a variety of useful methods that are
 * used by the code generator to branch based upon different handling for the
 * various types.
 *
 */
func (t*TTypeDesc)IsVoid() bool{
	return false
}

func (t*TTypeDesc)IsBaseType() bool{
	return false
}

func (t*TTypeDesc)IsString() bool{
	return false
}

func (t*TTypeDesc)IsBool() bool{
	return false
}

func (t*TTypeDesc)IsTypedef() bool{
	return false
}

func (t*TTypeDesc)IsEnum() bool{
	return false
}

func (t*TTypeDesc)IsStruct() bool{
	return false
}

func (t*TTypeDesc)IsXception() bool{
	return false
}

func (t*TTypeDesc)IsContainer() bool{
	return false
}

func (t*TTypeDesc)IsList() bool{
	return false
}

func (t*TTypeDesc)IsSet() bool{
	return false
}

func (t*TTypeDesc)IsMap() bool{
	return false
}

func (t*TTypeDesc)IsService() bool{
	return false
}

func (t *TTypeDesc)GetProgram() (*TProgram){
	return t.Program_
}

// Return a string that uniquely identifies this type
// from any other thrift type in the world, as far as
// TDenseProtocol is concerned.
// We don't cache this, which is a little sloppy,
// but the compiler is so fast that it doesn't really matter.
//virtual std::string get_fingerprint_material() const = 0;
func (t *TTypeDesc)GetFingerPrintMaterial() string{
	return fmt.Sprintf("invalid implatmentation")
}

func (t *TTypeDesc)HasFingerprint() bool {
	for _, v := range t.Fingerprint_{
		if v != 0 {
			return true
		}
	}

	return false;
}

func Byte2Hex(b byte) string{
	var rv string;
	rv += string(Nybble2Xdigit(b >> 4));
	rv += string(Nybble2Xdigit(b & 0x0f));
	return rv;
}

// This function will break (maybe badly) unless 0 <= num <= 16.
func Nybble2Xdigit(num uint8) byte{
	if (num < 10) {
		return (byte)('0' + num);
	} else {
		return (byte)('A' + num - 10);
	}
}

func (t*TTypeDesc) GetName() string {
	return t.Name_
}


func (t*TTypeDesc) SetName(name string){
	t.Name_ = name
}

func (t*TTypeDesc)GetFingerprint() []uint8{
	return t.Fingerprint_
}

func (t*TTypeDesc)GenerateFingerprint(){
	fmt.Printf("generating fingerprint for %s\n", t.GetName());
	material := t.GetFingerprintMaterial();
	h := md5.New()
	h.Write([]byte(material))
}

func (t*TTypeDesc)GetFingerprintMaterial() string{
	return "ddd"
}

func (t*TTypeDesc)GetAnnotations() map[string]string{
	return t.GetAnnotations()
}

func (t*TTypeDesc)SetAnnotations(annotation map[string]string) {
	t.Annotations_ = annotation
}

func (t*TTypeDesc)SetAnnotation(key, val string){
	t.Annotations_[key] = val
}

//
// 每一个类型的接口
//
type ICafType interface {
	GetDoc() string
	HasDoc() bool
	GetName() string

	GetProgram() (*TProgram)
	GetFingerprintMaterial() string
	GenerateFingerprint()
	HasFingerprint() bool
	GetFingerprint() []uint8
	GetAnnotations() map[string]string
    SetAnnotations(annotation map[string]string)
	SetAnnotation(key, val string)

	// 类型定义
	IsVoid() bool
	IsBaseType() bool
	IsString() bool
	IsBool() bool
	IsTypedef() bool
	IsEnum() bool
	IsStruct() bool
	IsXception() bool
	IsContainer() bool
	IsList() bool
	IsSet() bool
	IsMap() bool
	IsService() bool
}

func GenerateBaseFingerprint(cafType ICafType) {
	if (!cafType.HasFingerprint()) {
		fmt.Printf("generating fingerprint for %s\n", cafType.GetName());
		material := cafType.GetFingerprintMaterial();
		h := md5.New()
		h.Write([]byte(material))
		h.Sum(cafType.GetFingerprint())
	}
}

func GetBinaryFingerprint(cafType ICafType) ([]uint8){
	if (cafType.HasFingerprint()) { // lazy fingerprint generation, right now only used with the c++ generator
		cafType.GenerateFingerprint()
	}

	return cafType.GetFingerprint()
}

func GetAsciiFingerprint(cafType ICafType) string{
	if (cafType.HasFingerprint()) { // lazy fingerprint generation, right now only used with the c++ generator
		cafType.GenerateFingerprint()
	}

	var rv string
	for _, v := range cafType.GetFingerprint() {
		rv += Byte2Hex(v)
	}

	return rv
}

/*
func
t_type* t_type::get_true_type() {
t_type* type = this;
while (type->is_typedef()) {
type = ((t_typedef*)type)->get_type();
}
return type;
}
*/