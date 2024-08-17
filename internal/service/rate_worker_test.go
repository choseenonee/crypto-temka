package service

import (
	"crypto-temka/internal/models"
	"crypto-temka/internal/utils"
	"crypto-temka/pkg/config"
	"crypto-temka/pkg/database"
	"fmt"
	"github.com/guregu/null"
	"github.com/jmoiron/sqlx"
	"testing"
	"time"
)

var userID int
var firstRateID int
var secondRateID int

func deleteDependencies() {
	config.InitConfig()

	db := database.MustGetDB()

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(`DELETE FROM users WHERE id = $1`, userID)

	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			panic(err.Error() + rbErr.Error())
		}
		panic(err.Error())
	}

	_, err = tx.Exec(`DELETE FROM rates WHERE id = $1`, firstRateID)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			panic(err.Error() + rbErr.Error())
		}
		panic(err.Error())
	}

	_, err = tx.Exec(`DELETE FROM rates WHERE id = $1`, secondRateID)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			panic(err.Error() + rbErr.Error())
		}
		panic(err.Error())
	}

	err = tx.Commit()
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			panic(err.Error() + rbErr.Error())
		}
		panic(err.Error())
	}
}

func initDependencies() *sqlx.DB {
	config.InitConfig()

	db := database.MustGetDB()

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	defer tx.Rollback()

	row := tx.QueryRow(`INSERT INTO users (email, phone_number, hashed_password, status, refer_id, properties) 
VALUES ('popa', 'piska', '123', 'verified', null, '"string"') RETURNING id`)

	err = row.Scan(&userID)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			panic(err.Error() + rbErr.Error())
		}
		panic(err.Error())
	}

	row = tx.QueryRow(`INSERT INTO rates (title, profit, min_lock_days, commission, properties) VALUES 
                                                                             ('first', 10, 7, 1, '"string"') RETURNING id`)

	err = row.Scan(&firstRateID)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			panic(err.Error() + rbErr.Error())
		}
		panic(err.Error())
	}

	row = tx.QueryRow(`INSERT INTO rates (title, profit, min_lock_days, commission, properties) VALUES 
                                                                             ('second', 200, 7, 1, '"string"') RETURNING id`)

	err = row.Scan(&secondRateID)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			panic(err.Error() + rbErr.Error())
		}
		panic(err.Error())
	}

	err = tx.Commit()
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			panic(err.Error() + rbErr.Error())
		}
		panic(err.Error())
	}

	return database.MustGetDB()
}

func TestWorker(t *testing.T) {
	TestWorker_NoRatesToProcess(t)
	fmt.Println("no rates to process completed successfully!")
	TestWorker_OneToProcess(t)
	fmt.Println("one rate to process completed successfully!")
	TestWorker_TwoToProcessOneToOutcome(t)
	fmt.Println("two rates to process one to outcome completed successfully!")
	TestWorker_TwoToProcessOneToOutcomeOneNdc0SecondNdc1(t)
	fmt.Println("two rates to process one to outcome one ndc 0 second ndc 1 completed successfully!")
	TestWorker_TwoToProcessOneToOutcomeOneNdc0SecondNdc1(t)
	fmt.Println("two rates to process one to outcome one ndc < earned second earned < ndc < outcome completed successfully!")
	TestWorker_TwoToProcessOneToOutcomeOneNdcLessThanEarnedSecondNdcMoreThanEarnedButLessThanOutcome(t)
	fmt.Println("two rates to process one to outcome one ndc < earned second earned < ndc < earned + outcome completed successfully!")
	TestWorker_TwoToProcessOneToOutcomeOneNdcLessThanEarnedSecondNdcMoreThanEarnedPlusOutcome(t)
	fmt.Println("one rate to process one to outcome one earned + outcome< ndc < deposit completed successfully!")
}

