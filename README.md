# Example Golang ssh client

1. Define alias in `~/.ssh/config`

```
Host example
User ubuntu
Hostname 10.11.11.11
IdentityFile ~/.ssh/id_rsa
ForwardAgent yes
```

2. Run command `whoami` on remote host with alias `example`.

```golang
	conn := createConnection("example")
	runCommand("whoami", conn)
  defer conn.Close()
```


# Install

```
go mod download
```
