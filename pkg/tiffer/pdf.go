package tiffer

import (
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"github.com/nu7hatch/gouuid"
	"github.com/pkg/errors"
	"log"
	"os/exec"
	loader "../loader"
)

//create pdf from received email content
func Create_pdf(informations map[string]string, acc loader.Account_Informations) ([]string, error) {
	var err error
	var file_path string
	file_path = Create_folder(acc.Fax_email_connection_host.String, acc.Fax_email_connection_username.String)

	pdf := gofpdf.New("P", "pt", "A4", "")
	pdf.SetAutoPageBreak(true, 40)
	pdf.SetMargins(30, 0, 30)
	testText := "tly with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum."
	pdf.SetFooterFunc(func() {
		pdf.SetY(-25)
		pdf.SetFont("times", "I", 8)
		pdf.CellFormat(0, 10, testText,
			"", 1, "C", false, 0, "")
	})
	pdf.AddPage()

	pdf.SetTopMargin(30)
	pdf.SetFont("times", "", 10)
	pdf.Image(("./test_logos/logo_test.png"), 30, 30, 180, 70, false, "", 0, "")
	pdf.Image(("./test_logos/cover_test.png"), 450, 30, 120, 70, false, "", 0, "")

	//display umlaute in right form
	tr := pdf.UnicodeTranslatorFromDescriptor("")
	pdf.Ln(20)
	pdf.SetX(100)
	pdf.SetY(150)
	text := "FROM: " + informations["from"]
	pdf.Cell(10, 10, tr(text))

	pdf.SetX(100)
	pdf.SetY(200)
	text = "TO: " + informations["to"]
	pdf.Cell(10, 10, tr(text))

	pdf.SetX(20)
	pdf.SetY(250)
	text = "SUBJECT: " + informations["subject"]
	pdf.Cell(10, 10, tr(text))
	pdf.SetX(30)
	pdf.SetY(300)
	pdf.SetFont("times", "", 12)
	pdf.MultiCell(0, 20, tr(informations["message"]), "0", "", false)

	page_count := pdf.PageCount()
	for i := 0; i <= page_count; i++ {
		fmt.Println("pagecount: ", pdf.PageCount())
		pdf.SetPage(i)
		//other way doesn't work
		if pdf.PageCount() == 1 {
			pdf.Line(30, 290, 565.28, 290)
			pdf.Line(30, 290, 30, 800)
			pdf.Line(565.28, 290, 565.28, 800)
			pdf.Line(30, 800, 565.28, 800)
		} else {
			if i == 1 {
				pdf.Line(30, 290, 565.28, 290)
				pdf.Line(30, 290, 30, 800)
				pdf.Line(565.28, 290, 565.28, 800)
				pdf.Line(30, 800, 565.28, 800)

			} else {
				pdf.Line(30, 30, 565.28, 30)
				pdf.Line(30, 30, 30, 800)
				pdf.Line(565.28, 30, 565.28, 800)
				pdf.Line(30, 800, 565.28, 800)

			}
		}

	}

	//create uuid for unique pdf name (temp file - later this file will be removed because of a merge with other files)
	uuid, _ := uuid.NewV4()
	tmp_filename := uuid.String() + ".pdf"
	result := []string{}
	result = append(result, file_path, tmp_filename)
	file_path += "/" + tmp_filename
	err = pdf.OutputFileAndClose(file_path)
	if err != nil {
		return nil, errors.Errorf("[create_pdf] failed to store created message pdf: \n%v", err)
	}
	return result, nil
}

//INFO: function is tested on linux
//need to be root -> execute with sudo
//merge created pdf with attachments from mail and create a new one
//key in map is used as paht and contains all associated files in this path as value
func Merge_pdf(filenames map[string][]string) (string, error) {
	command := ""
	tmppath := ""
	var extraCmds []string
	for path, name := range filenames {
		tmppath = path
		for i := 0; i < len(name); i++ {
			command = path + "/" + name[i]
			extraCmds = append(extraCmds, command)
		}
	}
	uuid, _ := uuid.NewV4()
	tmp_filename := uuid.String() + ".pdf"
	tmppath += "/" + tmp_filename
	extraCmds = append(extraCmds, tmppath)
	//use pdfunite to merge pdf
	s, err := exec.Command("pdfunite", extraCmds...).Output()
	reslt := string(s)
	log.Println(reslt)
	if err != nil {
		return "", errors.Errorf("[merge_pdf] failed to execute pdfunite command: \n%v", err)
	}
	tif_file, err := create_tif(tmppath)
	if err != nil {
		return "", errors.Errorf("[merge_pdf] failed to create tif file from created pdf: \n%v", err)
	}
	return tif_file, nil

}

