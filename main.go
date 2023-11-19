package main

import (
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Employee struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

func main() {

	employees := []Employee{
		{
			Id:   1,
			Name: "Mike",
			Role: "engineer",
		},
		{
			Id:   2,
			Name: "Will",
			Role: "stsm",
		},
		{
			Id:   3,
			Name: "Den",
			Role: "architect",
		},
		{
			Id:   4,
			Name: "Gem",
			Role: "manager",
		},
		{
			Id:   5,
			Name: "Ru",
			Role: "lead-engineer",
		},
		{
			Id:   6,
			Name: "Jo",
			Role: "engineer",
		},
		{
			Id:   7,
			Name: "Nirai",
			Role: "engineer",
		},
		{
			Id:   8,
			Name: "Mant",
			Role: "engineer",
		},
		{
			Id:   9,
			Name: "Hima",
			Role: "engineer",
		},
		{
			Id:   10,
			Name: "Dan",
			Role: "engineer",
		},
	}

	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"https://carbon-nextjs-ucc.apps.o1-373093.cp.fyre.ibm.com"}

	r.GET("/employees", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"employees": employees,
		})
	})

	r.GET("/employees/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))

		var targetEmployee Employee
		for _, currentEmployee := range employees {
			if currentEmployee.Id == id {
				targetEmployee = currentEmployee
			}
		}

		if (targetEmployee != Employee{}) {
			c.JSON(200, gin.H{
				"employee": targetEmployee,
			})
		} else {
			c.JSON(404, gin.H{"message": "Not found"})
		}
	})

	r.Use(cors.New(config))
	r.Run()
}
