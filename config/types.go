package config

type Env struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Containers struct {
	Name  string `json:"name"`
	Image string `json:"image"`
	Env   []Env  `json:"env"`
}

type ScaledJobStruct struct {
	Metadata   map[string]any `json:"metadata"`
	Containers []Containers   `json:"containers"`
}

type AuthenticationSecret struct {
	Name      string `json:"name"`
	Parameter string `json:"parameter"`
	Value     string `json:"value"`
}
