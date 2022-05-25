# turbotx-backend

Note you do not need the blockchain to use this. It uses all delegates commands

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

Build the program  
`go build .`

Run the program  
`systemctl start turbotx-backend`

Make sure to install the frontend so you can view the turbo tx id https://github.com/X-CASH-official/turbotx-frontend
