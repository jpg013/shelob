package proxy

import (
	"log"
	"shelob/db"
)

type UnknownProxyPort struct {
	FilePath string
	OCRText  string
	Port     int32
}

func (u *UnknownProxyPort) Insert() int64 {
	sql := "INSERT INTO unknown_proxy_port(file_path, ocr_text) VALUES(?,?)"
	stmt, err := db.Conn.Prepare(sql)

	if err != nil {
		log.Fatal(err)
	}

	res, err := stmt.Exec(
		u.FilePath,
		u.OCRText,
	)

	if err != nil {
		log.Fatal(err)
	}

	id, err := res.LastInsertId()

	if err != nil {
		log.Fatal(err)
	}

	return id
}
