package service

import (
	"crypto-temka/internal/models"
	"crypto-temka/internal/utils"
	"crypto-temka/pkg/log"
	"database/sql"
	"errors"
	"fmt"
	"github.com/guregu/null"
	"github.com/jmoiron/sqlx"
	"math"
	"time"
)

type worker struct {
	db     *sqlx.DB
	logger *log.Logs
}

func InitUserRateWorker(db *sqlx.DB, logger *log.Logs) {
	wrk := worker{
		db:     db,
		logger: logger,
	}

	go wrk.startWorker()
}

func waitRate() {
	// todo: debug mode
	//time.Sleep(time.Hour * 12)
	time.Sleep(time.Second * 5)
}

func (w *worker) startWorker() {
	var shouldWait = false
	for {
		if shouldWait {
			shouldWait = false
			waitRate()
		}
		tx, err := w.db.Begin()
		if err != nil {
			w.logger.Error(err.Error())
			continue
		}

		// todo: check last_updated правильность
		row := tx.QueryRow(`SELECT ur.id, user_id, rate_id, lock, last_updated, opened, deposit, earned_pool, 
       next_day_charge, outcome_pool, token, profit FROM users_rates ur JOIN rates r on r.id = ur.rate_id 
                                                   WHERE last_updated < current_date ORDER BY last_updated LIMIT 1 FOR UPDATE`)

		var ur models.UserRate
		var lastUpdated time.Time
		var nextDayCharge null.Float
		var profit int
		err = row.Scan(&ur.ID, &ur.UserID, &ur.RateID, &ur.Lock, &lastUpdated, &ur.Opened, &ur.Deposit, &ur.EarnedPool,
			&nextDayCharge, &ur.OutcomePool, &ur.Token, &profit)
		if err != nil {
			rbErr := tx.Rollback()
			if errors.Is(err, sql.ErrNoRows) {
				shouldWait = true
				//waitRate()
			}
			if rbErr != nil {
				w.logger.Error(fmt.Sprintf("err: %v, rbErr: %v", err, rbErr))
				continue
			}
			w.logger.Error(err.Error())
			continue
		}

		// earned_pool -> outcome_pool
		// todo: могут быть интересные вещи при получении даты оттуда
		ur.Lock = utils.DateOnly(ur.Lock)
		todayDate := utils.DateOnly(time.Now())
		if todayDate.Sub(ur.Lock) >= 0 {
			fmt.Println("1")
			if todayDate.Sub(ur.Opened)%(time.Hour*24*7) == 0 {
				fmt.Println("2")

				ur.OutcomePool += ur.EarnedPool
				ur.EarnedPool = 0

				_, err = tx.Exec(`UPDATE users_rates SET outcome_pool = $1, earned_pool = 0 WHERE id = $2`,
					ur.OutcomePool+ur.EarnedPool, ur.ID)

				if err != nil {
					rbErr := tx.Rollback()
					if rbErr != nil {
						w.logger.Error(fmt.Sprintf("err: %v, rbErr: %v", err, rbErr))
						continue
					}
					w.logger.Error(err.Error())
					continue
				}
			}
		}

		// increase money. nextDayCharge развилка
		if nextDayCharge.Valid {
			if nextDayCharge.Float64 > 0 {
				ur.EarnedPool += nextDayCharge.Float64
			} else if nextDayCharge.Float64 < 0 {
				if math.Abs(nextDayCharge.Float64) > ur.EarnedPool {
					nextDayCharge.Float64 += ur.EarnedPool
					ur.EarnedPool = 0
					if math.Abs(nextDayCharge.Float64) > ur.OutcomePool {
						nextDayCharge.Float64 += ur.OutcomePool
						ur.OutcomePool = 0
						if !utils.FloatEquals(nextDayCharge.Float64, 0) {
							ur.Deposit += nextDayCharge.Float64
						}
					} else {
						ur.OutcomePool += nextDayCharge.Float64
						nextDayCharge.Float64 = 0
					}
				} else {
					ur.EarnedPool += nextDayCharge.Float64
					nextDayCharge.Float64 = 0
				}
			}

			_, err = tx.Exec(`UPDATE users_rates SET next_day_charge = null WHERE id = $1`, ur.ID)

			if err != nil {
				rbErr := tx.Rollback()
				if rbErr != nil {
					w.logger.Error(fmt.Sprintf("err: %v, rbErr: %v", err, rbErr))
					continue
				}
				w.logger.Error(err.Error())
				continue
			}
		} else {
			ur.EarnedPool += (ur.EarnedPool + ur.OutcomePool + ur.Deposit) * (float64(profit) / 100)
		}

		// update db to current money changes
		_, err = tx.Exec(`UPDATE users_rates SET deposit = $1, earned_pool = $2, outcome_pool = $3 WHERE id = $4`,
			ur.Deposit, ur.EarnedPool, ur.OutcomePool, ur.ID)
		if err != nil {
			rbErr := tx.Rollback()
			if rbErr != nil {
				w.logger.Error(fmt.Sprintf("err: %v, rbErr: %v", err, rbErr))
				continue
			}
			w.logger.Error(err.Error())
			continue
		}

		// change last updated
		_, err = tx.Exec(`UPDATE users_rates SET last_updated = current_date WHERE id = $1`, ur.ID)
		if err != nil {
			rbErr := tx.Rollback()
			if rbErr != nil {
				w.logger.Error(fmt.Sprintf("err: %v, rbErr: %v", err, rbErr))
				continue
			}
			w.logger.Error(err.Error())
			continue
		}

		err = tx.Commit()
		if err != nil {
			rbErr := tx.Rollback()
			if rbErr != nil {
				w.logger.Error(fmt.Sprintf("err: %v, rbErr: %v", err, rbErr))
				continue
			}
			w.logger.Error(err.Error())
			continue
		}
	}
}
