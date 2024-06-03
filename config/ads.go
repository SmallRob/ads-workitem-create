package config

type ADS struct {
	OrganizationUrl     string `mapstructure:"organizationurl" json:"organizationurl" yaml:"organizationurl"`
	PersonalaccessToken string `mapstructure:"personalaccesstoken" json:"personalaccesstoken" yaml:"personalaccesstoken"`
	// 项目名称
	ProjectName string `mapstructure:"projectname" json:"projectname" yaml:"projectname"`
	// Api版本6.0
	ApiVersion string `mapstructure:"api_version" json:"api_version" yaml:"api_version"`
	// 默认返回迭代数量
	IterationsCount int `mapstructure:"iterationscount" json:"iterationscount" yaml:"iterationscount"`

	AppName string `mapstructure:"app_name" json:"app_name" yaml:"app_name"`
	AppUrl  string `mapstructure:"app_url" json:"app_url" yaml:"app_url"`
}
