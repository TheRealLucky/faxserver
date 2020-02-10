package main

import (
	"../../pkg/config"
	"../../pkg/loader"
	log "../../pkg/logger"
	mailer "../../pkg/mailer"
	"context"
	"database/sql"
	"fmt"
	uuid "github.com/nu7hatch/gouuid"
	"golang.org/x/sync/semaphore"
	"time"
	sender "../../pkg/sender"
)

var Config *config.Configuration
var Database *sql.DB

func main() {
	log.Info("starting application")
	Config, err := config.InitConfig("../../conf.yaml")
	if err != nil {
		panic(err)
	}

	Database, err := config.InitDatabase(Config)
	if err != nil {
		panic(err)
	}

	sem := semaphore.NewWeighted(10)
	ctx := context.TODO()


	var start_time time.Time
	var end_time time.Time

	running := true
	for running {
		start_time = time.Now()
		//get default_settings
		default_settings, err := loader.Load_default_settings(Database)
		if err != nil {
			fmt.Println(err)
		}

		fax_send_mode_default := default_settings["fax"]["send_mode"]["text"]
		if len(fax_send_mode_default) == 0 {
			fax_send_mode_default = append(fax_send_mode_default, "direct")
		}

		fax_allowed_extension_default := default_settings["fax"]["allowed_extension"]
		if len(fax_allowed_extension_default[""]) == 0 {
			val := []string{".pdf", ".tiff", "tif"}
			fax_allowed_extension_default[""] = val
		}

		//get accounts informations
		acc_list, err := loader.Get_account_informations(Database)
		if err != nil {
			log.Error("failed to get account informations: %v", err)
		} else {
			log.Info("got account informations list")
		}
		cnt := 1
		for _, acc_info := range acc_list {
			if err := sem.Acquire(ctx, 1); err != nil {
				log.Error("Failed to acquire semaphore: %v", err)
				break
			}
			go func() {
				cnt++
				uuid, _ := uuid.NewV4()
				fmt.Println("starttttt: ", uuid)
				//load domain_setting and replace it with default_settings
				domain_settings, err1 := loader.Load_domain_settings(Database, default_settings, acc_info.Domain_uuid.String)
				//str, _ := json.Marshal(domain_settings)
				//fmt.Println(string(str))
				if err1 != nil {
					log.Error("failed to load domain settings: %v",err1)
				}

				//	setVariables(e)
				//TODO: do i need this? in php code these variables are set but never used
				fax_send_mode := domain_settings["fax"]["conver_font"]["text"]
				if len(fax_send_mode) == 0 {
					fax_send_mode = fax_send_mode_default
				}
				fax_allowed_extension := domain_settings["fax"]["allowed_extension"]
				if len(fax_allowed_extension[""]) == 0 {
					fax_allowed_extension = fax_allowed_extension_default
				}
				//connect to host with user and password
				//fetch eamils from server
				//create pdf file and merge these to a tif file
				tif_file, fax_numbers, err := mailer.Get_emails(acc_info)
				fmt.Println(tif_file)
				fmt.Println(fax_numbers)
				if err != nil {
					log.Error("failed to get mails: %v",err)
				} else {
				err = sender.Send_fax(Database, acc_info, domain_settings, tif_file, fax_numbers)
					if err != nil {
						fmt.Println("send fax error")
						log.Error("failed to send fax: \n%v", err)
						panic(err)
					}
					fmt.Println("last one")

				}

				fmt.Println("done")
				if cnt == 10 {
					fmt.Println("it's 10, yeah")
					//	time.Sleep(time.Minute * 10)
				}
				defer sem.Release(1)
				//out[i] = collatzSteps(i + 1)
			}()

		}

		//wait until the interval is over
		fmt.Println("come to end")
		end_time = time.Now()
		diff := end_time.Sub(start_time)
		if int(diff.Seconds()) < Config.Interval {
			sleepy := Config.Interval - int(diff.Seconds())
			log.Info("sleep for %s seconds", string(sleepy))
			time.Sleep(time.Duration(sleepy) * time.Second)

		}

	}


	// testing loader package
	/*
	log.Info("loading default settings")
	default_settings, err := loader.Load_default_settings(Database)
	if err != nil {
		panic(err)
	}
	fmt.Println(default_settings)
	*/
	/*
	log.Info("loading account informations")
	account, err := loader.Get_account_informations(Database)
	if err != nil {
		panic(err)
	}
	fmt.Println(account)
	*/
	//log.Info("getting user uuid")
	//uuid, domain, err := loader.Get_user_uuid(Database, )

	//create debug folder with generated pdf's
	/*tiffer.Create_folder("x","u")
	m := make(map[string]string)
	m["from"]="KundeA"
	m["to"]="KundeB"
	m["subject"]="xsx"
	file, err := tiffer.Create_pdf(m,"x","u")
	if err != nil {
		fmt.Println("failed")
	}
	fmt.Println(file)
	*/

	//testcode on linux
	/*
	ka := make(map[string][]string)
	t := make([]string,1)
	l := make([]string,1)
	t = append(t,"d8a9bea9-640d-423e-5de0-f6f517d1d529.pdf")
	t = append(t, "f3d72fa4-9584-486c-4877-133980b81860.pdf")
	l = append(l, "dummy.pdf","sample.pdf")
	ka["./debug/mailboxes/x/u/2020-01-02"]= t
	//ka["./testlogos/"] = l
	_, err = tiffer.Merge_pdf(ka)
	if err != nil {   fmt.Println(err)}*/

}