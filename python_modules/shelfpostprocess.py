import imagezmq
import numpy as np

class PostProcess:
    
    def __init__(self, imagezmqSocket, mqttTopic, pipelineName):
        self.pipelineNumber = np.random.randint(0, 255)
        self.data = {"imagezmqSocket": imagezmqSocket, "mqttTopic": mqttTopic, "pipelineName": pipelineName}
        self.topic = self.data["pipelineName"]+"/"+str(self.pipelineNumber)
        self.sender = imagezmq.ImageSender(connect_to = self.data["imagezmqSocket"])
        return
    
    def shelfPostProcess(self,frame):
        pipeline_data = {"Name": self.data["pipelineName"], "Number":self.pipelineNumber}

        roi_list = []
        with frame.data() as mat:
            for roi in frame.regions():
                roi_dict = {"label": roi.label_id(), "confidence": roi.confidence(),
                "rect": roi.rect()}
                roi_list.append(roi_dict)
            pipeline_data["data"] =  roi_list
            self.sender.send_image(self.topic,  mat[:,:,:3], pipeline_data)
        return True
