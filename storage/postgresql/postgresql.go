package postgresql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/KekemonBS/dhtcrawler/crawler/models"
)

//DbImpl stores db connection pointer and additional info to work with db
type DbImpl struct {
	storage *sql.DB
}

//New creates db instance
func New(db *sql.DB) *DbImpl {
	return &DbImpl{
		storage: db,
	}
}

//Create adpends one row to shares table
func (db DbImpl) Create(ctx context.Context, share models.Share) error {
	query := `INSERT INTO shares (name, shareSize, fileTree, magnetLink) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING;`
	res, err := db.storage.ExecContext(ctx, query,
		share.Name,
		share.Size,
		share.FileTree,
		share.MagnetLink,
	)
	if err != nil {
		return fmt.Errorf("store error: %w", err)
	}

	_, err = res.RowsAffected()
	if err != nil {
		return fmt.Errorf("create rows affected: %w", err)
	}
	return nil
}

//DeleteByID deletes one row from shares table
func (db DbImpl) DeleteByID(ctx context.Context, uuid string) error {
	query := `DELETE FROM shares WHERE where id = $1;`
	res, err := db.storage.ExecContext(ctx, query,
		uuid,
	)
	if err != nil {
		return fmt.Errorf("store error: %w", err)
	}

	_, err = res.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete rows affected: %w", err)
	}
	return nil
}

//ReadByID reads one row from shares table
func (db DbImpl) ReadByID(ctx context.Context, uuid string) (models.Share, error) {
	query := `SELECT name, shareSize, fileTree, magnetLink FROM shares WHERE id = $1;`
	res, err := db.storage.QueryContext(ctx, query,
		uuid,
	)
	defer res.Close()
	resShare := models.Share{}
	ok := res.Next()
	if !ok {
		return models.Share{}, fmt.Errorf("no matching rows left")
	}
	err = res.Scan(&resShare.Name,
		&resShare.Size,
		&resShare.FileTree,
		&resShare.MagnetLink)
	if err != nil {
		return models.Share{}, fmt.Errorf("read error: %w", err)
	}
	return resShare, nil
}

//ReadAll reads all entries from shares table
func (db DbImpl) ReadAll(ctx context.Context) ([]models.Share, error) {
	query := `SELECT name, shareSize, fileTree, magnetLink FROM shares`
	res, err := db.storage.QueryContext(ctx, query)
	if err != nil {
		return []models.Share{}, fmt.Errorf("query error: %w", err)
	}
	defer res.Close()
	resShares := []models.Share{}
	for res.Next() {
		resShare := models.Share{}
		err = res.Scan(&resShare.Name,
			&resShare.Size,
			&resShare.FileTree,
			&resShare.MagnetLink)
		resShares = append(resShares, resShare)
		if err != nil {
			return []models.Share{}, fmt.Errorf("read error: %w", err)
		}
	}
	return resShares, nil
}

//ReadPage reads page from shares table nth page with defined size
func (db DbImpl) ReadPage(ctx context.Context, size, n int, queryShares string) (models.SharesPage, error) {
	//Read page
	offset := size * (n - 1)
	queryShares = "%" + queryShares + "%"
	query := `SELECT name, shareSize, fileTree, magnetLink FROM shares WHERE name LIKE $3 LIMIT $1 OFFSET $2;`
	res, err := db.storage.QueryContext(ctx, query, size, offset, queryShares)
	if err != nil {
		return models.SharesPage{}, fmt.Errorf("query error: %w", err)
	}
	defer res.Close()
	resShares := []models.Share{}
	for res.Next() {
		resShare := models.Share{}
		err = res.Scan(&resShare.Name,
			&resShare.Size,
			&resShare.FileTree,
			&resShare.MagnetLink)
		if err != nil {
			return models.SharesPage{}, fmt.Errorf("read error: %w", err)
		}
		resShares = append(resShares, resShare)
	}

	//Count total res
	countQuery := `SELECT COUNT(name) FROM shares WHERE name LIKE $1;`
	resCount, err := db.storage.QueryContext(ctx, countQuery, queryShares)
	if err != nil {
		return models.SharesPage{}, fmt.Errorf("countQuery error: %w", err)
	}
	defer resCount.Close()
	var count int
	resCount.Next()
	err = resCount.Scan(&count)
	if err != nil {
		return models.SharesPage{}, fmt.Errorf("countQuery error: %w", err)
	}

	resPage := models.SharesPage{
		Total:   count,
		Results: resShares,
	}

	return resPage, nil
}
