package core

import "fmt"

type Environment []map[string]any

func (env *Environment) push(v map[string]any) Environment {
	*env = append(*env, v)
	return *env
}

func (env *Environment) pop() map[string]any {
	l := len(*env)
	if l == 0 {
		panic("Cannot pop from an empty environment")
	}

	res := (*env)[l-1]
	*env = (*env)[:l-1]
	return res
}

func (env Environment) findVar(name string) (any, error) {
	for i := len(env) - 1; i >= 0; i-- {
		if val, ok := env[i][name]; ok {
			return val, nil
		}
	}

	return nil, fmt.Errorf("Variable %s not found", name)
}

func (env *Environment) setVar(name string, value any) error {
	for i := len(*env) - 1; i >= 0; i-- {
		if _, ok := (*env)[i][name]; ok {
			(*env)[i][name] = value
			return nil
		}
	}

	return fmt.Errorf("Variable %s not found", name)
}

func (env *Environment) declareVar(name string, value any) {
	l := len(*env)
	(*env)[l-1][name] = value
}

func DefaultEnvironment() Environment {
	env := make(Environment, 1)
	env[0] = make(map[string]any)
	env[0]["print"] = Function{
		Arity: 1,
		Call: func(_ Environment, args []any) (any, error) {
			fmt.Println(args[0])
			return nil, nil
		},
	}

	return env
}
