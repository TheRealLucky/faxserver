package loader

import "gopkg.in/guregu/null.v3"

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
