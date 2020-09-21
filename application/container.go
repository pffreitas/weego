package application

import (
	"reflect"
)

type provider struct {
	producedType reflect.Type
}

type container struct {
	providers map[reflect.Type]provider
	instances map[reflect.Type][]reflect.Value
}

type param struct {
	Type    reflect.Type
	IsSlice bool
}

func newContainer() *container {
	container := &container{
		map[reflect.Type]provider{},
		map[reflect.Type][]reflect.Value{},
	}

	v := reflect.ValueOf(container)
	container.instances[v.Type()] = []reflect.Value{v}

	return container
}

func (c *container) provide(constructor interface{}) {

	fnVal := reflect.ValueOf(constructor)
	fnType := fnVal.Type()

	var params []param
	var args []reflect.Value
	var producedType reflect.Type
	var instance []reflect.Value

	if fnType.Kind() == reflect.Struct {
		params = getParamsFromFields(fnType)
		args = getArgs(c, params)
		producedType = fnType
		instance = instantiateStruct(fnType, args)
	} else {
		params = getParamsFromFuncIn(fnType)
		args = getArgs(c, params)
		producedType = fnType.Out(0)
		instance = fnVal.Call(args)
	}

	//c.providers[producedType] = provider{producedType}
	c.instances[producedType] = instance
}

func instantiateStruct(fnType reflect.Type, args []reflect.Value) []reflect.Value {
	ptr := reflect.New(fnType).Elem()
	for i := 0; i < ptr.NumField(); i++ {
		field := ptr.Field(i)
		field.Set(args[i])
	}
	return []reflect.Value{ptr}
}

func (c *container) invoke(fn interface{}) interface{} {

	fnVal := reflect.ValueOf(fn)
	fnType := fnVal.Type()

	inParams := getParamsFromFuncIn(fnType)
	inArgs := getArgs(c, inParams)

	ret := fnVal.Call(inArgs)
	return ret[0].Interface()
}

func getParamsFromFields(fnType reflect.Type) []param {
	var inParams []param

	for i := 0; i < fnType.NumField(); i++ {
		inParamType := fnType.Field(i).Type
		inParams = append(inParams, param{
			inParamType,
			inParamType.Kind() == reflect.Slice,
		})
	}

	return inParams
}

func getParamsFromFuncIn(fnType reflect.Type) []param {
	numIn := fnType.NumIn()
	var inParams []param

	for i := 0; i < numIn; i++ {
		inParamType := fnType.In(i)
		inParams = append(inParams, param{
			inParamType,
			inParamType.Kind() == reflect.Slice,
		})
	}

	return inParams
}

func getArgs(c *container, inParams []param) []reflect.Value {
	var inArgs []reflect.Value

	for _, param := range inParams {
		if param.IsSlice {
			ins := reflect.MakeSlice(param.Type, 0, 0)
			for k, v := range c.instances {
				if k == param.Type.Elem() {
					ins = reflect.Append(ins, v[0])
				} else if param.Type.Elem().Kind() == reflect.Interface && k.Implements(param.Type.Elem()) {
					ins = reflect.Append(ins, v[0])
				}
			}

			inArgs = append(inArgs, ins)
		} else {
			inArgs = append(inArgs, c.instances[param.Type][0])
		}
	}

	return inArgs
}

func (c *container) injectConfig(config weegoConfig) {
	c.instances[config.configType] = []reflect.Value{config.configValue}
}
