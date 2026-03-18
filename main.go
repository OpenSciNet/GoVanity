/*
 * ____________________________________________________________________________
 * |                                                                          |
 * |     ██████╗  ██████╗ ██╗   ██╗ █████╗ ███╗   ██╗██╗████████╗██╗   ██╗    |
 * |    ██╔════╝ ██╔═══██╗██║   ██║██╔══██╗████╗  ██║██║╚══██╔══╝╚██╗ ██╔╝    |
 * |    ██║  ███╗██║   ██║██║   ██║███████║██╔██╗ ██║██║   ██║    ╚████╔╝     |
 * |    ██║   ██║██║   ██║╚██╗ ██╔╝██╔══██║██║╚██╗██║██║   ██║     ╚██╔╝      |
 * |    ╚██████╔╝╚██████╔╝ ╚████╔╝ ██║  ██║██║ ╚████║██║   ██║      ██║       |
 * |     ╚═════╝  ╚═════╝   ╚═══╝  ╚═╝  ╚═╝╚═╝  ╚═══╝╚═╝   ╚═╝      ╚═╝       |
 * |__________________________________________________________________________|
 * |                                                                          |
 * |  PROJECT:	GoVanity Vanity Miner v0.1                                    |
 * |  AUTHOR:	@nocachy                                                      |
 * |  DATE:     March 13 2026                                                 |
 * |  STATUS:	Active | Entropy Pool: Verified                               |
 * |  LICENSE:	Business Source License 1.1									  |
 * |  NOTE: 	The 'bech32' sub-package is licensed under MIT				  |
 * |__________________________________________________________________________|
 * |                                                                          |
 * |  "In mathematics we trust; in entropy we verify."                        |
 * |__________________________________________________________________________|
 */
package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"govanity/bech32m"
	"math"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/txscript"
	//"github.com/btcsuite/btcutil/base58"
	"github.com/cosmos/btcutil/bech32"
	"github.com/nbd-wtf/go-nostr/nip19"
	"golang.org/x/crypto/ripemd160"
)

func hash160(b []byte) []byte {
	sha := sha256.Sum256(b)
	h := ripemd160.New()
	h.Write(sha[:])
	return h.Sum(nil)
}

