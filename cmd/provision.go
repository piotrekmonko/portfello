/*
Copyright © 2023 Piotr Mońko <piotrek.monko@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/lithammer/shortuuid/v4"
	"github.com/piotrekmonko/portfello/pkg/auth"
	"github.com/piotrekmonko/portfello/pkg/conf"
	"github.com/piotrekmonko/portfello/pkg/dao"
	"github.com/piotrekmonko/portfello/pkg/logz"
	"github.com/spf13/cobra"
	"time"
)

// provisionCmd represents the provision command
var provisionCmd = &cobra.Command{
	Use:     "provision",
	Aliases: []string{"prov"},
	Short:   "Add objects to systems",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		conf := conf.New()
		log := logz.NewLogger(&conf.Logging)
		db, dbq, err := dao.NewDAO(ctx, log, conf.DatabaseDSN)
		if err != nil {
			return
		}
		defer db.Close()
		authProvider, err := auth.NewProvider(ctx, log, conf, dbq)
		if err != nil {
			return
		}
		authService := auth.New(authProvider)

		if email, _ := cmd.Flags().GetString("user"); email != "" {
			if err := provisionUser(cmd.Context(), authService, email); err != nil {
				cobra.CheckErr(err)
			}
		}

		if testData, _ := cmd.Flags().GetBool("test"); testData {
			if err := provisionTestData(cmd.Context(), dbq); err != nil {
				cobra.CheckErr(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(provisionCmd)

	provisionCmd.Flags().StringP("user", "u", "", "Add a new user")
	provisionCmd.Flags().BoolP("test", "t", false, "Add example test data to database")
}

func provisionUser(ctx context.Context, authService *auth.Service, email string) error {
	user, err := authService.CreateUser(ctx, email, email, auth.Roles{auth.RoleUser})
	if err != nil {
		return err
	}

	fmt.Printf("User created: %+v\n", user)
	return nil
}

func provisionTestData(ctx context.Context, dbq *dao.DAO) error {
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
