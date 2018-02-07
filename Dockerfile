FROM alpine
MAINTAINER Quentin McGaw <quentin.mcgaw@gmail.com>
RUN apk add --no-cache openvpn curl unzip
RUN mkdir /pia
WORKDIR /pia
RUN curl https://www.privateinternetaccess.com/openvpn/openvpn.zip > openvpn.zip && unzip openvpn.zip && rm openvpn.zip
RUN apk del curl unzip
COPY script.sh ./
RUN chmod +x script.sh
ENV REGION="Romania"
ENTRYPOINT ["/pia/script.sh"]