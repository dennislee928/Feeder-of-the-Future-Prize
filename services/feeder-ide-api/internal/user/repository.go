package user

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/feeder-platform/feeder-ide-api/internal/database"
	"github.com/google/uuid"
)

// Repository 用戶資料存取介面
type Repository interface {
	// User CRUD
	CreateUser(user *User) error
	GetUserByID(id string) (*User, error)
	GetUserByEmail(email string) (*User, error)
	UpdateUser(user *User) error

	// OAuth
	CreateOrUpdateOAuth(oauth *UserOAuth) error
	GetOAuthByProvider(provider, providerUserID string) (*UserOAuth, error)
	GetOAuthByUserID(userID string) ([]*UserOAuth, error)

	// Subscription
	CreateSubscription(sub *Subscription) error
	GetSubscriptionByID(id string) (*Subscription, error)
	GetActiveSubscriptionByUserID(userID string) (*Subscription, error)
	GetSubscriptionsByProviderID(provider, providerSubscriptionID string) ([]*Subscription, error)
	UpdateSubscription(sub *Subscription) error

	// Quota
	CreateOrUpdateQuota(quota *UserQuota) error
	GetQuotaByUserID(userID string) (*UserQuota, error)
	UpdateQuota(quota *UserQuota) error

	// Payment
	CreatePayment(payment *Payment) error
	GetPaymentByID(id string) (*Payment, error)
	GetPaymentsByUserID(userID string) ([]*Payment, error)
	UpdatePayment(payment *Payment) error
}

// PostgresUserRepository PostgreSQL 實作
type PostgresUserRepository struct {
	db *sql.DB
}

// NewPostgresUserRepository 建立新的 PostgreSQL user repository
func NewPostgresUserRepository() (*PostgresUserRepository, error) {
	if database.DB == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	return &PostgresUserRepository{
		db: database.DB,
	}, nil
}

