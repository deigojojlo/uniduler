package api
import (
	"net/http"
    "github.com/gin-gonic/gin"
    "github.com/jackc/pgx/v5"
	"uniduler/utils"
    "context"
	"strings"

	"fmt"
)

func GetEvents(c *gin.Context, conn *pgx.Conn) {
	/* get all arguments */
	name := c.Query("name");
	groups := c.Query("groups"); // Like maths1, mathinfo1 ...
	year := c.Query("year"); // L1 L2 L3 M1 M2


	if name == "" || groups == "" || year == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Need required argument : name, groups, year",
		})
		return
	}

	name = strings.ToUpper(name) + "%"
	groups = strings.ToUpper(groups)
	year = strings.ToUpper(year)

	rows, _ := conn.Query(context.Background(),"SELECT * FROM events WHERE UPPER(name) LIKE $1 AND UPPER(groups)=$2 AND UPPER(year)=$3",name,groups,year)
	
	var events []utils.Event
	for rows.Next() {
		fmt.Print("&")
		var event utils.Event
		err := rows.Scan(
			&event.Name,
			&event.Groups,
			&event.Summary,
			&event.Location,
			&event.Start,
			&event.End,
			&event.Year,
			&event.DayOfTheWeek,
			&event.Type,
			&event.Parcours)
		if err != nil {continue}
		events = append(events,event)
	}
	c.IndentedJSON(http.StatusOK,events)
}


/*
	List of Maths,Computer science, eco,law ...
*/
func GetFormation(c *gin.Context, conn *pgx.Conn) {
	// no argument

	rows, _ := conn.Query(context.Background(),"SELECT DISTINCT parcours FROM events");
	var res []string
	for rows.Next() {
		var formation string
		err := rows.Scan(&formation)
		if err != nil {continue}
		res = append(res,formation)
	}
	c.IndentedJSON(http.StatusOK,res)
}

/*
	List of year in a formation like L1 L2 L3 ...
*/
func GetYear(c *gin.Context, conn *pgx.Conn) {
	// one argument
	formation := c.Query("formation")
	if formation == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Need required argument : formation",
		})
		return
	}

	rows, _ := conn.Query(context.Background(),"SELECT DISTINCT year FROM events WHERE parcours=$1",formation)

	var res []string
	for rows.Next() {
		var year string
		err := rows.Scan(&year)
		if err != nil {continue}

		res = append(res,year)
	}
	c.IndentedJSON(http.StatusOK,res)
}

/*
	List of groups in a year formation
*/
func GetGroups(c *gin.Context, conn *pgx.Conn) {
	// two argument
	formation := c.Query("formation")
	if formation == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Need required argument : formation",
		})
		return
	}
	year := c.Query("year")
	if year == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Need required argument : year",
		})
		return
	}

	rows, _ := conn.Query(context.Background(),"SELECT DISTINCT groups FROM events WHERE parcours=$1 AND year=$2",formation,year)

	var res []string
	for rows.Next() {
		var groups string
		err := rows.Scan(&groups)
		if err != nil {continue}
		res = append(res,groups)
	}
	c.IndentedJSON(http.StatusOK,res)
}

/*
	List of subject in a subject of a year of formation
*/
func GetSubject(c *gin.Context, conn *pgx.Conn) {
	// three argument
	formation := c.Query("formation")
	if formation == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Need required argument : formation",
		})
		return
	}
	year := c.Query("year")
	if year == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Need required argument : year",
		})
		return
	}
	groups := c.Query("groups")
	if groups == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Need required argument : groups",
		})
		return
	}

	rows, _ := conn.Query(context.Background(),"SELECT DISTINCT name FROM events WHERE parcours=$1 AND year=$2 AND groups=$3",formation,year,groups)
	var res []string
	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {continue}
		res = append(res,name)
	}
	c.IndentedJSON(http.StatusOK,res)
}