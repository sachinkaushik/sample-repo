import cv2
from collections import namedtuple
import time
import sys
import math, random
import imagezmq
import numpy as np


class Person:
    """
    Store the data of the people for tracking
    """
    def __init__(self, p_id, in_time):
        self.id = p_id
        self.counted = False
        self.gone = False
        self.in_time = in_time
        self.out_time = None
        self.looking = 0
        self.positive = 0
        self.negative = 0
        self.neutral = 0
        self.sentiment = ''

class CentroidInit:
    """
    Store centroid details of the face detected for tracking
    """
    def __init__(self, p_id, point, gone_count):
        self.id = p_id
        self.point = point
        self.gone_count = gone_count

class Centroid:
    """
    Store centroid details of the face detected for tracking
    """
    def __init__(self):
        self.centroids = []
        self.tracked_person = []
        self.person_id = 0
        self.MAX_FRAME_GONE = 20
        self.INTEREST_COUNT_TIME = 5
        self.interested= 0
        self.not_interested= 0
        self.CENTROID_DISTANCE = 150

    def remove_centroid(self,p_id):
        """
        Remove the centroid from the "centroids" list when the person is out of the frame and
        set the person.gone variable as true

        :param p_id: ID of the person whose centroid data has to be deleted
        :return: None
        """
        for idx, centroid in enumerate(self.centroids):
            if centroid.id is p_id:
                del self.centroids[idx]
                break

        if self.tracked_person[p_id]:
            self.tracked_person[p_id].gone = True
            self.tracked_person[p_id].out_time = time.time()

    def add_centroid(self, point):
        """
        Add the centroid of the object to the "centroids" list

        :param point: Centroid point to be added
        :return: None
        """
        centroid = CentroidInit(self.person_id, point, gone_count=0)
        person = Person(self.person_id, time.time())
        self.centroids.append(centroid)
        self.tracked_person.append(person)
        self.person_id = self.person_id + 1
    
    def closest_centroid(self, point):
        """
        Find the closest centroid

        :param point: Coordinate of the point for which the closest centroid point has to be detected
        :return p_idx: Id of the closest centroid
                dist: Distance of point from the closest centroid
        """
        p_idx = 0
        dist = sys.float_info.max

        for idx, centroid in enumerate(self.centroids):
            _point = centroid.point
            dx = point.x - _point.x
            dy = point.y - _point.y
            _dist = math.sqrt(dx*dx + dy*dy)
            if _dist < dist:
                dist = _dist
                p_idx = centroid.id

        return [p_idx, dist]

    def update_centroid(self, points, looking, sentiment, fps):
        """
        Update the centroid data in the centroids list and check whether the person is  or not interested

        :param points: List of centroids of the faces detected
        :param looking: List of bool values indicating if the person is looking at the camera or not
        :param sentiment: List containing the mood of the people looking at the camera
        :param fps: FPS of the input stream
        :return: None
        """
        if len(points) is 0:
            for idx, centroid in enumerate(self.centroids):
                centroid.gone_count += 1
                if centroid.gone_count > self.MAX_FRAME_GONE:
                    self.remove_centroid(centroid.id)

        if not self.centroids:
            for idx, point in enumerate(points):
                self.add_centroid(point)
        else:
            checked_points = len(points) * [None]
            checked_points_dist = len(points) * [None]
            for idx, point in enumerate(points):
                p_id, dist = self.closest_centroid(point)
                if dist > self.CENTROID_DISTANCE:
                    continue

                if p_id in checked_points:
                    p_idx = checked_points.index(p_id)
                    if checked_points_dist[p_idx] > dist:
                        checked_points[p_idx] = None
                        checked_points_dist[p_idx] = None

                checked_points[idx] = p_id
                checked_points_dist[idx] = dist
            for centroid in self.centroids:
                if centroid.id in checked_points:
                    p_idx = checked_points.index(centroid.id)
                    centroid.point = points[p_idx]
                    centroid.gone_count = 0
                else:
                    centroid.gone_count += 1
                    if centroid.gone_count > self.MAX_FRAME_GONE:
                        self.remove_centroid(centroid.id)
            for idx in range(len(checked_points)):
                if checked_points[idx] is None:
                    self.add_centroid(points[idx])
                else:
                    if looking[idx] is True:
                        self.tracked_person[checked_points[idx]].sentiment = sentiment[idx]
                        self.tracked_person[checked_points[idx]].looking += 1
                        if sentiment[idx] == "happy" or sentiment[idx] == "surprise":
                            self.tracked_person[checked_points[idx]].positive += 1
                        elif sentiment[idx] == 'sad' or sentiment[idx] == 'anger':
                            self.tracked_person[checked_points[idx]].negative += 1
                        elif sentiment[idx] == 'neutral':
                            self.tracked_person[checked_points[idx]].neutral += 1
                    else:
                        self.tracked_person[checked_points[idx]].sentiment = "Not looking"

            for person in self.tracked_person:
                if person.counted is False:
                    positive = person.positive + person.neutral

                    # If the person is looking at the camera for specified time
                    # and his mood is positive, increment the interested variable
                    if (person.looking > fps * self.INTEREST_COUNT_TIME) and (positive > person.negative):
                        self.interested += 1
                        person.counted = True

                        # If the person is gone out of the frame, increment the not_ variable
                        if person.gone is True:
                            self.not_interested= 1
                            person.counted = True
        return self.centroids, self.tracked_person

