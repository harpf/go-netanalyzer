# Go NetAnalyzer

**Go NetAnalyzer** is a CLI tool to perform diagnostics across OSI layers 1–4 using SNMP, ICMP, and TCP tools.

---

## 🧱 Structure

This tool currently supports the following commands:

---

## 🧪 Layer 1: Physical Layer

### `linkstatus [host] [community] [ifIndex]`
- Queries SNMP OID `1.3.6.1.2.1.2.2.1.8.X`
- Returns link state:
  - `1 = up`
  - `2 = down`
  - `3 = testing`
- **Example:**
  ```bash
  netanalyzer linkstatus 192.168.1.1 public 2
  ```

### `interfacespeed [host] [community] [ifIndex]`
- SNMP OID: `1.3.6.1.2.1.2.2.1.5.X`
- Reports interface speed in bits/second (32-bit limit ~4Gbps)
- **Example:**
  ```bash
  netanalyzer interfacespeed 192.168.1.1 public 2
  ```

### `highspeed [host] [community] [ifIndex]`
- SNMP OID: `1.3.6.1.2.1.31.1.1.1.15.X`
- Reports interface speed in Mbps (64-bit support for >4Gbps)
- **Example:**
  ```bash
  netanalyzer highspeed 192.168.1.1 public 2
  ```

---

## 🧪 Layer 2: Data Link Layer

### `mactable [host] [community]`
- Walks OID `1.3.6.1.2.1.17.4.3.1.2`
- Shows MAC addresses mapped to switch ports (bridge forwarding table)
- **Example:**
  ```bash
  netanalyzer mactable 192.168.1.1 public
  ```

### `arptable [host] [community]`
- Walks OID `1.3.6.1.2.1.4.22.1.2`
- Displays IP-to-MAC address mappings (ARP table)
- **Example:**
  ```bash
  netanalyzer arptable 192.168.1.1 public
  ```

### `stpinfo [host] [community]`
- Walks OID `1.3.6.1.2.1.17.2.15`
- Returns STP port states:
  - `1 = Disabled`
  - `2 = Blocking`
  - `3 = Listening`
  - `4 = Learning`
  - `5 = Forwarding`
  - `6 = Broken`
- **Example:**
  ```bash
  netanalyzer stpinfo 192.168.1.1 public
  ```

---

## 🧪 Layer 3: Network Layer

### `ping [host]`
- ICMP ping with 4 echo requests
- Requires root/Admin privileges on some systems
- **Example:**
  ```bash
  netanalyzer ping 8.8.8.8
  ```

### `traceroute [host]`
- ICMP-based path tracing to destination
- Shows intermediate hops and response times
- **Example:**
  ```bash
  netanalyzer traceroute google.com
  ```

### `dnslookup [host]`
- Resolves DNS name to IP addresses (A/AAAA)
- **Example:**
  ```bash
  netanalyzer dnslookup example.com
  ```

### `ipinfo [host]`
- Shows IP and hostname for a given address or domain
- **Example:**
  ```bash
  netanalyzer ipinfo 8.8.8.8
  ```

---

## 🧪 Layer 4: Transport Layer

### `tcpscan [host] [start-port] [end-port]`
- Performs TCP port scan in range
- **Example:**
  ```bash
  netanalyzer tcpscan 192.168.1.1 20 80
  ```

### `tcpbanner [host] [port]`
- Reads TCP banner from service (HTTP, SMTP, FTP...)
- **Example:**
  ```bash
  netanalyzer tcpbanner mail.example.com 25
  ```

### `tcpservice [host] [port]`
- Checks if TCP service is available without banner grab
- **Example:**
  ```bash
  netanalyzer tcpservice 192.168.1.1 443
  ```

---

## 🧰 Usage

```bash
# Build
go build -o netanalyzer.exe ./main.go

# Run
./netanalyzer.exe [command] [args]
```

---

## 📦 Dependencies

- [gosnmp](https://github.com/gosnmp/gosnmp)
- [cobra](https://github.com/spf13/cobra)

---

## 📄 License
