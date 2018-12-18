package main

import (
    "fmt"
    "database/sql"
    "os"
    _ "github.com/go-sql-driver/mysql"
)

func main() {
    fmt.Println("Go MySQL")
    conn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_IP"), os.Getenv("DB_NAME"))
    db, err := sql.Open("mysql", conn)
    
    if err != nil {
        panic(err.Error())
    }
    
    defer db.Close()
    
    q := fmt.Sprintf("INSERT INTO filme ( filme_name ) VALUES ( '%s' )", os.Getenv("VALUE"))
    insert, err := db.Query(q)
    
    if err != nil {
        panic(err.Error())
    }
    
    defer insert.Close()  
}