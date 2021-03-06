# EVM-BABBLE DEMOS
Deploying **evm-babble** and **babble** side by side

**The following scripts were only tested on Ubuntu 16.04**

## Dependencies

### Node.js

As part of the demos, we use javascript to interact with Smart Contracts. This  
allows us to reuse some popular libraries that were developed to work with Ethereum.  
Node.js allows us to run javascript in the console.

```bash
# install node version manager
curl -o- https://raw.githubusercontent.com/creationix/nvm/v0.33.5/install.sh | bash
# use nvm to intall stable version of node
nvm install node stable
```

## Demo

The **demo** shows how one might setup a crowd funding campaign on a trusted  
network of Babble nodes. We show how to use a Smart Contract on evm-babble to  
distribute and automate the logic that will securely receive contributions for a  
crowd funding campaign and transfer the funds to the beneficiary if and only if  
the funding goal is met.

## Docker

Launch a set of Docker containers to setup a local evm-babble testnet. 

Obviously this requires [Docker](https://docker.com). Follow the link to find installation instructions.

```bash
[...]/evm-babble/demo/demo-docker$ make  # create testnet
[...]/evm-babble/demo/demo-docker$ make demo # run through a demo scenario
[...]/evm-babble/demo/demo-docker$ make stop # stop and remove all resources
```

## AWS

Setup a testnet in AWS using the [Terraform](https://www.terraform.io/) utility.

This is a more complicated scenario. Please contact us if you need help.  
You need an AWS account and an authentication key. 

There are two main parts to this procedure:

    1. Use the AWS console to create a base image.
    2. Use terraform scripts to launch a certain number of nodes in a testnet  
       and start babble and evm-babble on them. 

1. Create an AWS Image with babble and evm-babble binaries

This step cannot really be automated

We could automate the deployment of babble and evm-babble binaries to each  
instance but it would be very slow since these files are large. So the idea is to  
manually create an image (snapshot) of a machine configured with babble and evm-babble  
preinstalled. We can then use Terraform to create other machines based on that image.  
This process makes it a lot faster to bootstrap new testnets but requires a manual  
step everytime there is a new build for babble or evm-babble

Our approach to this is to keep an Ubuntu 16.04 instance in AWS that will serve as  
a template. When we want to test a new build for babble or evm-babble, we copy the  
binaries into that instance using **scp** and we take a snapshot of the the template  
instance. We then copy the resulting snapshot's ID into our Terraform scripts (example.tf).

2. Use scripts to deploy the testnet and execute demos  

```bash
[...]/evm-babble/demo/demo-terraform$ make "nodes=12"
[...]/evm-babble/demo/demo-terraform$ make demo #run a demo scenario
# ssh into a node directly. From there you can look at logs or system resources
[...]/evm-babble/demo/demo-terraform$ ssh -i babble.pem ubuntu@[public ip] 
[...]/evm-babble/demo/demo-terraform$ make destroy #destroy resources
```

