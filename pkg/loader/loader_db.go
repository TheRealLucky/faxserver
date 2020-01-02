package loader

import (
	log "../logger"
	"database/sql"
	"github.com/pkg/errors"
)

func getDefaultSettingsDB(db *sql.DB) (map[string]map[string]map[string][]string, error) {
	query := "select default_setting_name, default_setting_category, default_setting_subcategory, default_setting_value " +
		"       from v_default_settings " +
		"      where default_setting_enabled='true';"

	log.Info("executing query on database")
	res, err := db.Query(query)
	if err != nil {
		return nil, errors.Errorf("failed to execute query: ", err)
	}

	result := make(map[string]map[string]map[string][]string)

	for res.Next() {
		var name string
		var category string
		var subcategory string
		var value string

		err = res.Scan(&name, &category, &subcategory, &value)
		if result[category] == nil {
			result[category] = make(map[string]map[string][]string)
		}
		if result[category][subcategory] == nil {
			result[category][subcategory] = make(map[string][]string)
		}
		if len(subcategory) == 0 {
			if name == "array" {
				if result[category] == nil {
					result[category] = make(map[string]map[string][]string)
					result[category][""][""] = append(result[category][""][""], value)
				}
			} else {
				result[category][""][name] = append(result[category][""][name], value)
			}
		} else {
			if name == "array" {
				result[category][subcategory][""] = append(result[category][subcategory][""], value)

			} else {
				result[category][subcategory][name] = append(result[category][subcategory][name], value)
			}
		}
	}
	return result, nil
}

func getDomainSettingsDB(db *sql.DB, ) {

}