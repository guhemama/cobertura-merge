package main

import (
	"encoding/xml"
	"fmt"
	"os"
)

type Metrics struct {
	LineRate   float64 `xml:"line-rate,attr"`
	BranchRate float64 `xml:"branch-rate,attr"`
	Complexity int     `xml:"complexity,attr"`
}

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

type Source struct {
	Path string `xml:",chardata"`
}

type Package struct {
	Name    string  `xml:"name,attr"`
	Classes []Class `xml:"classes>class"`
	Metrics
}

type Class struct {
	Name     string   `xml:"name,attr"`
	Filename string   `xml:"filename,attr"`
	Methods  []Method `xml:"methods>method"`
	Lines    []Line   `xml:"lines>line"`
	Metrics
}

type Method struct {
	Name      string `xml:"name,attr"`
	Signature string `xml:"signature,attr"`
	Lines     []Line `xml:"lines>line"`
	Metrics
}

type Line struct {
	Number int `xml:"number,attr"`
	Hits   int `xml:"hits,attr"`
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run merge_cobertura.go output.xml input1.xml input2.xml ...")
		os.Exit(1)
	}

	// Create a new coverage report to store the merged data
	mergedCoverage := &Coverage{}

	// Loop through input files starting from the third argument (index 2)
	for _, filePath := range os.Args[2:] {
		// Read the input Cobertura XML file
		xmlData, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", filePath, err)
			continue
		}

		// Parse the XML data into a Coverage structure
		var coverage Coverage
		err = xml.Unmarshal(xmlData, &coverage)
		if err != nil {
			fmt.Printf("Error parsing XML from file %s: %v\n", filePath, err)
			continue
		}

		// Merge the data from the current report into the mergedCoverage
		mergeCoverageReports(mergedCoverage, &coverage)
	}

	// Recalculate the metrics from the merged reports
	mergedCoverage.recalculateMetrics()

	// Encode the merged coverage data back to XML
	mergedXML, err := xml.MarshalIndent(mergedCoverage, "", "  ")
	if err != nil {
		fmt.Println("Error encoding merged coverage data:", err)
		os.Exit(1)
	}
	// Prepend the XML header
	output := []byte(xml.Header + string(mergedXML))

	// Write the merged coverage data to the output file
	outputFilePath := os.Args[1]
	err = os.WriteFile(outputFilePath, output, 0644)
	if err != nil {
		fmt.Printf("Error writing merged coverage data to %s: %v\n", outputFilePath, err)
		os.Exit(1)
	}

	fmt.Printf("Merged coverage data written to %s\n", outputFilePath)
}

func mergeCoverageReports(mergedCoverage, currentCoverage *Coverage) {
	// Loop through sources in the current report
	for _, currentSource := range currentCoverage.Sources {
		// Check if a source exists with the same path in mergedCoverage
		existingSource := findSourceByPath(mergedCoverage, currentSource.Path)

		if existingSource == nil {
			// If the source doesn't exist, add it to mergedCoverage
			mergedCoverage.Sources = append(mergedCoverage.Sources, currentSource)
		}
	}

	// Loop through the packages in the current report
	for _, currentPackage := range currentCoverage.Packages {
		// Check if a package with the same name already exists in mergedCoverage
		existingPackage := findPackageByName(mergedCoverage, currentPackage.Name)

		if existingPackage == nil {
			// If the package doesn't exist, add it to the mergedCoverage
			mergedCoverage.Packages = append(mergedCoverage.Packages, currentPackage)
		} else {
			// If the package already exists, merge the class data
			mergeClasses(existingPackage, &currentPackage)
		}
	}
}

func findSourceByPath(coverage *Coverage, sourcePath string) *Source {
	for _, source := range coverage.Sources {
		if source.Path == sourcePath {
			return &source
		}
	}
	return nil
}

func findPackageByName(coverage *Coverage, packageName string) *Package {
	for _, pkg := range coverage.Packages {
		if pkg.Name == packageName {
			return &pkg
		}
	}
	return nil
}

