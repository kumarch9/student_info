package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var tmpl *template.Template

// type studentdate struct {
// 	syear, smonth, sdate uint
// }
type student struct {
	Id         string
	Fname      string
	Lname      string
	Age        string
	Fathername string
	Address    string
	State      string
	Country    string
	Course     string
	Batch      string
	Rollnum    string
}

type MsgAll struct {
	Status  bool      `bson:"status" json:"status"`
	Message string    `bson:"message" json:"message"`
	Data    []student `bson:"data" json:"data"`
}

type Msg struct {
	Status  bool        `bson:"status" json:"status"`
	Message string      `bson:"message" json:"message"`
	Data    interface{} `bson:"data" json:"data"`
}

type UserMessage struct {
	EmailId string `bson:"emailid" json:"emailid"`
	Message string `bson:"message" json:"message"`
}

var s student

func init() {
	tmpl = template.Must(template.ParseGlob("./templates/*.html"))
}

var ptrStudentDB *sql.DB

func sqlCon() *sql.DB {
	db, err := sql.Open("mysql", "root:@(127.0.0.1:3306)/studentdb?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	return db
}
func main() {
	//to access sub file in templates
	fs := http.FileServer(http.Dir("./assets/"))
	http.Handle("/assets/", http.StripPrefix("/assets", fs))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/insert", insertHandler)
	http.HandleFunc("/get", getHanler)
	http.HandleFunc("/put", putHandler)
	http.HandleFunc("/update", updateHandler)
	http.HandleFunc("/delete", deleteHanler)
	http.HandleFunc("/contact", contactHanler)
	http.ListenAndServe(":8055", nil)
}
func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "index.html", nil)
}

func insertHandler(w http.ResponseWriter, r *http.Request) {
	defer ptrStudentDB.Close()
	if r.Method != "POST" {
		tmpl.ExecuteTemplate(w, "insert.html", nil)
		return
	}

	stud := student{
		Fname:      r.FormValue("stdfname"),
		Lname:      r.FormValue("stdfname"),
		Age:        r.FormValue("stdage"),
		Fathername: r.FormValue("stdfathername"),
		Address:    r.FormValue("stdaddress"),
		State:      r.FormValue("stdstate"),
		Country:    r.FormValue("stdcountry"),
		Course:     r.FormValue("stdcourse"),
		Batch:      r.FormValue("stdbatch"),
		Rollnum:    r.FormValue("stdrollnum"),
	}

	ptrStudentDB = sqlCon()
	result, err := ptrStudentDB.Exec("insert into student (id,Fname,Lname,Age,Fathername,Address,State,Country,Course,Batch,Rollnum) value(null,?,?,?,?,?,?,?,?,?,?)", stud.Fname, stud.Lname, stud.Age, stud.Fathername, stud.Address, stud.State, stud.Country, stud.Course, stud.Batch, stud.Rollnum)
	if err != nil {
		Resp := Msg{false, err.Error(), nil}
		tmpl.ExecuteTemplate(w, "index.html", Resp)
		log.Println(err)
		return
	}
	log.Println(result)
	Resp := Msg{true, "Data Inserted", nil}
	tmpl.ExecuteTemplate(w, "insert.html", Resp)

}

func getHanler(w http.ResponseWriter, r *http.Request) {
	defer ptrStudentDB.Close()

	D_data := []student{}
	ptrStudentDB = sqlCon()
	rows, err := ptrStudentDB.Query("select * from student order by id DESC;")
	if err != nil {
		Resp := Msg{false, err.Error(), nil}
		tmpl.ExecuteTemplate(w, "index.html", Resp)
		log.Println(err)
		return
	}

	for rows.Next() {
		if er := rows.Scan(&s.Id, &s.Fname, &s.Lname, &s.Age, &s.Fathername, &s.Address, &s.State, &s.Country, &s.Course, &s.Batch, &s.Rollnum); er != nil {
			log.Println("err in rows scan ", er)
		}
		D_data = append(D_data, s)

	}
	log.Println("student:", s)
	log.Println("d_data :", D_data)
	rows.Close()
	Resp := Msg{true, "Data fetched", D_data}
	tmpl.ExecuteTemplate(w, "get.html", Resp)

}

