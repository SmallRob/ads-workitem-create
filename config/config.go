package config

type Configuration struct {
	App         App         `mapstructure:"app" json:"app" yaml:"app"`
	ADS         ADS         `mapstructure:"ads" json:"ads" yaml:"ads"`
	ProductItem ProductItem `mapstructure:"productitem" json:"productitem" yaml:"productitem"`
	Task        Task        `mapstructure:"task" json:"task" yaml:"task"`
	Log         Log         `mapstructure:"log" json:"log" yaml:"log"`
	Database    Database    `mapstructure:"database" json:"database" yaml:"database"`
	Jwt         Jwt         `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
	Redis       Redis       `mapstructure:"redis" json:"redis" yaml:"redis"`
}
