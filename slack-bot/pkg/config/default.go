package config

// DefaultConfig with some common values
var DefaultConfig = Config{
	StoragePath: "./storage/",
	Logger: Logger{
		File:  "./edith.log",
		Level: "info",
	},
	OpenWeather: OpenWeather{
		Units: "metric",
	},
}
