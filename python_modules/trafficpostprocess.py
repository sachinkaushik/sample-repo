import cv2
import numpy as np
import imagezmq
import random
from server.common.utils import logging

logger = logging.get_logger('shopper_count_duration', is_static=True)

class PostProcess:
    
    def __init__(self, imagezmqSocket="tcp://server:5555", mqttTopic="shopper-count-duration", pipelineName="shopper-count-duration"):
        logger.info("starting")
        self.data = {"imagezmqSocket": imagezmqSocket, "mqttTopic": mqttTopic, "pipelineName": pipelineName}
        self.labels = ['person']
        self.last_correct_count = [0] * len(self.labels)
        self.total_count = [0] * len(self.labels)
        self.current_count = [0] * len(self.labels)
        self.changed_count = [False] * len(self.labels)
        self.candidate_count = [0] * len(self.labels)
        self.candidate_confidence = [0] * len(self.labels)
        self.mog = cv2.createBackgroundSubtractorMOG2()
        self.CONF_CANDIDATE_CONFIDENCE = 3
        self.thresh = 0.45
        self.flag = 0
        self.topic = self.data["mqttTopic"]+"/"+str(random.randrange(2001,3000))
        self.sender = imagezmq.ImageSender(connect_to = self.data["imagezmqSocket"])
        return

    def trafficPostProcess(self, frame):
        if self.flag == 0:
            width = frame.video_info().width
            height = frame.video_info().height
            self.accumulated_frame = np.zeros((int(height), int(width)), np.uint8)
            self.flag = 1

        with frame.data() as mat:
            for roi in frame.regions():
                if roi.confidence() > self.thresh:
                    label = roi.label_id() - 1
                    self.current_count[label] += 1
            
            for i in range(len(self.labels)):
                if self.candidate_count[i] == self.current_count[i]:
                    self.candidate_confidence[i] += 1
                else:
                    self.candidate_confidence[i] = 0
                    self.candidate_count[i] = self.current_count[i]

                if self.candidate_confidence[i] == self.CONF_CANDIDATE_CONFIDENCE:

                    self.candidate_confidence[i] = 0
                    self.changed_count[i] = True
                else:
                    continue
                if self.current_count[i] > self.last_correct_count[i]:
                    self.total_count[i] += self.current_count[i] - self.last_correct_count[i]
                self.last_correct_count[i] = self.current_count[i]

            ##### Heatmap generation
            # Convert to grayscale
            gray = cv2.cvtColor(mat, cv2.COLOR_BGR2GRAY)
            # Remove the background
            fgbgmask = self.mog.apply(gray)
            # Threshold the image
            thresh = 2
            max_value = 2
            threshold_frame = cv2.threshold(fgbgmask, thresh, max_value, cv2.THRESH_BINARY)[1]
            # Add threshold image to the accumulated image
            self.accumulated_frame = cv2.add(threshold_frame, self.accumulated_frame)
            colormap_frame = cv2.applyColorMap(self.accumulated_frame, cv2.COLORMAP_HOT)
            self.frame = cv2.addWeighted(mat[:,:,:3], 0.6, colormap_frame, 0.4, 0)
            self.frame = cv2.putText(self.frame, 'use-case='+self.data["pipelineName"], (10,30), cv2.FONT_HERSHEY_SIMPLEX, 1, (255,0,0), thickness=3)
            for idx, label in enumerate(self.labels):
                cv2.putText(self.frame, "person current count:"+str(self.current_count[idx]), (10, 60),cv2.FONT_HERSHEY_SIMPLEX, 1, (255, 255, 255), 1)
                cv2.putText(self.frame, "person total count:"+str(self.total_count[idx]), (10, 80),cv2.FONT_HERSHEY_SIMPLEX, 1, (255, 255, 255), 1)
            self.sender.send_image(self.topic, cv2.resize(self.frame[:,:,:3], (480,360)))
        self.current_count = [0] * len(self.labels)
        return True
