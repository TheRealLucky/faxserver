package sender

import (
	"../loader"
	log "../logger"
	"database/sql"
	"fmt"
	"github.com/nu7hatch/gouuid"
	"github.com/pkg/errors"
)

func send_fax(db *sql.DB, acc_info loader.Account_Informations, domain_settings map[string]map[string]map[string][]string, tif_file string, fax_numbers []string) error {
	log.Info("sending fax")
	var mailto_address_user, mailto_address_fax, fax_prefix, mailto_address string
	mailfrom_address := ""
	var err error

	user_uuid, domain_name, err := loader.GetUserUuid(db, acc_info.Domain_uuid.String, acc_info.Fax_uuid)
	if err != nil {
		log.Error("failed to get user uuid: %v", err)
		return err
	}

	//TODO: if superadmin or admin -> php code
	fax_info, err := loader.Get_assigned_fax_extensions(db, acc_info.Domain_uuid.String, acc_info.Fax_uuid, user_uuid)
	if err != nil {
		log.Error("failed to get assigned fax extensions from database")
		return err
	}

	if domain_settings["fax"]["smtp_from"]["var"] != nil {
		mailfrom_address = domain_settings["fax"]["smtp_from"]["var"][0]
	} else {
		mailfrom_address = domain_settings["email"]["smtp_from"]["var"][0]
	}

	mailto_address_fax, fax_prefix, err = loader.GetFaxAddressAndPrefix(db, acc_info.Fax_uuid)
	if err != nil {
		log.Error("failed to get mailto_address and fax_prefix from database")
		return err
	}

	mailto_address_user, err = loader.Get_address_user(db, user_uuid)
	if err != nil {
		log.Error("failed to get mailto_address_user from database")
	}

	if mailto_address_fax != "" && mailto_address_user != mailto_address_fax {
		mailto_address = mailto_address_fax + "," + mailto_address_user
	} else {
		mailto_address = mailto_address_user
	}
	fmt.Println(fax_prefix)
	fmt.Println(mailto_address)
	fmt.Println(mailfrom_address)

	//create dial string
	dial_string := fmt.Sprintf("for_fax=1, accountcode=%s, ship_h_X-accountcode=%s, domain_uuid=%s, " +
		"domain_name=%s, origination_caller_id_name=%s, origination_caller_id_number=%s, fax_ident=%s, " +
		"fax_header=%s, fax_file=%s,", fax_info.Accountcode, fax_info.Accountcode, acc_info.Domain_uuid.String,
		domain_name, fax_info.Fax_caller_id_name, fax_info.Fax_caller_id_number, fax_info.Fax_caller_id_number,
		fax_info.Fax_caller_id_name, tif_file)

	for _, fax_number := range fax_numbers {
		//TODO: fax_split_dtmf
		fax_dtmf := ""
		routearray, err := OutboundRouteToBridge(db, acc_info.Domain_uuid.String, fax_number)
		if err != nil {
			return errors.Errorf("[send_fax] failed to call outbound route to bridge")
		}
		var fax_uri, fax_variables string
		if len(routearray) == 0 {
			fax_uri = "users/" + fax_number + "@" + domain_name
			fax_variables = ""
		} else {
			fax_uri = routearray[0]
			fax_variables = ""
			for _, element := range domain_settings["fax"]["variale"][""] {
				fax_variables += element

			}

		}
		wav_file := ""
		err = faxEnqueue(db, acc_info.Fax_uuid, tif_file, wav_file, mailto_address, fax_uri, fax_dtmf, dial_string)
		if err != nil {
			log.Error("failed to enqueue fax")
			return err
		}
	}
	return nil
}

func faxEnqueue(db *sql.DB, faxUUID, tif, wav, mailToAddress, faxURI, faxDTMF, dialString string) error {
	taskUUID, _ := uuid.NewV4()

	dialString = dialString + "fax_task_uuid='" + taskUUID.String() + "',"

	err := insertIntoFaxQueue(db, taskUUID.String(), faxUUID, tif, wav, faxURI, dialString, faxDTMF, mailToAddress)
	if err != nil {
		log.Error("failed to insert task into database")
		return err
	}

	return nil
}