package ginCore

type AppConfig struct {
	Name       string `json:"name"`
	Port       int    `json:"port"`
	RunMode    string `json:"run_mode"`
	PublicPath string `json:"public_path"`
}

type Configure interface {
	GetAppConfig() *AppConfig
	GetDBConfig() *DBConfig
	GetCacheConfig() *CacheConfig
	GetLocalConfig() *DBConfig
}
