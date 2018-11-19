package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"text/template"
)

// default variables
const (
	TmplFileName       = "./src/tmpl.mustache"
	CSVFileNameYoume   = "./src/youme.csv"
	CSVFileNameLady    = "./src/lady.csv"
	CSVFileNameGentle  = "./src/gentle.csv"
	DistFileNameYoume  = "./../../inc/contents-youme.php"
	DistFileNameLady   = "./../../inc/contents-lady.php"
	DistFileNameGentle = "./../../inc/contents-gentle.php"
)

type fileName struct {
	Tfn string
	Cfn string
	Dfn string
}

var fns = []fileName{
	fileName{Tfn: TmplFileName, Cfn: CSVFileNameYoume, Dfn: DistFileNameYoume},
	fileName{Tfn: TmplFileName, Cfn: CSVFileNameLady, Dfn: DistFileNameLady},
	fileName{Tfn: TmplFileName, Cfn: CSVFileNameGentle, Dfn: DistFileNameGentle},
	fileName{Tfn: "./src/tmpl_sp.mustache", Cfn: CSVFileNameYoume, Dfn: "./../../../sp/xmas/inc/contents-youme.php"},
	fileName{Tfn: "./src/tmpl_sp.mustache", Cfn: CSVFileNameLady, Dfn: "./../../../sp/xmas/inc/contents-lady.php"},
	fileName{Tfn: "./src/tmpl_sp.mustache", Cfn: CSVFileNameGentle, Dfn: "./../../../sp/xmas/inc/contents-gentle.php"},
}

func main() {
	for _, fn := range fns {
		generateRoopedFile(fn.Tfn, fn.Cfn, fn.Dfn)
	}
}

func generateRoopedFile(tfn string, cfn string, dfn string) {
	checkFileExists(tfn)
	checkFileExists(cfn)
	t := createTmplate(tfn)
	yd := parseCSVFile(cfn)
	yw := createWriter(dfn)
	defer yw.Close()
	executeTemplate(t, yw, yd, dfn)
}

func executeTemplate(t *template.Template, w *os.File, ms []map[string]string, fn string) {
	if err := t.Execute(w, ms); err != nil {
		log.Fatal(fn + "への書き込みに失敗しました。\n" + err.Error())
	}
	fmt.Println(fn + "を作成しました。")
}

func checkFileExists(fn string) {
	if _, err := os.Stat(fn); err != nil {
		log.Fatal(fn + "が配置されていません。")
	}
}

func createTmplate(fn string) *template.Template {
	fbs, err := ioutil.ReadFile(fn)
	if err != nil {
		log.Fatal(fn + "を開けませんでした。\n" + err.Error())
	}

	var hbs []byte
	hbs = append(hbs, strings.TrimSpace(string(fbs))...)
	hbs = append(hbs, "\n"...)

	tbs := []byte("{{ range . }}\n")
	tbs = append(tbs, string(hbs)...)
	tbs = append(tbs, "{{ end }}"...)

	t, err := template.New("template").Parse(string(tbs))
	if err != nil {
		log.Fatal(fn + "の解析に失敗しました。\n" + err.Error())
	}

	return t
}

func parseCSVFile(fn string) []map[string]string {
	const ABC = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	cf, err := os.Open(fn)
	if err != nil {
		log.Fatal(fn + "を開けませんでした。\n" + err.Error())
	}
	defer cf.Close()

	r := csv.NewReader(cf)
	r.FieldsPerRecord = -1
	record, err := r.ReadAll()
	if err != nil {
		log.Fatal(fn + "の解析に失敗しました。\n" + err.Error())
	}

	skipLine := 1
	var dataSet []map[string]string
	for idx, items := range record {
		if skipLine-1 >= idx {
			continue
		}

		data := make(map[string]string)
		for key, item := range items {
			item = strings.Replace(item, "\n", "<br/>", -1)
			item = strings.TrimSpace(item)
			data[string(ABC[key])] = item
		}

		sn := "0" + strconv.Itoa(idx+1-skipLine)
		data["SN"] = sn[len(sn)-2:]
		dataSet = append(dataSet, data)
	}

	return dataSet
}

func createWriter(fn string) *os.File {
	w, err := os.Create(fn)
	if err != nil {
		log.Fatal(fn + "の作成に失敗しました。\n" + err.Error())
	}
	return w
}
