#!/usr/bin/env bash
curl -X POST http://localhost:17000 -d "reset
white"
curl -X POST http://localhost:17000 -d "figure 0.25 0.25
update"
sleep 1
curl -X POST http://localhost:17000 -d "move 0.25 0.25
update"
sleep 1
curl -X POST http://localhost:17000 -d "move -0.5 -0.5
update"



