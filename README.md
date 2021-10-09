## Prerequisites
- Go 1.16+

## Tests and running
This project contains a few tests that can show the potential of this simulation program

### Simulate a single match
- The test generates 2 teams with randoms skills and pits them against each other
- Run "cd gamesim"
- "go test -v -run TestSim_Start"

### Simulate a tournament of 3 stages (8 competing teams)
- "cd gamesim"
- "go test -v -run TestTourneySim_Advance"