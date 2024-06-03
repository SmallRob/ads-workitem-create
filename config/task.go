package config

type Task struct {
	AssignedTo string `mapstructure:"assignedto" json:"assignedto" yaml:"assignedto"`
	// 初始估计
	OriginalEstimate string `mapstructure:"originalestimate" json:"originalestimate" yaml:"originalestimate"`
	// 已完成工作
	CompletedWork string `mapstructure:"completedwork" json:"completedwork" yaml:"completedwork"`

	// 工时类型
	WorkType string `mapstructure:"worktype" json:"worktype" yaml:"worktype"`
	// 工时收费状态
	BillableStatus string `mapstructure:"billablestatus" json:"billablestatus" yaml:"billablestatus"`
	// 出差
	Travel string `mapstructure:"travel" json:"travel" yaml:"travel"`
	// 出差地点
	TravelLocation string `mapstructure:"travellocation" json:"travellocation" yaml:"travellocation"`
	//BillableStatus  string `mapstructure:"app_url" json:"app_url" yaml:"app_url"`
	State string `mapstructure:"state" json:"state" yaml:"state"`
}
