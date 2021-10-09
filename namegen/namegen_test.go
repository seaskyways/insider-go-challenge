package namegen

import "testing"

func TestGenerate(t *testing.T) {
	t.Log(Generate())
}

func TestGenerateTeamName(t *testing.T) {
	t.Log(GenerateTeamName())
}
