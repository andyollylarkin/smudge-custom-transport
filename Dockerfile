# This is a multi-stage Dockerfile. The first part executes a build in a Golang
# container, and the second retrieves the binary from the build container and
# inserts it into a "scratch" image.

# Part 1: Compile the binary in a containerized Go environment
#
FROM golang:1.20 as build

COPY . /build


RUN cd /build; CGO_ENABLED=0 GOOS=linux go build -o /smudge ./smudge/smudge.go


# Part 3: Build the Smudge image proper
#
FROM alpine as image

COPY --from=build /smudge .

EXPOSE 9999

CMD ["/smudge"]
