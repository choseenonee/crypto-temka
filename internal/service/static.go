package service

import (
	"context"
	"crypto-temka/internal/models"
	"crypto-temka/internal/repository"
	"crypto-temka/pkg/config"
	"crypto-temka/pkg/log"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"math/rand"
	"os"
	"sync"
	"time"
)

type static struct {
	sync.RWMutex
	metricsSet     models.MetricsSet
	metricsSetFile *os.File

	outcomes     []models.Outcome
	outcomeMutex sync.RWMutex

	logger *log.Logs
	repo   repository.Static
}

var once sync.Once
var st *static

func InitStatic(repo repository.Static, logger *log.Logs, metricsSetFile *os.File) Static {
	once.Do(func() {
		metricsSetJSON, err := io.ReadAll(metricsSetFile)

		var metricsSet models.MetricsSet
		err = json.Unmarshal(metricsSetJSON, &metricsSet)
		if err != nil {
			panic(fmt.Sprintf("Error unmarshalling json from static/metrics.json file, err: %v", err.Error()))
		}

		st = &static{
			logger:         logger,
			repo:           repo,
			metricsSet:     metricsSet,
			metricsSetFile: metricsSetFile,
		}

		go st.serveOutcomesGeneration()
	})

	return st
}

var tokens = []string{"BitCoin", "Ethereum", "LiteCoin", "DogeCoin", "Dash", "BitcoinCash", "Zcash", "Ripple", "TRON",
	"Stellar", "BinanceCoin", "TRON_TRC20", "Ethereum_ERC20"}

func generateOutcome() models.Outcome {
	var otc models.Outcome
	otc.Amount = float64(rand.Intn(viper.GetInt(config.OutcomeAmountMax)+1)) + rand.Float64()
	otc.Token = tokens[rand.Intn(len(tokens))]
	otc.UserID = viper.GetInt(config.OutcomeUserIDMin) + rand.Intn(viper.GetInt(config.OutcomeUserIDMax)-viper.GetInt(config.OutcomeUserIDMin))
	return otc
}

func waitOutcome() {
	t := time.Second * time.Duration(viper.GetInt(config.OutcomeTickerMin)+
		rand.Intn(viper.GetInt(config.OutcomeTickerMax)-viper.GetInt(config.OutcomeTickerMin)+1))
	time.Sleep(t)
}

func (s *static) serveOutcomesGeneration() {
	if len(s.outcomes) == 0 {
		s.outcomeMutex.Lock()
		for i := 0; i < viper.GetInt(config.OutcomesAmount); i++ {
			s.outcomes = append(s.outcomes, generateOutcome())
		}
		s.outcomeMutex.Unlock()
	}
	for {
		waitOutcome()
		s.outcomeMutex.Lock()
		newOutcome := generateOutcome()
		outcomes := []models.Outcome{newOutcome}
		s.outcomes = append(outcomes, s.outcomes[:len(s.outcomes)-1]...)
		s.outcomeMutex.Unlock()
	}
}

func (s *static) CreateReview(ctx context.Context, rc models.ReviewCreate) (int, error) {
	id, err := s.repo.CreateReview(ctx, rc)
	if err != nil {
		s.logger.Error(err.Error())
		return 0, err
	}

	return id, nil
}

func (s *static) GetReviews(ctx context.Context, page, perPage int) ([]models.Review, error) {
	reviews, err := s.repo.GetReviews(ctx, page, perPage)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	return reviews, nil
}

func (s *static) UpdateReview(ctx context.Context, r models.Review) error {
	err := s.repo.UpdateReview(ctx, r)
	if err != nil {
		s.logger.Error(err.Error())
		return err
	}

	return nil
}

func (s *static) DeleteReview(ctx context.Context, id int) error {
	err := s.repo.DeleteReview(ctx, id)
	if err != nil {
		s.logger.Error(err.Error())
		return err
	}

	return nil
}

func (s *static) SetMetrics(m models.MetricsSet) error {
	mJSON, _ := json.Marshal(m)

	err := s.metricsSetFile.Truncate(0)
	if err != nil {
		return err
	}

	_, err = s.metricsSetFile.WriteAt(mJSON, 0)
	if err != nil {
		return err
	}

	s.Lock()
	defer s.Unlock()

	s.metricsSet = m

	return nil
}

func (s *static) GetMetrics() models.Metrics {
	s.RLock()
	defer s.RUnlock()

	return models.Metrics{
		MetricsSet:   s.metricsSet,
		IncomeSubOut: s.metricsSet.AlltimeIncome - s.metricsSet.AlltimeOut,
	}
}

func (s *static) CreateCase(ctx context.Context, cc models.CaseCreate) (int, error) {
	id, err := s.repo.CreateCase(ctx, cc)
	if err != nil {
		s.logger.Error(err.Error())
		return 0, err
	}

	return id, nil
}

func (s *static) UpdateCase(ctx context.Context, cu models.Case) error {
	err := s.repo.UpdateCase(ctx, cu)
	if err != nil {
		s.logger.Error(err.Error())
		return err
	}

	return nil
}

func (s *static) DeleteCase(ctx context.Context, id int) error {
	err := s.repo.DeleteCase(ctx, id)
	if err != nil {
		s.logger.Error(err.Error())
		return err
	}

	return nil
}

func (s *static) GetCase(ctx context.Context, id int) (models.Case, error) {
	c, err := s.repo.GetCase(ctx, id)
	if err != nil {
		s.logger.Error(err.Error())
		return models.Case{}, err
	}

	return c, nil
}

func (s *static) GetCases(ctx context.Context, page, perPage int) ([]models.Case, error) {
	cases, err := s.repo.GetCases(ctx, page, perPage)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	return cases, nil
}

func (s *static) GetOutcome() []models.Outcome {
	s.outcomeMutex.RLock()
	defer s.outcomeMutex.RUnlock()
	return s.outcomes
}