func putHandler(w http.ResponseWriter, r *http.Request) {
	defer ptrStudentDB.Close()
	D_data := []student{}
	ptrStudentDB = sqlCon()
	sId := r.URL.Query().Get("id")
	rows, err := ptrStudentDB.Query("select * from student where id=?", sId)
	if err != nil {
		Resp := Msg{false, err.Error(), nil}
		tmpl.ExecuteTemplate(w, "get.html", Resp)
		log.Println("error in query of sql:", err)
		return
	}

	for rows.Next() {
		if er := rows.Scan(&s.Id, &s.Fname, &s.Lname, &s.Age, &s.Fathername, &s.Address, &s.State, &s.Country, &s.Course, &s.Batch, &s.Rollnum); er != nil {
			log.Println("err in rows scan ", er)
		}
		D_data = append(D_data, s)
	}
	log.Println("s", s)
	log.Println("d_data", D_data)
	rows.Close()
	Resp := Msg{true, "Update Data fetched", D_data}
	tmpl.ExecuteTemplate(w, "put.html", Resp)

}

func deleteHanler(w http.ResponseWriter, r *http.Request) {
	defer ptrStudentDB.Close()
	ptrStudentDB = sqlCon()
	sid := r.URL.Query().Get("id")
	rows, err := ptrStudentDB.Prepare("delete from student where id=?")
	if err != nil {
		log.Println(err)
		return
	}
	rows.Exec(sid)
	log.Println("Deleted data.")
	//http.Redirect(w, r, "/get", http.StatusMovedPermanently)
	tmpl.ExecuteTemplate(w, "/get", nil)

}

func updateHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		stud := student{
			Id:         r.FormValue("stdid"),
			Fname:      r.FormValue("stdfname"),
			Lname:      r.FormValue("stdfname"),
			Age:        r.FormValue("stdage"),
			Fathername: r.FormValue("stdfathername"),
			Address:    r.FormValue("stdaddress"),
			State:      r.FormValue("stdstate"),
			Country:    r.FormValue("stdcountry"),
			Course:     r.FormValue("stdcourse"),
			Batch:      r.FormValue("stdbatch"),
			Rollnum:    r.FormValue("stdrollnum"),
		}
		log.Println("ready in update handler to cmd sql :", stud)
		ptrStudentDB = sqlCon()
		result, err := ptrStudentDB.Prepare("update student set Fname=?,Lname=?,Age=?,Fathername=?,Address=?,State=?,Country=?,Course=?,Batch=?,Rollnum=? where id=?")
		if err != nil {
			Resp := Msg{false, err.Error(), stud}
			tmpl.ExecuteTemplate(w, "get.html", Resp)
			log.Println("error in executing update cmd:", err)
			return
		}
		result.Exec(stud.Fname, stud.Lname, stud.Age, stud.Fathername, stud.Address, stud.State, stud.Country, stud.Course, stud.Batch, stud.Rollnum, stud.Id)
		log.Println("after executed cmd the result:", result)
		//Resp := Msg{true, "Data Updated", nil}
		fmt.Println("update successful.")
	}

	defer ptrStudentDB.Close()
	time.Sleep(3000 * time.Millisecond)
	//tmpl.ExecuteTemplate(w, "index.html", Resp)
	http.Redirect(w, r, "/get", http.StatusMovedPermanently)
}

func contactHanler(w http.ResponseWriter, r *http.Request) {
	ptrStudentDB = sqlCon()
	if r.Method != "POST" {
		//Resp := Msg{false, "Not POST Method Call.", nil}
		tmpl.ExecuteTemplate(w, "contact.html", nil)
		return
	}
	msg := UserMessage{
		EmailId: r.FormValue("mailid"),
		Message: r.FormValue("sendtxt"),
	}
	ptrStudentDB = sqlCon()
	result, err := ptrStudentDB.Prepare("insert into message (id,email,msg) values (?,?,?);")
	if err != nil {
		Resp := Msg{false, err.Error(), nil}
		tmpl.ExecuteTemplate(w, "contact.html", Resp)
		log.Println(err)
		return
	}
	result.Exec("Null", msg.EmailId, msg.Message)
	Resp := Msg{true, "Message Sent.", nil}
	defer ptrStudentDB.Close()
	tmpl.ExecuteTemplate(w, "contact.html", Resp)

}
