package models

import "strings"

type ModelString string

func (m ModelString) Model() string {
	return strings.Split(string(m), ":")[0]
}

func (m ModelString) Dest() string {
	return strings.Split(string(m), ":")[1]
}

func (m ModelString) Spec() string {
	return strings.Split(string(m), ":")[2]
}

func (m ModelString) Rev() string {
	return strings.Split(string(m), ":")[2]
}
