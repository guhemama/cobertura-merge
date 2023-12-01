package models

type Class struct {
	Name     string   `xml:"name,attr"`
	Filename string   `xml:"filename,attr"`
	Methods  []Method `xml:"methods>method"`
	Lines    []Line   `xml:"lines>line"`
	Metrics
}

func (class *Class) recalculateMetrics() {
}
