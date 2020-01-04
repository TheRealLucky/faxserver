package sender

import (
	"database/sql"
	log "../logger"
)



func insertIntoFaxQueue(db *sql.DB, taskUUID, faxUUID, faxFile, wavFile, faxURI, dialString, faxDTMF, replyAddress string) error {
	statement := "INSERT INTO v_fax_tasks( fax_task_uuid, fax_uuid, task_next_time, task_lock_time, " +
		"             task_fax_file, task_wav_file, task_uri, task_dial_string, task_dtmf, task_interrupted, " +
		"             task_status, task_no_answer_counter, task_no_answer_retry_counter, task_retry_counter, " +
		"             task_reply_address, task_description) " +
		"              VALUES ($1, $2, NOW() at time zone 'utc', NULL, $3, $4, $5, $6, $7, 'false', 0, 0, 0, 0, $8, '');"

	_, err := db.Exec(statement, taskUUID, faxUUID, faxFile, wavFile, faxURI, dialString, faxDTMF, replyAddress)
	if err != nil {
		log.Error("failed to perform insert statement")
		return err
	}
	return nil

}