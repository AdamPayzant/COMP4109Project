# This is a workaround for what is likely a gap in knowledge
# Go packages are tricky and I'm still learning
cp ./smvsclient $(go env GOPATH)/src/smvsclient
cp ./smvshost $(go env GOPATH)/src/smvshost
cp ./smvsserver $(go env GOPATH)/src/smvsserver

go get smvsclient
go install smvsclient
go get smvshost
go install smvshost
go get smvsserver
go install smvsserver