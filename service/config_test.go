package service

import (
	"io/ioutil"
	"math/big"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestConfig(t *testing.T) {
	c := &States{
		[]StateConfig{
			{
				ChainID:     big.NewInt(1),
				GenesisFile: "abc",
			},
			{
				ChainID:     big.NewInt(1),
				GenesisFile: "",
			},
		},
	}
	out, err := yaml.Marshal(c)
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile("testdata/config", out, 0600)
	if err != nil {
		t.Fatal(err)
	}
}

func TestLoadConfig(t *testing.T) {
	c := &States{}
	out, err := ioutil.ReadFile("testdata/config")
	if err != nil {
		t.Fatal(err)
	}
	err = yaml.Unmarshal(out, c)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v", c)
}
