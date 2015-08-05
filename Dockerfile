FROM golang
RUN apt-get update && \
        apt-get -y install unzip
RUN wget https://github.com/getqujing/qtunnel/archive/master.zip
RUN unzip master.zip && cd qtunnel* && make && mv bin/qtunnel /tmp/qtunnel
EXPOSE 8080
CMD /tmp/qtunnel -listen=:8080 -secret="u6H?6fTKmBbp&wuu" -crypto=aes256cfb -backend=106.185.42.18:8100 -clientmode=false
