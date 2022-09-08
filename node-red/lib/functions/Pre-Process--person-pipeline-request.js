// name: Pre-Process: person pipeline request
// outputs: 1
// initialize: // Code added here will be run once\n// whenever the node is started.\n
// finalize: // Code added here will be run when the\n// node is being stopped or re-deployed.\n
// info: 
msg.endpoint = msg.pipelineList[msg.pipeline];
var pipelineNumber = flow.get("pipeline_number") || 0;
pipelineNumber = pipelineNumber + 1;
flow.set("pipeline_number", pipelineNumber);
msg.payload = {
    "source": {
        "uri": msg.video,
        "type": "uri"
    },
    "destination": {
        "metadata": msg.metadataDest
    },
    "parameters": {
        "zmq-config": {
            "imagezmqSocket": msg.zmqSocket,
            "mqttTopic": msg.pipeline,
            "pipelineNumber": pipelineNumber,
            "pipelineName": msg.pipeline
        }
    }
};
return msg;