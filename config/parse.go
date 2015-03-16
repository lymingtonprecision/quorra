package config

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func (c *Config) VmMemoryInMB() int64 {
	if len(c.VM.Memory) == 0 {
		return 1024
	}

	re := regexp.MustCompile("(?i)^(\\d+)(G|M|K)B$")

	if !re.MatchString(c.VM.Memory) {
		panic(fmt.Errorf("unable to parse '%s' as an amount of memory", c.VM.Memory))
	}

	m := re.FindAllStringSubmatch(strings.ToUpper(c.VM.Memory), -1)

	s, err := strconv.Atoi(m[0][1])
	if err != nil {
		panic(err)
	}

	var x int64
	switch {
	case m[0][2] == "G":
		x = int64(s) * int64(1024)
	case m[0][2] == "M":
		x = int64(s)
	case m[0][2] == "K":
		x = int64(s) / 1024
	}

	return x
}
