package sender

import (
	log "../logger"
	"database/sql"
	"regexp"
	"strings"
)

func OutboundRouteToBridge(db *sql.DB, domainUUID, destinationNumber string) ([]string, error) {
	destinationNumber = strings.Trim(destinationNumber,  " \t\n\r")
	bridgeArray := make([]string, 0)

	// check for a valid destination number
	pattern := "^[\\*\\+0-9]*$"
	match, err := regexp.MatchString(pattern, destinationNumber)
	if err != nil {
		log.Error("failed to match destination number '%v' with regex pattern: %v")
		return nil, err
	}

	if match == false {
		log.Warn("destination number '%v' does not match regex pattern: %v", destinationNumber, pattern)
		return bridgeArray, nil
	}

	dialPlanInfos, err := GetDialPlans(db, domainUUID)
	if err != nil {
		log.Error("failed to retrieve dial plan information from database")
		return nil, err
	}

	regexMatches := make([]string, 5)

	for _, plan := range dialPlanInfos {
		regexMatch := false
		for _, detail := range plan.Info {
			if detail.DetailTag == "condition" && detail.DetailType == "destination_number" {
				rex := regexp.MustCompile(detail.DetailData)
				matches := rex.FindAllStringSubmatch(destinationNumber, -1)
				if match == true {
					regexMatch = true
					regexMatches[0] = matches[0][1]
					regexMatches[1] = matches[1][1]
					regexMatches[2] = matches[2][1]
					regexMatches[3] = matches[3][1]
					regexMatches[4] = matches[4][1]
				}else {
					regexMatch = false
				}
			}
		}
		if regexMatch == true {
			for _, detail := range plan.Info {
				if detail.DetailTag == "action" && detail.DetailType == "bridge" && detail.DetailData != "${enum_auto_route}" {
					detailData := strings.Replace(detail.DetailData, "$1", regexMatches[0], -1)
					detailData = strings.Replace(detail.DetailData, "$2", regexMatches[1], -1)
					detailData = strings.Replace(detail.DetailData, "$3", regexMatches[3], -1)
					detailData = strings.Replace(detail.DetailData, "$4", regexMatches[4], -1)
					detailData = strings.Replace(detail.DetailData, "$5", regexMatches[1], -1)
					bridgeArray = append(bridgeArray, detailData)
					if plan.Continue == "false" {
						return bridgeArray, nil
					}
				}
			}
		}
	}
	return bridgeArray, nil
}