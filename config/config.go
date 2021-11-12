package config

func InitializeConfig() {
	NewLoggerService()
	ConnectDatabase()
	InitSessionStore()
	ConnectNats()
}
