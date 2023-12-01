package models

type Package struct {
	Name    string  `xml:"name,attr"`
	Classes []Class `xml:"classes>class"`
	Metrics
}
