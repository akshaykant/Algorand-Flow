package transaction

import (
	"encoding/base64"
	"fmt"
	"github.com/algorand/go-algorand-sdk/types"
	"log"

	"github.com/algorand/go-algorand-sdk/client/algod"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/mnemonic"
	"github.com/algorand/go-algorand-sdk/transaction"

)

type Algo interface {
	GenerateAccount() (*Account, error)
	CreateAsset(string, string, uint64, uint32, *Account, *Account, *Account, *Account) (uint64, error)
	OptIn(uint64) error
	SendAsset(*Account, uint64, uint64) error
	AtomicSwap(*Account, uint64, uint64,  uint64, uint64) error
}

type algo struct {
	address      *Account
	apiAddress   string
	apiKey       string
	minFee       uint64
	seedAlgo     uint64
	l            *log.Logger
}

func New(address *Account, apiAddress, apiKey string, minFee, seedAlgo uint64, l *log.Logger) Algo {
	return &algo{
		address:      address,
		apiAddress:   apiAddress,
		apiKey:       apiKey,
		minFee:       minFee,
		seedAlgo:     seedAlgo,
		l:            l,
	}
}


func (a *algo) GenerateAccount() (*Account, error) {
	account := crypto.GenerateAccount()
	paraphrase, err := mnemonic.FromPrivateKey(account.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("generateAccount: error generating account: %s", err)
	}

	return &Account{
		AccountAddress:     account.Address.String(),
		PrivateKey:         string(account.PrivateKey),
		SecurityPassphrase: paraphrase,
	}, nil
}

func (a *algo) CreateAsset(tokenName string, tokenUnit string, totalSupply uint64, decimalSupply uint32, managerAddress, reserveAddress, freezeAddress, clawbackAddress  *Account) (uint64, error) {

	var headers []*algod.Header
	headers = append(headers, &algod.Header{Key: "X-API-Key", Value: a.apiKey})
	algodClient, err := algod.MakeClientWithHeaders(a.apiAddress, "", headers)
	if err != nil {
		return 0, fmt.Errorf("createAsset: error connecting to algo: %s", err)
	}

	txParams, err := algodClient.SuggestedParams()
	if err != nil {
		return 0, fmt.Errorf("createAsset: error getting suggested tx params: %s", err)
	}

	genID := txParams.GenesisID
	genHash := txParams.GenesisHash
	firstValidRound := txParams.LastRound
	lastValidRound := firstValidRound + 1000

	// Create an asset
	// Set parameters for asset creation transaction
	creator := a.address
	assetName := tokenName
	unitName := tokenUnit
	assetURL := "https://asseturl.com"
	assetMetadataHash := "thisIsSomeLength32HashCommitment"
	defaultFrozen := false
	decimals := decimalSupply
	totalIssuance := totalSupply
	manager := managerAddress
	reserve := reserveAddress
	freeze := freezeAddress
	clawback := clawbackAddress
	note := []byte(nil)
	txn, err := transaction.MakeAssetCreateTxn(creator.AccountAddress, a.minFee, firstValidRound, lastValidRound, note,
		genID, base64.StdEncoding.EncodeToString(genHash), totalIssuance, decimals, defaultFrozen, manager.AccountAddress, reserve.AccountAddress, freeze.AccountAddress, clawback.AccountAddress,
		unitName, assetName, assetURL, assetMetadataHash)
	if err != nil {
		return 0, fmt.Errorf("createAsset: failed to make asset: %s", err)
	}
	fmt.Printf("Asset created AssetName: %s\n", txn.AssetConfigTxnFields.AssetParams.AssetName)

	privateKey, err := mnemonic.ToPrivateKey(a.address.SecurityPassphrase)
	if err != nil {
		return 0, fmt.Errorf("createAsset: error getting private key from mnemonic: %s", err)
	}

	txid, stx, err := crypto.SignTransaction(privateKey, txn)
	if err != nil {
		return 0, fmt.Errorf("createAsset: failed to sign transaction: %s", err)
	}
	a.l.Printf("Signed txid: %s", txid)
	// Broadcast the transaction to the network
	txHeaders := append([]*algod.Header{}, &algod.Header{Key: "Content-Type", Value: "application/x-binary"})
	sendResponse, err := algodClient.SendRawTransaction(stx, txHeaders...)
	if err != nil {
		return 0, fmt.Errorf("createAsset: failed to send transaction: %s", err)
	}

	// Wait for transaction to be confirmed
	waitForConfirmation(algodClient, sendResponse.TxID)

	// Retrieve asset ID by grabbing the max asset ID
	// from the creator account's holdings.
	act, err := algodClient.AccountInformation(a.address.AccountAddress)
	if err != nil {
		return 0, fmt.Errorf("createAsset: failed to get account information: %s", err)
	}

	assetID := uint64(0)
	for i := range act.AssetParams {
		if i > assetID {
			assetID = i
		}
	}

	a.l.Printf("createAsset: asset ID from AssetParams: %d", assetID)
	// Retrieve asset info.
	assetInfo, err := algodClient.AssetInformation(assetID)
	if err != nil {
		return 0, fmt.Errorf("createAsset: error getting asset info: %s", err)
	}

	a.l.Printf("createAsset: assets info: %+v", assetInfo)

	return assetID, nil
}

