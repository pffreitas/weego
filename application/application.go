package application

type WeegoApplication struct {
	container *container
}

func New(app interface{}) WeegoApplication {
	container := newContainer()

	config, err := processConfig(app)
	if err == nil {
		container.injectConfig(config)
	}

	return container.invoke(newWeegoApplication).(WeegoApplication)
}

func newWeegoApplication(container *container) WeegoApplication {
	return WeegoApplication{container}
}

func (wa *WeegoApplication) Provide(constructor interface{}) {
	wa.container.provide(constructor)
}

func (wa *WeegoApplication) Invoke(fn interface{}) {
	wa.container.invoke(fn)
}