func TestWorker_NoRatesToProcess(t *testing.T) {
	db := initDependencies()
	defer deleteDependencies()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	defer tx.Rollback()

	fur := models.UserRate{
		UserRateCreate: models.UserRateCreate{
			UserID:  userID,
			RateID:  firstRateID,
			Lock:    utils.DateOnly(time.Now().Add(time.Hour * 24 * 7)),
			Deposit: 111,
			Token:   "sol",
		},
		Opened:      utils.DateOnly(time.Now()),
		EarnedPool:  333,
		OutcomePool: 444,
	}

	sur := models.UserRate{
		UserRateCreate: models.UserRateCreate{
			UserID:  userID,
			RateID:  secondRateID,
			Lock:    utils.DateOnly(time.Now().Add(time.Hour * 24 * 14)),
			Deposit: 777,
			Token:   "usd",
		},
		Opened:      utils.DateOnly(time.Now()),
		EarnedPool:  888,
		OutcomePool: 999,
	}

	row := tx.QueryRow(`INSERT INTO users_rates (user_id, rate_id, lock, last_updated, opened, deposit,
                        earned_pool, next_day_charge, outcome_pool, token) VALUES ($1, $2, $3, $4, $5,
                                                                                   $6, $7, $8, $9, $10) RETURNING id`,
		fur.UserID, fur.RateID, fur.Lock, utils.DateOnly(time.Now()), utils.DateOnly(time.Now()), fur.Deposit,
		fur.EarnedPool, nil, fur.OutcomePool, fur.Token)

	err = row.Scan(&fur.ID)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			panic(err.Error() + rbErr.Error())
		}
		panic(err.Error())
	}

	row = tx.QueryRow(`INSERT INTO users_rates (user_id, rate_id, lock, last_updated, opened, deposit,
                        earned_pool, next_day_charge, outcome_pool, token) VALUES ($1, $2, $3, $4, $5,
                                                                                   $6, $7, $8, $9, $10) RETURNING id`,
		sur.UserID, sur.RateID, sur.Lock, utils.DateOnly(time.Now()), utils.DateOnly(time.Now()), sur.Deposit,
		sur.EarnedPool, 0, sur.OutcomePool, sur.Token)

	err = row.Scan(&sur.ID)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			panic(err.Error() + rbErr.Error())
		}
		panic(err.Error())
	}

	err = tx.Commit()
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			panic(err.Error() + rbErr.Error())
		}
		panic(err.Error())
	}

	time.Sleep(time.Second * 7)

	row = db.QueryRow(`SELECT id, user_id, rate_id, lock, last_updated, opened, deposit, earned_pool, next_day_charge,
      outcome_pool, token FROM users_rates WHERE id = $1`, fur.ID)

	var nfur models.UserRate
	var nfurLU time.Time
	var nfurNDC null.Float
	err = row.Scan(&nfur.ID, &nfur.UserID, &nfur.RateID, &nfur.Lock, &nfurLU, &nfur.Opened, &nfur.Deposit, &nfur.EarnedPool,
		&nfurNDC, &nfur.OutcomePool, &nfur.Token)
	if err != nil {
		panic(err.Error())
	}

	nfur.Opened = utils.DateOnly(nfur.Opened)
	nfur.Lock = utils.DateOnly(nfur.Lock)
	if fur != nfur || utils.DateOnly(nfurLU) != utils.DateOnly(time.Now()) || nfurNDC.Valid {
		panic(fmt.Sprintf("nfur != fur. fur: %v, nfur: %v, nfurLU: %v, nfurNDC: %v", fur, nfur, nfurLU, nfurNDC))
	}

	row = db.QueryRow(`SELECT id, user_id, rate_id, lock, last_updated, opened, deposit, earned_pool, next_day_charge,
      outcome_pool, token FROM users_rates WHERE id = $1`, sur.ID)

	var nsur models.UserRate
	var nsurLU time.Time
	var nsurNDC null.Float
	err = row.Scan(&nsur.ID, &nsur.UserID, &nsur.RateID, &nsur.Lock, &nsurLU, &nsur.Opened, &nsur.Deposit, &nsur.EarnedPool,
		&nsurNDC, &nsur.OutcomePool, &nsur.Token)
	if err != nil {
		panic(err.Error())
	}

	nsur.Opened = utils.DateOnly(nsur.Opened)
	nsur.Lock = utils.DateOnly(nsur.Lock)
	if sur != nsur || utils.DateOnly(nsurLU) != utils.DateOnly(time.Now()) || !nsurNDC.Valid {
		panic(fmt.Sprintf("nfur != fur. fur: %v, nfur: %v, nfurLU: %v, nfurNDC: %v", sur, nsurLU, nfurLU, nsurNDC))
	}

	_, err = db.Exec(`DELETE FROM users_rates WHERE id = $1`, fur.ID)
	if err != nil {
		panic(err.Error())
	}
	_, err = db.Exec(`DELETE FROM users_rates WHERE id = $1`, sur.ID)
	if err != nil {
		panic(err.Error())
	}

	return
}

