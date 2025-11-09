@echo off
echo Starting MeshComm Client Node...
cd cmd/meshcomm
go run . -port 6001 -nick Client_%random% -api-port 3002 -same_string meshcomm