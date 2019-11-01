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

	inParams := getInParams(fnType)
	inArgs := getInArgs(c, inParams)

	producedType := fnType.Out(0)
	c.providers[producedType] = provider{
		producedType,
	}

	instance := fnVal.Call(inArgs)
	c.instances[producedType] = instance
}

// Invoke .
func (c *container) invoke(fn interface{}) interface{} {

	fnVal := reflect.ValueOf(fn)
	fnType := fnVal.Type()

	inParams := getInParams(fnType)
	inArgs := getInArgs(c, inParams)

	ret := fnVal.Call(inArgs)
	return ret[0].Interface()
}

func getInParams(fnType reflect.Type) []param {
	numIn := fnType.NumIn()
	inParams := []param{}

	for i := 0; i < numIn; i++ {
		inParamType := fnType.In(i)
		inParams = append(inParams, param{
			inParamType,
			inParamType.Kind() == reflect.Slice,
		})
	}

	return inParams
}

func getInArgs(c *container, inParams []param) []reflect.Value {
	inArgs := []reflect.Value{}

	for _, param := range inParams {
		if param.IsSlice {
			ins := reflect.MakeSlice(param.Type, 0, 0)
			for k, v := range c.instances {
				if k.Implements(param.Type.Elem()) {
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