// 2 тарифа 1 обработать но не сунуть ернед в ауткам
func TestWorker_OneToProcess(t *testing.T) {
	db := initDependencies()
	defer deleteDependencies()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	defer tx.Rollback()

	fur := models.UserRate{
		UserRateCreate: models.UserRateCreate{
			UserID:  userID,
			RateID:  firstRateID,
			Lock:    utils.DateOnly(time.Now().Add(time.Hour * 24 * 7)),
			Deposit: 111,
			Token:   "sol",
		},
		Opened:      utils.DateOnly(time.Now()),
		EarnedPool:  333,
		OutcomePool: 444,
	}

	sur := models.UserRate{
		UserRateCreate: models.UserRateCreate{
			UserID:  userID,
			RateID:  secondRateID,
			Lock:    utils.DateOnly(time.Now().Add(time.Hour * 24 * 14)),
			Deposit: 777,
			Token:   "usd",
		},
		Opened:      utils.DateOnly(time.Now()),
		EarnedPool:  888,
		OutcomePool: 999,
	}

	row := tx.QueryRow(`INSERT INTO users_rates (user_id, rate_id, lock, last_updated, opened, deposit,
                         earned_pool, next_day_charge, outcome_pool, token) VALUES ($1, $2, $3, $4, $5, 
                                                                                    $6, $7, $8, $9, $10) RETURNING id`,
		fur.UserID, fur.RateID, fur.Lock, utils.DateOnly(time.Now()), utils.DateOnly(time.Now()), fur.Deposit,
		fur.EarnedPool, nil, fur.OutcomePool, fur.Token)

	err = row.Scan(&fur.ID)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			panic(err.Error() + rbErr.Error())
		}
		panic(err.Error())
	}

	row = tx.QueryRow(`INSERT INTO users_rates (user_id, rate_id, lock, last_updated, opened, deposit,
                         earned_pool, next_day_charge, outcome_pool, token) VALUES ($1, $2, $3, $4, $5, 
                                                                                    $6, $7, $8, $9, $10) RETURNING id`,
		sur.UserID, sur.RateID, sur.Lock, utils.DateOnly(time.Now()).Add(-time.Hour*24), utils.DateOnly(time.Now()), sur.Deposit,
		sur.EarnedPool, nil, sur.OutcomePool, sur.Token)

	err = row.Scan(&sur.ID)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			panic(err.Error() + rbErr.Error())
		}
		panic(err.Error())
	}

	err = tx.Commit()
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			panic(err.Error() + rbErr.Error())
		}
		panic(err.Error())
	}

	time.Sleep(time.Second * 7)

	row = db.QueryRow(`SELECT id, user_id, rate_id, lock, last_updated, opened, deposit, earned_pool, next_day_charge, 
       outcome_pool, token FROM users_rates WHERE id = $1`, fur.ID)

	var nfur models.UserRate
	var nfurLU time.Time
	var nfurNDC null.Float
	err = row.Scan(&nfur.ID, &nfur.UserID, &nfur.RateID, &nfur.Lock, &nfurLU, &nfur.Opened, &nfur.Deposit, &nfur.EarnedPool,
		&nfurNDC, &nfur.OutcomePool, &nfur.Token)
	if err != nil {
		panic(err.Error())
	}

	nfur.Opened = utils.DateOnly(nfur.Opened)
	nfur.Lock = utils.DateOnly(nfur.Lock)
	if fur != nfur || utils.DateOnly(nfurLU) != utils.DateOnly(time.Now()) || nfurNDC.Valid {
		panic(fmt.Sprintf("nfur != fur. fur: %v, nfur: %v, nfurLU: %v, nfurNDC: %v", fur, nfur, nfurLU, nfurNDC))
	}

	row = db.QueryRow(`SELECT id, user_id, rate_id, lock, last_updated, opened, deposit, earned_pool, next_day_charge, 
       outcome_pool, token FROM users_rates WHERE id = $1`, sur.ID)

	var nsur models.UserRate
	var nsurLU time.Time
	var nsurNDC null.Float
	err = row.Scan(&nsur.ID, &nsur.UserID, &nsur.RateID, &nsur.Lock, &nsurLU, &nsur.Opened, &nsur.Deposit, &nsur.EarnedPool,
		&nsurNDC, &nsur.OutcomePool, &nsur.Token)
	if err != nil {
		panic(err.Error())
	}

	nsur.Opened = utils.DateOnly(nsur.Opened)
	nsur.Lock = utils.DateOnly(nsur.Lock)
	sur.EarnedPool = 888 + (888+999+777)*2
	if sur != nsur {
		panic(fmt.Sprintf("nsur != sur. sur: %v, nsur: %v", sur, nsur))
	}
	if utils.DateOnly(nsurLU) != utils.DateOnly(time.Now()) || nsurNDC.Valid {
		panic(fmt.Sprintf("nsurLU: %v, nsurNDC: %v", nsurLU, nsurNDC))
	}

	_, err = db.Exec(`DELETE FROM users_rates WHERE id = $1`, fur.ID)
	if err != nil {
		panic(err.Error())
	}
	_, err = db.Exec(`DELETE FROM users_rates WHERE id = $1`, sur.ID)
	if err != nil {
		panic(err.Error())
	}

	return
}

