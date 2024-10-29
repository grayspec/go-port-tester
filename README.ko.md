
# go-port-tester

`go-port-tester`는 여러 서버의 지정된 포트를 HTTP 요청을 통해 확인하는 간단한 Go 프로젝트입니다.   
네트워크 관리자나 여러 서버와 포트에서 실행 중인 서비스의 접근 가능성을 확인해야 하는 개발자에게 유용한 도구입니다.

[English](README.md)

이 프로젝트는 MIT 라이선스에 따라 라이선스가 부여되었습니다.   
자세한 내용은 [LICENSE](LICENSE) 파일을 참조하세요.

---

### 주요 기능
- 🌐 **멀티 서버, 멀티 포트 검사**: 여러 IP와 포트를 동시에 테스트
- ⚡ **빠르고 효율적인 동시 처리**: 대규모 네트워크에 최적화된 동시성 조절 기능
- 📊 **쉬운 결과 분석**: 결과를 깔끔한 CSV 형식으로 저장하여 분석 용이
- 🛠️ **크로스 플랫폼 호환성**: Windows, macOS, Linux 지원

---

### 목차
- [필수 조건](#필수-조건)
- [설치](#설치)
    - [Go 설치](#go-설치)
    - [Make 설치](#make-설치)
- [빌드 방법](#빌드-방법)
- [애플리케이션 실행](#애플리케이션-실행)
- [사용 예시](#사용-예시)

---

### 필수 조건

- **Go** (버전 1.20 이상)
- **GNU Make**

### 설치

#### Go 설치
- **Windows**: [Go 공식 웹사이트](https://golang.org/dl/)에서 다운로드 및 설치
- **macOS**: Homebrew를 사용하여 설치:
  ```bash
  brew install go
  ```
- **Linux**: [Go 공식 웹사이트](https://golang.org/doc/install)의 설명을 따르거나 `apt` 등의 패키지 매니저를 사용하여 설치:
  ```bash
  sudo apt update
  sudo apt install -y golang
  ```

#### Make 설치
Linux와 macOS에는 Make가 기본적으로 설치되어 있으며, Windows에서는 다음 방법으로 설치할 수 있습니다:

- **Windows**: Chocolatey를 통해 설치:
  ```powershell
  choco install make
  ```
  또는 Scoop을 통해 설치:
  ```powershell
  scoop install make
  ```

- **macOS**: Homebrew로 설치:
  ```bash
  brew install make
  ```

- **Linux**: Make는 보통 기본 설치되어 있으나, 설치가 필요하면 다음 명령어를 실행:
  ```bash
  sudo apt install -y make
  ```

---

### 빌드 방법

크로스 플랫폼 빌드를 위해 제공된 `Makefile`을 사용하여 빌드합니다.

1. **레포지토리 클론**:
   ```bash
   git clone https://github.com/your-username/go-port-tester.git
   cd go-port-tester
   ```

2. **모든 플랫폼 빌드**:
   ```bash
   make all
   ```

3. **플랫폼별 빌드**:
    - **Windows**:
      ```bash
      make build-windows
      ```
    - **macOS**:
      ```bash
      make build-macos
      ```
    - **Linux**:
      ```bash
      make build-linux
      ```

빌드된 바이너리는 `build` 디렉토리에 OS별로 저장됩니다.

---

### 애플리케이션 실행

#### 서버 (테스트 대상)
서버는 지정된 포트에서 간단한 HTTP 서비스를 호스팅하여 테스트를 지원합니다.

```bash
./build/<platform>/server --config server.csv
```

#### 클라이언트 (포트 테스터)
클라이언트는 지정된 IP와 포트에 요청을 보내 포트 상태를 확인합니다.

```bash
./build/<platform>/client --servers servers.csv --output result.csv --timeout 2 --concurrency 5
```

---

### 사용 예시

**서버 목록 CSV 예시** (`servers.csv`):
```csv
IP,Port
127.0.0.1,8080
192.168.1.10,8081
```
