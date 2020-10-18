package iniconfig

func Watch(filename string, listener Listener) (*Config, error) {
	listener.listen(filename)
	return InitConfig(filename)
}
