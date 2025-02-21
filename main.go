package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/blocto/solana-go-sdk/client"
	"github.com/blocto/solana-go-sdk/program/sysprog"
	"github.com/blocto/solana-go-sdk/types"
	"log"
	"os"
)

// pastikan disini saat ngecreate account, private dan public key disimpan as solana key
func CreateGetAccount(name string) *types.Account {
	if _, err := os.Stat(fmt.Sprintf("%s.xx", name)); err != nil {
		account := types.NewAccount()

		data, err := json.Marshal(account.PrivateKey)
		if err != nil {
			log.Fatalf("Failed to encode private key: %v", err)
		}

		if err := os.WriteFile(fmt.Sprintf("%s.xx", name), data, 0600); err != nil {
			log.Fatalf("Failed to save key file: %v", err)
		}

		fmt.Println("Created New Account:", name)
		fmt.Println("Public Key:", account.PublicKey.ToBase58())
		fmt.Println("Private Key:", string(account.PrivateKey))

		return &account
	}

	data, err := os.ReadFile(fmt.Sprintf("%s.xx", name))
	if err != nil {
		log.Fatalf("Failed to read key file: %v", err)
	}

	var privateKey []byte
	if err := json.Unmarshal(data, &privateKey); err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}

	account, err := types.AccountFromBytes(privateKey)
	if err != nil {
		log.Fatalf("Failed to load account: %v", err)
	}

	log.Println("account:", account.PublicKey.ToBase58())

	return &account

}

func main() {
	ctx := context.Background()
	c := client.NewClient("https://api.devnet.solana.com")
	version, err := c.GetVersion(ctx)
	if err != nil {
		log.Fatalf("Gagal mendapatkan versi Solana: %v", err)
	}
	fmt.Println("Versi Solana", version.SolanaCore)

	accountFajrul := CreateGetAccount("fajrul")
	accountWidia := CreateGetAccount("widia")

	balance, err := c.GetBalance(ctx, accountFajrul.PublicKey.ToBase58())
	if err != nil {
		log.Fatalf("Gagal mendapatkan saldo: %v", err)
	}
	fmt.Println("Saldo SOL Fajrul:", balance)

	// Balance fajrul jika kosong tambahin
	// atau lewat sini ngisinya https://faucet.solana.com/
	if balance == 0 {
		if _, err = c.RequestAirdrop(ctx, accountFajrul.PublicKey.ToBase58(), 1000000); err != nil {
			log.Fatalf("Gagal mendapatkan airdrop: %v", err)
		}
	}

	// Balance widia
	balance, err = c.GetBalance(ctx, accountWidia.PublicKey.ToBase58())
	if err != nil {
		log.Fatalf("Gagal mendapatkan saldo: %v", err)
	}
	fmt.Println("Saldo SOL Widia:", balance)

	recent, err := c.GetLatestBlockhash(ctx)
	if err != nil {
		log.Fatalf("Gagal mendapatkan recent blockhash: %v", err)
	}

	log.Println("Latest Blockhash:", recent)

	message := types.NewMessage(types.NewMessageParam{
		FeePayer:        accountFajrul.PublicKey,
		RecentBlockhash: recent.Blockhash,
		Instructions: []types.Instruction{sysprog.Transfer(sysprog.TransferParam{
			From:   accountFajrul.PublicKey,
			To:     accountWidia.PublicKey,
			Amount: 100000,
		})},
	})

	tx, err := types.NewTransaction(types.NewTransactionParam{
		Message: message,
		Signers: []types.Account{*accountFajrul},
	})
	if err != nil {
		log.Fatalf("Gagal membuat NewTransaction: %v", err)
	}

	txHash, err := c.SendTransaction(context.Background(), tx)
	if err != nil {
		log.Fatalf("Gagal mengirim transaksi: %v", err)
	}
	log.Printf("Transaksi berhasil dengan hash: %s", txHash)

}
