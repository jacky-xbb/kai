package types

// SystemConfig contains the system configuration
type SystemConfig struct {
	SwapDir   string `required:"true"`
	Logfile   string `required:"false"`
	Port      string `required:"true"`
	DBDriver  string `required:"true"`
	MongoHost string `required:"false"`
	NGpu      int    `required:"true"`
	Workers   int    `required:"true"`
}

// YoloConfig contains the yolo configuration
type YoloConfig struct {
	DataCfg    string  `required:"true"`
	CfgFile    string  `required:"true"`
	WeightFile string  `required:"true"`
	Thresh     float64 `required:"true"`
	HierThresh float64 `required:"true"`
}

// ServerConfig contains the system and yolo configuration
type ServerConfig struct {
	ServerName string `required:"true"`

	System SystemConfig `required:"true"`

	Yolo YoloConfig `required:"true"`
}
