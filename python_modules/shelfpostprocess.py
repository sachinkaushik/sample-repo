import cv2
import imagezmq
import random
from server.common.utils import logging
LABELS = [
                "person",
                "bicycle",
                "car",
                "motorbike",
                "aeroplane",
                "bus",
                "train",
                "truck",
                "boat",
                "traffic light",
                "fire hydrant",
                "stop sign",
                "parking meter",
                "bench",
                "bird",
                "cat",
                "dog",
                "horse",
                "sheep",
                "cow",
                "elephant",
                "bear",
                "zebra",
                "giraffe",
                "backpack",
                "umbrella",
                "handbag",
                "tie",
                "suitcase",
                "frisbee",
                "skis",
                "snowboard",
                "sports ball",
                "kite",
                "baseball bat",
                "baseball glove",
                "skateboard",
                "surfboard",
                "tennis racket",
                "bottle",
                "wine glass",
                "cup",
                "fork",
                "knife",
                "spoon",
                "bowl",
                "banana",
                "apple",
                "sandwich",
                "orange",
                "broccoli",
                "carrot",
                "hot dog",
                "pizza",
                "donut",
                "cake",
                "chair",
                "sofa",
                "pottedplant",
                "bed",
                "diningtable",
                "toilet",
                "tvmonitor",
                "laptop",
                "mouse",
                "remote",
                "keyboard",
                "cell phone",
                "microwave",
                "oven",
                "toaster",
                "sink",
                "refrigerator",
                "book",
                "clock",
                "vase",
                "scissors",
                "teddy bear",
                "hair drier",
                "toothbrush"
    ]
logger = logging.get_logger('shelf_object_count', is_static=True)
class PostProcess:
    
    def __init__(self,  imagezmqSocket, mqttTopic, pipelineName):
        logger.info("starting")
        self.data = {"imagezmqSocket": imagezmqSocket, "mqttTopic": mqttTopic, "pipelineName": pipelineName}
        self.thresh = 0.14
        self.labels = LABELS
        self.last_correct_count = [0] * len(self.labels)
        self.total_count = [0] * len(self.labels)
        self.current_count = [0] * len(self.labels)
        self.changed_count = [False] * len(self.labels)
        self.candidate_count = [0] * len(self.labels)
        self.candidate_confidence = [0] * len(self.labels)
        self.CONF_CANDIDATE_CONFIDENCE = 6
        self.topic = self.data["mqttTopic"]+"/"+str(random.randrange(0,1000))
        self.sender = imagezmq.ImageSender(connect_to = self.data["imagezmqSocket"])
        self.all_objects_current_count = 0
        self.all_objects_total_count = 0
        return
    
    def shelfPostProcess(self,frame):
        with frame.data() as mat:
            for roi in frame.regions():
                if roi.confidence() > self.thresh:
                    label = roi.label_id()
                    if self.labels[label] != "person" and self.labels[label] != "toothbrush" and self.labels[label] != "knife":
                        self.current_count[label] += 1
                        bbox = roi.rect()
                        cv2.putText(mat, self.labels[label],(int(bbox.x + bbox.w/2), int(bbox.y + bbox.h/2)),cv2.FONT_HERSHEY_SIMPLEX, 1, (255, 255, 255), 1)
                        cv2.rectangle(mat, (bbox.x, bbox.y), (bbox.x + bbox.w, bbox.y + bbox.h), (0, 255, 0), 1)

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

            for idx, label in enumerate(self.labels):
                self.all_objects_current_count += self.current_count[idx]
                self.all_objects_total_count += self.total_count[idx]

            mat = cv2.putText(mat, 'use-case='+self.data["pipelineName"], (10,30), cv2.FONT_HERSHEY_SIMPLEX, 1, (255,0,0), thickness=3)
            for idx, label in enumerate(self.labels):
                cv2.putText(mat, " current count:"+str(self.all_objects_current_count), (10, 60),cv2.FONT_HERSHEY_SIMPLEX, 1, (255, 255, 255), 1)
                cv2.putText(mat, " total count:"+str(self.all_objects_total_count), (10, 80),cv2.FONT_HERSHEY_SIMPLEX, 1, (255, 255, 255), 1)
            self.sender.send_image(self.topic, cv2.resize(mat[:,:,:3], (480,360)))
        self.current_count = [0] * len(self.labels)
        self.all_objects_total_count = 0
        self.all_objects_current_count = 0
        return True