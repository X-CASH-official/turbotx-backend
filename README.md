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
`git clone https://github.com/X-CASH-official/turbotx-backend.git && cd turbotx-backend`

copy the systemd file  
`cp -a turbotx-backend.service /lib/systemd/system/ && sudo systemctl daemon-reload`

Install redis  
`sudo add-apt-repository ppa:redislabs/redis -y && sudo apt-get update && sudo apt install redis`

Setup local EMPTY xcash-rpc-wallet connected to the 5 seed nodes on ports 18285-18289
`screen -dmS network_data_node_1 /root/xcash-official/xcash-core/build/release/bin/xcash-wallet-rpc --wallet-file delegate-wallet-1 --password password --rpc-bind-port 18285 --confirm-external-bind --daemon-address dpops-test-internal-1.xcash.foundation:18281 --disable-rpc-login --trusted-daemon && screen -dmS network_data_node_2 /root/xcash-official/xcash-core/build/release/bin/xcash-wallet-rpc --wallet-file delegate-wallet-2 --password password --rpc-bind-port 18286 --confirm-external-bind --daemon-address dpops-test-internal-2.xcash.foundation:18281 --disable-rpc-login --trusted-daemon && screen -dmS network_data_node_3 /root/xcash-official/xcash-core/build/release/bin/xcash-wallet-rpc --wallet-file delegate-wallet-3 --password password --rpc-bind-port 18287 --confirm-external-bind --daemon-address dpops-test-internal-3.xcash.foundation:18281 --disable-rpc-login --trusted-daemon && screen -dmS network_data_node_4 /root/xcash-official/xcash-core/build/release/bin/xcash-wallet-rpc --wallet-file delegate-wallet-4 --password password --rpc-bind-port 18288 --confirm-external-bind --daemon-address dpops-test-internal-4.xcash.foundation:18281 --disable-rpc-login --trusted-daemon && screen -dmS network_data_node_5 /root/xcash-official/xcash-core/build/release/bin/xcash-wallet-rpc --wallet-file delegate-wallet-5 --password password --rpc-bind-port 18289 --confirm-external-bind --daemon-address dpops-test-internal-5.xcash.foundation:18281 --disable-rpc-login --trusted-daemon`

Build the program  
`go build .`

Run the program  
`systemctl start turbotx-backend`

