package provision

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/lithammer/shortuuid/v4"
	"github.com/piotrekmonko/portfello/pkg/auth"
	"github.com/piotrekmonko/portfello/pkg/dao"
	"github.com/piotrekmonko/portfello/pkg/logz"
	"github.com/spf13/cobra"
	"time"
)

// Provisioner handles user and test data provisioning.
type Provisioner struct {
	log  logz.Logger
	db   *dao.DAO
	auth *auth.Service
}

func NewProvisioner(log *logz.Log, db *dao.DAO, auth *auth.Service) *Provisioner {
	return &Provisioner{
		log:  log.Named("provisioner"),
		db:   db,
		auth: auth,
	}
}

func (p *Provisioner) HandleCommand(cmd *cobra.Command) error {
	ctx := cmd.Context()

	if email, _ := cmd.Flags().GetString("user"); email != "" {
		if err := p.User(ctx, p.auth, email); err != nil {
			return err
		}
	}

	if testData, _ := cmd.Flags().GetBool("test"); testData {
		if err := p.TestData(ctx, p.db); err != nil {
			return err
		}
	}

	return nil
}

func (p *Provisioner) User(ctx context.Context, authService *auth.Service, email string) error {
	user, err := authService.CreateUser(ctx, email, email, auth.Roles{auth.RoleUser})
	if err != nil {
		return err
	}

	fmt.Printf("User created: %+v\n", user)
	return nil
}

func (p *Provisioner) TestData(ctx context.Context, dbq *dao.DAO) error {
	q, rollbacker, err := dbq.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer rollbacker()

	expenses := []float64{100.50, -20.12, -34.14, -3.1415, 60.0, -10.50, -23.99}
	testUserID := gofakeit.Email()

	testWallet := &dao.WalletInsertParams{
		ID:        shortuuid.New(),
		UserID:    testUserID,
		Balance:   0.0,
		Currency:  "USD",
		CreatedAt: time.Now().Add(time.Hour * time.Duration(-1*len(expenses))).UTC(),
	}
	if err := q.WalletInsert(ctx, testWallet); err != nil {
		return err
	}

	for i, amount := range expenses {
		testExpense := &dao.ExpenseInsertParams{
			ID:          shortuuid.New(),
			WalletID:    testWallet.ID,
			Amount:      amount,
			Description: dao.NilStr(gofakeit.HackerPhrase()),
			CreatedAt:   time.Now().Add(time.Hour * time.Duration(-1*len(expenses)-i)).UTC(),
		}
		if err := q.ExpenseInsert(ctx, testExpense); err != nil {
			return err
		}
	}

	return q.Commit(ctx)
}
