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
	fmt.Print("Name :",name)
	fmt.Print("\nYear :",year)
	fmt.Print("\nGroups :",groups)

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