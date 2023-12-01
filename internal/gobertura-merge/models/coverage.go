package models

import "encoding/xml"

type Coverage struct {
	XMLName  xml.Name  `xml:"coverage"`
	Sources  []Source  `xml:"sources>source"`
	Packages []Package `xml:"packages>package"`
	Metrics
	LinesCovered    int `xml:"lines-covered,attr"`
	LinesValid      int `xml:"lines-valid,attr"`
	BranchesCovered int `xml:"branches-covered,attr"`
	BranchesValid   int `xml:"branches-valid,attr"`
}
