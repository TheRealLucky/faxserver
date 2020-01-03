package sender

import (
	log "../logger"
	"database/sql"
	"github.com/pkg/errors"
)

type DialPlanInfo struct {
	UUID		string
	DetailTag	string
	DetailType	string
	Continue	string
}

type DialPlanDetailsInfo struct {
	DetailsTag  string
	DetailsType string
	DetailsData string
}

func GetDialPlanInfo(db *sql.DB, domainUUID string) ([]DialPlanInfo, error) {
	var result []DialPlanInfo

	query := "select dialplan_uuid, dialplan_detail_tag, dialplan_detail_type, dialplan_continue " +
		"       from v_dialplans " +
		"      where (domain_uuid = $1 or domain_uuid is null) " +
		"        and app_uuid='8c914ec3-9fc0-8ab5-4cda-6c9288bdc9a3' " +
		"        and dialplan_enabled = 'true' " +
		"      order by dialplan_order asc"

	res, err := db.Query(query, domainUUID)
	if err != nil {
		log.Error("failed to perform query")
		return nil, errors.Errorf("[outbound_route_to_bridge] failed to get dialplan infos \n%v", err)
	}

	for res.Next() {
		info := DialPlanInfo{}
		err = res.Scan(&info.UUID, &info.DetailTag, &info.DetailType, & info.Continue)
		if err != nil {
			log.Error("failed to scan result set")
			return nil, err
		}
		result = append(result, info)
	}
	return result, nil
}