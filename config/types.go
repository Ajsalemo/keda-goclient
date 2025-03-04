package config

type Env struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type ScaledJobStruct struct {
	Metadata map[string]any `json:"metadata"`
	Env      []Env          `json:"env"`
}
