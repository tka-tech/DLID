package dlidparser

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

const ColoradoIssuerId string = "636020"
const ConnecticutIssuerId string = "636006"
const IllinoisIssuerId string = "636035"
const MassachusettsIssuerId string = "636002"
const SouthCarolinaIssuerId string = "636005"
const TennesseeIssuerId string = "636053"

func parseV1(data string, issuer string) (license *DLIDLicense, err error) {

	start, end, err := dataRangeV1(data)

	if issuer == IllinoisIssuerId {

		// Illinois are the worst offenders so far in terms of mangling the DLID
		// spec.  They store name, licence number, expiry date and date of birth
		// as expected, but then go all-out crazy and encrypt everything else.
		// This means that the data range exceeds the size of the licence data
		// string.  We have to treat Illinois as a special case.
		end = len(data) - 1
	}

	if end >= len(data) {
		err = errors.New("Payload location does not exist in data")
	}

	payload := data[start:end]

	if err != nil {
		return
	}

	license, err = parseDataV1(payload, issuer)

	if err != nil {
		return
	}

	return
}

func dataRangeV1(data string) (start int, end int, err error) {

	start, err = strconv.Atoi(data[21:25])

	if err != nil {
		err = errors.New("Data contains malformed payload location")
		return
	}

	end, err = strconv.Atoi(data[25:29])

	if err != nil {
		err = errors.New("Data contains malformed payload length")
		return
	}

	end += start

	return
}

func parseDataV1(licenceData string, issuer string) (license *DLIDLicense, err error) {

	// Version 1 of the DLID card spec was published in 2000.  As of 2012, it is
	// the version used in Colorado.

	license = new(DLIDLicense)
	// We want to strip off the "DL" chunk identifier, but every other state has
	// managed to screw this up too.  Rather than handle this on a
	// state-by-state basis, we'll check to see what's at the target location
	// and handle it appropriately.
	if strings.HasPrefix(licenceData, "DL") || strings.HasPrefix(licenceData, "ID") {
		if strings.HasPrefix(licenceData, "DL") {
			license.SetDocumentType("DL")
		}

		if strings.HasPrefix(licenceData, "ID") {
			license.SetDocumentType("ID")
		}
		// POMG!  They actually got it right!
		licenceData = licenceData[2:]
	} else if strings.HasPrefix(licenceData, "L") {
		license.SetDocumentType("DL")
		// Either the guys in South Carolina can't count or they don't consider
		// the "DL" header part of the licence data.  In either case, their
		// offset is off by one.
		licenceData = licenceData[1:]
	} else {

		// Honestly, the spec really isn't that hard to follow.  I have no idea
		// why just about every implementation gets it wrong.  Massachusetts,
		// Connecticut and Pennsylvania don't include the "DL" chunk header in
		// at least some of their licenses.
		//
		// This else block is here just so I can grumble about badly-implemented
		// specs.
	}

	components := strings.Split(licenceData, "\n")

	license.SetIssuerId(issuer)
	license.SetIssuerName(issuers[issuer])

	// Country is always USA for V1 licenses
	license.SetCountry("USA")

	for component := range components {

		if len(components[component]) < 3 {
			continue
		}

		identifier := components[component][0:3]
		data := components[component][3:]

		data = strings.Trim(data, " ")

		switch identifier {
		case "DAR":
			license.SetVehicleClass(data)

		case "DAS":
			license.SetRestrictionCodes(data)

		case "DAT":
			license.SetEndorsementCodes(data)

		case "DAA":

			// Early versions of the Colorado implementation screwed up the
			// delimiter - they use a space instead of the specified comma.

			separator := " "

			if strings.Index(data, separator) == -1 {
				separator = ","
			}

			names := strings.Split(data, separator)

			// According to the spec, names are ordered LAST,FIRST,MIDDLE.
			// However, the geniuses in the Colorado and Tennessee DMVs order it
			// FIRST,MIDDLE,LAST.  We'll use the issuer ID number to
			// identify Colorado and adjust appropriately.  Issuer IDs can
			// be found here:
			//
			// http://www.aamva.org/IIN-and-RID/

			if issuer == ColoradoIssuerId || issuer == TennesseeIssuerId {

				// Colorado's backwards formatting style...
				license.SetFirstName(names[0])

				if len(names) > 2 {
					license.SetMiddleNames(names[1 : len(names)-1])
					license.SetLastName(names[len(names)-1])
				} else if len(names) > 1 {
					license.SetLastName(names[1])
				}
			} else {

				// Everyone else, hopefully.
				license.SetLastName(names[0])

				if len(names) > 1 {
					license.SetFirstName(names[1])

					if len(names) > 2 {
						license.SetMiddleNames(names[2:])
					}
				}
			}

		case "DAE":
			license.SetNameSuffix(data)

		case "DAL":

			// Colorado screws up again: they omit the *required* DAG field and
			// substitute the optional DAL field in older licences.
			fallthrough

		case "DAG":
			license.SetStreet(data)

		case "DAN":

			// Again, old Colorado licences ignore the spec.
			fallthrough

		case "DAI":
			license.SetCity(data)

		case "DAO":

			// Colorado strikes again.  Honestly, what is the point in having a
			// spec if you don't follow it?
			fallthrough

		case "DAJ":
			license.SetState(data)

		case "DAP":
			// More Colorado shenanigans.
			fallthrough

		case "DAK":

			// Colorado uses the 5-digit zip code.  South Carolina uses the
			// 5 digit zip code plus the +4 extension all smooshed together
			// into one long string.  Massachusetts uses the 5 digit zip
			// plus the +4 extension separated by "-".  The zip is
			// apparently never written like that and always uses "+" as a
			// separator.  Who knows what other states managed to
			// accomplish.  At this point your dedicated programmer admits
			// defeat in trying to untangle the incredible mess implemented
			// in this single field; we'll just show the zip as it is
			// stored.
			license.SetPostal(strings.Trim(data, " "))

		case "DAQ":
			license.SetCustomerId(data)

		case "DBA":
			license.SetExpiryDate(parseDateV1(data))

		case "DBB":
			license.SetDateOfBirth(parseDateV1(data))

		case "DBC":

			// Sex can be stored as M/F if it uses the DLID code.  It could
			// also be stored as 0/1/2/9 if it uses the ANSI D-20 codes,
			// available here:
			//
			// http://www.aamva.org/ANSI-D20-Standard-for-Traffic-Records-Systems/

			switch data {
			case "M":
				fallthrough
			case "1":
				license.SetSex(DriverSexMale)
			case "F":
				fallthrough
			case "2":
				license.SetSex(DriverSexFemale)
			default:
				license.SetSex(DriverSexNone)
			}

		case "DBD":
			license.SetIssueDate(parseDateV1(data))

		case "DBK":

			// Optional and probably not available
			license.SetSocialSecurityNumber(data)
		}
	}

	return
}

func parseDateV1(data string) time.Time {

	year, err := strconv.Atoi(data[:4])

	if err != nil {
		return time.Unix(0, 0)
	}

	month, err := strconv.Atoi(data[4:6])

	if err != nil {
		return time.Unix(0, 0)
	}

	day, err := strconv.Atoi(data[6:8])

	if err != nil {
		return time.Unix(0, 0)
	}

	location, err := time.LoadLocation("UTC")

	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, location)
}
