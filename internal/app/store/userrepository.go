package store

import (
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/GritselMaks/postgreSQL-api-server/internal/app/model"
)

type ArticlesRepository struct {
	store *Store
}

func (r *ArticlesRepository) Save(a *model.Articles) (*model.Articles, error) {
	if len(a.Title) > 200 || len(a.FullText) > 1000 || len(a.URLFoto) > 3 {
		return nil, errors.New("andeclared len title, text article, or count foto")
	}
	if err := r.store.Db.QueryRow(
		"INSERT INTO articles (title, full_text, price, data, urlfoto) values ($1,$2,$3,$4,$5) RETURNING id",
		a.Title,
		a.FullText,
		a.Price,
		a.Data,
		a.URLFoto,
	).Scan(&a.ID); err != nil {
		return nil, err
	}

	return a, nil
}

//
func (r *ArticlesRepository) ShowArticle(article model.Articles, vars map[string]string) error {

	return r.store.Db.QueryRow(fmt.Sprintf("select title, full_text, price, urlfoto from articles where id=%s", vars["id"])).
		Scan(&article.Title, &article.FullText, &article.Price, &article.URLFoto)

}

func (r *ArticlesRepository) ShowArticles(v url.Values) ([]model.Articles, error) {
	var insert *sql.Rows
	var err error
	s := v.Get("sort")
	sortParam := strings.Split(s, ",")
	fmt.Println(sortParam[0])
	for i, s := range sortParam {
		if s == "price" {
			sortParam[i] = "price DESC"
		}
		if s == "-price" {
			sortParam[i] = "price ASC"
		}
		if s == "-data" {
			sortParam[i] = "data_at ASC"
		}
		if s == "data" {
			sortParam[i] = "data_at DESC"
		}
	}

	fmt.Println(sortParam)
	insert, err = r.store.Db.Query(fmt.Sprintf("SELECT * FROM articles ORDER BY %s, %s;", sortParam[0], sortParam[1]))

	if err != nil {
		fmt.Println("errors query from DB")
		return nil, errors.New("errors query from DB")
	}
	defer insert.Close()

	articles := []model.Articles{}
	for insert.Next() {
		var article model.Articles
		if err := insert.Scan(&article.ID, &article.Title, &article.FullText, &article.Price, &article.URLFoto, &article.Data); err != nil {
			fmt.Println("errors Scan articles")
			return nil, errors.New("errors Scan articles")
		}
		articles = append(articles, article)
	}
	return articles, nil
}
