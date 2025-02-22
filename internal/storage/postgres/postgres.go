package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type Storage struct {
	db *pgx.Conn
}

type Post struct {
	ID        int64
	NanoID    string
	Message   string
	CreatedAt time.Time
}

func New(dsn string) (*Storage, error) {
	const op = "storage.postgres.New"

	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	s := &Storage{db: conn}

	q := `
	CREATE TABLE IF NOT EXISTS post (
		id SERIAL PRIMARY KEY,
	    nanoid VARCHAR(10) NOT NULL,
		message TEXT NOT NULL,
	    created_at TIMESTAMPTZ DEFAULT NOW()
	);
	`

	_, err = s.db.Exec(context.Background(), q)
	if err != nil {
		s.Close()
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return s, nil
}

func (s *Storage) SavePost(nanoid, message string) (int64, time.Time, error) {
	const op = "storage.postgres.SavePost"

	q := `
    INSERT INTO post (nanoid, message)
    VALUES ($1, $2)
    RETURNING id, created_at;
    `

	var id int64
	var createdAt time.Time
	err := s.db.QueryRow(context.Background(), q, nanoid, message).Scan(&id, &createdAt)
	if err != nil {
		return 0, time.Time{}, fmt.Errorf("%s: %w", op, err)
	}

	return id, createdAt, nil
}

func (s *Storage) GetAllPost() ([]Post, error) {
	const op = "storage.postgres.GetAllPost"

	q := `SELECT id, nanoid, message, created_at FROM post;`

	rows, err := s.db.Query(context.Background(), q)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		err := rows.Scan(&p.ID, &p.NanoID, &p.Message, &p.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		posts = append(posts, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return posts, nil
}

func (s *Storage) Close() {
	if s.db != nil {
		s.db.Close(context.Background())
	}
}
