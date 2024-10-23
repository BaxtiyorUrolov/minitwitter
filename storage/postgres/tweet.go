package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5"
	"strings"
	"twitter/api/models"
	"twitter/pkg/logger"
	"twitter/storage"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type tweetRepo struct {
	db  *pgxpool.Pool
	log logger.Logger
}

func NewTweetRepo(db *pgxpool.Pool, log logger.Logger) storage.ITweetStorage {
	return &tweetRepo{
		db:  db,
		log: log,
	}
}

func (t *tweetRepo) CreateTweet(ctx context.Context, tweet models.Tweet) (string, error) {
	uid := uuid.New()

	_, err := t.db.Exec(ctx, `
		INSERT INTO tweets (id, user_id, content, media, created_at) 
		VALUES ($1, $2, $3, $4, NOW())
		`,
		uid.String(),
		tweet.UserID,
		tweet.Content,
		tweet.Media,
	)
	if err != nil {
		t.log.Error("error while inserting tweet", logger.Error(err))
		return "", err
	}

	return uid.String(), nil
}

func (t *tweetRepo) GetTweet(ctx context.Context, tweetID string) (models.Tweet, error) {

	var content, media sql.NullString
	var updatedAt sql.NullTime

	query := `
		SELECT t.id, t.content, t.media, t.created_at, t.updated_at,
		       t.views_count,
		       (SELECT COUNT(*) FROM likes l WHERE l.tweet_id = t.id) AS likes_count 
		FROM tweets t
		WHERE t.id = $1;
	`

	var tweet models.Tweet
	err := t.db.QueryRow(ctx, query, tweetID).Scan(
		&tweet.ID,
		&content,
		&media,
		&tweet.CreatedAt,
		&updatedAt,
		&tweet.ViewsCount,
		&tweet.LikesCount,
	)
	if err != nil {
		t.log.Error("Error while getting tweet", logger.Error(err))
		return models.Tweet{}, err
	}

	if content.Valid {
		tweet.Content = content.String
	}

	if media.Valid {
		tweet.Media = media.String
	}

	if updatedAt.Valid {
		tweet.UpdatedAt = updatedAt.Time
	}

	return tweet, nil
}

func (t *tweetRepo) GetTweetList(ctx context.Context, limit, offset int, search string) ([]models.Tweet, error) {

	var content, media sql.NullString
	var updatedAt sql.NullTime

	query := `
		SELECT t.id, t.content, t.media, t.created_at, t.updated_at,
		       t.views_count, 
		       (SELECT COUNT(*) FROM likes l WHERE l.tweet_id = t.id) AS likes_count
		FROM tweets t
		WHERE t.deleted_at = 0
	`

	if search != "" {
		search = "%" + search + "%"
		query += " AND t.content ILIKE $3"
	}

	query += ` LIMIT $1 OFFSET $2`

	var rows pgx.Rows
	var err error
	if search != "" {
		rows, err = t.db.Query(ctx, query, limit, offset, search)
	} else {
		rows, err = t.db.Query(ctx, query, limit, offset)
	}

	if err != nil {
		t.log.Error("Error while getting tweet list", logger.Error(err))
		return nil, err
	}
	defer rows.Close()

	var tweets []models.Tweet
	for rows.Next() {
		var tweet models.Tweet
		err := rows.Scan(
			&tweet.ID,
			&content,
			&media,
			&tweet.CreatedAt,
			&updatedAt,
			&tweet.ViewsCount,
			&tweet.LikesCount,
		)
		if err != nil {
			t.log.Error("Error while scanning tweet", logger.Error(err))
			continue
		}
		tweets = append(tweets, tweet)
	}

	return tweets, nil
}

func (t *tweetRepo) DeleteTweet(ctx context.Context, tweetID string) error {
	_, err := t.db.Exec(ctx, `DELETE FROM tweets WHERE id = $1`, tweetID)
	if err != nil {
		t.log.Error("error while deleting tweet", logger.Error(err))
		return err
	}
	return nil
}

func (t *tweetRepo) ListTweetsByUser(ctx context.Context, userID string) ([]models.Tweet, error) {
	rows, err := t.db.Query(ctx, `SELECT t.id, t.content, t.media, t.created_at, t.updated_at,
		       t.views_count, 
		       (SELECT COUNT(*) FROM likes l WHERE l.tweet_id = t.id) AS likes_count 
		FROM tweets t
		WHERE t.user_id = $1;`, userID)
	if err != nil {
		t.log.Error("error while listing tweets", logger.Error(err))
		return nil, err
	}
	defer rows.Close()

	var tweets []models.Tweet
	for rows.Next() {
		tweet := models.Tweet{}
		err = rows.Scan(
			&tweet.ID,
			&tweet.Content,
			&tweet.Media,
			&tweet.CreatedAt,
			&tweet.UpdatedAt,
			&tweet.ViewsCount,
			&tweet.LikesCount,
		)
		if err != nil {
			t.log.Error("error while scanning tweet", logger.Error(err))
			continue
		}
		tweets = append(tweets, tweet)
	}
	return tweets, nil
}

func (t *tweetRepo) UpdateTweet(ctx context.Context, tweet models.Tweet) error {
	setClauses := []string{}
	args := []interface{}{}
	argID := 1

	if tweet.Content != "" {
		setClauses = append(setClauses, fmt.Sprintf("content = $%d", argID))
		args = append(args, tweet.Content)
		argID++
	}

	if tweet.Media != "" {
		setClauses = append(setClauses, fmt.Sprintf("media = $%d", argID))
		args = append(args, tweet.Media)
		argID++
	}

	if len(setClauses) == 0 {
		return nil
	}

	args = append(args, tweet.ID)
	query := fmt.Sprintf("UPDATE tweets SET %s WHERE id = $%d", strings.Join(setClauses, ", "), argID)

	_, err := t.db.Exec(ctx, query, args...)
	if err != nil {
		t.log.Error("error while updating tweet", logger.Error(err))
		return err
	}

	return nil
}

func (t *tweetRepo) IncrementTweetViews(ctx context.Context, userID, tweetID string) error {
	var count int
	checkQuery := `SELECT COUNT(1) FROM views WHERE user_id = $1 AND tweet_id = $2`
	err := t.db.QueryRow(ctx, checkQuery, userID, tweetID).Scan(&count)
	if err != nil {
		t.log.Error("Error while checking if tweet is viewed", logger.Error(err))
		return err
	}

	if count == 0 {
		// 2. Views jadvaliga yozuv qo'shamiz
		insertQuery := `INSERT INTO views (user_id, tweet_id) VALUES ($1, $2)`
		_, err = t.db.Exec(ctx, insertQuery, userID, tweetID)
		if err != nil {
			t.log.Error("Error while inserting into views", logger.Error(err))
			return err
		}

		// 3. Tvitning views_count ni oshiramiz
		updateQuery := `UPDATE tweets SET views_count = views_count + 1 WHERE id = $1`
		_, err = t.db.Exec(ctx, updateQuery, tweetID)
		if err != nil {
			t.log.Error("Error while updating tweet views count", logger.Error(err))
			return err
		}
	}

	return nil
}
