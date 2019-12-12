package glox

type Environment struct {
	Enclosing *Environment
	Values    map[string]interface{}
}

func (env *Environment) define(name string, value interface{}) {
	if env.Values == nil {
		env.Values = make(map[string]interface{})
	}
	env.Values[name] = value
}

func (env *Environment) assign(name Token, value interface{}) {
	if env.Values != nil {
		if _, hasKey := env.Values[name.Lexeme]; hasKey {
			env.Values[name.Lexeme] = value
		}
	}

	if env.Enclosing != nil {
		env.Enclosing.assign(name, value)
	}
}

func (env *Environment) get(name Token) interface{} {
	//fmt.Println(env)
	if env.Values != nil {
		if _, hasKey := env.Values[name.Lexeme]; hasKey {
			return env.Values[name.Lexeme]
		}
	}
	if env.Enclosing != nil {
		return env.Enclosing.get(name)
	}

	return nil
}
