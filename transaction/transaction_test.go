package transaction

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestCreateUSDV(t *testing.T) {
	creatorUSDV := Account{
		AccountAddress:     "OA27TLE463PL4C6ZB4C5FBMGFYSH7QNOB7IHONN24GWO7TUX73VUPEU6CM",
		SecurityPassphrase: "buffalo dry basic lake corn glide warfare blue clip bulk salmon potato useless industry business happy neck write word diamond pigeon spray forget abandon much",
	}
	manageAddress := Account{
		AccountAddress:     "CFBVIN6Q3CW45MMCQ6VQ5VAUYVSN3A3MC5LVZUL4F3ATSMXRQ62QV3RX6Y",
	}
	reserveAddress := Account{
		AccountAddress:     "CFBVIN6Q3CW45MMCQ6VQ5VAUYVSN3A3MC5LVZUL4F3ATSMXRQ62QV3RX6Y",
	}

	freezeAddress := Account{
		AccountAddress:     "CFBVIN6Q3CW45MMCQ6VQ5VAUYVSN3A3MC5LVZUL4F3ATSMXRQ62QV3RX6Y",
	}

	clawbackAddress := Account{
		AccountAddress:     "CFBVIN6Q3CW45MMCQ6VQ5VAUYVSN3A3MC5LVZUL4F3ATSMXRQ62QV3RX6Y",
	}

	var buf bytes.Buffer
	l := log.New(&buf, "logger:", log.Lshortfile)
	a := New(&creatorUSDV,
		"https://testnet-algorand.api.purestake.io/ps1",
		"LDV76UoaH15icurAUz6Hd3CvmfQpKRZj8CkoYUM2",
		1000,
		100,
		l)

	assetID, err := a.CreateAsset("USDV",
		"V-USD",
		uint64(100000000),	//1 million
		uint32(2),
		&manageAddress,
		&reserveAddress,
		&freezeAddress,
		&clawbackAddress)
	fmt.Printf("%d, %+v", assetID, err)
}


func TestOptIn(t *testing.T) {
	OptInAccount := Account{
		AccountAddress:     "MEPXYD7V6TYDH6F3GYRMODCY45K2ORLMQVNDOQCKC3MG5MQ7DG5E5QILIQ",
		SecurityPassphrase: "music scrap fantasy bone card page output squeeze civil remove cup gloom aware drip relax assist identify quit sort music next bird outdoor abandon helmet",
	}

	assetID := uint64(10741965)  //USDV

	var buf bytes.Buffer
	l := log.New(&buf, "logger:", log.Lshortfile)
	a := New(&OptInAccount,
		"https://testnet-algorand.api.purestake.io/ps1",
		"LDV76UoaH15icurAUz6Hd3CvmfQpKRZj8CkoYUM2",
		1000,
		100,
		l)

	err := a.OptIn(assetID)
	fmt.Printf("%+v", err)
}

func TestSendAsset(t *testing.T) {

	fromAccount := Account{
		AccountAddress:     "OA27TLE463PL4C6ZB4C5FBMGFYSH7QNOB7IHONN24GWO7TUX73VUPEU6CM",
		SecurityPassphrase: "buffalo dry basic lake corn glide warfare blue clip bulk salmon potato useless industry business happy neck write word diamond pigeon spray forget abandon much",
	}

	toAccount := Account{
		AccountAddress:     "MEPXYD7V6TYDH6F3GYRMODCY45K2ORLMQVNDOQCKC3MG5MQ7DG5E5QILIQ",
	}

	assetID := uint64(10741965)  //USDV

	amount := uint64(10000)  //100 USDV

	var buf bytes.Buffer
	l := log.New(&buf, "logger:", log.Lshortfile)
	a := New(&fromAccount,
		"https://testnet-algorand.api.purestake.io/ps1",
		"LDV76UoaH15icurAUz6Hd3CvmfQpKRZj8CkoYUM2",
		1000,
		100,
		l)

	err := a.SendAsset(&toAccount, assetID, amount)
	fmt.Printf("%+v", err)

	if err == nil{
		assert.Equal(t, "", "", "No Error")
	}

}


