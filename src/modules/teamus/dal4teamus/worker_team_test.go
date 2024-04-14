package dal4teamus

import (
	"context"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/mocks4dalgo/mocks4dal"
	"github.com/golang/mock/gomock"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/core4teamus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/dto4teamus"
	"github.com/sneat-co/sneat-go-backend/src/modules/teamus/models4teamus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type fooModuleTeamData struct {
	Int1 int
	Str1 string
}

func (fooModuleTeamData) Validate() error {
	return nil
}

func TestRunModuleTeamWorker(t *testing.T) {
	ctx := context.Background()
	user := &facade.AuthUser{ID: "user1"}
	request := dto4teamus.TeamRequest{TeamID: "team1"}
	const moduleID = "test_module"
	assertTxWorker := func(ctx context.Context, tx dal.ReadwriteTransaction, params *ModuleTeamWorkerParams[*fooModuleTeamData]) (err error) {
		if err := params.GetRecords(ctx, tx, params.UserID); err != nil {
			return err
		}
		assert.NotNil(t, params)
		assert.NotNil(t, params.TeamModuleEntry)
		assert.NotNil(t, params.TeamModuleEntry.Record)
		assert.NotNil(t, params.TeamModuleEntry.Data)
		assert.NotNil(t, params.TeamModuleEntry.Record.Data())
		return nil
	}
	facade.GetDatabase = func(ctx context.Context) dal.DB {
		ctrl := gomock.NewController(t)
		db := mocks4dal.NewMockDatabase(ctrl)
		//var db2 dal.DB
		//db2.RunReadwriteTransaction()
		db.EXPECT().RunReadwriteTransaction(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, worker dal.RWTxWorker, options ...dal.TransactionOption) error {
			tx := mocks4dal.NewMockReadwriteTransaction(ctrl)
			tx.EXPECT().GetMulti(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, records []dal.Record) error {
				for _, record := range records {
					record.SetError(nil)
					if record.Key().Collection() == "teams" {
						teamData := record.Data().(*models4teamus.TeamDbo)
						teamData.CreatedAt = time.Now()
						teamData.CreatedBy = "test"
						teamData.IncreaseVersion(teamData.CreatedAt, teamData.CreatedBy)
						teamData.Type = core4teamus.TeamTypeFamily
						teamData.CountryID = "UK"
						teamData.Status = dbmodels.StatusActive
						teamData.UserIDs = []string{user.ID}
					}
				}
				return nil
			})
			return worker(ctx, tx)
		})
		return db
	}
	err := RunModuleTeamWorker(ctx, user, request, moduleID, new(fooModuleTeamData), assertTxWorker)
	assert.Nil(t, err)
	//type args[D TeamModuleData] struct {
	//	ctx      context.Context
	//	user     facade.User
	//	request  dto4teamus.TeamRequest
	//	moduleID string
	//	worker   func(ctx context.Context, tx dal.ReadwriteTransaction, teamWorkerParams *ModuleTeamWorkerParams[D]) (err error)
	//}
	//type testCase[D TeamModuleData] struct {
	//	name    string
	//	args    args[D]
	//	wantErr bool
	//}
	//tests := []testCase[ /* TODO: Insert concrete types here */ ]{
	//	// TODO: Add test cases.
	//}
	//for _, tt := range tests {
	//	t.Run(tt.name, func(t *testing.T) {
	//		if err := RunModuleTeamWorker(tt.args.ctx, tt.args.user, tt.args.request, tt.args.moduleID, tt.args.worker); (err != nil) != tt.wantErr {
	//			t.Errorf("RunModuleTeamWorker() error = %v, wantErr %v", err, tt.wantErr)
	//		}
	//	})
	//}
}

func TestRunModuleTeamWorkerTx(t *testing.T) {
	ctx := context.Background()
	user := &facade.AuthUser{ID: "user1"}
	request := dto4teamus.TeamRequest{TeamID: "team1"}
	const moduleID = "test_module"
	//assertTxWorker := func(ctx context.Context, tx dal.ReadwriteTransaction, teamWorkerParams *ModuleTeamWorkerParams[*fooModuleTeamData]) (err error) {
	//	assert.NotNil(t, teamWorkerParams)
	//	assert.NotNil(t, teamWorkerParams.TeamModuleEntry)
	//	assert.NotNil(t, teamWorkerParams.TeamModuleEntry.Record)
	//	assert.NotNil(t, teamWorkerParams.TeamModuleEntry.Data)
	//	assert.NotNil(t, teamWorkerParams.TeamModuleEntry.Record.Data())
	//	return nil
	//}
	assert.Panics(t, func() {
		_ = RunModuleTeamWorkerTx(ctx, nil, user, request, moduleID, new(fooModuleTeamData), nil)
	})
}
