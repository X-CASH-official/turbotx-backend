# turbotx-backend

Note you do not need the blockchain to use this. It uses all delegates commands

You can download the latest release and run it as well as the [frontend](https://github.com/X-CASH-official/turbotx-frontend) locally if you want, or use the [official website](http://162.55.235.87/)

# How to build from source

install go  
download the latest go version from https://go.dev/doc/install
 
untar it  
`tar -xf go* && rm go*.tar.gz && mv go /usr/local/`
 
edit the path  
`sudo nano ~/.profile`
 
add this line to the end of the file  
`export PATH=$PATH:/usr/local/go/bin`
 
save the file  
`source ~/.profile`
 
verify the install  
`go version`

Install redis  
`sudo add-apt-repository ppa:redislabs/redis && sudo apt-get update && sudo apt-get install redis && sudo systemctl enable --now redis-server`

Install  
`git clone https://github.com/X-CASH-official/turbotx-backend.git && cd turbotx-backend`

copy the systemd file  
`cp -a turbotx-backend.service /lib/systemd/system/ && sudo systemctl daemon-reload`

Install redis  
`sudo add-apt-repository ppa:redislabs/redis -y && sudo apt-get update && sudo apt install redis`

Build the program  
`make clean ; make release`

Run the program  
`systemctl start turbotx-backend`

Make sure to install the [frontend](https://github.com/X-CASH-official/turbotx-frontend) so you can view the turbo tx id
