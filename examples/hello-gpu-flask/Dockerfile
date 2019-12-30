FROM nvidia/cuda:10.1-base

# Install Python & pip
RUN apt-get update && apt-get install -y --no-install-recommends \
  python3 python3-pip \
  && \
  apt-get clean && \
  apt-get autoremove && \
  rm -rf /var/lib/apt/lists/*
RUN pip3 install setuptools

# Set an arbitrary /app workdir.
WORKDIR /app

# Install production dependencies.
# We copy the requirements.txt separately so that we don't reinstall the
# dependencies everytime we modify the rest of the source code.
COPY requirements.txt requirements.txt
RUN pip3 install -r requirements.txt

# Copy local code to the container image.
COPY *.py ./

# Launch the app on port 80.
ENV PORT 80

# Run the app using gunicorn.
CMD exec gunicorn --bind :$PORT --workers 1 --threads 8 main:app