#!/bin/bash

# Load environment variables from .env file for start project without docker
if [ -f .env ]; then
    # Read .env file line by line and export variables safely
    while IFS= read -r line; do
        # Skip empty lines and comments
        if [[ -z "$line" || "$line" =~ ^[[:space:]]*# ]]; then
            continue
        fi
        
        # Remove leading/trailing whitespace
        line=$(echo "$line" | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')
        
        # Skip empty lines after trimming
        if [[ -z "$line" ]]; then
            continue
        fi
        
        # Export the variable
        export "$line"
    done < .env
    
    echo "Environment variables loaded from .env"
    echo "Imported variables:"
    echo "=================="
    
    # Display imported variables
    while IFS= read -r line; do
        # Skip empty lines and comments
        if [[ -z "$line" || "$line" =~ ^[[:space:]]*# ]]; then
            continue
        fi
        
        # Remove leading/trailing whitespace
        line=$(echo "$line" | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')
        
        # Skip empty lines after trimming
        if [[ -z "$line" ]]; then
            continue
        fi
        
        # Extract variable name and value
        var_name=$(echo "$line" | cut -d'=' -f1)
        var_value=$(echo "$line" | cut -d'=' -f2-)
        
        # Display variable (mask sensitive values)
        # Debug: echo "Checking: $var_name"
        if [[ "$var_name" =~ (PASSWORD|SECRET|KEY|TOKEN|API_KEY) ]]; then
            # Show first 3 characters and mask the rest
            if [[ ${#var_value} -gt 6 ]]; then
                masked_value="${var_value:0:3}***${var_value: -3}"
            else
                masked_value="***"
            fi
            echo "$var_name=$masked_value"
        else
            echo "$var_name=$var_value"
        fi
    done < .env
    
    echo "=================="
else
    echo "Warning: .env file not found. Using default values."
fi
