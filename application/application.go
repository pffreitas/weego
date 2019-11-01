package application

// WeegoApplication .
type WeegoApplication struct {
	container *container
}

// New .
func New(app interface{}) *WeegoApplication {
	container := newContainer()

	config, err := processConfig(app)
	if err == nil {
		container.injectConfig(config)
	}

	return &WeegoApplication{
		container,
	}
}

// Run .
func (wa *WeegoApplication) Run() {

}
