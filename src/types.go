package main

type Auth struct {
	Type   string         `json:"type"`
	ApiKey []KeyValueType `json:"apikey"`
	Bearer []KeyValueType `json:"bearer"`
	Basic  []KeyValueType `json:"basic"`
}

type KeyValueType struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type"`
}
