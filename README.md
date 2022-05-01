# turbotx-backend

install go  
download the latest go version from https://go.dev/doc/install
 
untar it  
`tar -xf go* && rm go*.tar.gz`
 
edit the path  
`sudo nano ~/.profile`
 
add this line to the end of the file  
`export PATH=$PATH:/usr/local/go/bin`
 
save the file  
`source ~/.profile`
 
verify the install  
`go version`

Install  
`go get -u github.com/X-CASH-official/turbotx-backend`

copy the systemd file  
`cp -a turbotx-backend.service /lib/systemd/system/ && sudo systemctl daemon-reload`

Install redis  
`sudo add-apt-repository ppa:redislabs/redis -y && sudo apt-get update && sudo apt install redis`

Setup local EMPTY xcash-rpc-wallet (you do not need a local xcashd)  
`/root/xcash-official/xcash-core/build/release/bin/xcash-wallet-rpc --wallet-file wallet1 --password password --rpc-bind-port 18285 --confirm-external-bind --daemon-port 18281 --disable-rpc-login --trusted-daemon`

