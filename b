go mod tidy
go fmt ./...
go build
rm /home/gokul/Workspace/bin/note 2> /dev/null
mv note /home/gokul/Workspace/bin