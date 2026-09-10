package event

// Status_Lgr 状态信息 (lagrange)
type Status_Lgr struct {
	Status `mapstructure:"sender,squash"`

	AppInitialized bool `json:"app_initialized" mapstructure:"app_initialized"`
	AppEnabled     bool `json:"app_enabled" mapstructure:"app_enabled"`
	AppGood        bool `json:"app_good" mapstructure:"app_good"`
}
