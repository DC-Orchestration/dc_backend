# 🔐 Secure Collaborative Document Platform

A microservice-based, real-time collaborative document editing and signing platform tailored for regulated industries such as banking, insurance, and legal. Built with strong encryption, blockchain-auditable logs, and AI-assisted legal compliance.

---

## 📦 Overview

This platform consists of three major services:

- **Main Backend (Go + SQLite)** – Core authentication, document handling, real-time sync, and signature logic.
- **AI Service (FastAPI)** – NLP-based clause detection and compliance rule validation.
- **Blockchain Gateway (Go/FastAPI)** – Anchors audit logs to Hyperledger for tamperproof traceability.

---

## 🚀 Features

- 🔄 **Real-time collaborative editing**
- 🔐 **End-to-end encryption**
- 🧑‍⚖️ **Role-based access control**
- 📝 **eIDAS-compliant digital signatures**
- 📜 **Immutable blockchain audit log**
- 🤖 **AI-powered compliance checks**
- 🌉 **Integration with core banking/ERP systems**

---

## 🧱 Architecture

### Services

| Service             | Stack                  | Description                                     |
|---------------------|------------------------|-------------------------------------------------|
| Main Backend        | Go + SQLite            | Handles auth, docs, permissions, sync, signing |
| AI Compliance       | FastAPI (Python)       | NLP clause analysis and rule validation        |
| Blockchain Gateway  | Go or FastAPI          | Anchors logs to Hyperledger/Ethereum L2        |
| Frontend            | React + TipTap/Slate   | Rich-text editor with encryption & sync        |

---

## 🧪 Setup Instructions

### Prerequisites

- Go 1.20+
- Python 3.10+
- Node.js 18+
- SQLite (included with Python/Go)
- Hyperledger Fabric (for audit gateway)

### Running the Services

#### 1. Main Backend (Go + SQLite)

```bash
cd main-service
go run main.go
