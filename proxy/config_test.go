package proxy

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"math/big"
	"testing"
)

func TestConfig(t *testing.T) {
	c := &Config{
		StateConfig: StateConfig{
			ChainIDs: []*big.Int{big.NewInt(1), big.NewInt(2)},
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
	c := &Config{}
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