// 2 тарифа 2 обработать, первый сунуть, второй не сунуть
func TestWorker_TwoToProcessOneToOutcome(t *testing.T) {
	db := initDependencies()
	defer deleteDependencies()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	defer tx.Rollback()

	fur := models.UserRate{
		UserRateCreate: models.UserRateCreate{
			UserID:  userID,
			RateID:  firstRateID,
			Lock:    utils.DateOnly(time.Now().Add(-time.Hour * 24 * 7)),
			Deposit: 111,
			Token:   "sol",
		},
		Opened:      utils.DateOnly(time.Now()).Add(-time.Hour * 24 * 14),
		EarnedPool:  333,
		OutcomePool: 444,
	}

	sur := models.UserRate{
		UserRateCreate: models.UserRateCreate{
			UserID:  userID,
			RateID:  secondRateID,
			Lock:    utils.DateOnly(time.Now().Add(time.Hour * 24 * 14)),
			Deposit: 777,
			Token:   "usd",
		},
		Opened:      utils.DateOnly(time.Now()),
		EarnedPool:  888,
		OutcomePool: 999,
	}

	row := tx.QueryRow(`INSERT INTO users_rates (user_id, rate_id, lock, last_updated, opened, deposit,
                         earned_pool, next_day_charge, outcome_pool, token) VALUES ($1, $2, $3, $4, $5, 
                                                                                    $6, $7, $8, $9, $10) RETURNING id`,
		fur.UserID, fur.RateID, fur.Lock, utils.DateOnly(time.Now()).Add(-time.Hour*24), fur.Opened, fur.Deposit,
		fur.EarnedPool, nil, fur.OutcomePool, fur.Token)

	err = row.Scan(&fur.ID)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			panic(err.Error() + rbErr.Error())
		}
		panic(err.Error())
	}

	row = tx.QueryRow(`INSERT INTO users_rates (user_id, rate_id, lock, last_updated, opened, deposit,
                         earned_pool, next_day_charge, outcome_pool, token) VALUES ($1, $2, $3, $4, $5, 
                                                                                    $6, $7, $8, $9, $10) RETURNING id`,
		sur.UserID, sur.RateID, sur.Lock, utils.DateOnly(time.Now()).Add(-time.Hour*24), utils.DateOnly(time.Now()), sur.Deposit,
		sur.EarnedPool, nil, sur.OutcomePool, sur.Token)

	err = row.Scan(&sur.ID)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			panic(err.Error() + rbErr.Error())
		}
		panic(err.Error())
	}

	err = tx.Commit()
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			panic(err.Error() + rbErr.Error())
		}
		panic(err.Error())
	}

	time.Sleep(time.Second * 7)

	row = db.QueryRow(`SELECT id, user_id, rate_id, lock, last_updated, opened, deposit, earned_pool, next_day_charge, 
       outcome_pool, token FROM users_rates WHERE id = $1`, fur.ID)

	var nfur models.UserRate
	var nfurLU time.Time
	var nfurNDC null.Float
	err = row.Scan(&nfur.ID, &nfur.UserID, &nfur.RateID, &nfur.Lock, &nfurLU, &nfur.Opened, &nfur.Deposit, &nfur.EarnedPool,
		&nfurNDC, &nfur.OutcomePool, &nfur.Token)
	if err != nil {
		panic(err.Error())
	}

	nfur.Opened = utils.DateOnly(nfur.Opened)
	nfur.Lock = utils.DateOnly(nfur.Lock)
	fur.EarnedPool = (111 + 333 + 444) * 0.1
	fur.OutcomePool = 444 + 333
	if fur != nfur {
		if utils.FloatEquals(fur.EarnedPool, nfur.EarnedPool) {
			fur.EarnedPool = nfur.EarnedPool
			if fur != nfur {
				panic(fmt.Sprintf("nfur != fur. fur: %v, nfur: %v", fur, nfur))
			}
		} else {
			panic(fmt.Sprintf("nfur != fur. fur: %v, nfur: %v", fur, nfur))
		}
	}
	if utils.DateOnly(nfurLU) != utils.DateOnly(time.Now()) || nfurNDC.Valid {
		panic(fmt.Sprintf("nsurLU: %v, nsurNDC: %v", nfurLU, nfurNDC))
	}

	row = db.QueryRow(`SELECT id, user_id, rate_id, lock, last_updated, opened, deposit, earned_pool, next_day_charge, 
       outcome_pool, token FROM users_rates WHERE id = $1`, sur.ID)

	var nsur models.UserRate
	var nsurLU time.Time
	var nsurNDC null.Float
	err = row.Scan(&nsur.ID, &nsur.UserID, &nsur.RateID, &nsur.Lock, &nsurLU, &nsur.Opened, &nsur.Deposit, &nsur.EarnedPool,
		&nsurNDC, &nsur.OutcomePool, &nsur.Token)
	if err != nil {
		panic(err.Error())
	}

	nsur.Opened = utils.DateOnly(nsur.Opened)
	nsur.Lock = utils.DateOnly(nsur.Lock)
	sur.EarnedPool = 888 + (888+999+777)*2
	if sur != nsur {
		panic(fmt.Sprintf("nsur != sur. sur: %v, nsur: %v", sur, nsur))
	}
	if utils.DateOnly(nsurLU) != utils.DateOnly(time.Now()) || nsurNDC.Valid {
		panic(fmt.Sprintf("nsurLU: %v, nsurNDC: %v", nsurLU, nsurNDC))
	}

	_, err = db.Exec(`DELETE FROM users_rates WHERE id = $1`, fur.ID)
	if err != nil {
		panic(err.Error())
	}
	_, err = db.Exec(`DELETE FROM users_rates WHERE id = $1`, sur.ID)
	if err != nil {
		panic(err.Error())
	}

	return
}

