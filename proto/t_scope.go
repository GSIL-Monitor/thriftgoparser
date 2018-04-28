package proto

import "fmt"

/**
 * This represents a variable scope used for looking up predefined types and
 * services. Typically, a scope is associated with a t_program. Scopes are not
 * used to determine code generation, but rather to resolve identifiers at
 * parse time.
 *
 */

type TScope struct{
	// Map of names to types
	Types_ map[string] ICafType

	// Map of names to constants
	Constants_ map[string] *TConst

	// Map of names to services
	Services_ map[string] *TService

}

func  NewScope()  *TScope{
	ret := new(TScope)
	ret.Types_ = make(map[string]ICafType)
	ret.Constants_ = make(map[string]*TConst)
	ret.Services_ = make(map[string]*TService)
	return ret
}

func (s *TScope)AddType(name string, t ICafType) {
	s.Types_[name] = t
}

func (s *TScope)GetType(name string) ICafType{
	return s.Types_[name];
}

func (s *TScope)AddService(name string, service *TService){
	s.Services_[name] = service
}

func (s *TScope)GetService(name string) *TService{
	return s.Services_[name]
}

func (s *TScope)AddConstant(name string, constant* TConst) {
	_, ok := s.Constants_[name]
	if (ok){
		panic("Enum " + name + " is already defined!")
	} else {
		s.Constants_[name] = constant;
	}
}

func (s *TScope)GetConstant(name string) *TConst{
	return s.Constants_[name]
}

func (s *TScope)Print() {
	for k, v := range s.Types_{
		fmt.Printf("%s => %s\n", k, v.GetName())
	}
}

func (s *TScope)ResolveConstValue(const_val *TConstValue, ttype ICafType) {
	switch conV := ttype.(type){
	case *TMap:
		{
			for k,v := range const_val.GetMap() {
				s.ResolveConstValue(k, conV.GetKeyType())
				s.ResolveConstValue(v, conV.GetKeyType())
			}
		}
	case *TList:
		{
			for _, v := range const_val.GetMap() {
				s.ResolveConstValue(v, conV.GetElemType())
			}
		}
	case *TSet:
		{
			for _,v := range const_val.GetMap()  {
				s.ResolveConstValue(v, conV.GetElemType())
			}
		}
	case *TStruct:
		{
			for k,v := range const_val.GetMap(){
				f := conV.GetFieldByName(k.GetString());
				if (nil == f) {
					panic(fmt.Sprintf("No field named %s was found in struct of type %s", k.GetString(), conV.GetName()))
				}
				s.ResolveConstValue(v, f.GetType())
			}
		}
	case *TEnum:
		{
			{
				if (const_val.GetType() == CV_IDENTIFIER) {
					const_val.SetEnum(conV)
				} else {
					// enum constant with non-identifier value. set the enum and find the
					// value's name.
					emVal := conV.GetConstantByValue(const_val.GetInt())
					if (emVal == nil) {
						panic(fmt.Sprintf("Couldn't find a named value in enum %s for value :%d", conV.GetName(), const_val.GetInt()))
					}
					const_val.SetIdentifier(conV.GetName() + "." + emVal.GetName())
					const_val.SetEnum(conV)
				}
			}
		}
		default:
			{
				if const_val.GetType() != CV_IDENTIFIER {
					constant := s.GetConstant(const_val.GetIdentifier())
					if constant == nil {
						panic("No enum value or constant found named \"" + const_val.GetIdentifier() + "\"!");
					}

					// Resolve typedefs to the underlying type
					/*t_type* const_type = constant.GetType().get_true_type();

					if (const_type->is_base_type()) {
						switch (((t_base_type*)const_type)->get_base()) {
							case t_base_type::TYPE_I16:
							case t_base_type::TYPE_I32:
							case t_base_type::TYPE_I64:
							case t_base_type::TYPE_BOOL:
							case t_base_type::TYPE_BYTE:
							const_val->set_integer(constant->get_value()->get_integer());
							break;
							case t_base_type::TYPE_STRING:
							const_val->set_string(constant->get_value()->get_string());
							break;
							case t_base_type::TYPE_DOUBLE:
							const_val->set_double(constant->get_value()->get_double());
							break;
							case t_base_type::TYPE_VOID:
							throw "Constants cannot be of type VOID";
						}
					} else if (const_type->is_map()) {
						const std::map<t_const_value*, t_const_value*>& map = constant->get_value()->get_map();
						std::map<t_const_value*, t_const_value*>::const_iterator v_iter;

						const_val->set_map();
						for (v_iter = map.begin(); v_iter != map.end(); ++v_iter) {
							const_val->add_map(v_iter->first, v_iter->second);
						}
					} else if (const_type->is_list()) {
						const std::vector<t_const_value*>& val = constant->get_value()->get_list();
						std::vector<t_const_value*>::const_iterator v_iter;

						const_val->set_list();
						for (v_iter = val.begin(); v_iter != val.end(); ++v_iter) {
							const_val->add_list(*v_iter);
						}
					}
					*/
				}
			}
	}
}
