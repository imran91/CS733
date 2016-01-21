# CS733
Versioned File Server

#Installation

#Command Specification
1.	Write: create a file, or update the file’s contents if it already exists.
	
	write <filename> <numbytes> [<exptime>]\r\n
    <content bytes>\r\n
	
	The server responds with the following:
	
	OK <version>\r\n

2.	Read: Given a filename, retrieve the corresponding file:
	
	read <filename>\r\n
	
	The server responds with the following format (or one of the errors described later)
	
	CONTENTS <version> <numbytes> <exptime> \r\n
	<content bytes>\r\n

3.	Compare and swap. This replaces the old file contents with the new content
	provided the version is still the same.
	
	cas <filename> <version> <numbytes> [<exptime>]\r\n
	<content bytes>\r\n
	
	The server responds with the new version if successful 
	
	OK <version>\r\n

4.	Delete file
	
	delete <filename>\r\n
	
	The server response (if successful)
	
	OK\r\n

#Errors

1.	ERR_VERSION <newversion>\r\n (the contents were not updated because of a
	version mismatch. The latest version is returned)

2.	ERR_FILE_NOT_FOUND\r\n (the filename doesn’t exist)

3.	ERR_CMD_ERR\r\n (the command is not formatted correctly)

4.	ERR_INTERNAL\r\n (any other error you wish to report that is not covered by the
	rest (optional))