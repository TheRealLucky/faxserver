package loader

import (
	log "../logger"
	"database/sql"
)

func Get_account_informations(db *sql.DB) ([]Account_Informations, error) {
	query := "select fax_uuid, domain_uuid, fax_extension, fax_email, fax_pin_number, fax_caller_id_name, " +
		"            fax_caller_id_number, fax_email_connection_type, fax_email_connection_host, " +
		"            fax_email_connection_port, fax_email_connection_security, fax_email_connection_validate, " +
		"            fax_email_connection_username, fax_email_connection_password, fax_email_connection_mailbox, " +
		"            fax_email_outbound_subject_tag, fax_email_outbound_authorized_senders, fax_send_greeting " +
		"       from v_fax " +
		"      where fax_email_connection_host != '' " +
		"        and fax_email_connection_host is not null " +
		"        and fax_email_connection_type = 'imap'"

	res, err := db.Query(query)
	if err != nil {
		log.Error("failed to execute query: %v", err)
		return nil, err
	}
	var acc_list []Account_Informations
	for res.Next() {
		acc := Account_Informations{}
		err := res.Scan(&acc.Fax_uuid, &acc.Domain_uuid, &acc.Fax_extension, &acc.Fax_email, &acc.Fax_pin_number, &acc.Fax_caller_id_name, &acc.Fax_caller_id_number, &acc.Fax_email_connection_type, &acc.Fax_email_connection_host, &acc.Fax_email_connection_port, &acc.Fax_email_connection_security, &acc.Fax_email_connection_validate, &acc.Fax_email_connection_username, &acc.Fax_email_connection_password, &acc.Fax_email_connection_mailibox, &acc.Fax_email_outboud_subject_tag, &acc.Fax_email_outbound_authorized_senders, &acc.Fax_send_greeting)
		if err != nil {
			log.Error("failed to scan result set: %v", err)
			return nil, err
		}
		acc_list = append(acc_list, acc)
	}
	return acc_list, nil
}


func Get_user_uuid(db *sql.DB, domain_uuid string, fax_uuid string) (string, string, error) {
	query := "select user_uuid, domain_name " +
		"       from v_fax_users as fu, v_domains as d " +
		"      where fu.domain_uuid = d.domain_uuid " +
		"        and fu.fax_uuid = $1 " +
		"        and fu.domain_uuid = $2"

	res := db.QueryRow(query, domain_uuid, fax_uuid)

	var user_uuid, domain_name string
	err := res.Scan(&user_uuid, &domain_name)
	if err != nil {
		log.Error("failed to scan result set: %v", err)
		return "", "", err
	} else {
		return user_uuid, domain_name, nil
	}
}

func Get_fax_address_and_prefix(db *sql.DB, fax_uuid string) (string, string, error) {
	query := "select fax_email, fax_prefix " +
		"       from v_fax " +
		"      where fax_uuid = $1"

	res := db.QueryRow(query, fax_uuid)
	var fax_email, fax_prefix string
	err := res.Scan(&fax_email, &fax_prefix)
	if err != nil {
		log.Error("failed to scan result set")
		return "", "", err
	}
	return fax_email, fax_prefix, nil
}

func Get_address_user(db *sql.DB, user_uuid string) (string, error) {
	query := "select contact_uuid " +
		"       from v_users " +
		"      where user_uuid = $1"

	res := db.QueryRow(query, user_uuid)
	var contact_uuid string
	err := res.Scan(&contact_uuid)
	if err != nil {
		return "", nil
	}
	query = "select email_address " +
		"      from v_contact_emails " +
		"     where contact_uuid = $1 " +
		"     order by email_primary desc;"
	res = db.QueryRow(query, contact_uuid)
	var email_address string
	err = res.Scan(&email_address)
	if err != nil {
		return "", err
	}
	return email_address, nil
}



