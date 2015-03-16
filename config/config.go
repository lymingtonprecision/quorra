package config

type NetworkConfig struct {
	Name       string
	Vlan       int
	Allocation string
	IpRange    string
}

type Config struct {
	ApiUrl   string
	Username string
	Password string

	ExtCert string
	ExtKey  string

	Datacenter string

	Default struct {
		Datastore string
		Folder    string
		Host      string
		OVA       string
	}

	VM struct {
		Datastore string
		Folder    string
		Memory    string
	}

	DataVolume struct {
		Datastore string
		Folder    string
	}

	Network struct {
		Private NetworkConfig
		Public  NetworkConfig
	}
}

func FromDefaultSources() (*Config, error) {
	c, err := ParseDefaultFiles()
	if err != nil {
		return nil, err
	}

	if err := c.ParseEnv(); err != nil {
		return nil, err
	}

	return c, nil
}
