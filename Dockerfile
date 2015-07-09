FROM scratch
ADD anomi /anomi 
EXPOSE 8080 
CMD ["/anomi"]
