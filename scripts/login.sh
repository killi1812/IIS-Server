!#/bin/bash

curl -X POST http://localhost:5555/login -d "username=admin&password=password"

curl -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFkbWluIiwiZXhwIjoxNzQxNTM3MDc5fQ.J48uOcbl8QePmecJ1tTOMOLYkAgC1CgSKaWIkUX6H9o" http://localhost:5555/secure

curl -X POST http://localhost:5555/weather -H "Content-Type: text/xml" --data \
    '<?xml version="1.0"?>
<methodCall>
    <methodName>GetTemp</methodName>
<params>
        <param><value><string>Zagreb</string></value></param>
</params>
</methodCall>' | xmllint --format -

curl -X POST http://localhost:5555/upload/rng \
    -F "file=@data.xml"
