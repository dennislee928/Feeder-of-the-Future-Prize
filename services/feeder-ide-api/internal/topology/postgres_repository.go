package topology

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/feeder-platform/feeder-ide-api/internal/database"
	"github.com/google/uuid"
)

// PostgresRepository PostgreSQL 實作
type PostgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository 建立新的 PostgreSQL repository
func NewPostgresRepository() (*PostgresRepository, error) {
	if database.DB == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	return &PostgresRepository{
		db: database.DB,
	}, nil
}

func (r *PostgresRepository) Create(topology *Topology) error {
	if topology.ID == "" {
		topology.ID = uuid.New().String()
	}

	now := time.Now()
	if topology.CreatedAt.IsZero() {
		topology.CreatedAt = now
	}
	topology.UpdatedAt = now

	// 序列化 nodes 和 lines 為 JSONB
	nodesJSON, err := json.Marshal(topology.Nodes)
	if err != nil {
		return fmt.Errorf("failed to marshal nodes: %w", err)
	}

	linesJSON, err := json.Marshal(topology.Lines)
	if err != nil {
		return fmt.Errorf("failed to marshal lines: %w", err)
	}

	query := `
		INSERT INTO topologies (id, user_id, name, description, profile_type, nodes, lines, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err = r.db.Exec(query,
		topology.ID,
		topology.UserID,
		topology.Name,
		topology.Description,
		topology.ProfileType,
		nodesJSON,
		linesJSON,
		topology.CreatedAt,
		topology.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create topology: %w", err)
	}

	return nil
}

func (r *PostgresRepository) GetByID(id string) (*Topology, error) {
	return r.GetByIDAndUserID(id, nil)
}

func (r *PostgresRepository) GetByIDAndUserID(id string, userID *string) (*Topology, error) {
	var query string
	var args []interface{}

	if userID == nil {
		// 無用戶ID檢查（允許訪問任何拓樸，用於 demo 模式）
		query = `SELECT id, user_id, name, description, profile_type, nodes, lines, created_at, updated_at
		         FROM topologies WHERE id = $1`
		args = []interface{}{id}
	} else {
		// 檢查用戶ID（註冊用戶只能訪問自己的拓樸）
		query = `SELECT id, user_id, name, description, profile_type, nodes, lines, created_at, updated_at
		         FROM topologies WHERE id = $1 AND (user_id = $2 OR user_id IS NULL)`
		args = []interface{}{id, *userID}
	}

	var topology Topology
	var nodesJSON, linesJSON []byte
	var userIDPtr sql.NullString

	err := r.db.QueryRow(query, args...).Scan(
		&topology.ID,
		&userIDPtr,
		&topology.Name,
		&topology.Description,
		&topology.ProfileType,
		&nodesJSON,
		&linesJSON,
		&topology.CreatedAt,
		&topology.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrTopologyNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get topology: %w", err)
	}

	if userIDPtr.Valid {
		userIDStr := userIDPtr.String
		topology.UserID = &userIDStr
	}

	// 反序列化 nodes 和 lines
	if err := json.Unmarshal(nodesJSON, &topology.Nodes); err != nil {
		return nil, fmt.Errorf("failed to unmarshal nodes: %w", err)
	}

	if err := json.Unmarshal(linesJSON, &topology.Lines); err != nil {
		return nil, fmt.Errorf("failed to unmarshal lines: %w", err)
	}

	return &topology, nil
}

func (r *PostgresRepository) Update(id string, topology *Topology) error {
	topology.UpdatedAt = time.Now()

	// 序列化 nodes 和 lines
	nodesJSON, err := json.Marshal(topology.Nodes)
	if err != nil {
		return fmt.Errorf("failed to marshal nodes: %w", err)
	}

	linesJSON, err := json.Marshal(topology.Lines)
	if err != nil {
		return fmt.Errorf("failed to marshal lines: %w", err)
	}

	query := `
		UPDATE topologies
		SET name = $1, description = $2, profile_type = $3, nodes = $4, lines = $5, updated_at = $6
		WHERE id = $7
	`

	result, err := r.db.Exec(query,
		topology.Name,
		topology.Description,
		topology.ProfileType,
		nodesJSON,
		linesJSON,
		topology.UpdatedAt,
		id,
	)

	if err != nil {
		return fmt.Errorf("failed to update topology: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrTopologyNotFound
	}

	topology.ID = id
	return nil
}

func (r *PostgresRepository) Delete(id string) error {
	query := `DELETE FROM topologies WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete topology: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrTopologyNotFound
	}

	return nil
}

func (r *PostgresRepository) List() ([]*Topology, error) {
	return r.ListByUserID(nil)
}

func (r *PostgresRepository) ListByUserID(userID *string) ([]*Topology, error) {
	var query string
	var args []interface{}

	if userID == nil {
		// 列出所有拓樸（demo 模式）
		query = `SELECT id, user_id, name, description, profile_type, nodes, lines, created_at, updated_at
		         FROM topologies ORDER BY created_at DESC`
		args = []interface{}{}
	} else {
		// 列出用戶的拓樸
		query = `SELECT id, user_id, name, description, profile_type, nodes, lines, created_at, updated_at
		         FROM topologies WHERE user_id = $1 ORDER BY created_at DESC`
		args = []interface{}{*userID}
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query topologies: %w", err)
	}
	defer rows.Close()

	topologies := []*Topology{}
	for rows.Next() {
		var topology Topology
		var nodesJSON, linesJSON []byte
		var userIDPtr sql.NullString

		err := rows.Scan(
			&topology.ID,
			&userIDPtr,
			&topology.Name,
			&topology.Description,
			&topology.ProfileType,
			&nodesJSON,
			&linesJSON,
			&topology.CreatedAt,
			&topology.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan topology: %w", err)
		}

		if userIDPtr.Valid {
			userIDStr := userIDPtr.String
			topology.UserID = &userIDStr
		}

		// 反序列化 nodes 和 lines
		if err := json.Unmarshal(nodesJSON, &topology.Nodes); err != nil {
			return nil, fmt.Errorf("failed to unmarshal nodes: %w", err)
		}

		if err := json.Unmarshal(linesJSON, &topology.Lines); err != nil {
			return nil, fmt.Errorf("failed to unmarshal lines: %w", err)
		}

		topologies = append(topologies, &topology)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return topologies, nil
}

func (r *PostgresRepository) CountByUserID(userID *string) (int, error) {
	var query string
	var args []interface{}

	if userID == nil {
		// 統計所有拓樸（demo 模式）
		query = `SELECT COUNT(*) FROM topologies WHERE user_id IS NULL`
		args = []interface{}{}
	} else {
		// 統計用戶的拓樸
		query = `SELECT COUNT(*) FROM topologies WHERE user_id = $1`
		args = []interface{}{*userID}
	}

	var count int
	err := r.db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count topologies: %w", err)
	}

	return count, nil
}

