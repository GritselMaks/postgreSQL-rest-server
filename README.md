# postgreSQL-rest-server

API service for storing and serving articles. Articles are stored in a database. 
The service provides an API that runs on top of HTTP in JSON format.

Services have 3 method: get list of articles, get one article and post article.

Articles consist of fields:
 * ID
 * Title
 * FullText
 * Price
 * URLFoto
 * Date
    
1. Titles up to 200 characters. 
2. Full description no more than 1000 characters
3. URLPhoto array containing no more than 3 links. First link is main.
    
## API

|Method         | Path           | Operation  |
| :-----------: | :------------: | :--------: | 
| GET           | /articles?sort=?     | "SELECT * FROM articles ORDER BY ? "|
| GET           | /article/id?fields=? | "SELECT title,price, ? FROM articles WHERE id=|
| POST          | /articles/new        | "INSERT INTO articles VALUES (data.values)|

Get list of articles

 1. Sort by fields price and date.
 2. Sort fields are required 
 3. "-" before fields is used to sort the data returned in ascending order.
 
Get one article 
 
 1. Manualy return name,price and main photo from article
 2. Optional fields can be requested by passing the fields parameter("full_text"): description, links to all photos. 
 
 
# Get start
Open The Terminal

  1. git clone https://github.com/GritselMaks/postgreSQL-rest-server
  2. Build the Api by
       command :: go build main.go
       
Open different terminal
  1. Start postgres DB server
  2. If it's need create database by
        command :: migrate -path migrations -database "postgres://localhost/articles?sslmode=disable" up
  
@then in first terminal run ./myapp

Default port=5432
