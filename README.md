# ⚡ GoVanity

![Build Status](https://github.com/OpenSciNet/GoVanity/actions/workflows/go.yml/badge.svg)

**GoVanity** is a high-performance, multi-threaded vanity address generator designed for decentralized ecosystems. Built by **OpenSciNet**, it allows users to mine custom-patterned addresses for Bitcoin (SegWit & Taproot), Nostr (NIP-19), and Cosmos-based chains with cryptographic precision.

> In mathematics we trust; in entropy we verify

---

## 🛠 Features

- **Multi-Chain Support:** Generate vanity addresses for:
  - **Bitcoin:** Native SegWit (`bc1q`) and Taproot (`bc1p`).
  - **Nostr:** Public keys (`npub1`).
  - **Cosmos:** Custom HRP support (e.g., `osmo`, `cosmos`, `stars`).
- **High-Performance Parallelism:** Automatically scales to use all available CPU cores.
- **Advanced Matching:** Support for `prefix`, `suffix`, and `any` position matching.
- **Luck Estimation:** Real-time probability calculation based on current hash rate and pattern difficulty.

---

## 🚀 Quick Start

### Installation
Ensure you have [Go](https://go.dev/) installed (v1.21+ recommended).

```bash
git clone https://github.com/OpenSciNet/GoVanity
cd govanity
go build -o govanity .
````

### Usage Examples

**Mine a Bitcoin Taproot address starting with "scinet":**

```bash
./govanity -t bc1p -o prefix scinet
```

**Mine a Cosmos address ending in "99" (using Osmosis HRP):**

```bash
./govanity -t cosmos -hrp osmo -o suffix 99
```

**Find a Nostr `npub` containing the string "cool":**

```bash
./govanity -t nostr -o any cool
```

-----

## 📊 Command Line Flags

| Flag | Default | Description |
| :--- | :--- | :--- |
| `-t` | `cosmos` | Address layout: `bc1q`, `bc1p`, `cosmos`, `nostr` |
| `-hrp` | `cosmos` | The Human Readable Part (Cosmos layouts only) |
| `-o` | `prefix` | Matching mode: `prefix`, `suffix`, or `any` |

-----

## ⚖️ Licensing

This project utilizes a **Dual-License** structure to balance protocol compatibility with ecosystem protection:

1.  **GoVanity Core:** Licensed under the **Business Source License 1.1 (BSL)**.
      - Non-commercial and personal use is permitted.
      - Transitions to the **MIT License** on **December 31, 2030**.
2.  **bech32 Package:** Located in `/bech32m`, this protocol implementation is licensed under **MIT** for maximum compatibility with the Bitcoin/Cosmos ecosystem.

*Refer to the [LICENSE](https://github.com/orgs/OpenSciNet/GoVanity/LICENSE) and [bech32/LICENSE](https://github.com/orgs/OpenSciNet/GoVanity/bech32m/LICENSE) files for full legal text.*

-----

## 🛡️ Security

GoVanity uses `crypto/rand` for initial entropy. Private keys are generated locally and never leave your machine.

**Warning:** Always run vanity miners on trusted hardware. Never share your generated private keys with anyone.

-----

## 🤝 Contributing

We welcome contributions\! Please see [CONTRIBUTING.md](https://github.com/orgs/OpenSciNet/GoVanity/CONTRIBUTING.md) for guidelines on how to submit improvements or report issues.

**Author:** @nocachy | **Project:** [OpenSciNet](https://github.com/orgs/OpenSciNet)
