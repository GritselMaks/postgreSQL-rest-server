package store

import (
	"errors"
	"fmt"

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

func (r *ArticlesRepository) ShowArticle(id string, fields string) (model.Articles, error) {

	article := model.Articles{}
	var err error
	if fields != "" {
		err = r.store.Db.QueryRow(fmt.Sprintf("SELECT title,%s,price,url_foto FROM articles WHERE id=%s", fields, id)).
			Scan(&article.Title, &article.FullText, &article.Price, &article.URLFoto)
	} else {
		err = r.store.Db.QueryRow(fmt.Sprintf("SELECT title,price,url_foto FROM articles WHERE id=%s", id)).
			Scan(&article.Title, &article.Price, &article.URLFoto)
	}

	return article, err

}

func (r *ArticlesRepository) ShowArticles(sort []string) ([]model.Articles, error) {

	insert, err := r.store.Db.Query(fmt.Sprintf("SELECT * FROM articles ORDER BY %s, %s;", sort[0], sort[1]))
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