func main() {
	layout := flag.String("t", "cosmos", "address layout (bc1q, bc1p, cosmos, nostr)")
	hrp := flag.String("hrp", "cosmos", "bech32 hrp for cosmos only")
	suffixMode := flag.String("o", "prefix", "mine prefix, suffix or any")
	//
	flag.Usage = func() {
		fmt.Printf("\n✨ GoVanity Vanity Miner\n")
		fmt.Printf("Usage: govanity [FLAGS] <pattern>\n\n")
		fmt.Printf("Example: govanity -t bc1q -o prefix lucky\n")
		fmt.Printf("Example: govanity -t cosmos -hrp osmo -o suffix 99\n\n")
		fmt.Printf("Flags:\n")
		flag.PrintDefaults()
	}
	flag.Parse()
	//	check arguments length == 1
	if flag.NArg() < 1 {
		fmt.Println("❌ Error: Missing pattern.")
		fmt.Println("💡 Hint: You must provide a pattern to mine at the end of the command.")
		fmt.Println("👉 Try: govanity -t nostr 0000")
		fmt.Println("\nRun 'govanity h' for full help.")
		return
	} else if flag.NArg() > 1 {
		fmt.Printf("❌ Error: Too many arguments detected.\n")
		fmt.Println("💡 Hint: If your pattern contains spaces, wrap it in \"quotes\".")
		fmt.Println("👉 Proper usage: govanity [flags] <single_pattern>")
		return
	}
	//	init
	target := strings.ToLower(flag.Arg(0))
	numWorkers := runtime.NumCPU()
	var found int32
	var attempts uint64
	start := time.Now()
	mode := "prefix"
	switch *suffixMode {
		case "suffix":
			mode = "suffix"
		case "any":
			mode = "any"
		default:
			mode = "prefix"
	}
	fmt.Printf("Mining %s (%s '%s') using %d workers\n", *layout, mode, target, numWorkers)
	var wg sync.WaitGroup
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			privBytes := make([]byte, 32)
			rand.Read(privBytes)
			localAttempts := uint64(0)
			for atomic.LoadInt32(&found) == 0 {
				priv, pub := btcec.PrivKeyFromBytes(privBytes)
				pubBytes := pub.SerializeCompressed()
				var addr string
				switch *layout {
					case "cosmos":
						h160 := hash160(pubBytes)
						converted, _ := bech32.ConvertBits(h160, 8, 5, true)
						addr, _ = bech32.Encode(*hrp, converted)
					case "bc1q":
						h160 := hash160(pubBytes)
						converted, _ := bech32.ConvertBits(h160, 8, 5, true)
						data := append([]byte{0}, converted...)
						addr, _ = bech32.Encode("bc", data)
					case "bc1p":
						//ikey := pub.SerializeCompressed()[1:33]				//	internal key
						tapKey := txscript.ComputeTaprootKeyNoScript(pub)		//	Taproot tweak (no script tree)
						xonly := schnorr.SerializePubKey(tapKey)				//	x-only pubkey (32 bytes)
						converted, _ := bech32.ConvertBits(xonly, 8, 5, true)
						data := append([]byte{1}, converted...)
						addr = bech32m.EncodeBech32m("bc", data)
					case "nostr":
						pubHex := hex.EncodeToString(pub.SerializeCompressed()[1:])
						addr, _ = nip19.EncodePublicKey(pubHex)
					//case "zec":
					//	h160 := hash160(pubBytes)
					//	version := []byte{0x1C, 0xB8} // Zcash t1
					//	payload := append(version, h160...)
					//	check := sha256.Sum256(payload)
					//	check2 := sha256.Sum256(check[:])
					//	full := append(payload, check2[:4]...)
					//	addr = base58.Encode(full)
					default:
						fmt.Println("unknown layout")
						return
				}
				var check string
				if *layout == "cosmos" {
					check = addr[len(*hrp)+1:]
				} else if strings.HasPrefix(addr, "bc1") {
					check = addr[4:]
				} else if strings.HasPrefix(addr, "npub1") {
					check = addr[5:]
				} else if strings.HasPrefix(addr, "t1") {
					check = addr[2:]
				} else {
					check = addr
				}
				var match bool
				switch *suffixMode {
					case "any":
						match = strings.Contains(check, target)
					case "suffix":
						match = strings.HasSuffix(check, target)
					default:
						match = strings.HasPrefix(check, target)
				}
				if match && atomic.CompareAndSwapInt32(&found, 0, 1) {
				fmt.Println("\n\nSUCCESS")
					fmt.Println("address:", addr)
					fmt.Printf("privkey: %x\n", priv.Serialize())
				return
				}
				for i := 31; i >= 0; i-- {
					privBytes[i]++
					if privBytes[i] != 0 {
						break
					}
				}
				localAttempts++
				if localAttempts%1000 == 0 {
					atomic.AddUint64(&attempts, 1000)
				}
			}
		}()
	}
	go func() {
		p := 1.0 / math.Pow(32, float64(len(target)))
		for atomic.LoadInt32(&found) == 0 {
			time.Sleep(3 * time.Second)
			a := atomic.LoadUint64(&attempts)
			pPerKey := p
			if *suffixMode == "any" {
				pPerKey = 1 - math.Pow(1-p, 35)
			}
			// luck = 1 - exp(a * ln(1 - pPerKey))
			// For very small p, ln(1 - p) is approximately -p
			luck := 1.0 - math.Exp(float64(a) * math.Log(1.0-pPerKey))
			speed := float64(a) / time.Since(start).Seconds()
			fmt.Printf("\r%.0f keys/sec | %d tested | %.2f%%", speed, a, luck*100)
			}
	}()
	wg.Wait()
}