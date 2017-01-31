#!/bin/sh
# prints file in $1 as go byte slice literal, nicely formatted
exec python -c "print(', '.join(['0x%x'%ord(x) for x in open('"$1"', 'rb').read()]))" | par -w78
