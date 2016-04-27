package main

type CreateRuleRequest struct {
	Type    string `json:"type"`
	Trigger string `json:"trigger"`
	Target  string `json:"target"`
	Action  string `json:"action"`
}
