# CS733
Versioned File Server

#Installation
1.Download a copy using:	go get github.com/imran91/CS733

2.Change directory to : cd $GOPATH/src/github.com/imran91/CS733/Assignment1

3.Run the test script using : go test

#Command Specification
1.	Write: create a file, or update the file’s contents if it already exists. <br/>
	write &lt;filename&gt; &lt;numbytes&gt; [&lt;exptime&gt;]\r\n<br/>
    	&lt;content bytes&gt;\r\n<br/>

	The server responds with the following:<br/>	
	OK &lt;version&gt;\r\n<br/>

2.	Read: Given a filename, retrieve the corresponding file:<br/>
    read &lt;filename&gt;\r\n<br/>
	
	The server responds with the following format (or one of the errors described later)<br/>
	CONTENTS &lt;version&gt; &lt;numbytes&gt; &lt;exptime&gt; \r\n<br/>
	&lt;content bytes&gt;\r\n<br/>

3.	Compare and swap. This replaces the old file contents with the new content
	provided the version is still the same.<br/>
	cas &lt;filename&gt; &lt;version&gt; &lt;numbytes&gt; [&lt;exptime&gt;]\r\n<br/>
	&lt;content bytes&gt;\r\n<br/>
	
	The server responds with the new version if successful <br/>
	OK &lt;version&gt;\r\n<br/>

4.	Delete file<br/>
	delete &lt;filename&gt;\r\n<br/>
	
	The server response (if successful)<br/>
	OK\r\n

#Errors

1.	ERR_VERSION &lt;newversion&gt;\r\n (the contents were not updated because of a version mismatch. The latest version is returned)
2.	ERR_FILE_NOT_FOUND\r\n (the filename doesn’t exist)
3.	ERR_CMD_ERR\r\n (the command is not formatted correctly)
4.	ERR_INTERNAL\r\n (remaining errors)
