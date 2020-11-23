package dlidparser

import (
	"errors"
	"strings"
)

func parseV4(data string, issuer string) (license *DLIDLicense, err error) {

	start, end, err := dataRangeV2(data)

	if end > len(data) {
		err = errors.New("Payload location does not exist in data")
	}

	payload := data[start:end]

	if err != nil {
		return
	}

	license, err = parseDataV4(payload, issuer)

	if err != nil {
		return
	}

	return
}

func parseDataV4(licenceData string, issuer string) (license *DLIDLicense, err error) {

	// Version 4 of the DLID card spec was published in 2009.

	if !strings.HasPrefix(licenceData, "DL") && !strings.HasPrefix(licenceData, "ID") {
		err = errors.New("Missing header in licence data chunk")
		return
	}

	license = new(DLIDLicense)

	if strings.HasPrefix(licenceData, "DL") {
		license.SetDocumentType("DL")
	}

	if strings.HasPrefix(licenceData, "ID") {
		license.SetDocumentType("ID")
	}

	licenceData = licenceData[2:]

	components := strings.Split(licenceData, "\n")

	license.SetIssuerId(issuer)
	license.SetIssuerName(issuers[issuer])

	var dateOfBirth string
	var expiryDate string
	var issueDate string

	for component := range components {

		if len(components[component]) < 3 {
			continue
		}

		identifier := components[component][0:3]
		data := components[component][3:]

		data = strings.Trim(data, " ")

		switch identifier {
		case "DCA":
			license.SetVehicleClass(data)

		case "DCB":
			license.SetRestrictionCodes(data)

		case "DCD":
			license.SetEndorsementCodes(data)

		case "DCS":
			license.SetLastName(data)

		case "DCU":
			license.SetNameSuffix(data)

		case "DAC":
			license.SetFirstName(data)

		case "DAD":
			names := strings.Split(data, ",")
			license.SetMiddleNames(names)

		case "DCG":
			license.SetCountry(data)

		case "DAG":
			license.SetStreet(data)

		case "DAI":
			license.SetCity(data)

		case "DAJ":
			license.SetState(data)

		case "DAK":
			license.SetPostal(data)

		case "DAQ":
			license.SetCustomerId(data)

		case "DBA":
			expiryDate = data

		case "DBB":
			dateOfBirth = data

		case "DBC":
			switch data {
			case "1":
				license.SetSex(DriverSexMale)
			case "2":
				license.SetSex(DriverSexFemale)
			default:
				license.SetSex(DriverSexNone)
			}

		case "DBD":
			issueDate = data
		}
	}

	if license.Country() == "" {
		license.SetCountry("USA") // set to USA if empty
	}

	// At this point we should know the country and the postal code (both are
	// mandatory fields) so we can undo the desperate mess the standards body
	// made of the postal code field.

	if license.Country() == "USA" && len(license.Postal()) > 0 {

		// Another change to the postal code field!  Surprise!  This time the
		// standards guys trimmed the field down to 9 characters, which makes
		// sense because US zip codes are only 9 digits long.  Canadian post
		// codes are only 6 characters.  Why was the original spec 11 digits?
		// Because the standards guys are *nuts*.
		//
		// We will extract the 5-digit zip and the +4 section.  If the +4 is all
		// zeros we can discard it.

		if len(license.Postal()) > 5 {
			zip := license.Postal()[:5]
			plus4 := license.Postal()[5:9]

			if plus4 == "0000" {
				license.SetPostal(zip)
			} else {
				license.SetPostal(zip + "+" + plus4)
			}
		}
	}

	// Now we can parse the dates, too.
	if len(license.Country()) > 0 {
		license.SetDateOfBirth(parseDateV3(dateOfBirth, license.Country()))
		license.SetExpiryDate(parseDateV3(expiryDate, license.Country()))
		license.SetIssueDate(parseDateV3(issueDate, license.Country()))
	}

	return
}
