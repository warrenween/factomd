package wsapi_test

import (
	/*"encoding/json"
	  "fmt"
	  "github.com/FactomProject/factomd/common/entryBlock"
	  "github.com/FactomProject/factomd/common/interfaces"
	  "github.com/FactomProject/factomd/common/primitives"
	  "github.com/FactomProject/factomd/receipts"
	  "github.com/FactomProject/web"
	  "net/http"
	  "strings"*/

	"encoding/hex"
	"testing"

	"github.com/FactomProject/factomd/testHelper"
	. "github.com/FactomProject/factomd/wsapi"
)

func TestHandleV2FactoidACK(t *testing.T) {
	state := testHelper.CreateAndPopulateTestState()
	blocks := testHelper.CreateFullTestBlockSet()

	for _, block := range blocks {
		for _, tx := range block.FBlock.GetTransactions() {
			req := AckRequest{}
			txID := tx.GetHash().String()
			req.TxID = txID

			r, jError := HandleV2FactoidACK(state, req)

			if jError != nil {
				t.Errorf("%v", jError)
				continue
			}

			resp, ok := r.(*FactoidTxStatus)
			if ok == false {
				t.Error("Invalid response type returned")
				continue
			}

			if resp.TxID != txID {
				t.Error("Invalid TxID returned")
			}
			if resp.Status != AckStatusDBlockConfirmed {
				t.Error("Invalid status returned")
			}

			req = AckRequest{}
			h, err := tx.MarshalBinary()
			if err != nil {
				t.Errorf("%v", err)
				continue
			}
			req.FullTransaction = hex.EncodeToString(h)

			r, jError = HandleV2FactoidACK(state, req)

			if jError != nil {
				t.Errorf("%v", jError)
				continue
			}

			resp, ok = r.(*FactoidTxStatus)
			if ok == false {
				t.Error("Invalid response type returned")
				continue
			}

			if resp.TxID != txID {
				t.Error("Invalid TxID returned")
			}
			if resp.Status != AckStatusDBlockConfirmed {
				t.Error("Invalid status returned")
			}
		}
	}
}

func TestHandleV2EntryACK(t *testing.T) {

}