// 2 тарифа 2 обработать, первый сунуть, второй не сунуть, перый ndc = 0, второй ndc = 1 (вместо 200%)
func TestWorker_TwoToProcessOneToOutcomeOneNdc0SecondNdc1(t *testing.T) {
	db := initDependencies()
	defer deleteDependencies()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	defer tx.Rollback()

	fur := models.UserRate{
		UserRateCreate: models.UserRateCreate{
			UserID:  userID,
			RateID:  firstRateID,
			Lock:    utils.DateOnly(time.Now().Add(-time.Hour * 24 * 7)),
			Deposit: 111,
			Token:   "sol",
		},
		Opened:      utils.DateOnly(time.Now()).Add(-time.Hour * 24 * 14),
		EarnedPool:  333,
		OutcomePool: 444,
	}

	sur := models.UserRate{
		UserRateCreate: models.UserRateCreate{
			UserID:  userID,
			RateID:  secondRateID,
			Lock:    utils.DateOnly(time.Now().Add(time.Hour * 24 * 14)),
			Deposit: 777,
			Token:   "usd",
		},
		Opened:      utils.DateOnly(time.Now()),
		EarnedPool:  888,
		OutcomePool: 999,
	}

	row := tx.QueryRow(`INSERT INTO users_rates (user_id, rate_id, lock, last_updated, opened, deposit,
                         earned_pool, next_day_charge, outcome_pool, token) VALUES ($1, $2, $3, $4, $5, 
                                                                                    $6, $7, $8, $9, $10) RETURNING id`,
		fur.UserID, fur.RateID, fur.Lock, utils.DateOnly(time.Now()).Add(-time.Hour*24), fur.Opened, fur.Deposit,
		fur.EarnedPool, 0, fur.OutcomePool, fur.Token)

	err = row.Scan(&fur.ID)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			panic(err.Error() + rbErr.Error())
		}
		panic(err.Error())
	}

	row = tx.QueryRow(`INSERT INTO users_rates (user_id, rate_id, lock, last_updated, opened, deposit,
                         earned_pool, next_day_charge, outcome_pool, token) VALUES ($1, $2, $3, $4, $5, 
                                                                                    $6, $7, $8, $9, $10) RETURNING id`,
		sur.UserID, sur.RateID, sur.Lock, utils.DateOnly(time.Now()).Add(-time.Hour*24), utils.DateOnly(time.Now()), sur.Deposit,
		sur.EarnedPool, 1, sur.OutcomePool, sur.Token)

	err = row.Scan(&sur.ID)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			panic(err.Error() + rbErr.Error())
		}
		panic(err.Error())
	}

	err = tx.Commit()
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			panic(err.Error() + rbErr.Error())
		}
		panic(err.Error())
	}

	time.Sleep(time.Second * 7)

	row = db.QueryRow(`SELECT id, user_id, rate_id, lock, last_updated, opened, deposit, earned_pool, next_day_charge, 
       outcome_pool, token FROM users_rates WHERE id = $1`, fur.ID)

	var nfur models.UserRate
	var nfurLU time.Time
	var nfurNDC null.Float
	err = row.Scan(&nfur.ID, &nfur.UserID, &nfur.RateID, &nfur.Lock, &nfurLU, &nfur.Opened, &nfur.Deposit, &nfur.EarnedPool,
		&nfurNDC, &nfur.OutcomePool, &nfur.Token)
	if err != nil {
		panic(err.Error())
	}

	nfur.Opened = utils.DateOnly(nfur.Opened)
	nfur.Lock = utils.DateOnly(nfur.Lock)
	fur.EarnedPool = 0
	fur.OutcomePool = 444 + 333
	if fur != nfur {
		if utils.FloatEquals(fur.EarnedPool, nfur.EarnedPool) {
			fur.EarnedPool = nfur.EarnedPool
			if fur != nfur {
				panic(fmt.Sprintf("nfur != fur. fur: %v, nfur: %v", fur, nfur))
			}
		} else {
			panic(fmt.Sprintf("nfur != fur. fur: %v, nfur: %v", fur, nfur))
		}
	}
	if utils.DateOnly(nfurLU) != utils.DateOnly(time.Now()) || nfurNDC.Valid {
		panic(fmt.Sprintf("nsurLU: %v, nsurNDC: %v", nfurLU, nfurNDC))
	}

	row = db.QueryRow(`SELECT id, user_id, rate_id, lock, last_updated, opened, deposit, earned_pool, next_day_charge, 
       outcome_pool, token FROM users_rates WHERE id = $1`, sur.ID)

	var nsur models.UserRate
	var nsurLU time.Time
	var nsurNDC null.Float
	err = row.Scan(&nsur.ID, &nsur.UserID, &nsur.RateID, &nsur.Lock, &nsurLU, &nsur.Opened, &nsur.Deposit, &nsur.EarnedPool,
		&nsurNDC, &nsur.OutcomePool, &nsur.Token)
	if err != nil {
		panic(err.Error())
	}

	nsur.Opened = utils.DateOnly(nsur.Opened)
	nsur.Lock = utils.DateOnly(nsur.Lock)
	sur.EarnedPool = 888 + 1
	if sur != nsur {
		if utils.FloatEquals(sur.EarnedPool, nsur.EarnedPool) {
			sur.EarnedPool = nsur.EarnedPool
			if sur != nsur {
				panic(fmt.Sprintf("nsur != sur. sur: %v, nsur: %v", fur, nsur))
			}
			panic(fmt.Sprintf("nsur != sur. sur: %v, nsur: %v", sur, nsur))
		}
	}
	if utils.DateOnly(nsurLU) != utils.DateOnly(time.Now()) || nsurNDC.Valid {
		panic(fmt.Sprintf("nsurLU: %v, nsurNDC: %v", nsurLU, nsurNDC))
	}

	_, err = db.Exec(`DELETE FROM users_rates WHERE id = $1`, fur.ID)
	if err != nil {
		panic(err.Error())
	}
	_, err = db.Exec(`DELETE FROM users_rates WHERE id = $1`, sur.ID)
	if err != nil {
		panic(err.Error())
	}

	return

}

