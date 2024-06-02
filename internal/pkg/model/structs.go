package model

type SyntheticTestsConfig struct {
	HttpTests []SyntheticTestHttpOptions `yaml:"HttpTests"`
	GrpcTests []SyntheticTestGrpcOptions `yaml:"GrpcTests"`
}

type SyntheticTestHttpOptions struct {
	Name             string            `yaml:"Name"`
	Region           string            `yaml:"Region,omitempty"`
	Timeout          float64           `yaml:"Timeout,omitempty"`
	Url              string            `yaml:"Url,omitempty"`
	Method           string            `yaml:"Method,omitempty"` //GET, POST..
	Headers          map[string]string `yaml:"Headers,omitempty"`
	JwtToken         string            `yaml:"JwtToken,omitempty"`
	TickEverySec     int64             `yaml:"TickEverySec,omitempty"`
	SlackChannel     string            `yaml:"SlackChannel,omitempty"`
	RenotifyInterval int64             `yaml:"RenotifyInterval,omitempty"`

	RequestBody           string `yaml:"RequestBody,omitempty"`
	AssertionJsonPath     string `yaml:"AssertionJsonPath,omitempty"`
	AssertionJsonOperator string `yaml:"AssertionJsonOperator,omitempty"`
	AssertionTargetValue  string `yaml:"AssertionTargetValue,omitempty"`
}

type SyntheticTestGrpcOptions struct {
	Name             string            `yaml:"Name,omitempty"`
	Region           string            `yaml:"Region,omitempty"`
	Timeout          float64           `yaml:"Timeout,omitempty"`
	Url              string            `yaml:"Url,omitempty"`
	Method           string            `yaml:"Method,omitempty"` //GET, POST..
	Headers          map[string]string `yaml:"Headers,omitempty"`
	JwtToken         string            `yaml:"JwtToken,omitempty"`
	TickEverySec     int64             `yaml:"TickEverySec,omitempty"`
	SlackChannel     string            `yaml:"SlackChannel,omitempty"`
	RenotifyInterval int64             `yaml:"RenotifyInterval,omitempty"`

	Host                     string `yaml:"Host,omitempty"`
	Port                     int64  `yaml:"Port,omitempty"`
	ServiceToCheck           string `yaml:"ServiceToCheck,omitempty"`
	Message                  string `yaml:"Message,omitempty"`
	CompressedJsonDescriptor string `yaml:"CompressedJsonDescriptor,omitempty"`
}
