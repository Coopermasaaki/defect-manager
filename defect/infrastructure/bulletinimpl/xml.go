package bulletinimpl

import "encoding/xml"

type CvrfBA struct {
	XMLName            xml.Name           `xml:"cvrfdoc,omitempty"`
	Xmlns              string             `xml:"xmlns,attr"`
	XmlnsCvrf          string             `xml:"xmlns:cvrf,attr"`
	DocumentTitle      DocumentTitle      `xml:"DocumentTitle,omitempty"`
	DocumentType       string             `xml:"DocumentType"`
	DocumentPublisher  DocumentPublisher  `xml:"DocumentPublisher,omitempty"`
	DocumentTracking   DocumentTracking   `xml:"DocumentTracking,omitempty"`
	DocumentNotes      DocumentNotes      `xml:"DocumentNotes,omitempty"`
	DocumentReferences DocumentReferences `xml:"DocumentReferences,omitempty"`
	ProductTree        ProductTree        `xml:"ProductTree,omitempty"`
	Vulnerability      []Vulnerability    `xml:"Vulnerability,omitempty"`
}

type DocumentTitle struct {
	XMLName       xml.Name `xml:"DocumentTitle,omitempty"`
	XmlLang       string   `xml:"xml:lang,attr"`
	DocumentTitle string   `xml:",innerxml"`
}

type DocumentPublisher struct {
	XMLName          xml.Name `xml:"DocumentPublisher,omitempty"`
	Type             string   `xml:"Type,attr"`
	ContactDetails   string   `xml:"ContactDetails"`
	IssuingAuthority string   `xml:"IssuingAuthority"`
}

type DocumentTracking struct {
	XMLName            xml.Name        `xml:"DocumentTracking,omitempty"`
	Identification     Identification  `xml:"Identification,omitempty"`
	Status             string          `xml:"Status"`
	Version            string          `xml:"Version"`
	RevisionHistory    RevisionHistory `xml:"RevisionHistory,omitempty"`
	InitialReleaseDate string          `xml:"InitialReleaseDate"`
	CurrentReleaseDate string          `xml:"CurrentReleaseDate"`
	Generator          Generator       `xml:"Generator,omitempty"`
}

type Identification struct {
	XMLName xml.Name `xml:"Identification,omitempty"`
	Id      string   `xml:"ID"`
}

type RevisionHistory struct {
	XMLName  xml.Name   `xml:"RevisionHistory,omitempty"`
	Revision []Revision `xml:"Revision,omitempty"`
}

type Revision struct {
	XMLName     xml.Name `xml:"Revision,omitempty"`
	Number      string   `xml:"Number"`
	Date        string   `xml:"Date"`
	Description string   `xml:"Description"`
}

type Generator struct {
	XMLName xml.Name `xml:"Generator,omitempty"`
	Engine  string   `xml:"Engine"`
	Date    string   `xml:"Date"`
}

type DocumentNotes struct {
	XMLName xml.Name `xml:"DocumentNotes,omitempty"`
	Note    []Note   `xml:"Note,omitempty"`
}

type Note struct {
	XMLName xml.Name `xml:"Note,omitempty"`
	Title   string   `xml:"Title,attr"`
	Type    string   `xml:"Type,attr"`
	Ordinal string   `xml:"Ordinal,attr"`
	XmlLang string   `xml:"xml:lang,attr"`
	Note    string   `xml:",innerxml"`
}

type DocumentReferences struct {
	XMLName      xml.Name       `xml:"DocumentReferences,omitempty"`
	CveReference []CveReference `xml:"Reference,omitempty"`
}

type CveReference struct {
	XMLName xml.Name `xml:"Reference,omitempty"`
	Type    string   `xml:"Type,attr"`
	CveUrl  []CveUrl `xml:"URL,omitempty"`
}

type CveUrl struct {
	XMLName xml.Name `xml:"URL,omitempty"`
	Url     string   `xml:",innerxml"`
}

type ProductTree struct {
	XMLName         xml.Name          `xml:"ProductTree,omitempty"`
	Xmlns           string            `xml:"xmlns,attr"`
	OpenEulerBranch []OpenEulerBranch `xml:"Branch,omitempty"`
}

type OpenEulerBranch struct {
	XMLName         xml.Name          `xml:"Branch,omitempty"`
	Type            string            `xml:"Type,attr"`
	Name            string            `xml:"Name,attr"`
	FullProductName []FullProductName `xml:"FullProductName,omitempty"`
}

type FullProductName struct {
	XMLName         xml.Name `xml:"FullProductName,omitempty"`
	ProductId       string   `xml:"ProductID,attr"`
	Cpe             string   `xml:"CPE,attr"`
	FullProductName string   `xml:",innerxml"`
}

type Vulnerability struct {
	XMLName         xml.Name        `xml:"Vulnerability,omitempty"`
	Ordinal         string          `xml:"Ordinal,attr"`
	Xmlns           string          `xml:"xmlns,attr"`
	CveNotes        CveNotes        `xml:"Notes,omitempty"`
	ReleaseDate     string          `xml:"ReleaseDate"`
	Bug             string          `xml:"Bug"`
	ProductStatuses ProductStatuses `xml:"ProductStatuses,omitempty"`
	Threats         Threats         `xml:"Threats,omitempty"`
	Remediations    Remediations    `xml:"Remediations,omitempty"`
}

type CveNotes struct {
	XMLName xml.Name `xml:"Notes,omitempty"`
	CveNote CveNote  `xml:"Note,omitempty"`
}

type CveNote struct {
	XMLName xml.Name `xml:"Note,omitempty"`
	Title   string   `xml:"Title,attr"`
	Type    string   `xml:"Type,attr"`
	Ordinal string   `xml:"Ordinal,attr"`
	XmlLang string   `xml:"xml:lang,attr"`
	Note    string   `xml:",innerxml"`
}

type ProductStatuses struct {
	XMLName xml.Name `xml:"ProductStatuses,omitempty"`
	Status  Status   `xml:"Status,omitempty"`
}

type Status struct {
	XMLName   xml.Name    `xml:"Status,omitempty"`
	Type      string      `xml:"Type,attr"`
	ProductId []ProductId `xml:"ProductID,omitempty"`
}

type ProductId struct {
	XMLName   xml.Name `xml:"ProductID,omitempty"`
	ProductId string   `xml:",innerxml"`
}

type Threats struct {
	XMLName xml.Name `xml:"Threats,omitempty"`
	Threat  Threat   `xml:"Threat,omitempty"`
}

type Threat struct {
	XMLName     xml.Name `xml:"Threat,omitempty"`
	Type        string   `xml:"Type,attr"`
	Description string   `xml:"Description"`
}

type Remediations struct {
	XMLName     xml.Name    `xml:"Remediations,omitempty"`
	Remediation Remediation `xml:"Remediation,omitempty"`
}

type Remediation struct {
	XMLName     xml.Name `xml:"Remediation,omitempty"`
	Type        string   `xml:"Type,attr"`
	Description string   `xml:"Description"`
	Date        string   `xml:"DATE"`
	Url         string   `xml:"URL"`
}