// 2 тарифа 2 обработать, первый сунуть, второй не сунуть, перый ndc < earned, второй earned < ndc < outcome
func TestWorker_TwoToProcessOneToOutcomeOneNdcLessThanEarnedSecondNdcMoreThanEarnedButLessThanOutcome(t *testing.T) {
	db := initDependencies()
	defer deleteDependencies()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	defer tx.Rollback()

	fur := models.UserRate{
		UserRateCreate: models.UserRateCreate{
			UserID:  userID,
			RateID:  firstRateID,
			Lock:    utils.DateOnly(time.Now().Add(-time.Hour * 24 * 7)),
			Deposit: 111,
			Token:   "sol",
		},
		Opened:      utils.DateOnly(time.Now()).Add(-time.Hour * 24 * 14),
		EarnedPool:  333,
		OutcomePool: 444,
	}

	sur := models.UserRate{
		UserRateCreate: models.UserRateCreate{
			UserID:  userID,
			RateID:  secondRateID,
			Lock:    utils.DateOnly(time.Now().Add(time.Hour * 24 * 14)),
			Deposit: 777,
			Token:   "usd",
		},
		Opened:      utils.DateOnly(time.Now()),
		EarnedPool:  888,
		OutcomePool: 999,
	}

	row := tx.QueryRow(`INSERT INTO users_rates (user_id, rate_id, lock, last_updated, opened, deposit,
                         earned_pool, next_day_charge, outcome_pool, token) VALUES ($1, $2, $3, $4, $5, 
                                                                                    $6, $7, $8, $9, $10) RETURNING id`,
		fur.UserID, fur.RateID, fur.Lock, utils.DateOnly(time.Now()).Add(-time.Hour*24), fur.Opened, fur.Deposit,
		fur.EarnedPool, -332, fur.OutcomePool, fur.Token)

	err = row.Scan(&fur.ID)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			panic(err.Error() + rbErr.Error())
		}
		panic(err.Error())
	}

	row = tx.QueryRow(`INSERT INTO users_rates (user_id, rate_id, lock, last_updated, opened, deposit,
                         earned_pool, next_day_charge, outcome_pool, token) VALUES ($1, $2, $3, $4, $5, 
                                                                                    $6, $7, $8, $9, $10) RETURNING id`,
		sur.UserID, sur.RateID, sur.Lock, utils.DateOnly(time.Now()).Add(-time.Hour*24), utils.DateOnly(time.Now()), sur.Deposit,
		sur.EarnedPool, -889, sur.OutcomePool, sur.Token)

	err = row.Scan(&sur.ID)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			panic(err.Error() + rbErr.Error())
		}
		panic(err.Error())
	}

	err = tx.Commit()
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			panic(err.Error() + rbErr.Error())
		}
		panic(err.Error())
	}

	time.Sleep(time.Second * 7)

	row = db.QueryRow(`SELECT id, user_id, rate_id, lock, last_updated, opened, deposit, earned_pool, next_day_charge, 
       outcome_pool, token FROM users_rates WHERE id = $1`, fur.ID)

	var nfur models.UserRate
	var nfurLU time.Time
	var nfurNDC null.Float
	err = row.Scan(&nfur.ID, &nfur.UserID, &nfur.RateID, &nfur.Lock, &nfurLU, &nfur.Opened, &nfur.Deposit, &nfur.EarnedPool,
		&nfurNDC, &nfur.OutcomePool, &nfur.Token)
	if err != nil {
		panic(err.Error())
	}

	nfur.Opened = utils.DateOnly(nfur.Opened)
	nfur.Lock = utils.DateOnly(nfur.Lock)
	fur.EarnedPool = 0
	fur.OutcomePool = 444 + 1
	if fur != nfur {
		if utils.FloatEquals(fur.EarnedPool, nfur.EarnedPool) {
			fur.EarnedPool = nfur.EarnedPool
			if fur != nfur {
				panic(fmt.Sprintf("nfur != fur. fur: %v, nfur: %v", fur, nfur))
			}
		} else {
			panic(fmt.Sprintf("nfur != fur. fur: %v, nfur: %v", fur, nfur))
		}
	}
	if utils.DateOnly(nfurLU) != utils.DateOnly(time.Now()) || nfurNDC.Valid {
		panic(fmt.Sprintf("nsurLU: %v, nsurNDC: %v", nfurLU, nfurNDC))
	}

	row = db.QueryRow(`SELECT id, user_id, rate_id, lock, last_updated, opened, deposit, earned_pool, next_day_charge, 
       outcome_pool, token FROM users_rates WHERE id = $1`, sur.ID)

	var nsur models.UserRate
	var nsurLU time.Time
	var nsurNDC null.Float
	err = row.Scan(&nsur.ID, &nsur.UserID, &nsur.RateID, &nsur.Lock, &nsurLU, &nsur.Opened, &nsur.Deposit, &nsur.EarnedPool,
		&nsurNDC, &nsur.OutcomePool, &nsur.Token)
	if err != nil {
		panic(err.Error())
	}

	nsur.Opened = utils.DateOnly(nsur.Opened)
	nsur.Lock = utils.DateOnly(nsur.Lock)
	sur.EarnedPool = 0
	sur.OutcomePool -= 1
	if sur != nsur {
		if utils.FloatEquals(sur.EarnedPool, nsur.EarnedPool) {
			sur.EarnedPool = nsur.EarnedPool
			if sur != nsur {
				panic(fmt.Sprintf("nsur != sur. sur: %v, nsur: %v", fur, nsur))
			}
			panic(fmt.Sprintf("nsur != sur. sur: %v, nsur: %v", sur, nsur))
		}
	}
	if utils.DateOnly(nsurLU) != utils.DateOnly(time.Now()) || nsurNDC.Valid {
		panic(fmt.Sprintf("nsurLU: %v, nsurNDC: %v", nsurLU, nsurNDC))
	}

	_, err = db.Exec(`DELETE FROM users_rates WHERE id = $1`, fur.ID)
	if err != nil {
		panic(err.Error())
	}
	_, err = db.Exec(`DELETE FROM users_rates WHERE id = $1`, sur.ID)
	if err != nil {
		panic(err.Error())
	}

	return

}

