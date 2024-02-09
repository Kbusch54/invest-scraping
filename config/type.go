package config

type Config struct {
	Persistence Persistence `mapstructure:"persistence"`
	Services    Services    `mapstructure:"services"`
	Server      Server      `mapstructure:"server"`
	Stream      Stream      `mapstructure:"stream"`
	Monitors    []*Monitor
}

type Persistence struct {
	MongoDB MongoDB `mapstructure:"mongodb"`
}

type MongoDB struct {
	URL      string `mapstructure:"url"`
	Database string `mapstructure:"database"`
}

type Services struct {
	InvestingServices InvestingServices `mapstructure:"investingservices"`
}

type InvestingServices struct {
	Host string `mapstructure:"host"`
}

type Server struct {
	Port string `mapstructure:"port"`
}
type Stream struct {
	Kafka Kafka `mapstructure:"kafka"`
}

type Kafka struct {
	Brokers []string `mapstructure:"brokers"`
}
type Monitor struct {
	Name        string `yaml:"name"`
	Key         string `yaml:"key"`
	Type        string `yaml:"type"`
	Topic       string `yaml:"topic"`
	IndexType   string `yaml:"indexType"`
	Symbol      string `yaml:"symbol"`
	Endpoint    string `yaml:"endpoint"`
	EndpointExt string `yaml:"endpointExt"`
	NameXpath   string `yaml:"nameXpath"`
	PriceXpath  string `yaml:"priceXpath"`
}

func (m *Monitor) TopicName() string {
	if m.Topic == "" {
		return m.Key
	} else {
		return m.Topic
	}
}
