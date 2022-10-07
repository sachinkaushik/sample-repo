import cv2
from collections import namedtuple
import imagezmq
import numpy as np

class PostProcess:
    def __init__(self, imagezmqSocket, mqttTopic, pipelineName):
        self.pipelineNumber = np.random.randint(0, 65535)
        self.data = {"imagezmqSocket": imagezmqSocket, "mqttTopic": mqttTopic, "pipelineName": pipelineName}
        self.Point = namedtuple("Point", "x,y")
        self.topic = self.data["pipelineName"]+"/"+str(self.pipelineNumber)
        self.sender = imagezmq.ImageSender(connect_to = self.data["imagezmqSocket"])
        return
    
    def facePostProcess(self, frame):
        return True
    
    def headPostProcess(self, frame):
        add_regions = []
        for roi in frame.regions():
            angle_p_fc = 0
            angle_y_fc = 0
            bbox = roi.rect()
            for tensor in roi.tensors():    
                if (tensor.layer_name() == "angle_p_fc"):  
                    angle_p_fc = angle_p_fc + tensor.data()
                elif (tensor.layer_name() == "angle_y_fc"):
                    angle_y_fc = angle_y_fc + tensor.data()

            if (angle_y_fc > -22.5) & (angle_y_fc < 22.5) & (angle_p_fc > -22.5) & (angle_p_fc < 22.5):
                add_regions.append([bbox.x ,bbox.y, bbox.w, bbox.h, "looking", 1.0])
            else:
                add_regions.append([bbox.x ,bbox.y, bbox.w, bbox.h, "looking", 0.0])

        for region in add_regions:
            frame.add_region(region[0], region[1], region[2], region[3], region[4], region[5])

        return True

    def emotionPostProcess(self, frame):
        self.sentiment = []
        self.looking = []
        self.frame_centroids = []
        self.bbox = []
        with frame.data() as mat:

            for roi in frame.regions():
                if roi.label() == "looking":
                    continue
                pt = roi.rect()
                x = pt.x + int(pt.w / 2)
                y = pt.y + int(pt.h / 2)
                point = self.Point(x, y)
                self.frame_centroids.append(point)
                self.bbox.append([pt.x, pt.y, pt.w, pt.h])

            for roi in frame.regions():
                bbox = roi.rect()
                if roi.label() != "looking":
                    continue
                x = bbox.x + int(bbox.w / 2)
                y = bbox.y + int(bbox.h / 2)
                
                for tensor in roi.tensors():
                    if (tensor.layer_name() == "prob_emotion"):
                        data = tensor.data()
                        emotions = ["neutral", "happy", "sad", "surprise", "anger"]
                        label_out = emotions[np.argmax(data)]
                        if roi.confidence() > 0.9:
                            self.looking.append(True)
                            self.sentiment.append(label_out)
                        else:
                            self.looking.append(False)
                            self.sentiment.append(-1)
            
            pipeline_data = {"frame_centroids": self.frame_centroids, "looking":self.looking, "sentiment":self.sentiment, "bbox":self.bbox}
            pipeline_data["Name"] =  self.data["pipelineName"]
            pipeline_data["Number"] = self.pipelineNumber
            
            self.sender.send_image(self.topic, mat[:,:,:3], pipeline_data)
        return True

    def __del__(self):
        self.sender.close()