// 1 тариф 1 обработать, сунуть, перый earned + outcome < ndc
func TestWorker_TwoToProcessOneToOutcomeOneNdcLessThanEarnedSecondNdcMoreThanEarnedPlusOutcome(t *testing.T) {
	db := initDependencies()
	defer deleteDependencies()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	defer tx.Rollback()

	fur := models.UserRate{
		UserRateCreate: models.UserRateCreate{
			UserID:  userID,
			RateID:  firstRateID,
			Lock:    utils.DateOnly(time.Now().Add(-time.Hour * 24 * 7)),
			Deposit: 111,
			Token:   "sol",
		},
		Opened:      utils.DateOnly(time.Now()).Add(-time.Hour * 24 * 14),
		EarnedPool:  333,
		OutcomePool: 444,
	}

	row := tx.QueryRow(`INSERT INTO users_rates (user_id, rate_id, lock, last_updated, opened, deposit,
                         earned_pool, next_day_charge, outcome_pool, token) VALUES ($1, $2, $3, $4, $5, 
                                                                                    $6, $7, $8, $9, $10) RETURNING id`,
		fur.UserID, fur.RateID, fur.Lock, utils.DateOnly(time.Now()).Add(-time.Hour*24), fur.Opened, fur.Deposit,
		fur.EarnedPool, -778, fur.OutcomePool, fur.Token)

	err = row.Scan(&fur.ID)
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			panic(err.Error() + rbErr.Error())
		}
		panic(err.Error())
	}

	err = tx.Commit()
	if err != nil {
		rbErr := tx.Rollback()
		if rbErr != nil {
			panic(err.Error() + rbErr.Error())
		}
		panic(err.Error())
	}

	time.Sleep(time.Second * 7)

	row = db.QueryRow(`SELECT id, user_id, rate_id, lock, last_updated, opened, deposit, earned_pool, next_day_charge, 
       outcome_pool, token FROM users_rates WHERE id = $1`, fur.ID)

	var nfur models.UserRate
	var nfurLU time.Time
	var nfurNDC null.Float
	err = row.Scan(&nfur.ID, &nfur.UserID, &nfur.RateID, &nfur.Lock, &nfurLU, &nfur.Opened, &nfur.Deposit, &nfur.EarnedPool,
		&nfurNDC, &nfur.OutcomePool, &nfur.Token)
	if err != nil {
		panic(err.Error())
	}

	nfur.Opened = utils.DateOnly(nfur.Opened)
	nfur.Lock = utils.DateOnly(nfur.Lock)
	fur.EarnedPool = 0
	fur.OutcomePool = 0
	fur.Deposit = 110
	if fur != nfur {
		if utils.FloatEquals(fur.EarnedPool, nfur.EarnedPool) {
			fur.EarnedPool = nfur.EarnedPool
			if fur != nfur {
				panic(fmt.Sprintf("nfur != fur. fur: %v, nfur: %v", fur, nfur))
			}
		} else {
			panic(fmt.Sprintf("nfur != fur. fur: %v, nfur: %v", fur, nfur))
		}
	}
	if utils.DateOnly(nfurLU) != utils.DateOnly(time.Now()) || nfurNDC.Valid {
		panic(fmt.Sprintf("nsurLU: %v, nsurNDC: %v", nfurLU, nfurNDC))
	}

	_, err = db.Exec(`DELETE FROM users_rates WHERE id = $1`, fur.ID)
	if err != nil {
		panic(err.Error())
	}

	return

}
