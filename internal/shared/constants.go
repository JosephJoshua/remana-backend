package shared

type AppEnv string

const (
	AppEnvDev        = AppEnv("development")
	AppEnvProduction = AppEnv("production")
)