func TestCreateEquity(t *testing.T) {
	creatorUSDV := Account{
		AccountAddress:     "QRS4TBQWFFB5M2TWOKIHOHUKNXE5UB4H6APMNM4M5KY7QZXEI2SHKEZJCY",
		SecurityPassphrase: "arm leisure wrestle oven kite thunder juice evil alert base dog useless visual fiscal candy diet possible whale actor next entry inquiry ecology absorb demand",
	}
	manageAddress := Account{
		AccountAddress:     "6WKGTTSI2NADMNVPC2KL4B5ZRYXO36IPOADZSL2XCBTSJB5TLWEJR66KEQ",
		SecurityPassphrase: "reward anxiety liberty black ritual unveil credit gadget hotel borrow record increase cattle mutual erode wide gap aerobic true resource hint tent cross abstract mix",
	}
	reserveAddress := Account{
		AccountAddress:     "6WKGTTSI2NADMNVPC2KL4B5ZRYXO36IPOADZSL2XCBTSJB5TLWEJR66KEQ",
		SecurityPassphrase: "reward anxiety liberty black ritual unveil credit gadget hotel borrow record increase cattle mutual erode wide gap aerobic true resource hint tent cross abstract mix",
	}

	freezeAddress := Account{
		AccountAddress:     "6WKGTTSI2NADMNVPC2KL4B5ZRYXO36IPOADZSL2XCBTSJB5TLWEJR66KEQ",
		SecurityPassphrase: "reward anxiety liberty black ritual unveil credit gadget hotel borrow record increase cattle mutual erode wide gap aerobic true resource hint tent cross abstract mix",
	}

	clawbackAddress := Account{
		AccountAddress:     "6WKGTTSI2NADMNVPC2KL4B5ZRYXO36IPOADZSL2XCBTSJB5TLWEJR66KEQ",
		SecurityPassphrase: "reward anxiety liberty black ritual unveil credit gadget hotel borrow record increase cattle mutual erode wide gap aerobic true resource hint tent cross abstract mix",
	}

	var buf bytes.Buffer
	l := log.New(&buf, "logger:", log.Lshortfile)
	a := New(&creatorUSDV,
		"https://testnet-algorand.api.purestake.io/ps1",
		"LDV76UoaH15icurAUz6Hd3CvmfQpKRZj8CkoYUM2",
		1000,
		100,
		l)

	assetID, err := a.CreateAsset("Microsoft",
		"i-MSFT",
		uint64(1000),
		uint32(0),
		&manageAddress,
		&reserveAddress,
		&freezeAddress,
		&clawbackAddress)
	fmt.Printf("%d, %+v", assetID, err)
}

func TestOptInShare(t *testing.T) {
	OptInAccount := Account{
		AccountAddress:     "MEPXYD7V6TYDH6F3GYRMODCY45K2ORLMQVNDOQCKC3MG5MQ7DG5E5QILIQ",
		SecurityPassphrase: "music scrap fantasy bone card page output squeeze civil remove cup gloom aware drip relax assist identify quit sort music next bird outdoor abandon helmet",
	}

	assetID := uint64(10740658)  //MSFT

	var buf bytes.Buffer
	l := log.New(&buf, "logger:", log.Lshortfile)
	a := New(&OptInAccount,
		"https://testnet-algorand.api.purestake.io/ps1",
		"LDV76UoaH15icurAUz6Hd3CvmfQpKRZj8CkoYUM2",
		1000,
		100,
		l)

	err := a.OptIn(assetID)
	fmt.Printf("%+v", err)
}

func TestAtomicSwap(t *testing.T) {

	account1 := Account{
		AccountAddress:     "MEPXYD7V6TYDH6F3GYRMODCY45K2ORLMQVNDOQCKC3MG5MQ7DG5E5QILIQ",
		SecurityPassphrase: "music scrap fantasy bone card page output squeeze civil remove cup gloom aware drip relax assist identify quit sort music next bird outdoor abandon helmet",
	}

	account2 := Account{
		AccountAddress:     "QRS4TBQWFFB5M2TWOKIHOHUKNXE5UB4H6APMNM4M5KY7QZXEI2SHKEZJCY",
		SecurityPassphrase: "arm leisure wrestle oven kite thunder juice evil alert base dog useless visual fiscal candy diet possible whale actor next entry inquiry ecology absorb demand",
	}

	assetID1 := uint64(10741965)
	assetID2 := uint64(10740658)

	amount1 := uint64(9000)		//90 USDV
	amount2 := uint64(1)		//1 MSFT

	var buf bytes.Buffer
	l := log.New(&buf, "logger:", log.Lshortfile)
	a := New(&account1,
		"https://testnet-algorand.api.purestake.io/ps1",
		"LDV76UoaH15icurAUz6Hd3CvmfQpKRZj8CkoYUM2",
		1000,
		100,
		l)

	err := a.AtomicSwap(&account2, assetID1, amount1, assetID2, amount2)
	fmt.Printf("%+v", err)

	if err == nil{
		assert.Equal(t, "", "", "No Error")
	}

}

