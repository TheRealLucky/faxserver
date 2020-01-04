package mailer

import (
	"crypto/tls"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
	loader "../loader"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/charset"
	"github.com/emersion/go-message/mail"
	Log "../logger"
	tiffer "../tiffer"

)

func Get_emails(acc loader.Account_Informations) (string, []string, error) {
	var fax_numbers []string
	var tif_file string
	if acc.Fax_email_connection_mailibox.String == "" {
		acc.Fax_email_connection_mailibox.String = "INBOX"
	}
	Log.Info("[get_emails] Connecting to server...")

	// Connect to server
	//TODO: check if there a opportunity to give mail typ as parameter (imap or pop)
	connection_host := acc.Fax_email_connection_host.String + ":" + acc.Fax_email_connection_port.String

	var c *client.Client
	var err error
	//connect with ssl or tls (depends on config)
	if acc.Fax_email_connection_security.String == "ssl" {
		tlsConn := tls.Config{}
		tlsConn.InsecureSkipVerify = true
		tlsConn.MinVersion = tls.VersionSSL30
		//tlsConn.MaxVersion = tls.VersionSSL30
		tlsConn.ServerName = acc.Fax_email_connection_host.String
		c, err = client.DialTLS(connection_host, &tlsConn)
	} else if acc.Fax_email_connection_security.String == "tls" {
		tlsConn := tls.Config{}
		tlsConn.InsecureSkipVerify = true
		tlsConn.MinVersion = tls.VersionTLS10
		tlsConn.MaxVersion = tls.VersionTLS13
		tlsConn.ServerName = acc.Fax_email_connection_host.String
		c, err = client.DialTLS(connection_host, &tlsConn)

	} else {
		c, err = client.Dial(connection_host)
	}
	if err != nil {
		return "", nil, errors.Errorf("[get_emails] failed to connect to host: \n%v", err)
	} else {
		Log.Info("[get_emails] connected to host")
	}

	// Don't forget to logout
	defer c.Logout()

	// Login
	if err := c.Login(acc.Fax_email_connection_username.String, acc.Fax_email_connection_password.String); err != nil {
		return "", nil, errors.Errorf("[get_emails] failed to login: \n%v", err)
	} else {
		Log.Info("[get_emails] Logged in")
	}

	mbox, err := c.Select(acc.Fax_email_connection_mailibox.String, false)
	if err != nil {
		return "", nil, errors.Errorf("[get_emails] failed to select deposit mailbox: \n%v", err)
	}

	if mbox.Messages == 0 {
		return "", nil, errors.Errorf("[get_emails] mailbox is empty")
	}

	seqSet := new(imap.SeqSet)
	seqSet.AddRange(1, mbox.Messages)
	//seqSet.AddNum(mbox.Messages)
	// Get the whole message body
	var section imap.BodySectionName
	items := []imap.FetchItem{section.FetchItem()}

	//fetch messages from mailbox
	messages := make(chan *imap.Message, 10)
	output_error := make(chan error)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		if err := c.Fetch(seqSet, items, messages); err != nil {
			output_error <- errors.Errorf("[get_emails] failed to fetch emails: %v", err)
		}
		output_error <- nil
		wg.Done()
	}()

	err = <-output_error
	if err != nil {
		return "", nil, err
	}
	wg.Wait()

	for msg := range messages {
		if msg == nil {
			Log.Fatal("[get_emails] Server didn't return message")
		}
		r := msg.GetBody(&section)
		if r == nil {
			Log.Fatal("[get_emails] Server didn't return message body")
		}

		//set charset var from the go-message package to support more charsets
		imap.CharsetReader = charset.Reader
		// Create a new mail reader
		mr, err := mail.CreateReader(r)
		if err != nil {
			return "", nil, errors.Errorf("[get_emails] failed to create mail reader: \n%v", err)
		}

		//information to set in fax
		informations_to_fax := make(map[string]string)

		// get some info about the message
		header := mr.Header
		if date, err := header.Date(); err == nil {
			log.Println("Date:", date)
		}
		if from, err := header.AddressList("From"); err == nil {
			if len(from) > 0 {
				if !check_if_authorized(acc.Fax_email_outbound_authorized_senders.String, from[0].Address) {
					log.Println(acc.Fax_email, acc.Fax_email_connection_host)
					Log.Error("[get_emails] not authorized to send a fax: \n%v", from[0].Address)
					continue
				}
				informations_to_fax["from"] = from[0].Address
				log.Println("From:", from)
			} else {
				continue
			}

		}
		if to, err := header.AddressList("To"); err == nil {
			informations_to_fax["to"] = to[0].Address
			log.Println("To:", to)
		}
		//TODO: check subject prefix (col in database)
		if subject, err := header.Subject(); err == nil {
			log.Println("Subject:", subject)
			//hope this is right way how to get fax_numbers
			fax_numbers = strings.Split(subject, ",")
		}

		//map to merge pdf's into one
		var filenames map[string]string
		//array which files have to be deleted after process
		var files_to_delete []string

		// Process each message's part
		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				break
			} else if err != nil {
				return "", nil, errors.Errorf("[get_emails] failed to extract mail")
			}

			switch h := p.Header.(type) {
			case *mail.InlineHeader:
				// This is the message's text (can be plain-text or HTML)
				b, _ := ioutil.ReadAll(p.Body)
				if informations_to_fax["message"] == "" {
					informations_to_fax["message"] = string(b)
				}
				//log.Println("Got text: ", string(b))
			case *mail.AttachmentHeader:
				// This is an attachment
				filename, _ := h.Filename()
				log.Println("Got attachment: ", filename)

				//create folder
				path := tiffer.Create_folder(acc.Fax_email_connection_host.String, acc.Fax_email_connection_username.String)
				tmp_path := path + "/" + filename

				//TODO: tif/tiff too?
				if strings.ToLower(filename[len(filename)-4:]) == ".pdf" {
					//copy attachment to folder
					out, err := os.Create(tmp_path)
					if err != nil {
						return "", nil, errors.Errorf("[get_emails] failed to create new folder: \n%v", err)
					}
					defer out.Close()
					//copy file to right folder
					_, err = io.Copy(out, p.Body)
					if err != nil {
						return "", nil, errors.Errorf("[get_emails] failed to copy pdf into new folder: \n%v", err)
					}
					if filenames == nil {
						filenames = make(map[string]string)
					}
					filenames[filename] = path
				}
			}
		}

		//create pdf from message with footer (without attachments)
		created_pdf, err := tiffer.Create_pdf(informations_to_fax, acc) //filenames,
		if err != nil {
			return "", nil, errors.Errorf("[get_emails] failed to create pdf: \n%v", err)
		} else {

			right_order := make(map[string][]string)
			path := created_pdf[0] + "/" + created_pdf[1]
			if filenames != nil {
				files_to_delete = append(files_to_delete, path)
				right_order[created_pdf[0]] = append(right_order[created_pdf[0]], created_pdf[1])
				for name, _ := range filenames {
					right_order[created_pdf[0]] = append(right_order[created_pdf[0]], name)
				}
				//create pdf with attachment and from this pdf create a tif file
				tif_file, err = tiffer.Merge_pdf(right_order)
				if err != nil {
					return "", nil, errors.Errorf("[get_emails] failed to create merged pdf: \n%v", err)
				}
			} else {
				//create a tif file such if there are no attachments
				tif_file, err = tiffer.Create_tif(path)
				if err != nil {
					return "", nil, errors.Errorf("[get_emails] failed to create tif file from pdf: \n%v", err)
				}
			}

		}
		//delete tmp files (from create_pdf)
		err = delete_paths(files_to_delete)
		if err != nil {
			return "", nil, errors.Errorf("[get_mails] failed to delte file; \n%v", err)
		}


	}
	return tif_file, fax_numbers, nil
}

//check if person who would like to send a fax is authorized
func check_if_authorized(str string, wish_sender string) bool {
	senders := strings.Split(str, ",")
	for _, i := range senders {
		if i == wish_sender {
			return true
		}
	}
	return false
}

//delete files from path
func delete_paths(paths []string) error {
	for i := 0; i < len(paths); i++ {
		err := os.Remove(paths[i])
		if err != nil {
			return errors.Errorf("[delete_paths] failed to delete file: \n%v", err)
		}
		fmt.Println("is deleted: ", paths[i])
	}
	return nil
}
