# Project
- go
- worker

```sh
#  init mod
$ go mod init app/axel/worker

$ go install app/axel/worker

$ go env -w GOBIN=/somewhere/else/bin
$ go build -o worker 


$ go install example/user/hello
$ go install .
$ go install

$ go install app/axel/worker

$ go mod tidy

```

# set path go 

echo "export GOPATH=/Users/axel/Documents/go" >> .zshrc


echo "export GOPATH=/Users/axel/.asdf/installs/golang/1.15.8/packages" >> .zshrc




# UI worker 
```sh

./workwebui -redis="redis://localhost:6379" -ns="work_namespace" -listen="5040"
./workwebui -redis="redis:6379" -ns="work_namespace" -listen=":5040"
./workwebui -redis="local:redis:6379" -ns="work_namespace" -listen=":5040"
./workwebui -redis="redis:6379" -ns="work_namespace" -listen="0.0.0.0:5040"



```
