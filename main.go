package main

import (
  "fmt"
	"database/sql"
  "encoding/json"
  "net/http"
  "html/template"
  "io/ioutil"
  "time"
	_ "github.com/mattn/go-sqlite3"
)
 var db  *sql.DB

type Thread struct{
  Board string
  ThreadNum string
  Sticky string
  PostCount int
  Subject string
  Date string
  Old int
  Replies []Post
}
type Post struct{
  PostNum int
  Name string
  Date string
  Image string
  Thumb string
  IName string
  Size string
  Comment string
  SubPosts []SubPost
}
type SubPost struct{
  Name string
  Date string
  Comment string
  PostNum int
}

type Inputs struct{
  Time string
}

func main(){
  db_init, err := sql.Open("sqlite3", "/home/ecverniy/Desktop/4t2/github/db/shitaba.db")
  db = db_init
  if err != nil {
    panic(err)
  }

  //init
  http.HandleFunc("/", generateRootForm)
  http.HandleFunc("/get", returnThreadSave)
  http.ListenAndServe(":8080", nil)

}

func generateRootForm(w http.ResponseWriter, r *http.Request){
  input_tmpl_str, err := ioutil.ReadFile("./inputs.template")
  if err != nil {
    panic(err)
  }
  inputs_tmpl, err := template.New("inputs").Parse(string (input_tmpl_str))
  if err != nil {
    panic(err)
  }
  inputs := Inputs{time.Now().Format("Mon Jan _2 15:04:05 2006")}
  err = inputs_tmpl.ExecuteTemplate(w,"inputs", inputs)
  if err != nil {
    panic(err)
  }
}

func returnThreadSave(w http.ResponseWriter, r *http.Request){
  thread_arr := r.URL.Query()["thread"]
  var thread string
  if thread_arr != nil{
    thread = thread_arr[0]
  }
  board_arr := r.URL.Query()["board"]
  var board string
  if board_arr != nil{
    board = board_arr[0]
  }
  threads := getThread(thread, board)
  json, err := json.Marshal(threads)
  if err != nil{
    panic(err)
  }
  fmt.Fprintf(w, string (json))
}

func getThread(thread string, board string) (Thread){
  thread_struct := Thread{}
  thread_props, err := db.Query(`SELECT sticky, postcount, subject, Date, old FROM threads
    WHERE board=? AND threadnum=?`, board, thread)
  if err != nil {
    panic(err)
  }
  thread_struct.Board = board
  thread_struct.ThreadNum = thread
  for thread_props.Next(){
    var sticky sql.NullString
    var postcount sql.NullInt64
    var subject sql.NullString
    var Date sql.NullString
    var old sql.NullInt64
    err = thread_props.Scan(&sticky, &postcount, &subject, &Date, &old)
    if err != nil {
      panic(err)
    }
    if sticky.Valid{
      thread_struct.Sticky = sticky.String
    }
    if postcount.Valid{
      thread_struct.PostCount = int (postcount.Int64)
    }
    if subject.Valid{
      thread_struct.Subject = subject.String
    }
    if Date.Valid{
      thread_struct.Date = Date.String
    }
    if old.Valid{
      thread_struct.Old = int (old.Int64)
    }
  }

  var post_arr []Post;
  post_props, err := db.Query(`SELECT PostNum, Name, Date, Image, Thumb,IName, Size, Comment
    FROM posts WHERE board=? AND threadnum=?`, board, thread)
  if err != nil {
    panic(err)
  }
  for post_props.Next(){
    post_struct := Post{}
    var PostNum sql.NullInt64
    var Name sql.NullString
    var Date sql.NullString
    var Image sql.NullString
    var Thumb sql.NullString
    var IName sql.NullString
    var Size sql.NullString
    var Comment sql.NullString
    err = post_props.Scan(&PostNum, &Name, &Date, &Image, &Thumb, &IName, &Size, &Comment)
    if err != nil {
      panic(err)
    }

    if PostNum.Valid{
      post_struct.PostNum = int (PostNum.Int64)
    }
    if Name.Valid{
      post_struct.Name = Name.String
    }
    if Date.Valid{
      post_struct.Date = Date.String
    }
    if Image.Valid{
      post_struct.Image = Image.String
    }
    if Thumb.Valid{
      post_struct.Thumb = Thumb.String
    }
    if IName.Valid{
      post_struct.IName = IName.String
    }
    if Size.Valid{
      post_struct.Size = Size.String
    }
    if Comment.Valid{
      post_struct.Comment = Comment.String
    }
    post_arr = append(post_arr, post_struct)
  }

  var subpost_arr []SubPost;
  subpost_props, err := db.Query("SELECT Name, Date, Comment, PostNum FROM SubPosts WHERE board=? AND threadnum=?", board, thread)
  if err != nil {
    panic(err)
  }
  for subpost_props.Next(){
    subpost_struct := SubPost{}
    var Name string
    var Date string
    var Comment string
    var PostNum int
    err = subpost_props.Scan(&Name, &Date, &Comment, &PostNum)
    if err != nil {
      panic(err)
    }
    subpost_struct.Name = Name;
    subpost_struct.Date = Date;
    subpost_struct.Comment = Comment;
    subpost_struct.PostNum = PostNum;
    subpost_arr = append(subpost_arr, subpost_struct)
  }

  for subpost_no := 0 ; subpost_no < len(subpost_arr) ; subpost_no++{
    for post_no := 0; post_no < len(post_arr) ; post_no++{
      if subpost_arr[subpost_no].PostNum == post_arr[post_no].PostNum{
        post_arr[post_no].SubPosts = append(post_arr[post_no].SubPosts, subpost_arr[subpost_no])
      }
    }
  }

  thread_struct.Replies = post_arr

  return thread_struct

  }
