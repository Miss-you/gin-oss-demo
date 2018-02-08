package main

import (
    "github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
    "net/http"
	"database/sql"
	"log"
	"strconv"
)

type Person struct {
	Id         int    `json:"id" form:"id"`
	FirstParam string `json:"first_param" form:"first_param"`
	LastParam  string `json:"last_param" form:"last_param"`
}

func main(){

    router := gin.Default()

	//init mysql conn pool
	//db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test?parseTime=true")
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/test?parseTime=true")
	defer db.Close()
	if err != nil{
		log.Fatalln(err)
	}

	db.SetMaxIdleConns(20)
	db.SetMaxOpenConns(20)

	if err := db.Ping(); err != nil{
		log.Fatalln(err)
	}
	//init mysql end

    router.GET("/ping", func(c *gin.Context) {
        c.String(http.StatusOK, "Hello World")
    })

	router.POST("/helloworld", func(c *gin.Context) {
		c.String(http.StatusOK, "welcome use the qcloud API gateway. If you see this message, it means that your signature is success!")
    })

	router.POST("/create", func(c *gin.Context) {
        firstparam := c.DefaultPostForm("firstparam", "yousa")
        lastparam := c.PostForm("lastparam")

		_, err := db.Exec(
			"INSERT INTO person (first_param, last_param) VALUES (?, ?)", firstparam, lastparam)
		
		if err != nil {
			log.Fatal(err)
		}

		c.String(http.StatusOK, "Hello %s %s", firstparam, lastparam)
    })

	router.POST("/getall", func(c *gin.Context) {
        //firstparam := c.DefaultPostForm("firstparam", "yousa")
        //lastparam := c.PostForm("lastparam")

		rows, err := db.Query("SELECT first_param, last_param FROM person")
		if err != nil {
			log.Fatal(err)
		}

		persons := make([]Person, 0)

		//read
		for rows.Next() {
			var person Person
			rows.Scan(&person.FirstParam, &person.LastParam)
			persons = append(persons, person)
		}

		c.JSON(http.StatusOK, gin.H{
			"persons": persons,
		})
    })

	router.POST("/update", func(c *gin.Context) {
        //firstparam := c.DefaultPostForm("firstparam", "yousa")
        //lastparam := c.PostForm("lastparam")
		cid := c.PostForm("id")
		id, err := strconv.Atoi(cid)
		firstparam := c.DefaultPostForm("firstparam", "yousa")
		lastparam := c.PostForm("lastparam")

		person := Person{Id: id, FirstParam: firstparam, LastParam: lastparam}
		req, err := db.Prepare("UPDATE person SET first_param=?, last_param=? WHERE id=?")
 		defer req.Close()
		if err != nil {
			log.Fatalln(err)
		}

		result, err := req.Exec(person.FirstParam, person.LastParam, person.Id)
		if err != nil {
			log.Fatalln(err)
		}

		result.RowsAffected()
		/*
		_, err := result.RowsAffected()
		if err != nil {
			log.Fatalln(err)
		}
		*/

		c.JSON(http.StatusOK, gin.H{
			"persons": "ok",
		})
    })

	router.POST("/delete", func(c *gin.Context) {
		firstparam := c.DefaultPostForm("firstparam", "yousa")

		_, err := db.Exec("DELETE FROM person WHERE first_param=?", firstparam)
		if err != nil {
			log.Fatalln(err)
		}

		c.JSON(http.StatusOK, gin.H{
			"persons": "ok",
		})
	})
	/*
	router.POST("/update", func(c *gin.Context) {
        //firstparam := c.DefaultPostForm("firstparam", "yousa")
        //lastparam := c.PostForm("lastparam")
		cid := c.PostForm("id")
		id, err := strconv.Atoi(cid)
		firstparam := c.DefaultPostForm("firstparam", "yousa")
		lastparam := c.PostForm("lastparam")

		person := Person{Id: id, FirstParam: firstparam, LastParam: lastparam}
		req, err := db.Prepare("UPDATE person SET first_param=?, last_param=? WHERE id=?")
 		defer req.Close()
		if err != nil {
			log.Fatalln(err)
		}

		result, err := req.Exec(person.FirstParam, person.LastParam, person.Id)
		if err != nil {
			log.Fatalln(err)
		}

		result.RowsAffected()

		c.JSON(http.StatusOK, gin.H{
			"persons": "ok",
		})
    })
	*/

    router.Run(":8080")
}
