package facade4userus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTxGetUsers(t *testing.T) {
	type args struct {
		ctx   context.Context
		tx    dal.ReadwriteTransaction
		users []dal.Record
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "nil tx",
			args: args{
				ctx:   context.Background(),
				tx:    nil,
				users: nil,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Nil(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, TxGetUsers(tt.args.ctx, tt.args.tx, tt.args.users), fmt.Sprintf("TxGetUsers(%v, %v, %v)", tt.args.ctx, tt.args.tx, tt.args.users))
		})
	}
}
