## Backend Server
An exercise in Masochism

-------------------------
- [Abstract](#abstract)
- [Setup](#setup)
- [Usage](#usage)

### Abstract
According to me, this was the most interesting task. That is chiefly because I couldn't reverse engineer the binary. The other two (one, idk, I am writing this readme on 14/08 so I have a day to complete the third task) tasks I did were mundane because I am reasonably comfortable with both Python and Web Scraping/ML. And while I have built backend servers in the past, they have been in Python. For this task I decided to learn Go, and it was an interesting experience.

-------------------------
### Setup
The backend server requires a database. I have used mongoDB because that is another thing I wanted to learn, and I am not a boomer to use SQL. To spin up the database, run
```bash
docker-compose up
```

Port 8069 is the port local machine can access mongodb at. 

To run the server, 
```bash
cd backend
go build 
./backend-api
```
The server listens on port 12345.

-------------------------
### Usage
The server allows 5 operations, 
1. Adding a Entry
The API enpoint is **/create_entry**. The request type is **POST**.
The request body is JSON.
It returns the _id of inserted record.
<p align="center"><img src="images/add-postman.png" width="500"></p>

2. Deleting an Entry
The API endpoint is **/delete**. The request type is **POST**. The request body is JSON. It deletes every record that matches the filter. It returns the number of deleted records.
<p align="center"><img src="images/del-postman.png" width="500"></p>

3. Filtering Entries
The API endpoint is **/filter/(by)/(value)**. The request type is **GET**. I did not make it a post request because we were only required to apply filters on one attribute.
<p align="center"><img src="images/filter-postman.png" width="500"></p>

4. Updating Entries
The API endpoint is **/update**. Request method is **POST**. It returns the number of records updated.
<p align="center"><img src="images/update-postman.png" width="500"></p>

5. Getting All Entries
The API endpoint is **/filter/(by)/(value)**. The request type is **GET**. 
<p align="center"><img src="images/get_all-postman.png" width="500"></p>