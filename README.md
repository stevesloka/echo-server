# echo-server
Simple echo server for demos

## Setup
- Args:
  - echotext: The text to specify in the "app" section of the request
  - response-delay: The delay to add to a http response (ex. `--response-delay=2s`)
  - listen-port: The port used to listen on (Defaults to 8080)
  
## Example

Viewing in a web browser you get a clean UI to visualize the request. 
If wanting to view via curl or a terminal, append a querystring param `?format=text` to specify a text based output.

#### HTML Output
![example output](img/output.png)

#### Text Output
![example text output](img/output-text.png)
