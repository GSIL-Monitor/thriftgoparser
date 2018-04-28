package proto

/**
 * A const is a constant value defined across languages that has a type and
 * a value. The trick here is that the declared type might not match the type
 * of the value object, since that is not determined until after parsing the
 * whole thing out.
 *
 */
type TConst struct{
	TDoc
	Type_ ICafType
	Name_ string
	Value_ *TConstValue
}

func NewConst(t ICafType,  name string, val *TConstValue) *TConst{
	var ret = new(TConst)
	ret.Type_ = t
	ret.Name_ = name
	ret.Value_ = val
	return ret
}

func (c *TConst)GetType() ICafType{
	return c.Type_
}

func (c *TConst)GetName() string{
	return c.Name_
}

func (c *TConst)GetValue() *TConstValue{
	return c.Value_
}