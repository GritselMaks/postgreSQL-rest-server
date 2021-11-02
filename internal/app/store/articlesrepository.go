package store

import (
	"errors"
	"fmt"
	"strings"

	"github.com/GritselMaks/postgreSQL-api-server/internal/app/model"
)

type ArticlesRepository struct {
	store *Store
}

// Save data in DB and return id rows
func (r *ArticlesRepository) Save(a *model.Articles) (*model.Articles, error) {
	if len(a.Title) > 200 || len(a.FullText) > 1000 || len(a.URLFoto) > 3 {
		return nil, errors.New("andeclared len title, text article, or count foto")
	}
	if err := r.store.db.QueryRow(
		"INSERT INTO articles (title, full_text, price, data, urlfoto) values ($1,$2,$3,$4,$5) RETURNING id",
		a.Title,
		a.FullText,
		a.Price,
		a.Date,
		a.URLFoto,
	).Scan(&a.ID); err != nil {
		return nil, err
	}

	return a, nil
}

// ShowArticle makes select from DB one row and adds to variable with type Article.
func (r *ArticlesRepository) ShowArticle(id string, fields string) (model.Articles, error) {

	article := model.Articles{}
	var err error
	if fields == "" {
		err = r.store.db.QueryRow(fmt.Sprintf("SELECT title,price,url_foto FROM articles WHERE id=%s", id)).
			Scan(&article.Title, &article.Price, &article.URLFoto)
		return article, err
	}
	if fields == "full_text" {
		err = r.store.db.QueryRow(fmt.Sprintf("SELECT title,%s,price,url_foto FROM articles WHERE id=%s", fields, id)).
			Scan(&article.Title, &article.FullText, &article.Price, &article.URLFoto)
		return article, err
	}
	return article, errors.New("not valid fields value")
}

// ShowArticles makes select from DB all rows and adds to array.
func (r *ArticlesRepository) ShowArticles(sort []string) ([]model.Articles, error) {
	if len(sort) != 2 {
		return nil, errors.New("not valid sort parametrs")
	}
	for _, s := range sort {
		if !strings.Contains(s, "price") && !strings.Contains(s, "date") {
			return nil, errors.New("not valid sort parametrs")
		}
	}

	insert, err := r.store.db.Query(fmt.Sprintf("SELECT * FROM articles ORDER BY %s, %s;", sort[0], sort[1]))
	if err != nil {
		return nil, errors.New("errors query from DB")
	}
	defer insert.Close()

	articles := []model.Articles{}
	for insert.Next() {
		var article model.Articles
		if err := insert.Scan(&article.ID, &article.Title, &article.FullText, &article.Price, &article.URLFoto, &article.Date); err != nil {
			return articles, err
		}
		articles = append(articles, article)
	}
	return articles, nil
}
