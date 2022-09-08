// name: Check EVAM Status
// outputs: 1
// initialize: // Code added here will be run once\n// whenever the node is started.\n
// finalize: // Code added here will be run when the\n// node is being stopped or re-deployed.\n
// info: 
flow.set("success", false);
msg.resetEndpoint = "http://" + env.get("HOST_IP") + ":9000/hooks/reset-evam";

const checkStatus = async function () {
    const timer = new Promise((resolve, reject) => {
        let count = 0;
        const check = setInterval (() => {
            count = count + 1;
            if (flow.get("success") === true) {
                resolve(true);
            }
            if (count === 6) {
                reject(false);
                clearInterval(check);
            }
        }, 1000);
    });

    try {
        await timer;
    } catch (ex) {
        return msg;
    }

};

return checkStatus(msg);