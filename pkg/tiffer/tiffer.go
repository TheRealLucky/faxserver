package tiffer

import (
	"fmt"
	"github.com/pkg/errors"
	"log"
	"os"
	"os/exec"
	"time"
)


//create a folder to store generated store generated files
//create folder if it doesn't exists
func Create_folder(host string, user string) string {
	dt := time.Now()
	//.format doesn't work fine. so i get the substring of the first 10 character to get the right date
	tmp := []rune(dt.String())
	folder_date := string(tmp[0:10])
	path := "./debug/mailboxes/" + host + "/" + user + "/" + folder_date
	//fmt.Println(path)
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}
	return path
}

//create tif file from pdf
func Create_tif(path string) (string, error) {
	//TODO: -g paramter
	//a := 8.3 * 204
	//b := 11.7 * 196
	tmppath := path[:len(path)-4]
	tmppath += ".tif"
	extraCmds := []string{"-sDEVICE=tiffg3", "-r204x196", "-dNOPAUSE",
		fmt.Sprintf("-sOutputFile=%s", tmppath), fmt.Sprintf("%s", path), "-c quit",
	}
	//use ghostscript (gs) to create this tif file
	s, err := exec.Command("gs", extraCmds...).Output()
	reslt := string(s)
	log.Println(reslt)
	if err != nil {
		return "", errors.Errorf("[create_tif] failed to execute gs command: \n%v", err)
	}
	return tmppath, nil
}
