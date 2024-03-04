#!/bin/bash

GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${GREEN}Generating keys...${NC}"

if [ ! -d "./certs" ]; then mkdir ./certs; fi 
	    if [ -f "./certs/app.rsa.key" ] && [ -f "./certs/app.rsa.pub" ]; then 
        echo -e "${BLUE}Keys already exist. Do you want to override them?${NC} (y/n)";
        read REPLY 
        if [ "$REPLY" != "y" ]; then 
            echo "Keys not overridden."; 
            exit 0; 
        else 
            echo -e "${RED}Overriding keys...${NC}"; 
            rm -f ./certs/app.rsa.key; 
            rm -f ./certs/app.rsa.pub; 
            openssl genpkey -algorithm RSA -out ./certs/app.rsa.key -pkeyopt rsa_keygen_bits:2048
	        openssl rsa -pubout -in ./certs/app.rsa.key -out ./certs/app.rsa.pub
            echo -e "${GREEN}Keys overridden.${NC}";
        fi; 
    fi;