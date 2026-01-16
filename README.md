


Create a postgresql database and put information in a .env file like this

```env
LOGIN=
PASSWD=
DBNAME=
```

your table need to like this
```
CREATE TABLE $DBNAME (name TEXT, groups TEXT, summary TEXT, location TEXT, startDate TEXT, endDate TEXT, year TEXT, dayOfTheWeek TEXT,type TEXT, parcours TEXT);
```