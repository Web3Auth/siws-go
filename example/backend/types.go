package main

type Payload struct {
	Domain string `json:"domain"`
	Address string `json:"address"`
	URI string `json:"uri"`
	Version string `json:"version"`
	Options map[string]interface{} `json:"options"`
}