// User CRUD
func (r *PostgresUserRepository) CreateUser(user *User) error {
	if user.ID == "" {
		user.ID = uuid.New().String()
	}

	now := time.Now()
	if user.CreatedAt.IsZero() {
		user.CreatedAt = now
	}
	user.UpdatedAt = now

	query := `
		INSERT INTO users (id, email, name, avatar_url, subscription_tier, subscription_status, subscription_expires_at, api_key, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := r.db.Exec(query,
		user.ID,
		user.Email,
		user.Name,
		user.AvatarURL,
		user.SubscriptionTier,
		user.SubscriptionStatus,
		user.SubscriptionExpiresAt,
		user.APIKey,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *PostgresUserRepository) GetUserByID(id string) (*User, error) {
	query := `SELECT id, email, name, avatar_url, subscription_tier, subscription_status, subscription_expires_at, api_key, created_at, updated_at
	          FROM users WHERE id = $1`

	var user User
	var namePtr, avatarURLPtr, apiKeyPtr sql.NullString
	var expiresAtPtr sql.NullTime

	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&namePtr,
		&avatarURLPtr,
		&user.SubscriptionTier,
		&user.SubscriptionStatus,
		&expiresAtPtr,
		&apiKeyPtr,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if namePtr.Valid {
		user.Name = &namePtr.String
	}
	if avatarURLPtr.Valid {
		user.AvatarURL = &avatarURLPtr.String
	}
	if apiKeyPtr.Valid {
		user.APIKey = &apiKeyPtr.String
	}
	if expiresAtPtr.Valid {
		user.SubscriptionExpiresAt = &expiresAtPtr.Time
	}

	return &user, nil
}

func (r *PostgresUserRepository) GetUserByEmail(email string) (*User, error) {
	query := `SELECT id, email, name, avatar_url, subscription_tier, subscription_status, subscription_expires_at, api_key, created_at, updated_at
	          FROM users WHERE email = $1`

	var user User
	var namePtr, avatarURLPtr, apiKeyPtr sql.NullString
	var expiresAtPtr sql.NullTime

	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&namePtr,
		&avatarURLPtr,
		&user.SubscriptionTier,
		&user.SubscriptionStatus,
		&expiresAtPtr,
		&apiKeyPtr,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if namePtr.Valid {
		user.Name = &namePtr.String
	}
	if avatarURLPtr.Valid {
		user.AvatarURL = &avatarURLPtr.String
	}
	if apiKeyPtr.Valid {
		user.APIKey = &apiKeyPtr.String
	}
	if expiresAtPtr.Valid {
		user.SubscriptionExpiresAt = &expiresAtPtr.Time
	}

	return &user, nil
}

func (r *PostgresUserRepository) UpdateUser(user *User) error {
	user.UpdatedAt = time.Now()

	query := `
		UPDATE users
		SET email = $1, name = $2, avatar_url = $3, subscription_tier = $4, subscription_status = $5,
		    subscription_expires_at = $6, api_key = $7, updated_at = $8
		WHERE id = $9
	`

	result, err := r.db.Exec(query,
		user.Email,
		user.Name,
		user.AvatarURL,
		user.SubscriptionTier,
		user.SubscriptionStatus,
		user.SubscriptionExpiresAt,
		user.APIKey,
		user.UpdatedAt,
		user.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// OAuth
func (r *PostgresUserRepository) CreateOrUpdateOAuth(oauth *UserOAuth) error {
	if oauth.ID == "" {
		oauth.ID = uuid.New().String()
	}

	now := time.Now()
	if oauth.CreatedAt.IsZero() {
		oauth.CreatedAt = now
	}
	oauth.UpdatedAt = now

	query := `
		INSERT INTO user_oauth (id, user_id, provider, provider_user_id, access_token, refresh_token, token_expires_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (provider, provider_user_id)
		DO UPDATE SET
			user_id = EXCLUDED.user_id,
			access_token = EXCLUDED.access_token,
			refresh_token = EXCLUDED.refresh_token,
			token_expires_at = EXCLUDED.token_expires_at,
			updated_at = EXCLUDED.updated_at
	`

	_, err := r.db.Exec(query,
		oauth.ID,
		oauth.UserID,
		oauth.Provider,
		oauth.ProviderUserID,
		oauth.AccessToken,
		oauth.RefreshToken,
		oauth.TokenExpiresAt,
		oauth.CreatedAt,
		oauth.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create or update oauth: %w", err)
	}

	return nil
}

func (r *PostgresUserRepository) GetOAuthByProvider(provider, providerUserID string) (*UserOAuth, error) {
	query := `SELECT id, user_id, provider, provider_user_id, access_token, refresh_token, token_expires_at, created_at, updated_at
	          FROM user_oauth WHERE provider = $1 AND provider_user_id = $2`

	var oauth UserOAuth
	var accessTokenPtr, refreshTokenPtr sql.NullString
	var expiresAtPtr sql.NullTime

	err := r.db.QueryRow(query, provider, providerUserID).Scan(
		&oauth.ID,
		&oauth.UserID,
		&oauth.Provider,
		&oauth.ProviderUserID,
		&accessTokenPtr,
		&refreshTokenPtr,
		&expiresAtPtr,
		&oauth.CreatedAt,
		&oauth.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("oauth not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get oauth: %w", err)
	}

	if accessTokenPtr.Valid {
		oauth.AccessToken = &accessTokenPtr.String
	}
	if refreshTokenPtr.Valid {
		oauth.RefreshToken = &refreshTokenPtr.String
	}
	if expiresAtPtr.Valid {
		oauth.TokenExpiresAt = &expiresAtPtr.Time
	}

	return &oauth, nil
}

func (r *PostgresUserRepository) GetOAuthByUserID(userID string) ([]*UserOAuth, error) {
	query := `SELECT id, user_id, provider, provider_user_id, access_token, refresh_token, token_expires_at, created_at, updated_at
	          FROM user_oauth WHERE user_id = $1`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query oauth: %w", err)
	}
	defer rows.Close()

	oauths := []*UserOAuth{}
	for rows.Next() {
		var oauth UserOAuth
		var accessTokenPtr, refreshTokenPtr sql.NullString
		var expiresAtPtr sql.NullTime

		err := rows.Scan(
			&oauth.ID,
			&oauth.UserID,
			&oauth.Provider,
			&oauth.ProviderUserID,
			&accessTokenPtr,
			&refreshTokenPtr,
			&expiresAtPtr,
			&oauth.CreatedAt,
			&oauth.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan oauth: %w", err)
		}

		if accessTokenPtr.Valid {
			oauth.AccessToken = &accessTokenPtr.String
		}
		if refreshTokenPtr.Valid {
			oauth.RefreshToken = &refreshTokenPtr.String
		}
		if expiresAtPtr.Valid {
			oauth.TokenExpiresAt = &expiresAtPtr.Time
		}

		oauths = append(oauths, &oauth)
	}

	return oauths, nil
}

// Subscription
func (r *PostgresUserRepository) CreateSubscription(sub *Subscription) error {
	if sub.ID == "" {
		sub.ID = uuid.New().String()
	}

	now := time.Now()
	if sub.CreatedAt.IsZero() {
		sub.CreatedAt = now
	}
	sub.UpdatedAt = now

	query := `
		INSERT INTO subscriptions (id, user_id, tier, status, payment_provider, payment_subscription_id,
		                           current_period_start, current_period_end, cancel_at_period_end, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	_, err := r.db.Exec(query,
		sub.ID,
		sub.UserID,
		sub.Tier,
		sub.Status,
		sub.PaymentProvider,
		sub.PaymentSubscriptionID,
		sub.CurrentPeriodStart,
		sub.CurrentPeriodEnd,
		sub.CancelAtPeriodEnd,
		sub.CreatedAt,
		sub.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create subscription: %w", err)
	}

	return nil
}

func (r *PostgresUserRepository) GetSubscriptionByID(id string) (*Subscription, error) {
	query := `SELECT id, user_id, tier, status, payment_provider, payment_subscription_id,
	                 current_period_start, current_period_end, cancel_at_period_end, created_at, updated_at
	          FROM subscriptions WHERE id = $1`

	var sub Subscription
	var providerPtr, subscriptionIDPtr sql.NullString
	var startPtr, endPtr sql.NullTime

	err := r.db.QueryRow(query, id).Scan(
		&sub.ID,
		&sub.UserID,
		&sub.Tier,
		&sub.Status,
		&providerPtr,
		&subscriptionIDPtr,
		&startPtr,
		&endPtr,
		&sub.CancelAtPeriodEnd,
		&sub.CreatedAt,
		&sub.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("subscription not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	if providerPtr.Valid {
		sub.PaymentProvider = &providerPtr.String
	}
	if subscriptionIDPtr.Valid {
		sub.PaymentSubscriptionID = &subscriptionIDPtr.String
	}
	if startPtr.Valid {
		sub.CurrentPeriodStart = &startPtr.Time
	}
	if endPtr.Valid {
		sub.CurrentPeriodEnd = &endPtr.Time
	}

	return &sub, nil
}

func (r *PostgresUserRepository) GetActiveSubscriptionByUserID(userID string) (*Subscription, error) {
	query := `SELECT id, user_id, tier, status, payment_provider, payment_subscription_id,
	                 current_period_start, current_period_end, cancel_at_period_end, created_at, updated_at
	          FROM subscriptions WHERE user_id = $1 AND status = 'active' ORDER BY created_at DESC LIMIT 1`

	var sub Subscription
	var providerPtr, subscriptionIDPtr sql.NullString
	var startPtr, endPtr sql.NullTime

	err := r.db.QueryRow(query, userID).Scan(
		&sub.ID,
		&sub.UserID,
		&sub.Tier,
		&sub.Status,
		&providerPtr,
		&subscriptionIDPtr,
		&startPtr,
		&endPtr,
		&sub.CancelAtPeriodEnd,
		&sub.CreatedAt,
		&sub.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // 沒有訂閱是正常的
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	if providerPtr.Valid {
		sub.PaymentProvider = &providerPtr.String
	}
	if subscriptionIDPtr.Valid {
		sub.PaymentSubscriptionID = &subscriptionIDPtr.String
	}
	if startPtr.Valid {
		sub.CurrentPeriodStart = &startPtr.Time
	}
	if endPtr.Valid {
		sub.CurrentPeriodEnd = &endPtr.Time
	}

	return &sub, nil
}

func (r *PostgresUserRepository) GetSubscriptionsByProviderID(provider, providerSubscriptionID string) ([]*Subscription, error) {
	query := `SELECT id, user_id, tier, status, payment_provider, payment_subscription_id,
	                 current_period_start, current_period_end, cancel_at_period_end, created_at, updated_at
	          FROM subscriptions WHERE payment_provider = $1 AND payment_subscription_id = $2`

	rows, err := r.db.Query(query, provider, providerSubscriptionID)
	if err != nil {
		return nil, fmt.Errorf("failed to query subscriptions: %w", err)
	}
	defer rows.Close()

	subscriptions := []*Subscription{}
	for rows.Next() {
		var sub Subscription
		var providerPtr, subscriptionIDPtr sql.NullString
		var startPtr, endPtr sql.NullTime

		err := rows.Scan(
			&sub.ID,
			&sub.UserID,
			&sub.Tier,
			&sub.Status,
			&providerPtr,
			&subscriptionIDPtr,
			&startPtr,
			&endPtr,
			&sub.CancelAtPeriodEnd,
			&sub.CreatedAt,
			&sub.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan subscription: %w", err)
		}

		if providerPtr.Valid {
			sub.PaymentProvider = &providerPtr.String
		}
		if subscriptionIDPtr.Valid {
			sub.PaymentSubscriptionID = &subscriptionIDPtr.String
		}
		if startPtr.Valid {
			sub.CurrentPeriodStart = &startPtr.Time
		}
		if endPtr.Valid {
			sub.CurrentPeriodEnd = &endPtr.Time
		}

		subscriptions = append(subscriptions, &sub)
	}

	return subscriptions, nil
}

func (r *PostgresUserRepository) UpdateSubscription(sub *Subscription) error {
	sub.UpdatedAt = time.Now()

	query := `
		UPDATE subscriptions
		SET tier = $1, status = $2, payment_provider = $3, payment_subscription_id = $4,
		    current_period_start = $5, current_period_end = $6, cancel_at_period_end = $7, updated_at = $8
		WHERE id = $9
	`

	result, err := r.db.Exec(query,
		sub.Tier,
		sub.Status,
		sub.PaymentProvider,
		sub.PaymentSubscriptionID,
		sub.CurrentPeriodStart,
		sub.CurrentPeriodEnd,
		sub.CancelAtPeriodEnd,
		sub.UpdatedAt,
		sub.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("subscription not found")
	}

	return nil
}

// Quota
func (r *PostgresUserRepository) CreateOrUpdateQuota(quota *UserQuota) error {
	if quota.ID == "" {
		quota.ID = uuid.New().String()
	}

	now := time.Now()
	if quota.CreatedAt.IsZero() {
		quota.CreatedAt = now
	}
	quota.UpdatedAt = now

	query := `
		INSERT INTO user_quotas (id, user_id, max_topologies, used_topologies, max_simulations_per_day,
		                         used_simulations_today, last_simulation_reset_date, can_use_3d_rendering,
		                         can_use_ai_prediction, can_use_advanced_security, can_access_api, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		ON CONFLICT (user_id)
		DO UPDATE SET
			max_topologies = EXCLUDED.max_topologies,
			used_topologies = EXCLUDED.used_topologies,
			max_simulations_per_day = EXCLUDED.max_simulations_per_day,
			used_simulations_today = EXCLUDED.used_simulations_today,
			last_simulation_reset_date = EXCLUDED.last_simulation_reset_date,
			can_use_3d_rendering = EXCLUDED.can_use_3d_rendering,
			can_use_ai_prediction = EXCLUDED.can_use_ai_prediction,
			can_use_advanced_security = EXCLUDED.can_use_advanced_security,
			can_access_api = EXCLUDED.can_access_api,
			updated_at = EXCLUDED.updated_at
	`

	_, err := r.db.Exec(query,
		quota.ID,
		quota.UserID,
		quota.MaxTopologies,
		quota.UsedTopologies,
		quota.MaxSimulationsPerDay,
		quota.UsedSimulationsToday,
		quota.LastSimulationResetDate,
		quota.CanUse3DRendering,
		quota.CanUseAIPrediction,
		quota.CanUseAdvancedSecurity,
		quota.CanAccessAPI,
		quota.CreatedAt,
		quota.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create or update quota: %w", err)
	}

	return nil
}

func (r *PostgresUserRepository) GetQuotaByUserID(userID string) (*UserQuota, error) {
	query := `SELECT id, user_id, max_topologies, used_topologies, max_simulations_per_day,
	                 used_simulations_today, last_simulation_reset_date, can_use_3d_rendering,
	                 can_use_ai_prediction, can_use_advanced_security, can_access_api, created_at, updated_at
	          FROM user_quotas WHERE user_id = $1`

	var quota UserQuota
	err := r.db.QueryRow(query, userID).Scan(
		&quota.ID,
		&quota.UserID,
		&quota.MaxTopologies,
		&quota.UsedTopologies,
		&quota.MaxSimulationsPerDay,
		&quota.UsedSimulationsToday,
		&quota.LastSimulationResetDate,
		&quota.CanUse3DRendering,
		&quota.CanUseAIPrediction,
		&quota.CanUseAdvancedSecurity,
		&quota.CanAccessAPI,
		&quota.CreatedAt,
		&quota.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // 沒有配額是正常的，會使用默認值
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get quota: %w", err)
	}

	return &quota, nil
}

func (r *PostgresUserRepository) UpdateQuota(quota *UserQuota) error {
	quota.UpdatedAt = time.Now()

	query := `
		UPDATE user_quotas
		SET max_topologies = $1, used_topologies = $2, max_simulations_per_day = $3,
		    used_simulations_today = $4, last_simulation_reset_date = $5, can_use_3d_rendering = $6,
		    can_use_ai_prediction = $7, can_use_advanced_security = $8, can_access_api = $9, updated_at = $10
		WHERE id = $11
	`

	result, err := r.db.Exec(query,
		quota.MaxTopologies,
		quota.UsedTopologies,
		quota.MaxSimulationsPerDay,
		quota.UsedSimulationsToday,
		quota.LastSimulationResetDate,
		quota.CanUse3DRendering,
		quota.CanUseAIPrediction,
		quota.CanUseAdvancedSecurity,
		quota.CanAccessAPI,
		quota.UpdatedAt,
		quota.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update quota: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("quota not found")
	}

	return nil
}

// Payment
func (r *PostgresUserRepository) CreatePayment(payment *Payment) error {
	if payment.ID == "" {
		payment.ID = uuid.New().String()
	}

	now := time.Now()
	if payment.CreatedAt.IsZero() {
		payment.CreatedAt = now
	}
	payment.UpdatedAt = now

	metadataJSON, _ := json.Marshal(payment.Metadata)

	query := `
		INSERT INTO payments (id, user_id, subscription_id, amount, currency, payment_provider,
		                      payment_provider_id, status, metadata, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	var subscriptionIDPtr *string
	if payment.SubscriptionID != nil {
		subscriptionIDPtr = payment.SubscriptionID
	}

	_, err := r.db.Exec(query,
		payment.ID,
		payment.UserID,
		subscriptionIDPtr,
		payment.Amount,
		payment.Currency,
		payment.PaymentProvider,
		payment.PaymentProviderID,
		payment.Status,
		metadataJSON,
		payment.CreatedAt,
		payment.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create payment: %w", err)
	}

	return nil
}

func (r *PostgresUserRepository) GetPaymentByID(id string) (*Payment, error) {
	query := `SELECT id, user_id, subscription_id, amount, currency, payment_provider,
	                 payment_provider_id, status, metadata, created_at, updated_at
	          FROM payments WHERE id = $1`

	var payment Payment
	var subscriptionIDPtr sql.NullString
	var metadataJSON []byte

	err := r.db.QueryRow(query, id).Scan(
		&payment.ID,
		&payment.UserID,
		&subscriptionIDPtr,
		&payment.Amount,
		&payment.Currency,
		&payment.PaymentProvider,
		&payment.PaymentProviderID,
		&payment.Status,
		&metadataJSON,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("payment not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	if subscriptionIDPtr.Valid {
		payment.SubscriptionID = &subscriptionIDPtr.String
	}

	if len(metadataJSON) > 0 {
		json.Unmarshal(metadataJSON, &payment.Metadata)
	}

	return &payment, nil
}

func (r *PostgresUserRepository) GetPaymentsByUserID(userID string) ([]*Payment, error) {
	query := `SELECT id, user_id, subscription_id, amount, currency, payment_provider,
	                 payment_provider_id, status, metadata, created_at, updated_at
	          FROM payments WHERE user_id = $1 ORDER BY created_at DESC`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query payments: %w", err)
	}
	defer rows.Close()

	payments := []*Payment{}
	for rows.Next() {
		var payment Payment
		var subscriptionIDPtr sql.NullString
		var metadataJSON []byte

		err := rows.Scan(
			&payment.ID,
			&payment.UserID,
			&subscriptionIDPtr,
			&payment.Amount,
			&payment.Currency,
			&payment.PaymentProvider,
			&payment.PaymentProviderID,
			&payment.Status,
			&metadataJSON,
			&payment.CreatedAt,
			&payment.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan payment: %w", err)
		}

		if subscriptionIDPtr.Valid {
			payment.SubscriptionID = &subscriptionIDPtr.String
		}

		if len(metadataJSON) > 0 {
			json.Unmarshal(metadataJSON, &payment.Metadata)
		}

		payments = append(payments, &payment)
	}

	return payments, nil
}

func (r *PostgresUserRepository) UpdatePayment(payment *Payment) error {
	payment.UpdatedAt = time.Now()

	metadataJSON, _ := json.Marshal(payment.Metadata)

	query := `
		UPDATE payments
		SET status = $1, metadata = $2, updated_at = $3
		WHERE id = $4
	`

	result, err := r.db.Exec(query,
		payment.Status,
		metadataJSON,
		payment.UpdatedAt,
		payment.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update payment: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("payment not found")
	}

	return nil
}

