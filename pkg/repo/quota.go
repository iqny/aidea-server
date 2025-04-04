package repo

import (
	"context"
	"database/sql"
	"encoding/json"
	model2 "github.com/mylxsw/aidea-server/pkg/repo/model"
	"time"

	"github.com/mylxsw/aidea-server/config"
	"github.com/mylxsw/asteria/log"
	"github.com/mylxsw/eloquent"
	"github.com/mylxsw/eloquent/query"
	"github.com/mylxsw/go-utils/array"
	"gopkg.in/guregu/null.v3"
)

// QuotaRepo 用户配额仓库
type QuotaRepo struct {
	db   *sql.DB
	conf *config.Config
}

// NewQuotaRepo create a new QuotaRepo
func NewQuotaRepo(db *sql.DB, conf *config.Config) *QuotaRepo {
	return &QuotaRepo{db: db, conf: conf}
}

// AddUserQuota 创建用户配额
func (repo *QuotaRepo) AddUserQuota(ctx context.Context, userID int64, quotaValue int64, endAt time.Time, note, paymentID string) (int64, error) {
	quota := model2.Quota{
		UserId:        userID,
		Quota:         quotaValue,
		Rest:          quotaValue,
		Note:          note,
		PaymentId:     paymentID,
		PeriodStartAt: NowInDate(),
		PeriodEndAt:   TimeInDate(endAt),
	}

	return model2.NewQuotaModel(repo.db).Save(ctx, quota.ToQuotaN(
		model2.FieldQuotaUserId,
		model2.FieldQuotaQuota,
		model2.FieldQuotaRest,
		model2.FieldQuotaNote,
		model2.FieldQuotaPaymentId,
		model2.FieldQuotaPeriodStartAt,
		model2.FieldQuotaPeriodEndAt,
	))
}

// TimeInDate 获取时间的日期部分
func TimeInDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).AddDate(0, 0, 1)
}

// NowInDate 获取当前时间的日期部分
func NowInDate() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}