func mergeClasses(existingPackage *Package, currentPackage *Package) {
	// Loop through the classes in the current package
	for _, currentClass := range currentPackage.Classes {
		// Check if a class with the same name already exists in existingPackage
		existingClass := findClassByName(existingPackage, currentClass.Name)

		if existingClass == nil {
			// If the class doesn't exist, add it to existingPackage
			existingPackage.Classes = append(existingPackage.Classes, currentClass)
		} else {
			// If the class already exists, merge the method data
			mergeMethods(existingClass, &currentClass)
			mergeClassLines(existingClass, &currentClass)
		}
	}
}

func findClassByName(packageData *Package, className string) *Class {
	for _, class := range packageData.Classes {
		if class.Name == className {
			return &class
		}
	}
	return nil
}

func mergeMethods(existingClass *Class, currentClass *Class) {
	// Loop through the methods in the current class
	for _, currentMethod := range currentClass.Methods {
		// Check if a method with the same name already exists in existingClass
		existingMethod := findMethodByName(existingClass, currentMethod.Name)

		if existingMethod == nil {
			// If the method doesn't exist, add it to existingClass
			existingClass.Methods = append(existingClass.Methods, currentMethod)
		} else {
			// If the method already exists, merge the lines data
			mergeMethodLines(existingMethod, &currentMethod)
		}
	}
}

func findMethodByName(class *Class, methodName string) *Method {
	for _, method := range class.Methods {
		if method.Name == methodName {
			return &method
		}
	}
	return nil
}

func mergeMethodLines(existingMethod *Method, currentMethod *Method) {
	// Create a map to quickly access line data by line number
	lineMap := make(map[int]*Line)
	for i := range existingMethod.Lines {
		lineMap[existingMethod.Lines[i].Number] = &existingMethod.Lines[i]
	}

	// Loop through the lines in the current class
	for i := range currentMethod.Lines {
		line := &currentMethod.Lines[i]
		existingLine, ok := lineMap[line.Number]

		if ok {
			// If the line already exists, add the hits
			existingLine.Hits += line.Hits
		} else {
			// If the line doesn't exist, add it to existingClass
			existingMethod.Lines = append(existingMethod.Lines, *line)
			lineMap[line.Number] = &existingMethod.Lines[len(existingMethod.Lines)-1]
		}
	}
}

func mergeClassLines(existingClass *Class, currentClass *Class) {
	// Create a map to quickly access line data by line number
	lineMap := make(map[int]*Line)
	for i := range existingClass.Lines {
		lineMap[existingClass.Lines[i].Number] = &existingClass.Lines[i]
	}

	// Loop through the lines in the current class
	for i := range currentClass.Lines {
		line := &currentClass.Lines[i]
		existingLine, ok := lineMap[line.Number]

		if ok {
			// If the line already exists, add the hits
			existingLine.Hits += line.Hits
		} else {
			// If the line doesn't exist, add it to existingClass
			existingClass.Lines = append(existingClass.Lines, *line)
			lineMap[line.Number] = &existingClass.Lines[len(existingClass.Lines)-1]
		}
	}
}

func (cov *Coverage) recalculateMetrics() {
	var coverageValidLines, coverageCoveredLines float64

	for i := range cov.Packages {
		pkg := &cov.Packages[i]
		var packageComplexity int
		var packageValidLines, packageCoveredLines float64

		for i := range pkg.Classes {
			class := &pkg.Classes[i]
			var classComplexity int
			var classValidLines, classCoveredLines float64

			for i := range class.Methods {
				method := &class.Methods[i]
				var methodValidLines, methodCoveredLines float64

				for i := range method.Lines {
					line := &method.Lines[i]
					methodValidLines += 1

					if line.Hits > 0 {
						methodCoveredLines += 1
					}
				}

				method.LineRate = methodCoveredLines / methodValidLines
				classValidLines += methodValidLines
				classCoveredLines += methodCoveredLines
				classComplexity += method.Complexity
			}

			class.LineRate = classCoveredLines / classValidLines
			class.Complexity = classComplexity
			packageValidLines += classValidLines
			packageCoveredLines += classCoveredLines
			packageComplexity += classComplexity
		}

		pkg.Complexity = packageComplexity
		pkg.LineRate = packageCoveredLines / packageValidLines
		coverageValidLines += packageValidLines
		coverageCoveredLines += packageCoveredLines
	}

	cov.LinesCovered = int(coverageCoveredLines)
	cov.LinesValid = int(coverageValidLines)
	cov.LineRate = coverageCoveredLines / coverageValidLines
}
