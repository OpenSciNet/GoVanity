# Contributing to GoVanity

First off, thank you for considering contributing to OpenSciNet! It's people like you that make decentralized infrastructure possible.

## ⚖️ Licensing Agreement
By contributing to this repository, you agree that:
1. Your contributions to the **core miner** will be licensed under the **BSL 1.1** until the Change Date (Dec 31, 2030), after which they will transition to **MIT**.
2. Your contributions to the `/bech32` directory will be licensed under the **MIT License**.

## 🚀 How to Contribute

### 1. Reporting Bugs
- Check the "Issues" tab to see if the bug has already been reported.
- If not, open a new issue. Include your OS, Go version, and the exact command that caused the error.

### 2. Suggesting Enhancements
- We are always looking for better hardware-backed entropy sources or faster EC multiplication logic.
- Open an issue titled "Feature Request" to discuss your idea before writing code.

### 3. Pull Requests
- Fork the repo and create your branch from `main`.
- Ensure your code is formatted with `go fmt`.
- Write a clear description of what your PR changes.
- **Header Check:** Ensure new files include the appropriate BSL or MIT header as defined in the project structure.

## 🛠 Development Environment
To build GoVanity locally:
```bash
git clone [https://github.com/openscinet/govanity](https://github.com/openscinet/govanity)
cd govanity
go build -o vanity .