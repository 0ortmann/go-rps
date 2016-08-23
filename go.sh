#!/bin/sh
curl -XPOST localhost:5000/create -d '{"name": "game1"}'
echo ""
curl -XPOST localhost:5000/play -d '{"game": "game1", "player": "leon", "action": "paper"}' 
echo ""
curl -XPOST localhost:5000/play -d '{"game": "game1", "player": "felix", "action": "rock"}'
echo ""
curl -XPOST localhost:5000/play -d '{"game": "game1", "player": "p1", "action": "scissor"}' 
echo ""
curl -XPOST localhost:5000/play -d '{"game": "game1", "player": "p2", "action": "rock"}' 
echo ""
curl -XPOST localhost:5000/play -d '{"game": "game1", "player": "p3", "action": "rock"}' 
echo ""
curl -XPOST localhost:5000/eval -d '{"game": "game1"}'
