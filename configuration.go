package main

//Configuration contains typed configuration from file
type Configuration struct {
	OutputDir string
}

//LoadConfiguration loads the configuration from file
func LoadConfiguration() Configuration {
	c := Configuration{}
	c.OutputDir = "/tmp/soffice"
	return c
}
