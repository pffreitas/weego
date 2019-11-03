package application

// WeegoApplication .
type WeegoApplication struct {
	container *container
}

// New .
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

// Provide .
func (wa *WeegoApplication) Provide(constructor interface{}) {
	wa.container.provide(constructor)
}

// Invoke .
func (wa *WeegoApplication) Invoke(fn interface{}) {
	wa.container.invoke(fn)
}
