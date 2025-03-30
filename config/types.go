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

type AuthenticationRefName struct {
	Name string `json:"name"`
}

type Triggers struct {
	Type              string                `json:"type"`
	Metadata          map[string]any        `json:"metadata"`
	AuthenticationRef AuthenticationRefName `json:"authenticationRef"`
}

type ScaledJobStruct struct {
	Name       string           `json:"name"`
	Triggers   []map[string]any `json:"triggers"`
	Containers []Containers     `json:"containers"`
}

type AuthenticationSecret struct {
	Name      string `json:"name"`
	Parameter string `json:"parameter"`
	Value     string `json:"value"`
}
