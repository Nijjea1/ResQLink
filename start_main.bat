@echo off
echo Starting MeshComm Main Node...
cd cmd/meshcomm
go run . -port 6000 -nick EMS_Main -api-port 3001 -same_string meshcomm