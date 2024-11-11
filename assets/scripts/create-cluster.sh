#!/bin/bash

CLUSTER_NAME="my-cluster"

# Check if Kind is installed
if ! command -v kind &> /dev/null
then
    echo "Kind is not installed. Please install Kind to proceed."
    exit 1
fi

# Check if the cluster already exists
if kind get clusters | grep -q "$CLUSTER_NAME"; then
    echo "Cluster '$CLUSTER_NAME' already exists. Skipping creation."
    exit 0
fi

# Create the cluster
echo "Creating Kind cluster '$CLUSTER_NAME'..."
kind create cluster --name "$CLUSTER_NAME"

if [ $? -eq 0 ]; then
    echo "Cluster '$CLUSTER_NAME' created successfully."
else
    echo "Failed to create cluster '$CLUSTER_NAME'."
    exit 1
fi
