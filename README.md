# CG4002: Laser Tag

This repository contains Group B7's code for CG4002 Capstone Project.

## External Communications
1. To build, you will need python3 and golang.
1. Setup `~/.ssh/config` with `xilinx` as the Ultra96
1. Use `make -j 2` to build and push to `xilinx`. You may modify the host name or target environment from arm64 to other supported platforms.
1. Start up python fpga server on `xilinx` with `pushd ~/ecomm/pynq && sudo -E python3 main.py`.
1. To run with mock eval server on `xilinx`, simply run `~/ecomm/ecomm`.
1. To run with eval server on `xilinx`, simply run `~/ecomm/ecomm -evalAddr='10.1.1.1:1234'`.
1. More runtime options can be found with `~/ecomm/ecomm --help`.

## Internal Communications
Please refer to the [internal communications README](iComm/README.md) under the `iComm` folder.

## Hardware Sensor
Please refer to the [hardware sensor README](Hardware/README.md) under the `Hardware` folder.

## Hardware AI
For the software AI development, please refer to the [software AI README](AI/sw_model/README.md) under the `AI/sw_model` folder.

For the FPGA AI implementation, please refer to the [hardware AI README](AI/hw_model/README.md) under the `AI/hw_model` folder.

For the PYNQ deployment, please refer to the [PYNQ README](AI/FPGA_deployment/README.md) under the `AI/FPGA_deployment` folder.

## Visualizer
All assets used to develop the AR visualizer is under the `Assets` folder.