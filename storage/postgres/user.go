package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"twitter/api/models"
	"twitter/pkg/logger"
	"twitter/storage"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepo struct {
	db  *pgxpool.Pool
	log logger.Logger
}

func NewUserRepo(db *pgxpool.Pool, log logger.Logger) storage.IUserStorage {
	return &userRepo{
		db:  db,
		log: log,
	}
}

// Create - Foydalanuvchi yaratish
func (u *userRepo) Create(ctx context.Context, createUser models.CreateUser) (string, error) {
	uid := uuid.New()

	_, err := u.db.Exec(ctx, `
		INSERT INTO users (id, name, user_name, email, password) 
		VALUES ($1, $2, $3, $4, $5)
		`,
		uid,
		createUser.Name,
		createUser.UserName,
		createUser.Email,
		createUser.Password,
	)
	if err != nil {
		u.log.Error("error is while inserting data", logger.Error(err))
	}

	return uid.String(), nil
}

func (u *userRepo) GetByID(ctx context.Context, pKey models.PrimaryKey) (models.User, error) {
	user := models.User{}

	var bio, profilePicture sql.NullString
	var updatedAT sql.NullTime

	query := `
		SELECT id, name, user_name, email, bio, profile_picture, created_at, updated_at FROM users WHERE id = $1 AND deleted_at = 0
	`
	if err := u.db.QueryRow(ctx, query, pKey.ID).Scan(
		&user.ID,
		&user.Name,
		&user.UserName,
		&user.Email,
		&bio,
		&profilePicture,
		&user.CreatedAt,
		&updatedAT,
	); err != nil {
		u.log.Error("error is while selecting user by id", logger.Error(err))
	}

	if bio.Valid {
		user.Bio = bio.String
	}

	if profilePicture.Valid {
		user.ProfilePicture = profilePicture.String
	}

	if updatedAT.Valid {
		user.UpdatedAt = updatedAT.Time
	}

	return user, nil
}

func (u *userRepo) IsUserNameExist(ctx context.Context, userName string) (bool, error) {
	var exists bool
	err := u.db.QueryRow(ctx, `
		SELECT EXISTS (SELECT 1 FROM users WHERE user_name = $1)
	`, userName).Scan(&exists)
	if err != nil {
		fmt.Println("error while checking user_name existence:", err)
		return false, err
	}

	return exists, nil
}

func (u *userRepo) Update(ctx context.Context, updateUser models.UpdateUser) error {
	setClauses := []string{}
	args := []interface{}{}
	argID := 1

	if updateUser.Name != "" {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", argID))
		args = append(args, updateUser.Name)
		argID++
	}

	if updateUser.UserName != "" {
		setClauses = append(setClauses, fmt.Sprintf("user_name = $%d", argID))
		args = append(args, updateUser.UserName)
		argID++
	}

	if updateUser.Email != "" {
		setClauses = append(setClauses, fmt.Sprintf("email = $%d", argID))
		args = append(args, updateUser.Email)
		argID++
	}

	if updateUser.Bio != "" {
		setClauses = append(setClauses, fmt.Sprintf("bio = $%d", argID))
		args = append(args, updateUser.Bio)
		argID++
	}

	if updateUser.ProfilePicture != "" {
		setClauses = append(setClauses, fmt.Sprintf("profile_picture = $%d", argID))
		args = append(args, updateUser.ProfilePicture)
		argID++
	}

	if len(setClauses) == 0 {
		return nil
	}

	args = append(args, updateUser.ID)

	query := fmt.Sprintf("UPDATE users SET %s WHERE id = $%d", strings.Join(setClauses, ", "), argID)

	_, err := u.db.Exec(ctx, query, args...)
	if err != nil {
		u.log.Error("error while updating user", logger.Error(err))
		return err
	}

	return nil
}

func (u *userRepo) GetList(ctx context.Context, request models.GetListRequest) (models.UsersResponse, error) {
	var (
		page                = request.Page
		offset              = (page - 1) * request.Limit
		query               string
		countQuery          string
		count               = 0
		users               = []models.User{}
		search              = request.Search
		bio, profilePicture sql.NullString
		updatedAT           sql.NullTime
	)

	// So'rovlar uchun bazaviy count so'rovi
	countQuery = `SELECT count(1) FROM users WHERE deleted_at = 0 `
	if search != "" {
		// Qidiruv so'zi bor bo'lsa, qidiruvni countQuery'ga qo'shamiz
		countQuery += fmt.Sprintf(` AND (name ILIKE '%%%s%%' OR user_name ILIKE '%%%s%%' OR email ILIKE '%%%s%%')`, search, search, search)
	}
	if err := u.db.QueryRow(ctx, countQuery).Scan(&count); err != nil {
		u.log.Error("error while scanning count", logger.Error(err))
		return models.UsersResponse{}, err
	}

	query = `SELECT id, name, user_name, email, bio, profile_picture, created_at, updated_at FROM users WHERE deleted_at = 0 `
	if search != "" {
		query += fmt.Sprintf(` AND (name ILIKE '%%%s%%' OR user_name ILIKE '%%%s%%' OR email ILIKE '%%%s%%')`, search, search, search)
	}
	query += ` LIMIT $1 OFFSET $2`

	rows, err := u.db.Query(ctx, query, request.Limit, offset)
	if err != nil {
		u.log.Error("error while selecting users", logger.Error(err))
		return models.UsersResponse{}, err
	}
	defer rows.Close()

	for rows.Next() {
		user := models.User{}
		if err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.UserName,
			&user.Email,
			&bio,
			&profilePicture,
			&user.CreatedAt,
			&updatedAT,
		); err != nil {
			u.log.Error("error while scanning user", logger.Error(err))
			return models.UsersResponse{}, err
		}

		if bio.Valid {
			user.Bio = bio.String
		}

		if profilePicture.Valid {
			user.ProfilePicture = profilePicture.String
		}

		if updatedAT.Valid {
			user.UpdatedAt = updatedAT.Time
		}

		users = append(users, user)
	}

	return models.UsersResponse{
		Users: users,
		Count: count,
	}, nil
}

func (u *userRepo) Delete(ctx context.Context, id models.PrimaryKey) error {

	uuid, err := uuid.Parse(id.ID)
	if err != nil {
		return err
	}

	_, err = u.db.Exec(ctx, `
		UPDATE users 
		SET deleted_at = extract(epoch from current_timestamp) 
		WHERE id = $1
	`, uuid)

	if err != nil {
		u.log.Error("error while soft deleting user", logger.Error(err))
		return err
	}

	return nil
}

func (u userRepo) GetUserCredentials(ctx context.Context, key string) (models.User, error) {
	user := models.User{}

	query := `
		SELECT id, password 
		FROM users 
		WHERE deleted_at = 0 
		AND (user_name = $1 OR email = $1)
	`

	err := u.db.QueryRow(ctx, query, key).Scan(&user.ID, &user.Password)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
