package sender


import (
	log "../logger"
	"database/sql"
)

type DialPlan struct {
	UUID		string
	Continue	string
	Info		[]DialPlanInfo
}

type DialPlanInfo struct {
	DetailTag	string
	DetailType	string
	DetailData  string
}

func GetDialPlans(db *sql.DB, domainUUID string) ([]DialPlan, error) {
	result := make([]DialPlan, 0)

	query := "select dialplan_uuid, dialplan_continue" +
		"       from v_dialplans" +
		"      where (domain_uuid = $1 or domain_uuid is null)" +
		"        and app_uuid='8c914ec3-9fc0-8ab5-4cda-6c9288bdc9a3'" +
		"        and dialplan_enabled = 'true'" +
		"      order by dialplan_order asc"

	res, err := db.Query(query, domainUUID)
	if err != nil {
		log.Error("failed to execute query, domain uuid was: %v", domainUUID)
		return nil, err
	}

	for res.Next() {
		plan := DialPlan{}
		plan.Info = make([]DialPlanInfo, 0)

		err = res.Scan(&plan.UUID, &plan.Continue)
		if err != nil {
			log.Error("failed to scan result set")
			return nil, err
		}

		query = "select dialplan_detail_tag, dialplan_detail_type, dialplan_detail_data" +
			"      from v_dialplan_details" +
			"     where dialplan_uuid = $1"

		newRes, err := db.Query(query, plan.UUID)
		if err != nil {
			log.Error("failed to execute query, planUUID was: %v", plan.UUID)
			return nil, err
		}

		for newRes.Next() {
			info := DialPlanInfo{}
			err = newRes.Scan(&info.DetailTag, &info.DetailType, &info.DetailData)
			if err != nil {
				log.Error("failed to scan result set")
				return nil, err
			}
			plan.Info = append(plan.Info, info)
		}
		result = append(result, plan)
	}
	return result, nil
}