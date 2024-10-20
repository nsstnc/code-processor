#!/bin/bash
if [[ "$1" == "python" || "$1" == "python3" ]]; then
    python3 "$2"
elif [[ "$1" == "cpp" ]]; then
    g++ "$2" -o output && ./output
elif [[ "$1" == "c" ]]; then
    gcc "$2" -o output && ./output
fi