DoxBuilder
==========

    doxbuilder [-c configuration.yml] [-p port]

## Requirements
* Libreoffice commandline tools (soffice)

## How to convert a docx 
Send a post request to <code>/convert?format=[format]</code> as multipart/form-data and set the file paramter to the docx file. You'll get the new document returned.

__Example in Curl__

    curl -i -X POST -H "Content-Type: multipart/form-data" -F "file=@filename.docx"  http://localhost:3000/convert?format=pdf
    
## Configuration
By default doxbuilder wil look for a configuration file in <code>./configuration.yml</code>.
You can run doxbuilder with a custom configuration file:

    doxbuilder -c <myfile>.yml
    
__Example Configuration:__

    # configuration.yml
    output_dir: /tmp
    allowed_formats:
      - pdf
      - odt