func (a *algo) OptIn(assetID uint64) error {
	var headers []*algod.Header
	headers = append(headers, &algod.Header{Key: "X-API-Key", Value: a.apiKey})
	algodClient, err := algod.MakeClientWithHeaders(a.apiAddress, "", headers)
	if err != nil {
		return fmt.Errorf("optin: error connecting to algo: %s", err)
	}

	txParams, err := algodClient.SuggestedParams()
	if err != nil {
		return fmt.Errorf("optin: error getting suggested tx params: %s", err)
	}

	note := []byte(fmt.Sprintf("Opting in from %s", a.address.AccountAddress))
	genID := txParams.GenesisID
	genHash := txParams.GenesisHash
	firstValidRound := txParams.LastRound
	lastValidRound := firstValidRound + 1000

	// Account opts in to receive asset
	txn, err := transaction.MakeAssetAcceptanceTxn(a.address.AccountAddress, a.minFee, firstValidRound,
		lastValidRound, note, genID, base64.StdEncoding.EncodeToString(genHash), assetID)
	if err != nil {
		return fmt.Errorf("optin: failed to send transaction MakeAssetAcceptanceTxn: %s", err)
	}

	privateKey, err := mnemonic.ToPrivateKey(a.address.SecurityPassphrase)
	if err != nil {
		return fmt.Errorf("optin: error getting private key from mnemonic: %s", err)
	}

	txid, stx, err := crypto.SignTransaction(privateKey, txn)
	if err != nil {
		return fmt.Errorf("optin: failed to sign transaction: %s", err)
	}

	fmt.Printf("Transaction ID: %s\n", txid)
	// Broadcast the transaction to the network
	txHeaders := append([]*algod.Header{}, &algod.Header{Key: "Content-Type", Value: "application/x-binary"})
	sendResponse, err := algodClient.SendRawTransaction(stx, txHeaders...)
	if err != nil {
		return fmt.Errorf("optin: failed to send transaction: %s", err)
	}

	a.l.Printf("optin: transaction ID raw: %s", sendResponse.TxID)

	// Wait for transaction to be confirmed
	waitForConfirmation(algodClient, sendResponse.TxID)

	act, err := algodClient.AccountInformation(a.address.AccountAddress)
	if err != nil {
		return fmt.Errorf("optin: failed to get account information: %s", err)
	}

	a.l.Printf("optin: account info: %+v", act.Assets[assetID])

	return nil
}

func (a *algo) SendAsset(to *Account, assetID uint64, total uint64) error {
	var headers []*algod.Header
	headers = append(headers, &algod.Header{Key: "X-API-Key", Value: a.apiKey})
	algodClient, err := algod.MakeClientWithHeaders(a.apiAddress, "", headers)
	if err != nil {
		return fmt.Errorf("sendAsset: error connecting to algo: %s", err)
	}

	txParams, err := algodClient.SuggestedParams()
	if err != nil {
		return fmt.Errorf("sendAsset: error getting suggested tx params: %s", err)
	}

	note := []byte("Transferring asset")
	genID := txParams.GenesisID
	genHash := txParams.GenesisHash
	firstValidRound := txParams.LastRound
	lastValidRound := firstValidRound + 1000

	// Send  1 of asset from Account to Account
	sender := a.address.AccountAddress
	recipient := to.AccountAddress
	amount := total
	closeRemainderTo := ""
	txn, err := transaction.MakeAssetTransferTxn(sender, recipient,
		closeRemainderTo, amount, a.minFee, firstValidRound, lastValidRound, note,
		genID, base64.StdEncoding.EncodeToString(genHash), assetID)
	if err != nil {
		return fmt.Errorf("sendAsset: failed to send transaction MakeAssetTransfer Txn: %s", err)
	}

	privateKey, err := mnemonic.ToPrivateKey(a.address.SecurityPassphrase)
	if err != nil {
		return fmt.Errorf("sendAsset: error getting private key from mnemonic: %s", err)
	}

	txid, stx, err := crypto.SignTransaction(privateKey, txn)
	if err != nil {
		return fmt.Errorf("sendAsset: failed to sign transaction: %s", err)
	}
	fmt.Printf("Transaction ID: %s", txid)
	// Broadcast the transaction to the network
	txHeaders := append([]*algod.Header{}, &algod.Header{Key: "Content-Type", Value: "application/x-binary"})
	sendResponse, err := algodClient.SendRawTransaction(stx, txHeaders...)
	if err != nil {
		return fmt.Errorf("sendAsset: failed to send transaction: %s", err)
	}
	fmt.Printf("Transaction ID raw: %s\n", sendResponse.TxID)

	// Wait for transaction to be confirmed
	waitForConfirmation(algodClient, sendResponse.TxID)

	act, err := algodClient.AccountInformation(to.AccountAddress)
	if err != nil {
		return fmt.Errorf("sendAsset: failed to get account information: %s", err)
	}

	a.l.Printf("sendAsset: account info: %v", act.Assets[assetID])
	return nil
}

