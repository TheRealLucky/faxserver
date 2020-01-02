package tiffer

import (
	"fmt"
	"os"
	"time"
)


//create a folder to store generated store generated files
//create folder if it doesn't exists
func create_folder(host string, user string) string {
	dt := time.Now()
	//.format doesn't work fine. so i get the substring of the first 10 character to get the right date
	tmp := []rune(dt.String())
	folder_date := string(tmp[0:10])
	path := "./debug/mailboxes/" + host + "/" + user + "/" + folder_date
	fmt.Println(path)
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}
	return path
}
