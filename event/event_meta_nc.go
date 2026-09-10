package event

// Status_Nc 状态信息 (NapCat 扩展)
type Status_Nc struct {
	Status `mapstructure:"sender,squash"`

	AppInitialized bool `json:"app_initialized" mapstructure:"app_initialized"`
	AppEnabled     bool `json:"app_enabled" mapstructure:"app_enabled"`
	AppGood        bool `json:"app_good" mapstructure:"app_good"`
}
