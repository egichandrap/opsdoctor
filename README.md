# 🛠️ OpsDoctor – DevOps Diagnostic CLI Tool

OpsDoctor adalah tools CLI berbasis **Golang + Cobra** untuk membantu DevOps, Backend Engineer, dan SRE melakukan berbagai diagnostic penting untuk server Linux/RHEL, aplikasi Spring Boot, TLS, API, dan log.

Tools ini dirancang supaya:
- Bisa dijalankan di local maupun server
- Bisa dipakai untuk debugging cepat
- Bisa dipakai DevOps untuk health check otomatis
- Bisa dipakai tim backend untuk analisa log dan masalah API
- Bisa di-extend sebagai tool open-source

---

# 🚀 Fitur

### 🌐 1. Network Diagnostic
- Ping host/IP
- HTTP status + latency check
- DNS resolve

### 🔐 2. TLS Checker
- Expiry date
- Organization / Issuer
- TLS protocol version
- Cipher suite
- Warning untuk certificate yang hampir expired

### 📊 3. Log Analyzer
- Slow log detection (ms threshold)
- Regex filter
- Extract timestamp dan latency
- JSON export
- Summary (top slow queries)

### ☕ 4. Spring Boot Version Security Check
- Parse `pom.xml`
- Deteksi versi Spring Boot
- Rekomendasi upgrade
- Flagging jika versi rawan CVE

### 🔍 5. API Tester
- Simple GET/POST
- Latency measurement
- Pretty JSON output
- JSON export

### 🏥 6. Service Checker
- Check beberapa microservices via YAML
- HTTP health endpoint checking
- Summary OK/WARN/FAILED
- JSON export

### 🎨 7. Global CLI Features
- `--verbose`
- `--output json`
- `--no-color`

---

# 🔧 Instalasi

## 1. Instal via Go (paling mudah)

Pastikan GOPATH/bin sudah masuk PATH :

```bash
echo 'export PATH=$PATH:$(go env GOPATH)/bin' >> ~/.zshrc
source ~/.zshrc
```

Install : 
```go
go install github.com/egichandrap/opsdoctor@latest
```

Jalankan :
```bash
opsdoctor --help
```

## 2. Instal via source (clone & build)
```terminaloutput
git clone https://github.com/egichandrap/opsdoctor
cd opsdoctor
```

Build untuk local OS :
```bash
go build -o opsdoctor
```

Atau cross-build (Linux/macOS/Windows) :
```bash
sh scripts/build.sh
```

Binary akan ada di : ```dist/opsdoctor```

Install ke sistem :
```bash
sudo cp dist/opsdoctor /usr/local/bin/
```


# 📁 Struktur Project

```text
opsdoctor/
│
├── main.go
├── go.mod
├── go.sum
│
├── cmd/
│   ├── root.go
│   ├── net.go
│   ├── tls.go
│   ├── log.go
│   ├── spring.go
│   ├── api.go
│   └── svc.go
│
├── internal/
│   ├── netscan/
│   ├── tlscheck/
│   ├── loganalyzer/
│   ├── springcheck/
│   ├── apitest/
│   └── svcchecker/
│
├── utils/
│   ├── color.go
│   ├── export.go
│   ├── verbose.go
│
└── scripts/
    └── build.sh
```

## 📘 Cara Menggunakan

Semua module mengikuti pola :

```bash
opsdoctor <module> <command> [args] [flags]
```

Global flags :

| Flag            | Fungsi        |
| --------------- | ------------- |
| `--verbose`     | output detail |
| `--output json` | JSON output   |
| `--no-color`    | disable warna |

## 🌐 Network Diagnostic

**PING** 
```bash
opsdoctor net ping google.com
```

**HTTP Check**
```bash
opsdoctor net http https://google.com
```

**DNS Lookup**
```bash
opsdoctor net dns google.com
```

## 🔐 TLS Checker

```bash
opsdoctor tls https://google.com
```

**Output berisi :**
* CN / Issuer
* Expiration date
* TLS version
* Cipher
* Warning jika certificate hampir expired

**JSON mode :**
```bash
opsdoctor --output json tls https://google.com
```

## 📊 Log Analyzer

**Slow log > 1000ms**
```bash
opsdoctor log analyze app.log --slow 1000
```

**Dengan pattern**
```bash
opsdoctor log analyze app.log --slow 500 --pattern "GET /api"
```

**Export JSON**
```bash
opsdoctor log analyze app.log --export output.json
```


## ☕ Spring Boot Security Version Check
```
opsdoctor spring check ./pom.xml
```

**Output :**

* Deteksi versi Spring Boot
* Rekomendasi upgrade
* CVE warning (jika masuk list vulnerable)

## 🔍 API Tester
GET
```bash
opsdoctor api get https://jsonplaceholder.typicode.com/posts/1
```

POST
```bash
opsdoctor api post https://reqres.in/api/users -d '{"name":"Egi"}'
```

## 🏥 Service Checker
**File: services.yaml**
```yaml
services:
  - name: user-service
    url: http://localhost:8080/health
  - name: payment-service
    url: http://localhost:9000/actuator/health
```

**Jalankan :**
```bash
opsdoctor svc check services.yaml
```

# 🎉 Selesai
**Jika kamu mau tambahan :**
- Badge GitHub Actions build status
- Banner ASCII saat CLI start
- Contoh screenshot output (warna-warni)
- Auto installer script (`install.sh`)
- Homebrew formula (`brew tap`)

Tinggal bilang **“lanjutkan”**.
