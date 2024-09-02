package dal4spaceus

import (
	"context"
	"fmt"
	"github.com/dal-go/dalgo/dal"
	"github.com/dal-go/mocks4dalgo/mocks4dal"
	"github.com/golang/mock/gomock"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/core4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dbo4spaceus"
	"github.com/sneat-co/sneat-go-backend/src/modules/spaceus/dto4spaceus"
	"github.com/sneat-co/sneat-go-core/facade"
	"github.com/sneat-co/sneat-go-core/models/dbmodels"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type fooModuleSpaceData struct {
	Int1 int
	Str1 string
}

func (fooModuleSpaceData) Validate() error {
	return nil
}

func TestRunModuleSpaceWorker(t *testing.T) {
	ctx := context.Background()
	user := &facade.AuthUserContext{ID: "user1"}
	request := dto4spaceus.SpaceRequest{SpaceID: "space1"}
	const moduleID = "test_module"
	assertTxWorker := func(ctx context.Context, tx dal.ReadwriteTransaction, params *ModuleSpaceWorkerParams[*fooModuleSpaceData]) (err error) {
		if err = params.GetRecords(ctx, tx); err != nil {
			return err
		}
		assert.NotNil(t, params)
		assert.NotNil(t, params.SpaceModuleEntry)
		assert.NotNil(t, params.SpaceModuleEntry.Record)
		assert.NotNil(t, params.SpaceModuleEntry.Data)
		assert.NotNil(t, params.SpaceModuleEntry.Record.Data())
		return nil
	}
	facade.GetDatabase = func(ctx context.Context) (dal.DB, error) {
		ctrl := gomock.NewController(t)
		db := mocks4dal.NewMockDatabase(ctrl)
		//var db2 dal.DB
		//db2.RunReadwriteTransaction()
		db.EXPECT().RunReadwriteTransaction(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, worker dal.RWTxWorker, options ...dal.TransactionOption) error {
			tx := mocks4dal.NewMockReadwriteTransaction(ctrl)
			tx.EXPECT().Get(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, record dal.Record) error {
				switch key := record.Key(); key.Collection() {
				case "spaces":
					record.SetError(nil)
					spaceDbo := record.Data().(*dbo4spaceus.SpaceDbo)
					spaceDbo.CreatedAt = time.Now()
					spaceDbo.CreatedBy = "test"
					spaceDbo.IncreaseVersion(spaceDbo.CreatedAt, spaceDbo.CreatedBy)
					spaceDbo.Type = core4spaceus.SpaceTypeFamily
					spaceDbo.CountryID = "UK"
					spaceDbo.Status = dbmodels.StatusActive
					spaceDbo.UserIDs = []string{user.ID}
				case "modules":
					record.SetError(nil)
				default:
					err := fmt.Errorf("unexpected Get(key=%+v)", key)
					record.SetError(err)
					return err
				}
				return nil
			})
			tx.EXPECT().GetMulti(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, records []dal.Record) error {
				for _, record := range records {
					switch key := record.Key(); key.Collection() {
					case "spaces":
						if err := record.Error(); err == nil {
							t.Error(fmt.Errorf("duplicate call to get space record"))
						}
					case "modules":
						record.SetError(nil)
					default:
						err := fmt.Errorf("unexpected GetMulti(key=%+v)", key)
						record.SetError(err)
					}
				}
				return nil
			})
			return worker(ctx, tx)
		})
		return db, nil
	}
	err := RunModuleSpaceWorker(ctx, user, request.SpaceID, moduleID, new(fooModuleSpaceData), assertTxWorker)
	assert.Nil(t, err)
	//type args[ModuleDbo SpaceModuleDbo] struct {
	//	ctx      context.Context
	//	user     facade4debtus.User
	//	request  dto4spaceus.SpaceRequest
	//	moduleID string
	//	worker   func(ctx context.Context, tx dal.ReadwriteTransaction, teamWorkerParams *ModuleSpaceWorkerParams[ModuleDbo]) (err error)
	//}
	//type testCase[ModuleDbo SpaceModuleDbo] struct {
	//	name    string
	//	args    args[ModuleDbo]
	//	wantErr bool
	//}
	//tests := []testCase[ /* TODO: Insert concrete types here */ ]{
	//	// TODO: Add test cases.
	//}
	//for _, tt := range tests {
	//	t.Run(tt.name, func(t *testing.T) {
	//		if err := RunModuleSpaceWorker(tt.args.ctx, tt.args.user, tt.args.request, tt.args.moduleID, tt.args.worker); (err != nil) != tt.wantErr {
	//			t.Errorf("RunModuleSpaceWorker() error = %v, wantErr %v", err, tt.wantErr)
	//		}
	//	})
	//}
}

func TestRunModuleSpaceWorkerTx(t *testing.T) {
	ctx := context.Background()
	user := &facade.AuthUserContext{ID: "user1"}
	request := dto4spaceus.SpaceRequest{SpaceID: "space1"}
	const moduleID = "test_module"
	//assertTxWorker := func(ctx context.Context, tx dal.ReadwriteTransaction, teamWorkerParams *ModuleSpaceWorkerParams[*fooModuleSpaceData]) (err error) {
	//	assert.NotNil(t, teamWorkerParams)
	//	assert.NotNil(t, teamWorkerParams.SpaceModuleEntry)
	//	assert.NotNil(t, teamWorkerParams.SpaceModuleEntry.Record)
	//	assert.NotNil(t, teamWorkerParams.SpaceModuleEntry.Data)
	//	assert.NotNil(t, teamWorkerParams.SpaceModuleEntry.Record.Data())
	//	return nil
	//}
	assert.Panics(t, func() {
		_ = RunModuleSpaceWorkerTx(ctx, nil, user, request, moduleID, new(fooModuleSpaceData), nil)
	})
}
