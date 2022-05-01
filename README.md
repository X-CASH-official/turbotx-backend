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
