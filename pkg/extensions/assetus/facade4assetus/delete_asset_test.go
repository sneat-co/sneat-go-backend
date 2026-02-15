package facade4assetus

import "testing"

func TestDeleteAssetTxWorker(t *testing.T) {
	if err := deleteAssetTxWorker(nil, nil, nil); err != nil {
		t.Fatalf("deleteAssetTxWorker() error = %v, want nil", err)
	}
}
