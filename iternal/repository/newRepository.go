package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/LaughG33k/TestTask2/iternal/model"
	"gopkg.in/reform.v1"
)

type News struct {
	Db *reform.DB
}

func (r *News) Edit(ctx context.Context, id int64, editModel model.NewsModel) error {

	args := make([]any, 0, 3)
	types := make([]string, 0, 3)
	var err error

	tx, err := r.Db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})

	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if editModel.Title != "" && len(editModel.Title) < 255 {
		args = append(args, editModel.Title)
		types = append(types, "title")
	}

	if editModel.Content != "" {
		args = append(args, editModel.Content)
		types = append(types, "content")
	}

	if editModel.Id > 0 {
		args = append(args, editModel.Id)
		types = append(types, "id")
	}

	if editModel.Id == 0 {
		if err = r.updateNewsCategory(ctx, id, editModel.Categories, tx); err != nil {
			return err
		}
	} else if editModel.Id > 0 {
		if err = r.updateNewsCategory(ctx, editModel.Id, editModel.Categories, tx); err != nil {
			return err
		}
	}

	if len(args) == 0 {
		if err = tx.Commit(); err != nil {
			return err
		}
		return nil
	}

	query := "update News set"

	for i := 0; i < len(args); i++ {
		query += fmt.Sprintf(" %s=$%d,", types[i], i+1)
	}

	if query[len(query)-1] == ',' {
		query = query[:len(query)-1]
	}

	query += fmt.Sprintf(" where id = %d;", id)

	_, err = tx.ExecContext(ctx, query, args...)

	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *News) updateNewsCategory(ctx context.Context, newsId int64, categories []int64, tx *reform.TX) error {

	if len(categories) == 0 {
		return nil
	}

	if _, err := tx.ExecContext(ctx, "delete from NewsCategory where NewsId = $1", newsId); err != nil {
		return err
	}

	query := "insert into NewsCategory (NewsId, CategoryId) values"

	for i := 0; i < len(categories); i++ {

		query += fmt.Sprintf(" (%d, %d),", newsId, categories[i])

	}

	query = query[:len(query)-1]

	if _, err := tx.ExecContext(ctx, query); err != nil {
		return err
	}

	return nil
}

func (r *News) GetNews(ctx context.Context) ([]*model.NewsModel, error) {

	var err error

	rows, err := r.Db.QueryContext(ctx, "select news.id, news.title, news.content, newscategory.categoryid from news left join newscategory on news.id = newscategory.newsid;")

	if err != nil {

		if err.Error() == "no rows in result" {
			return nil, nil
		}

		return nil, err
	}

	defer rows.Close()

	newsModelMap := make(map[int64]*model.NewsModel, 0)
	res := make([]*model.NewsModel, 0)

	for rows.Next() {

		var id interface{}
		var title any
		var content any
		var category any

		if err = rows.Scan(&id, &title, &content, &category); err != nil {
			continue
		}

		m, ok := newsModelMap[id.(int64)]

		if !ok {
			m = &model.NewsModel{Id: id.(int64), Title: title.(string), Content: content.(string)}
			m.Categories = make([]int64, 0)

			if category != nil {
				m.Categories = append(m.Categories, category.(int64))
			}
			res = append(res, m)

		} else {
			if category != nil {
				m.Categories = append(m.Categories, category.(int64))
			}
		}

		newsModelMap[id.(int64)] = m

	}

	if len(newsModelMap) == 0 {
		return nil, err
	}

	return res, nil

}
