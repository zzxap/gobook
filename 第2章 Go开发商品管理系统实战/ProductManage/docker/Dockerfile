FROM centos
RUN mkdir -p  /home/app/product/
WORKDIR /home/app/product/
ADD . /home/app/product/
EXPOSE 8090
RUN chmod 771 demo
CMD [ "sh", "-c", "nohup  ./demo & " ]