func (act *algo) AtomicSwap(act2 *Account, assetID1, amount1, assetID2, amount2 uint64) error {
	var headers []*algod.Header
	headers = append(headers, &algod.Header{Key: "X-API-Key", Value: act.apiKey})
	algodClient, err := algod.MakeClientWithHeaders(act.apiAddress, "", headers)
	if err != nil {
		return fmt.Errorf("optin: error connecting to algo: %s", err)
	}

	txParams, err := algodClient.SuggestedParams()
	if err != nil {
		return fmt.Errorf("optin: error getting suggested tx params: %s", err)
	}

	note1 := []byte(fmt.Sprintf("atomic swap to %s", act2.AccountAddress))
	note2 := []byte(fmt.Sprintf("atomic swap to %s", act.address.AccountAddress))
	genID := txParams.GenesisID
	genHash := txParams.GenesisHash
	firstValidRound := txParams.LastRound
	lastValidRound := firstValidRound + 1000
	closeRemainderTo := ""

	// Create unsigned transactions
	// from account 1 to account 2
	txn1, err := transaction.MakeAssetTransferTxn(act.address.AccountAddress, act2.AccountAddress,
		closeRemainderTo, amount1, act.minFee, firstValidRound, lastValidRound, note1,
		genID, base64.StdEncoding.EncodeToString(genHash), assetID1)
	if err != nil {
		return fmt.Errorf("Error creating transaction: %s\n", err)
	}

	// from account 2 to account 1
	txn2, err := transaction.MakeAssetTransferTxn(act2.AccountAddress, act.address.AccountAddress,
		closeRemainderTo, amount2, act.minFee, firstValidRound, lastValidRound, note2,
		genID, base64.StdEncoding.EncodeToString(genHash), assetID2)
	if err != nil {
		return fmt.Errorf("Error creating transaction: %s\n", err)
	}

	//Combine Transaction
	// compute group id and put it into each transaction
	gid, err := crypto.ComputeGroupID([]types.Transaction{txn1, txn2})
	txn1.Group = gid
	txn2.Group = gid

	//Sign Transaction by each account for authorisation
	privateKey1, err := mnemonic.ToPrivateKey(act.address.SecurityPassphrase)
	if err != nil {
		return fmt.Errorf("optin: error getting private key from mnemonic: %s", err)
	}

	privateKey2, err := mnemonic.ToPrivateKey(act2.SecurityPassphrase)
	if err != nil {
		return fmt.Errorf("optin: error getting private key from mnemonic: %s", err)
	}

	_, stx1, err := crypto.SignTransaction(privateKey1, txn1)
	if err != nil {
		return fmt.Errorf("Failed to sign transaction: %s\n", err)
	}
	_, stx2, err := crypto.SignTransaction(privateKey2, txn2)
	if err != nil {
		return fmt.Errorf("Failed to sign transaction: %s\n", err)
	}

	//Assemble transaction group
	var signedGroup []byte
	signedGroup = append(signedGroup, stx1...)
	signedGroup = append(signedGroup, stx2...)


	//send transaction group
	// Broadcast the transaction to the network
	txHeaders := append([]*algod.Header{}, &algod.Header{Key: "Content-Type", Value: "application/x-binary"})
	sendResponse, err := algodClient.SendRawTransaction(signedGroup, txHeaders...)
	if err != nil {
		return fmt.Errorf("optin: failed to send transaction: %s", err)
	}

	act.l.Printf("atomic swap: transaction ID raw: %s", sendResponse.TxID)

	// Wait for transaction to be confirmed
	waitForConfirmation(algodClient, sendResponse.TxID)

	return nil
}

// Function that waits for a given txId to be confirmed by the network
func waitForConfirmation(algodClient algod.Client, txID string) {
	for {
		pt, err := algodClient.PendingTransactionInformation(txID)
		if err != nil {
			fmt.Printf("waiting for confirmation... (pool error, if any): %s\n", err)
			continue
		}
		if pt.ConfirmedRound > 0 {
			fmt.Printf("Transaction "+pt.TxID+" confirmed in round %d\n", pt.ConfirmedRound)
			break
		}
		nodeStatus, err := algodClient.Status()
		if err != nil {
			fmt.Printf("error getting algod status: %s\n", err)
			return
		}
		algodClient.StatusAfterBlock(nodeStatus.LastRound + 1)
	}
}