#!/usr/bin/python
# -*- coding: utf-8 -*-

import json

f = open('data.json','r')
data = json.load(f)
print(data)
for i in data:
    print(i,':', data[i])
f.close()
