# sra-low-code
#### Low code implementation of *Intel Smart Retail Analytics*. 
It takes video input from one or more sources. Node-Red runtime is used for Low-Code execution and [**Intel Edge Video Analytics Microservice**](https://github.com/intel/edge-video-analytics-microservice) for producing inferences.

## Setup Steps

- Check the exposed ports section in `.env` file and change if they are not available to bind.
- Make sure `http_proxy`, `https_proxy` and `no_proxy` environment variables are set properly on your host.
- Run the following commands:
```bash

# Set other required environment variable
export HOST_IP=$(hostname -I | cut -d' ' -f1)

# Clone this repo and change dir
git clone https://github.com/intel-sandbox/sra-low-code.git
cd sra-low-code

# Build and spin-up containers
docker-compose build
docker-compose up -d
```

- Copy video input files to `./video_src` directory. It will be available to **Intel Edge Video Analytics Microservice** inside `/data` dir.

### Running and visualizing Smart Retail Analytics 

#### From Node-Red UI
 - Head to `http://your_host:1880` to access node-red UI.
 - Trigger the Smart Retail Analytics flow by triggering the inject node.
#### From Command line
 ```bash
 curl -X POST http://<your_host>:1880/inject/5daa1599299a7569
 ```
 
 This will start the flow and pipelines will start executing.
 - Visualize streams at `http://your_host:5000`.
 - Visualize pipeline stats in Grafana at `http://your_host:3000`.



