# Custom URL-shortener

## Short Description

A RESTful web service build on top of **Redis** and **Go**.

## Start a service

docker build . -t url-short

docker run -p 8080:8080 url-short

## Usage

### Endcode link

**POST** to http://localhost:8080/encode/

**Description:**

A handler that encodes your link into number with base 62 and length 10.

**Request**:  
type _json_  
{  
 "url":"https://codex.so/link-shortener",  
 "expires":"2021-10-10 11:11:11"
}

**Response**:  
type _json_  
{  
 "shortUrl":"https://localhost:8080/xxxxxxxxxx",  
 "success": true
}

### Redirect link

**GET** to ttp://localhost:8080/{shortLink}

**Description:**  
Simply redirects you to the previously encoded link, incresasing a view counter.

### Info link

**GET** to http://localhost:8080/{shortLink}/info

**Response**:

{  
 "shortUrl":{  
"id": 1,  
"url": "http://localhost:8080/xxxxxxxxxx",  
"expires": "2021-10-10 11:11:11",  
"visits": 1  
},  
 "success": true
}
