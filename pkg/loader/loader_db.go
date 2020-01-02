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
