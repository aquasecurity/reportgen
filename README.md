# Reportgen

Generate PDF reports for Aqua CSP image and host vulnerabilities.

## Build 

To build the container image, clone this project and run the following command:

```
docker build -t aqua-reports .
```

## Generate Reports

To generate a PDF report, you should have an Aqua CSP console, user and password to this console and an image that you would like to create a report for.

After creating the reportgen image, you should run the following command to generate a PDF report for a specific image:

```
docker run -it -v /tmp:/reports aqua-reports -server <Aqua Server URL> -user <User> -password <Password> -image <Image name> -registry <Registry Name> -output /reports/report.pdf
```

This command will generate a report file called "report.pdf" and put it under /reports directory in the container, which is mounted to "/tmp" directory on your host.
