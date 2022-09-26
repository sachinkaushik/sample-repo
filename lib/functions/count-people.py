// name: Count People : Python
count = 0
message = msg['payload']
src = message['source']
video = src.rsplit('/', 1)[-1]
for obj in message['objects']:
    if obj['detection']['label_id'] == 1:
        count = count + 1
msg['payload'] = { 'num_of_person': count, 'video_src': video }
return msg