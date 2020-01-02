package loader

import (
	log "../logger"
	"database/sql"
	"github.com/pkg/errors"
	"gopkg.in/guregu/null.v3"
)

type Account_Informations struct {
	Fax_uuid                              string
	Domain_uuid                           null.String
	Fax_extension                         null.String
	Fax_email                             null.String
	Fax_pin_number                        null.String
	Fax_caller_id_name                    null.String
	Fax_caller_id_number                  null.String
	Fax_email_connection_type             null.String
	Fax_email_connection_host             null.String
	Fax_email_connection_port             null.String
	Fax_email_connection_security         null.String
	Fax_email_connection_validate         null.String
	Fax_email_connection_username         null.String
	Fax_email_connection_password         null.String
	Fax_email_connection_mailibox         null.String
	Fax_email_outboud_subject_tag         null.String
	Fax_email_outbound_authorized_senders null.String
	Fax_send_greeting                     null.String
}

type Fax_Info struct {
	Fax_uuid             string
	Fax_extension        string
	Fax_caller_id_name   string
	Fax_caller_id_number string
	Accountcode          string
	Fax_send_greetings   string
}


func Load_default_settings(db *sql.DB) (map[string]map[string]map[string][]string, error) {
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

func Load_domain_settings(db *sql.DB, default_settings map[string]map[string]map[string][]string, domain_uuid string) (map[string]map[string]map[string][]string, error) {
	log.Info("loading domain settings")
	query := "select domain_setting_name, domain_setting_category, domain_setting_subcategory, domain_setting_value " +
		"                          from v_domain_settings " +
		"                         where domain_setting_enabled ='true' " +
		"                           and domain_uuid=$1"

	res, err := db.Query(query, domain_uuid)
	if err != nil {
		return nil, errors.Errorf("[load_domain_settings] failed to execute query \n%v", err)
	}

	for res.Next() {
		var name string
		var category string
		var subcategory string
		var value string
		err = res.Scan(&name, &category, &subcategory, &value)
		if err != nil {
			return nil, err
		}

		if default_settings == nil {
			default_settings = make(map[string]map[string]map[string][]string)
		}
		if default_settings[category] == nil {
			default_settings[category] = make(map[string]map[string][]string)
		}
		if default_settings[category][subcategory] == nil {
			default_settings[category][subcategory] = make(map[string][]string)
		}
		if len(subcategory) == 0 {
			if name == "array" {
				default_settings[category][""][""] = nil
				default_settings[category][""][""] = append(default_settings[category][""][""], value)
			} else {
				default_settings[category][""][name] = nil
				default_settings[category][""][name] = append(default_settings[category][""][name], value)
			}
		} else {
			if name == "array" {
				default_settings[category][subcategory][""] = nil
				default_settings[category][subcategory][""] = append(default_settings[category][subcategory][""], value)
			} else {
				default_settings[category][subcategory][name] = nil
				default_settings[category][subcategory][name] = append(default_settings[category][subcategory][name], value)
			}
		}
	}
	return default_settings, nil
}
