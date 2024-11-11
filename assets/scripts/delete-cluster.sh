#!/bin/bash

CLUSTER_NAME="my-cluster"

# Check if Kind is installed
if ! command -v kind &> /dev/null
then
    echo "Kind is not installed. Please install Kind to proceed."
    exit 1
fi

# Check if the cluster exists
if ! kind get clusters | grep -q "$CLUSTER_NAME"; then
    echo "Cluster '$CLUSTER_NAME' does not exist. Skipping deletion."
    exit 0
fi

# Delete the cluster
echo "Deleting Kind cluster '$CLUSTER_NAME'..."
kind delete cluster --name "$CLUSTER_NAME"

if [ $? -eq 0 ]; then
    echo "Cluster '$CLUSTER_NAME' deleted successfully."
else
    echo "Failed to delete cluster '$CLUSTER_NAME'."
    exit 1
fi
