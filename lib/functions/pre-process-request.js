// name: Pre-Process: person pipeline request
// outputs: 1
// initialize: // Code added here will be run once\n// whenever the node is started.\n
// finalize: // Code added here will be run when the\n// node is being stopped or re-deployed.\n
// info: 
const parameters = {
    "parameters" : {
        ...msg.modelInstanceId,
        ...msg.config
    }
}
msg.payload = { 
    ...msg.source,
    ...msg.destination,
    ...parameters
};

msg.endpoint = msg.pipelineList[msg.pipeline];
return msg;