// Quota 配额详情
type Quota struct {
	Id            int64     `json:"id"`
	UserId        int64     `json:"user_id"`
	Quota         int64     `json:"quota"`
	Rest          int64     `json:"rest"`
	Note          string    `json:"note"`
	PaymentId     string    `json:"payment_id"`
	PeriodStartAt time.Time `json:"period_start_at"`
	PeriodEndAt   time.Time `json:"period_end_at"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Expired       bool      `json:"expired"`
}

// GetUserQuotaDetails 获取用户配额详情
func (repo *QuotaRepo) GetUserQuotaDetails(ctx context.Context, userID int64) ([]Quota, error) {
	// 查询当前用户的所有配额
	q := query.Builder().
		Where(model2.FieldQuotaUserId, userID).
		Limit(100).
		OrderBy(model2.FieldQuotaId, "DESC")

	res, err := model2.NewQuotaModel(repo.db).Get(ctx, q)
	if err != nil {
		return nil, err
	}

	return array.Map(res, func(qn model2.QuotaN, _ int) Quota {
		return Quota{
			Id:            qn.Id.ValueOrZero(),
			UserId:        qn.UserId.ValueOrZero(),
			Quota:         qn.Quota.ValueOrZero(),
			Rest:          qn.Rest.ValueOrZero(),
			Note:          qn.Note.ValueOrZero(),
			PaymentId:     qn.PaymentId.ValueOrZero(),
			PeriodStartAt: qn.PeriodStartAt.ValueOrZero(),
			PeriodEndAt:   qn.PeriodEndAt.ValueOrZero(),
			CreatedAt:     qn.CreatedAt.ValueOrZero(),
			UpdatedAt:     qn.UpdatedAt.ValueOrZero(),
			Expired:       qn.PeriodEndAt.ValueOrZero().Before(time.Now()),
		}
	}), nil
}

// QuotaSummary 配额汇总
type QuotaSummary struct {
	Quota int64 `json:"quota"`
	Rest  int64 `json:"rest"`
	Used  int64 `json:"used"`
}

// GetUserQuota 获取用户配额
func (repo *QuotaRepo) GetUserQuota(ctx context.Context, userID int64) (*QuotaSummary, error) {
	q := query.Builder().
		Table(model2.QuotaTable()).
		Select(
			query.Raw("SUM(quota) AS quota"),
			query.Raw("SUM(rest) AS rest"),
		).
		Where(model2.FieldQuotaUserId, userID).
		Where(model2.FieldQuotaPeriodEndAt, ">", time.Now())

	quotas, err := eloquent.Query(ctx, repo.db, q, func(row eloquent.Scanner) (*QuotaSummary, error) {
		var quota QuotaSummary
		var v1, v2 sql.NullInt64
		if err := row.Scan(&v1, &v2); err != nil {
			return nil, err
		}

		quota.Quota = v1.Int64
		quota.Rest = v2.Int64

		quota.Used = quota.Quota - quota.Rest
		return &quota, nil
	})
	if err != nil {
		return nil, err
	}

	return quotas[0], nil
}

type QuotaUsedMeta struct {
	Models      []string `json:"models"`
	Tag         string   `json:"tag"`
	InputToken  int      `json:"input_token,omitempty"`
	OutputToken int      `json:"output_token,omitempty"`
	InputPrice  float64  `json:"input_price,omitempty"`
	OutputPrice float64  `json:"output_price,omitempty"`
	ReqPrice    int64    `json:"req_price,omitempty"`
	SearchPrice int64    `json:"search_price,omitempty"`
}

func NewQuotaUsedMeta(tag string, models ...string) QuotaUsedMeta {
	return QuotaUsedMeta{
		Models: models,
		Tag:    tag,
	}
}

// QuotaConsume 更新用户配额已使用量
func (repo *QuotaRepo) QuotaConsume(ctx context.Context, userID int64, used int64, meta QuotaUsedMeta) error {
	relatedQuotaIds := make(map[int64]int64)
	var debt int64

	err := eloquent.Transaction(repo.db, func(tx query.Database) error {
		usedVar := used
		// 查询当前可用配额
		q := query.Builder().
			Where(model2.FieldQuotaUserId, userID).
			Where(model2.FieldQuotaRest, ">", 0).
			Where(model2.FieldQuotaPeriodEndAt, ">", time.Now()).
			OrderBy(model2.FieldQuotaPeriodEndAt, "ASC")
		quotas, err := model2.NewQuotaModel(tx).Get(ctx, q)
		if err != nil {
			return err
		}

		for _, quota := range quotas {
			quotaID := quota.Id.ValueOrZero()
			rest := quota.Rest.ValueOrZero()
			if rest >= usedVar {
				relatedQuotaIds[quotaID] = usedVar
				// 当前配额足够，直接更新配额
				_, err := tx.ExecContext(ctx, "UPDATE quota SET rest = rest - ? WHERE id = ?", usedVar, quotaID)
				return err
			}

			relatedQuotaIds[quotaID] = rest

			// 当前配额不够，更新配额为 0
			_, err := tx.ExecContext(ctx, "UPDATE quota SET rest = 0 WHERE id = ?", quotaID)
			if err != nil {
				return err
			}

			// 更新已使用量
			usedVar -= rest
		}

		// 没有配额了，创建欠费记录
		if usedVar > 0 {
			debt = usedVar
			if _, err := model2.NewDebtModel(tx).Create(ctx, query.KV{
				model2.FieldDebtUserId: userID,
				model2.FieldDebtUsed:   usedVar,
			}); err != nil {
				return err
			}

			return nil
		}

		return nil
	})

	if err == nil {
		log.F(log.M{
			"user_id":   userID,
			"used":      used,
			"quota_ids": relatedQuotaIds,
			"debt":      debt,
			"meta":      meta,
		}).Info("user quota consumed")

		quotaIdsBytes, _ := json.Marshal(relatedQuotaIds)
		metaBytes, _ := json.Marshal(meta)

		if _, err := model2.NewQuotaUsageModel(repo.db).Save(ctx, model2.QuotaUsageN{
			UserId:   null.IntFrom(userID),
			Used:     null.IntFrom(used),
			QuotaIds: null.StringFrom(string(quotaIdsBytes)),
			Debt:     null.IntFrom(debt),
			Meta:     null.StringFrom(string(metaBytes)),
		}); err != nil {
			log.F(log.M{"user_id": userID, "err": err}).Error("save quota usage failed")
		}
	}

	return err
}

// GetQuotaStatisticsRecently 获取近期的配额使用统计
func (repo *QuotaRepo) GetQuotaStatisticsRecently(ctx context.Context, userId int64, days int64) ([]model2.QuotaStatistics, error) {
	q := query.Builder().
		Table(model2.QuotaStatisticsTable()).
		Where(model2.FieldQuotaStatisticsUserId, userId).
		Where(model2.FieldQuotaStatisticsCreatedAt, ">=", time.Now().AddDate(0, 0, -int(days))).
		OrderBy(model2.FieldQuotaStatisticsId, "DESC")

	res, err := model2.NewQuotaStatisticsModel(repo.db).Get(ctx, q)
	if err != nil {
		return nil, err
	}

	return array.Map(res, func(item model2.QuotaStatisticsN, _ int) model2.QuotaStatistics {
		return item.ToQuotaStatistics()
	}), nil
}

type QuotaUsage struct {
	model2.QuotaUsage
	QuotaMeta QuotaUsedMeta `json:"quota_meta,omitempty"`
}

// GetQuotaDetails 获取配额使用详情
func (repo *QuotaRepo) GetQuotaDetails(ctx context.Context, userId int64, startAt, endAt time.Time) ([]QuotaUsage, error) {
	q := query.Builder().
		Where(model2.FieldQuotaUsageUserId, userId).
		Where(model2.FieldQuotaUsageCreatedAt, ">=", startAt.Format("2006-01-02 15:04:05")).
		Where(model2.FieldQuotaUsageCreatedAt, "<", endAt.Format("2006-01-02 15:04:05")).
		OrderBy(model2.FieldQuotaUsageId, "DESC")

	res, err := model2.NewQuotaUsageModel(repo.db).Get(ctx, q)
	if err != nil {
		return nil, err
	}

	return array.Map(res, func(item model2.QuotaUsageN, _ int) QuotaUsage {
		var quotaMeta QuotaUsedMeta
		_ = json.Unmarshal([]byte(item.Meta.ValueOrZero()), &quotaMeta)
		return QuotaUsage{
			QuotaUsage: item.ToQuotaUsage(),
			QuotaMeta:  quotaMeta,
		}
	}), nil
}
