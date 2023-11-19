package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"net/http"

	"github.com/gin-gonic/gin"
)

type Employee struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Role       string `json:"role"`
	ProjectIDs []int  `json:"projectIds,omitempty"`
	Projects   string `json:"projects"`
}

type ProjectResponse struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"descrpiton"`
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

var PROJECT_SERVICE_URL = os.Getenv("PROJECT_SERVICE_URL")
var PROJECTS_SERVICE_PROJECT_PATH = "/projects"

func main() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	employees := []Employee{
		{
			Id:         1,
			Name:       "Mike",
			Role:       "engineer",
			ProjectIDs: []int{1, 2},
		},
		{
			Id:         2,
			Name:       "Will",
			Role:       "stsm",
			ProjectIDs: []int{1, 2},
		},
		{
			Id:         3,
			Name:       "Den",
			Role:       "architect",
			ProjectIDs: []int{3, 4},
		},
		{
			Id:         4,
			Name:       "Gem",
			Role:       "manager",
			ProjectIDs: []int{1, 3},
		},
		{
			Id:         5,
			Name:       "Ru",
			Role:       "lead-engineer",
			ProjectIDs: []int{2, 3},
		},
		{
			Id:         6,
			Name:       "Jo",
			Role:       "engineer",
			ProjectIDs: []int{1, 4},
		},
		{
			Id:         7,
			Name:       "Nirai",
			Role:       "engineer",
			ProjectIDs: []int{2, 4},
		},
		{
			Id:         8,
			Name:       "Mant",
			Role:       "engineer",
			ProjectIDs: []int{3, 4},
		},
		{
			Id:         9,
			Name:       "Hima",
			Role:       "engineer",
			ProjectIDs: []int{2, 4},
		},
	}

	r := gin.Default()
	r.Use(CORS())

	r.GET("/employees", func(c *gin.Context) {
		var employeeResponseList []Employee

		for _, employee := range employees {
			employee.setEmployeeProjects()
			employeeResponseList = append(employeeResponseList, employee)
		}
		c.JSON(200, gin.H{
			"employees": employeeResponseList,
		})
	})

	r.GET("/employees/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		var targetEmployee Employee

		employeeFound := false
		employeeIndex := 0
		for !employeeFound {
			currentEmployee := employees[employeeIndex]
			if currentEmployee.Id == id {
				targetEmployee = currentEmployee
				employeeFound = true
			}
			employeeIndex++
		}

		if employeeFound {
			targetEmployee.setEmployeeProjects()
			c.JSON(200, gin.H{
				"employee": targetEmployee,
			})
		} else {
			c.JSON(404, gin.H{"message": "Not found"})
		}
	})
	r.Run()
}

func (employee *Employee) setEmployeeProjects() {
	employeeProjects, err := getEmployeeProjects(*employee)
	if err != nil {
		fmt.Printf("Error contacting Projects service: %s\n", err)
		employee.Projects = "Currently unavailable"
	} else {
		employee.Projects = employeeProjects
	}
	employee.ProjectIDs = nil
}

func getEmployeeProjects(employee Employee) (projectList string, err error) {
	if PROJECT_SERVICE_URL == "" {
		PROJECT_SERVICE_URL = "http://localhost:4000"
	}

	var employeeProjectsList = make([]string, 0)
	for _, projectId := range employee.ProjectIDs {
		projectRequestUrl := fmt.Sprintf("%s%s/%d", PROJECT_SERVICE_URL, PROJECTS_SERVICE_PROJECT_PATH, projectId)
		fmt.Printf("Calling %s to fetch projects...\n", projectRequestUrl)
		res, err := http.Get(projectRequestUrl)
		if err != nil {
			fmt.Printf("Error fetching projects: %s\n", err)
			return "", err
		}
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			panic(err.Error())
		}
		var project ProjectResponse
		json.Unmarshal(body, &project)
		employeeProjectsList = append(employeeProjectsList, project.Name)
	}
	return strings.Join(employeeProjectsList[:], ", "), nil

}
