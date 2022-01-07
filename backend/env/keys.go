package env

type EnvVar struct {
	RedisAddresses string
}

var EnvVarKey EnvVar = EnvVar{
	RedisAddresses: "REDIS_ADDRESSES",
}
