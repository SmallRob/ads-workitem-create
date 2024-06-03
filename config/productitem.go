package config

type ProductItem struct {
	Title    string `mapstructure:"title" json:"title" yaml:"title"`
	AreaPath string `mapstructure:"areapath" json:"areapath" yaml:"areapath"`
	//IterationPath string `mapstructure:"iterationpath" json:"iterationpath" yaml:"iterationpath"`
	State string `mapstructure:"state" json:"state" yaml:"state"`
}
