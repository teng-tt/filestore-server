package config

type Server struct {
	System   System   `json:"system"`
	DBConfig DBConfig `json:"DBConfig"`
}
