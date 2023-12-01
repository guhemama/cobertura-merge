package models

type Line struct {
	Number int `xml:"number,attr"`
	Hits   int `xml:"hits,attr"`
}
