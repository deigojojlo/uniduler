package main

import (
	"uniduler/utils"
	"uniduler/api"
	"strings"

    "github.com/jackc/pgx/v5"
    "github.com/gin-gonic/gin"
)

func fetchData(conn *pgx.Conn){
	dat := utils.ReadCalendars()
	for _,line := range dat {
		link := "https://edt.math.univ-paris-diderot.fr/data/"+ line.Code +".ics"

		/* We get events */
		events := utils.Trunck(utils.Sort_events(utils.AddDate(utils.Parse(utils.Get(link)))))
		if events == nil {continue}
		if strings.Contains(line.Parcours,"zz") {continue}
		/* tmp */
		if events == nil {continue}
		/* We add some record */
		for _,element := range events {
			utils.AddYear(element,line.YearRaw)
			utils.AddGroups(element,line.Label)
			element.Parcours = line.Parcours
			utils.AddName(element)
			/* register events in the temporary database */

			utils.AddData(conn,element)
		}
	}
}

func main() {
	/* open database */
	conn,err := utils.Connect()
	if err != nil {
		panic(err)
	}
	/* to update data */
	fetchData(conn)

	router := gin.Default()
	router.GET("/uniduler/formation",func(c *gin.Context) {api.GetFormation(c,conn)}) // no args
	router.GET("/uniduler/year",func(c *gin.Context) {api.GetYear(c,conn)}) // formation args
	router.GET("/uniduler/groups",func(c *gin.Context) {api.GetGroups(c,conn)}) // formation + year
	router.GET("/uniduler/subject",func(c *gin.Context) {api.GetSubject(c,conn)}) // formation + year + groups
	router.GET("/uniduler/events",func(c *gin.Context) {api.GetEvents(c,conn)}) // name(subject) + year + groups
	


	router.Run("localhost:8080")
}
