# IPSEC Checker

Simple service which ping IPSec IPs and send info in to Telegram via [Horn](https://github.com/requilence/integram) if failed.
Also restart IPSec service if ping failed.

### Installing for develop purpose
1) Clone project
    ```bash
    git clone https://github.com/Rishats/ipsec_checker.git
    ```
2) Change folder
    ```bash
    cd ipsec_checker
    ```
3) Create .env file from .env.example
    ```bash
     cp .env.example .env
    ```

   4) Configure your .env
       ```
         APP_ENV=production-or-other
         INTEGRAM_WEBHOOK_URI=https://integram.org/webhook/cdgdw68sIpR
         SENTRY_DSN=your-dsn
         IP_LIST=127.0.0.1,8.8.8.8
         PROBE=3
       ```

### Develop use cases

#### Via go native:

Download dependency
```bash
go mod download
```

Build for linux
```bash
env GOOS=linux GOARCH=amd64 go build main.go
```

#### Via docker:
```bash
 docker build --target=build-env -t ipsec_checker .
 docker run -d --name "ipsec_checker" ipsec_checker
```


### Usage

#### Via Ansible playbooks
https://github.com/Rishats/ansible-ipsec_checker

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/Rishats/ywpti/tags). 

## Authors

* **Rishat Sultanov** - [Rishats](https://github.com/Rishats)

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
