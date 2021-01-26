// Copyright (C) 2021 The Dank Grinder authors.
//
// This source code has been released under the GNU Affero General Public
// License v3.0. A copy of this license is available at
// https://www.gnu.org/licenses/agpl-3.0.en.html

package config

import (
	"fmt"
	"os"
	"path"

	"gopkg.in/yaml.v2"
)

const (
	ShiftStateActive  = "active"
	ShiftStateDormant = "dormant"
)

type Config struct {
	Token              string             `yaml:"token"`
	ChannelID          string             `yaml:"channel_id"`
	Features           Features           `yaml:"features"`
	Compat             Compat             `yaml:"compatibility"`
	SuspicionAvoidance SuspicionAvoidance `yaml:"suspicion_avoidance"`
	Swarm              Swarm              `yaml:"swarm"`
}

type Swarm struct {
	Instances []Instance `yaml:"instances"`
}

type Instance struct {
	Token     string  `yaml:"token"`
	ChannelID string  `yaml:"channel_id"`
	Shifts    []Shift `yaml:"shifts"`
}

type Compat struct {
	PostmemeOpts    []string `yaml:"postmeme_options"`
	AllowedSearches []string `yaml:"allowed_searches"`
	Cooldown        Cooldown `yaml:"cooldown"`
	AutoSell        []string `yaml:"auto_sell"`
}

type Cooldown struct {
	Beg      int `yaml:"beg"`
	Fish     int `yaml:"fish"`
	Hunt     int `yaml:"hunt"`
	Postmeme int `yaml:"postmeme"`
	Search   int `yaml:"search"`
	Highlow  int `yaml:"highlow"`
	Margin   int `yaml:"margin"`
}

type Features struct {
	Commands     Commands `yaml:"commands"`
	AutoBuy      AutoBuy  `yaml:"auto_buy"`
	BalanceCheck bool     `yaml:"balance_check"`
	LogToFile    bool     `yaml:"log_to_file"`
	Debug        bool     `yaml:"debug"`
}

type AutoBuy struct {
	FishingPole  bool `yaml:"fishing_pole"`
	HuntingRifle bool `yaml:"hunting_rifle"`
	Laptop       bool `yaml:"laptop"`
}

type Commands struct {
	Fish bool `yaml:"fish"`
	Hunt bool `yaml:"hunt"`
}

type SuspicionAvoidance struct {
	Typing       Typing       `yaml:"typing"`
	MessageDelay MessageDelay `yaml:"message_delay"`
	Shifts       []Shift      `yaml:"shifts"`
}

type Typing struct {
	Base     int `yaml:"base"`     // A base duration in milliseconds.
	Speed    int `yaml:"speed"`    // Speed in keystrokes per minute.
	Variance int `yaml:"variance"` // A random value in milliseconds from [0,n) added to the base.
}

// MessageDelay is used to
type MessageDelay struct {
	Base     int `yaml:"base"`     // A base duration in milliseconds.
	Variance int `yaml:"variance"` // A random value in milliseconds from [0,n) added to the base.
}

// Shift indicates an application state (active or dormant) for a duration.
type Shift struct {
	State    string   `yaml:"state"`
	Duration Duration `yaml:"duration"`
}

// Duration is not related to a time.Duration. It is a structure used in a Shift
// type.
type Duration struct {
	Base     int `yaml:"base"`     // A base duration in seconds.
	Variance int `yaml:"variance"` // A random value in seconds from [0,n) added to the base.
}

// Load loads the config from the expected path.
func Load(dir string) (Config, error) {
	f, err := os.Open(path.Join(dir, "config.yml"))
	if err != nil {
		return Config{}, fmt.Errorf("error while opening config file: %v", err)
	}

	var cfg Config
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return Config{}, fmt.Errorf("error while decoding config: %v", err)
	}

	return cfg, nil
}

func (c Config) Validate() error {
	if c.Token == "" {
		return fmt.Errorf("token: no authorization token")
	}
	if c.ChannelID == "" {
		return fmt.Errorf("channel_id: no channel id")
	}
	if len(c.SuspicionAvoidance.Shifts) == 0 {
		return fmt.Errorf("suspicion_avoidance.shifts: no shifts, at least 1 is required")
	}
	if len(c.Compat.PostmemeOpts) == 0 {
		return fmt.Errorf("compatibility.postmeme: no compatibility options")
	}
	if len(c.Compat.AllowedSearches) == 0 {
		return fmt.Errorf("compatibility.allowed_searches: no compatibility options")
	}
	if c.Compat.Cooldown.Postmeme <= 0 {
		return fmt.Errorf("compatibility.cooldown.postmeme: value must be greater than 0")
	}
	if c.Compat.Cooldown.Hunt <= 0 {
		return fmt.Errorf("compatibility.cooldown.hunt: value must be greater than 0")
	}
	if c.Compat.Cooldown.Highlow <= 0 {
		return fmt.Errorf("compatibility.cooldown.highlow: value must be greater than 0")
	}
	if c.Compat.Cooldown.Fish <= 0 {
		return fmt.Errorf("compatibility.cooldown.fish: value must be greater than 0")
	}
	if c.Compat.Cooldown.Search <= 0 {
		return fmt.Errorf("compatibility.cooldown.search: value must be greater than 0")
	}
	if c.Compat.Cooldown.Beg <= 0 {
		return fmt.Errorf("compatibility.cooldown.beg: value must be greater than 0")
	}
	if c.Compat.Cooldown.Margin < 0 {
		return fmt.Errorf("compatibility.cooldown.margin: value must be greater than or equal to 0")
	}

	for _, shift := range c.SuspicionAvoidance.Shifts {
		if shift.State != ShiftStateActive && shift.State != ShiftStateDormant {
			return fmt.Errorf("invalid shift state: %v", shift.State)
		}
	}
	return nil
}
