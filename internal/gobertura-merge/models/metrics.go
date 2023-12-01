package models

type Metrics struct {
	LineRate   float64 `xml:"line-rate,attr"`
	BranchRate float64 `xml:"branch-rate,attr"`
	Complexity int     `xml:"complexity,attr"`
}
