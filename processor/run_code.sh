if [[ "$1" == "python" ]]; then
    python3 "$2"
elif [[ "$1" == "cpp" ]]; then
    g++ "$2" -o output && ./output
elif [[ "$1" == "c" ]]; then
    gcc "$2" -o output && ./output
fi