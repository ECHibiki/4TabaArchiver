package main

import (
  "fmt"
  "strconv"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)
 var db  *sql.DB

func main(){
  db_init, err := sql.Open("sqlite3", "/home/ecverniy/Desktop/4t2/github/db/shitaba.db")
  db = db_init
  if err != nil {
    panic(err)
  }


  thread, posts, subposts := getThread("1", "a")
  fmt.Println("thread", thread)
  fmt.Println("---")
  fmt.Println("posts", posts)
  fmt.Println("---")
  fmt.Println("subposts", subposts)
}

func getThread(thread string, board string) (map[string]string , []map[string]string , []map[string]string){
  thread_map := make(map[string]string)
  thread_props, err := db.Query(`SELECT sticky, postcount, subject, date, old FROM threads
    WHERE board=? AND threadnum=?`, board, thread)
  if err != nil {
    panic(err)
  }
  thread_map["board"] = board
  thread_map["threadnum"] = thread
  for thread_props.Next(){
    var sticky sql.NullString
    var postcount sql.NullInt64
    var subject sql.NullString
    var date sql.NullString
    var old sql.NullInt64
    err = thread_props.Scan(&sticky, &postcount, &subject, &date, &old)
    if err != nil {
      panic(err)
    }
    if sticky.Valid{
      thread_map["sticky"] = sticky.String
    }else{
      thread_map["sticky"] = ""
    }
    if postcount.Valid{
      thread_map["postcount"] = strconv.Itoa(int (postcount.Int64))
    }else{
      thread_map["postcount"] = ""
    }
    if subject.Valid{
      thread_map["subject"] = subject.String
    }else{
      thread_map["subject"] = ""
    }
    if date.Valid{
      thread_map["date"] = date.String
    }else{
      thread_map["date"] = ""
    }
    if old.Valid{
      thread_map["old"] = strconv.Itoa(int (old.Int64))
    }else{
      thread_map["old"] = ""
    }
  }

  var post_arr []map[string]string;
  post_props, err := db.Query(`SELECT postnum, name, date, image, thumb,iname, size, comment
    FROM posts WHERE board=? AND threadnum=?`, board, thread)
  if err != nil {
    panic(err)
  }
  for post_props.Next(){
    var post_map = make(map[string]string)
    var postnum int
    var name string
    var date string
    var image string
    var thumb string
    var iname string
    var size string
    var comment string
    err = post_props.Scan(&postnum, &name, &date, &image, &thumb, &iname, &size, &comment)
    if err != nil {
      panic(err)
    }
    post_map["postnum"] = strconv.Itoa(postnum)
    post_map["name"] = name
    post_map["date"] = date
    post_map["image"] = image
    post_map["thumb"] = thumb
    post_map["iname"] = iname
    post_map["size"] = size
    post_map["comment"] = comment
    post_arr = append(post_arr, post_map)
  }

  var subpost_arr []map[string]string;
  subpost_props, err := db.Query("SELECT name, date, comment FROM subposts WHERE board=? AND threadnum=?", board, thread)
  if err != nil {
    panic(err)
  }
  for subpost_props.Next(){
    var subpost_map = make(map[string]string)
    var name string
    var date string
    var comment string
    err = subpost_props.Scan(&name, &date, &comment)
    if err != nil {
      panic(err)
    }
    subpost_map["name"] = name;
    subpost_map["date"] = date;
    subpost_map["comment"] = comment;
    subpost_arr = append(subpost_arr, subpost_map)
  }


  return thread_map, post_arr, subpost_arr

  }
