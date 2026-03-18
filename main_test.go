package main

import (
	"encoding/hex"
	"testing"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/txscript"
    "github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/cosmos/btcutil/bech32"
	"govanity/bech32m"
)

// ---- TEST VECTORS ----

type testVector struct {
	priv    string
	nostr   string
	bc1q    string
	bc1p    string
	cosmos  string
}

var vectors = []testVector{
	{
		priv:   "9a6ad7dd18a5a68d0b90c77d6c16ee214062608a0e13fe07b6728f093748cbd4",
		nostr:  "npub18kgh0fjjwmykx83pztef4nucfdq4a2p2pjyyayqm4e3l4m7e7fpsddfkz3",
		bc1q:   "bc1q2apmutzzaupxjw2wzyatfsqpzzagruvkjcry4p",
		bc1p:   "bc1pxfh7l4ua2ermpsamyk4nk9tv4f4p6psamk8glzv8nfv0vpr3tfuqvphu4a",
		cosmos: "cosmos12apmutzzaupxjw2wzyatfsqpzzagruvkyza03v",
	},
	{
		priv:   "41e6118f80e851479030a4a3020d6f7586a52ac3d5ca9e0870cfb005bcf0b0e9",
		nostr:  "npub1n0str5v4y7z4dr6sngs35c0q26ngam2j25sxjr4qamyghwj5wlrsl5d05p",
		bc1q:   "bc1qjwmmzqzvhc2gu7lvafdh42y0atyg8crww7q24e",
		bc1p:   "bc1pz7j7hj9hm4aj4zaq8ljzxlmkx3xkafsnljp76kfrahc74f8ja27sk5d90r",
		cosmos: "cosmos1jwmmzqzvhc2gu7lvafdh42y0atyg8crwcy7p35",
	},
}

// ---- MAIN TEST ----

func TestAddressGeneration(t *testing.T) {
	for i, v := range vectors {
		privBytes, err := hex.DecodeString(v.priv)
		if err != nil {
			t.Fatalf("vector %d: invalid privkey", i)
		}

		_, pub := btcec.PrivKeyFromBytes(privBytes)
		pubBytes := pub.SerializeCompressed()

		// ---- HASH160 ----
		h160 := hash160(pubBytes)

		// ---- COSMOS ----
		t.Run("Cosmos_"+v.priv[:6], func(t *testing.T) {
			converted, _ := bech32.ConvertBits(h160, 8, 5, true)
			addr, _ := bech32.Encode("cosmos", converted)

			if addr != v.cosmos {
				t.Errorf("cosmos mismatch\n got: %s\nwant: %s", addr, v.cosmos)
			}
		})

		// ---- BITCOIN bc1q ----
		t.Run("BTC_bc1q_"+v.priv[:6], func(t *testing.T) {
			converted, _ := bech32.ConvertBits(h160, 8, 5, true)
			data := append([]byte{0}, converted...)
			addr, _ := bech32.Encode("bc", data)

			if addr != v.bc1q {
				t.Errorf("bc1q mismatch\n got: %s\nwant: %s", addr, v.bc1q)
			}
		})

		// ---- BITCOIN TAPROOT bc1p ----
		t.Run("BTC_bc1p_"+v.priv[:6], func(t *testing.T) {
			// Compute Taproot output key (BIP341 tweak)
			tapKey := txscript.ComputeTaprootKeyNoScript(pub)
			// Serialize as x-only (BIP340)
			xonly := schnorr.SerializePubKey(tapKey)
			converted, _ := bech32.ConvertBits(xonly, 8, 5, true)
			data := append([]byte{1}, converted...)
			addr := bech32m.EncodeBech32m("bc", data)

			if addr != v.bc1p {
				t.Errorf("bc1p mismatch\n got: %s\nwant: %s", addr, v.bc1p)
			}
		})

		// ---- NOSTR npub ----
		t.Run("Nostr_"+v.priv[:6], func(t *testing.T) {
			// Nostr uses 32-byte pubkey (no prefix)
			xonly := pubBytes[1:33]

			converted, _ := bech32.ConvertBits(xonly, 8, 5, true)
			addr, _ := bech32.Encode("npub", converted)

			if addr != v.nostr {
				t.Errorf("npub mismatch\n got: %s\nwant: %s", addr, v.nostr)
			}
		})
	}
}