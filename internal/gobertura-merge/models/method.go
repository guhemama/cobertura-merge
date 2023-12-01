package models

type Method struct {
	Name      string `xml:"name,attr"`
	Signature string `xml:"signature,attr"`
	Lines     []Line `xml:"lines>line"`
	Metrics
}
