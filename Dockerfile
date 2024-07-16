FROM scratch
ENTRYPOINT ["/cpu-limits-and-go"]
COPY cpu-limits-and-go /cpu-limits-and-go
