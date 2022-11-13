# Internal Comms

This folder contains the code for internal comms.

## Setup

Run `pip install -r requirements.txt`, which installs
the corresponding version of bluepy used for this project.

## Run
Before running, we first have to open a terminal and ssh into the xilinx server:

`ssh -L 8081:localhost:8081 xilinx`

Then, run `main.py` located within the iComm folder via 

`python main.py` within another terminal.

