# eval $(docker-machine env) && docker run -it --sig-proxy=true --rm  jogger


MAX_MEMORY_USAGE=128m
MAX_SCRIPT_TIME=1
MAX_STORAGE_SIZE=1G
MAX_OPEN_FILES=10


# Build the container
docker build -t jogger .

# Access the Python REPL
docker run -it --sig-proxy=true --rm  jogger

# Run a Python script in the container
docker run -it --rm --name testjoggerscript -v "$PWD":/usr/src/myapp -w /usr/src/myapp python:2.7 python 1_sample.py

# Run Python script with libraries
docker run -it --rm --name testjoggerscript -v "$PWD":/usr/src/myapp -w /usr/src/myapp -e PYTHONPATH="./libs/" python:2.7 python tests/2_sample_with_libs.py

# Networking disabled
docker run    --network="none" -i --rm --name testjoggerscript -v "$PWD":/usr/src/myapp -w /usr/src/myapp -e PYTHONPATH="./libs/" python:2.7    timeout 2 python tests/3_networking_disabled.py 2>&1
# returns urllib2.URLError: <urlopen error [Errno -2] Name or service not known>


# Can't create files. Storage size limit on container
# --storage-opt dm.basesize=$MAX_STORAGE_SIZE
gtimeout $MAX_SCRIPT_TIME   docker run -m $MAX_MEMORY_USAGE    --ulimit nofile=$MAX_OPEN_FILES:$MAX_OPEN_FILES  --network="none" -i --rm --name testjoggerscript -v "$PWD":/usr/src/myapp:ro -w /usr/src/myapp -e PYTHONPATH="./libs/" python:2.7   python tests/4_storage_limits.py 2>&1
# returns IOError: [Errno 30] Read-only file system: 'file1'



