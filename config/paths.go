package config

import (
	"path"
)

func (c *Config) VMStorePath() string {
	return path.Join(
		c.Datacenter,
		"vm",
		firstNonEmpty(c.VM.Folder, c.Default.Folder),
	)
}

func (c *Config) DataVolumeStorePath() string {
	return path.Join(
		c.Datacenter,
		"datastore",
		firstNonEmpty(c.DataVolume.Datastore, c.Default.Datastore),
	)
}

func (c *Config) DataVolumePath() string {
	return firstNonEmpty(c.DataVolume.Folder, c.Default.Folder, "/")
}

func firstNonEmpty(strings ...string) string {
	for _, s := range strings {
		if len(s) > 0 {
			return s
		}
	}

	return ""
}
