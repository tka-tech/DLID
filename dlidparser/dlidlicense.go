package dlidparser

import (
	"time"
)

type DriverSex int

const (
	DriverSexNone DriverSex = iota
	DriverSexMale
	DriverSexFemale
)

type DLIDLicense struct {
	firstName             string
	middleNames           []string
	lastName              string
	nameSuffix            string
	street                string
	city                  string
	state                 string
	country               string
	postal                string
	sex                   DriverSex
	socialSecurityNumber  string
	dateOfBirth           time.Time
	issuerId              string
	documentType          string
	issuerName            string
	expiryDate            time.Time
	issueDate             time.Time
	vehicleClass          string
	restrictionCodes      string
	endorsementCodes      string
	customerId            string
	documentDiscriminator string
}

func (d *DLIDLicense) SetFirstName(s string) {
	d.firstName = s
}

func (d *DLIDLicense) FirstName() string {
	return d.firstName
}

func (d *DLIDLicense) SetMiddleNames(s []string) {
	d.middleNames = s
}

func (d *DLIDLicense) MiddleNames() []string {
	return d.middleNames
}

func (d *DLIDLicense) SetLastName(s string) {
	d.lastName = s
}

func (d *DLIDLicense) LastName() string {
	return d.lastName
}

func (d *DLIDLicense) SetNameSuffix(s string) {
	d.nameSuffix = s
}

func (d *DLIDLicense) NameSuffix() string {
	return d.nameSuffix
}

func (d *DLIDLicense) SetStreet(s string) {
	d.street = s
}

func (d *DLIDLicense) Street() string {
	return d.street
}

func (d *DLIDLicense) SetCity(s string) {
	d.city = s
}

func (d *DLIDLicense) City() string {
	return d.city
}

func (d *DLIDLicense) SetState(s string) {
	d.state = s
}

func (d *DLIDLicense) State() string {
	return d.state
}

func (d *DLIDLicense) SetCountry(s string) {
	d.country = s
}

func (d *DLIDLicense) Country() string {
	return d.country
}

func (d *DLIDLicense) SetPostal(s string) {
	d.postal = s
}

func (d *DLIDLicense) Postal() string {
	return d.postal
}

func (d *DLIDLicense) SetSex(s DriverSex) {
	d.sex = s
}

func (d *DLIDLicense) Sex() DriverSex {
	return d.sex
}

func (d *DLIDLicense) SetSocialSecurityNumber(s string) {
	d.socialSecurityNumber = s
}

func (d *DLIDLicense) SocialSecurityNumber() string {
	return d.socialSecurityNumber
}

func (d *DLIDLicense) SetDateOfBirth(t time.Time) {
	d.dateOfBirth = t
}

func (d *DLIDLicense) DateOfBirth() time.Time {
	return d.dateOfBirth
}

func (d *DLIDLicense) IssuerId() string {
	return d.issuerId
}

func (d *DLIDLicense) SetDocumentType(s string) {
	d.documentType = s
}

func (d *DLIDLicense) DocumentType() string {
	return d.documentType
}

func (d *DLIDLicense) SetIssuerId(s string) {
	d.issuerId = s
}

func (d *DLIDLicense) IssuerName() string {
	return d.issuerName
}

func (d *DLIDLicense) SetIssuerName(s string) {
	d.issuerName = s
}

func (d *DLIDLicense) VehicleClass() string {
	return d.vehicleClass
}

func (d *DLIDLicense) SetVehicleClass(s string) {
	d.vehicleClass = s
}

func (d *DLIDLicense) RestrictionCodes() string {
	return d.restrictionCodes
}

func (d *DLIDLicense) SetRestrictionCodes(s string) {
	d.restrictionCodes = s
}

func (d *DLIDLicense) EndorsementCodes() string {
	return d.endorsementCodes
}

func (d *DLIDLicense) SetEndorsementCodes(s string) {
	d.endorsementCodes = s
}

func (d *DLIDLicense) CustomerId() string {
	return d.customerId
}

func (d *DLIDLicense) SetCustomerId(s string) {
	d.customerId = s
}

func (d *DLIDLicense) SetExpiryDate(t time.Time) {
	d.expiryDate = t
}

func (d *DLIDLicense) ExpiryDate() time.Time {
	return d.expiryDate
}

func (d *DLIDLicense) SetIssueDate(t time.Time) {
	d.issueDate = t
}

func (d *DLIDLicense) IssueDate() time.Time {
	return d.issueDate
}
