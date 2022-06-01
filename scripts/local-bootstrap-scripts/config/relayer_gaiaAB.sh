# Relayer home directory
RELAYER_DIR="${HOME}/archway/local/relayer_gaiaAB"

# Chain cluster configs
CHAIN1_CONFIG="config/gaiaA.sh"
CHAIN2_CONFIG="config/gaiaB.sh"

# Relayer default IBC timeout
TIMEOUT="20s"

# Paths (channels)
## This is a BASH array of arrays.
## Each array defines a single relayer path with the following values (by index):
##   0. "1->2" / "2->1". Defines a direction (chain 1 to chain 2 or vice versa);
##   1. Source port ID;
##   2. Destination port ID;
##   3. Order type ("ordered" / "unordered");
##   4. Version;
declare -a chat_pathAB=("1->2" "chat" "chat" "unordered" "chat-1")
# declare -a chat_pathBA=("2->1" "chat" "chat" "unordered" "chat-1")
# declare -a PATHS=("chat_pathAB" "chat_pathBA")
declare -a PATHS=("chat_pathAB")
