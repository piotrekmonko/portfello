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

	numTestData, _ := cmd.Flags().GetInt("num")

	if email, _ := cmd.Flags().GetString("user"); email != "" {
		if err := p.User(ctx, p.auth, email); err != nil {
			return err
		}
	}

	if testData, _ := cmd.Flags().GetBool("test"); testData {
		if err := p.TestData(ctx, p.db, numTestData); err != nil {
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

func (p *Provisioner) TestData(ctx context.Context, dbq *dao.DAO, numTestData int) error {
	q, rollbacker, err := dbq.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer rollbacker()

	for i := 0; i < numTestData; i++ {
		err = p.insertTestExpenses(ctx, q)
		if err != nil {
			return err
		}
	}

	return q.Commit(ctx)
}

func (p *Provisioner) insertTestExpenses(ctx context.Context, q dao.DBInterface) error {
	numExpenses := gofakeit.IntRange(0, 500)
	expensesEarliestDate := gofakeit.DateRange(time.Now().Add(time.Hour*24*time.Duration(numExpenses)), time.Now().UTC())
	testUserID := gofakeit.Email()

	testWallet := &dao.WalletInsertParams{
		ID:        shortuuid.New(),
		UserID:    testUserID,
		Balance:   0.0,
		Currency:  "USD",
		CreatedAt: expensesEarliestDate.UTC(),
	}
	if err := q.WalletInsert(ctx, testWallet); err != nil {
		return p.log.Errorw(ctx, err, "cannot insert a wallet", "args", testWallet)
	}

	for i := 0; i < numExpenses; i++ {
		testExpense := &dao.ExpenseInsertParams{
			ID:          shortuuid.New(),
			WalletID:    testWallet.ID,
			Amount:      gofakeit.Float64Range(-100, 100),
			Description: dao.NilStr(gofakeit.HackerPhrase()),
			CreatedAt:   gofakeit.DateRange(expensesEarliestDate.Add(time.Hour*24*time.Duration(i)), time.Now().UTC()),
		}
		if err := q.ExpenseInsert(ctx, testExpense); err != nil {
			return p.log.Errorw(ctx, err, "cannot insert an expense", "args", testExpense)
		}
	}

	p.log.Infof(ctx, "created user %s's wallet with %d expenses", testUserID, numExpenses)
	return nil
}
