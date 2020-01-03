package dao

import (
	"context"
	"database/sql"
	"log"
	"strings"

	"github.com/wsloong/go-crawler/global"

	"github.com/wsloong/go-crawler/model"
)

type jobDAO struct{}

var DefaultJob = jobDAO{}

func (j jobDAO) fields() string {
	return "id,name,city,company,education,jd,salary,welfare,work_exp,url,created_at"
}

func (j jobDAO) scanRow(row *sql.Row, job *model.JobInfo) error {
	return row.Scan(&job.ID, &job.Name, &job.City, &job.Company, &job.Education, &job.JD, &job.Salary,
		&job.Welfare, &job.WorkExp, &job.URL, &job.CreatedAt)
}

func (j jobDAO) Find(ctx context.Context, name, company string, offset, limit int) ([]*model.JobInfo, error) {
	var b strings.Builder

	b.WriteString("SELECT " + j.fields() + " FROM job_info")

	args := make([]interface{}, 0, 4)
	if name != "" {
		b.WriteString(" WHERE name LIKE ?")
		args = append(args, "%"+name+"%")
	}

	if company != "" {
		if len(args) > 0 {
			b.WriteString(" AND ")
		} else {
			b.WriteString(" WHERE ")
		}

		b.WriteString("company LIKE ?")
		args = append(args, "%"+company+"%")
	}

	b.WriteString(" LIMIT ?, ?")
	args = append(args, offset, limit)

	rows, err := masterDB.QueryContext(ctx, b.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	jobs := make([]*model.JobInfo, 0, limit)
	for rows.Next() {
		job := &model.JobInfo{}
		dest := []interface{}{
			&job.ID, &job.Name, &job.City, &job.Company, &job.Education, &job.JD, &job.Salary,
			&job.Welfare, &job.WorkExp, &job.URL, &job.CreatedAt,
		}
		err = rows.Scan(dest...)
		if err != nil {
			log.Println("job info scan error: ", err)
			continue
		}

		jobs = append(jobs, job)
	}

	return jobs, nil
}

func (j jobDAO) FindOne(ctx context.Context, id int) (*model.JobInfo, error) {
	job := &model.JobInfo{}

	strSql := "SELECT " + j.fields() + " FROM job_info WHERE id=?"
	row := masterDB.QueryRowContext(ctx, strSql, id)
	if err := j.scanRow(row, job); err != nil {
		if err == sql.ErrNoRows {
			return nil, global.ErrNoRecord
		}
		return nil, err
	}
	return job, nil
}

func (j jobDAO) FindOneByURL(ctx context.Context, jobURL string) (*model.JobInfo, error) {
	job := &model.JobInfo{}

	strSql := "SELECT " + j.fields() + " FROM job_info WHERE url=?"
	row := masterDB.QueryRowContext(ctx, strSql, jobURL)
	if err := j.scanRow(row, job); err != nil {
		if err == sql.ErrNoRows {
			return nil, global.ErrNoRecord
		}
		return nil, err
	}
	return job, nil
}

func (j jobDAO) Create(job *model.JobInfo) error {
	strSql := "INSERT INTO job_info(name,city,company,education,jd,salary,welfare,work_exp,url) VALUES(?,?,?,?,?,?,?,?,?)"
	result, err := masterDB.Exec(strSql, strings.TrimSpace(job.Name), strings.TrimSpace(job.City),
		strings.TrimSpace(job.Company), strings.TrimSpace(job.Education), strings.TrimSpace(job.JD),
		strings.TrimSpace(job.Salary), strings.TrimSpace(job.Welfare),
		strings.TrimSpace(job.WorkExp), job.URL)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	job.ID = id
	return nil

}