class PostProcess:
    
    def __init__(self, imagezmqSocket, mqttTopic, pipelineName):
        self.data = {"imagezmqSocket": imagezmqSocket, "mqttTopic": mqttTopic, "pipelineName": pipelineName}
        self.Point = namedtuple("Point", "x,y")
        self.topic = self.data["mqttTopic"]+"/"+str(random.randrange(1001,2000))
        self.sender = imagezmq.ImageSender(connect_to = self.data["imagezmqSocket"])
        self.centroid_obj = Centroid()
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
        with frame.data() as mat:

            for roi in frame.regions():
                if roi.label() == "looking":
                    continue
                pt = roi.rect()
                cv2.rectangle(mat, (pt.x, pt.y), (pt.x+pt.w, pt.y+pt.h), (255, 255, 0), 1)
                x = pt.x + int(pt.w / 2)
                y = pt.y + int(pt.h / 2)
                point = self.Point(x, y)
                self.frame_centroids.append(point)

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
                            cv2.putText(mat,'looking=True',(x,y),cv2.FONT_HERSHEY_COMPLEX,1,(255,0,0),thickness=1)
                        else:
                            self.looking.append(False)
                            self.sentiment.append(-1)
                            cv2.putText(mat,'looking=False',(x,y),cv2.FONT_HERSHEY_COMPLEX,1,(255,0,0),thickness=1)

            centroids, tracked_person = self.centroid_obj.update_centroid(self.frame_centroids, self.looking, self.sentiment, 150)
            mat = cv2.putText(mat, 'use-case='+self.data["pipelineName"], (10,30), cv2.FONT_HERSHEY_SIMPLEX, 1, (255,0,0), thickness=3)
            for centroid in centroids:
                cv2.rectangle(mat, (centroid.point.x, centroid.point.y), (centroid.point.x + 1, centroid.point.y + 1), (0, 255, 0), 1)
                cv2.putText(mat, "person:{}".format(centroid.id), (centroid.point.x + 1, centroid.point.y - 5),cv2.FONT_HERSHEY_SIMPLEX, 0.5, (255, 255, 255), 1)
                ht = 60
            for person in tracked_person:
                if person.gone is False:
                    message = "Person {} is {}".format(person.id, person.sentiment)
                    cv2.putText(mat, message, (10, ht), cv2.FONT_HERSHEY_COMPLEX, 1,(255, 255, 255), 1)
                    ht += 20
            self.sender.send_image(self.topic, cv2.resize(mat[:,:,:3], (480,360)))

        return True