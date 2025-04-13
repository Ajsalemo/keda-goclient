package config

import (
	apiv1 "k8s.io/api/core/v1"
)

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

type DeploymentStruct struct {
	Name       string            `json:"name"`
	Replicas   int32             `json:"replicas"`
	Containers []apiv1.Container `json:"containers"`
}

type ScaledJobStruct struct {
	Name       string           `json:"name"`
	Triggers   []map[string]any `json:"triggers"`
	Containers []Containers     `json:"containers"`
}

type ScaledObjectStruct struct {
	Name       string           `json:"name"`
	Triggers   []map[string]any `json:"triggers"`
	Containers []Containers     `json:"containers"`
}
type AuthenticationSecret struct {
	Name      string `json:"name"`
	Parameter string `json:"parameter"`
	Value     string `json:"value"`
}
