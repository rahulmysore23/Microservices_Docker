# Prerequisites
1. Install and setup docker on the system.
2. A Microsoft Azure account and a container within a storage account.

# How to Install/Use
1. Unzip the files.
2. Open microservice2\dockerfile.
3. Enter\Fill the details of your AZURE_STORAGE_ACCOUNT(name), AZURE_STORAGE_ACCESS_KEY, AZURE_STORAGE_CONTAINER(name).
4. Open cmd.
5. Cd to the folder which has the docker-compose.yml.
    ```sh
    $ cd <path_to_docker-compose.yml>
    ```
6. Run command - docker-compose build
    ```sh
    $ docker-compose build
    ```
7. Run command - docker-compose up
    ```sh
    $ docker-compose up
    ```
8. After a few mins the containers should be created and running.
9. Use the API documentation document for the list of APIs supported by the microservices.(under API_Docs folder open index.html)
 
### Example Senario

1. Make an API call to 'http://localhost:7070/upload' with 2 images as form data.
2. The Example_images folder has 2 example images to use as an input to the API.
3. The API will upload the images to the cloud with names - known_image.png, unknown_image.png
4. Then it will make a call to another microservice to verify if the unknown face from unknown_image is matching with the known face from known_image.


### Tech Used
* [Python] - For Microservice 1
* [GoLang] - For Microservice 2
* [Docker] - For deploying the services
* [Swagger] - For API Docs

### Frameworks Used
* [Django] - Webframework for python
* [Mux] - Web toolkit for GoLang
* [GORM] - ORM library for GoLang
* [Azure for go] - Azure SDK for GoLang

###### Note - For a successful response the image should be of type jpg and contain not more or less than 1 face. 

[//]: # (These are reference links used in the body of this note and get stripped out when the markdown processor does its job. There is no need to format nicely because it shouldn't be seen. Thanks SO - http://stackoverflow.com/questions/4823468/store-comments-in-markdown-syntax)

   [Python]: <https://www.python.org/>
   [GoLang]: <https://www.docker.com/>
   [Docker]: <https://golang.org/>
   [Django]: <https://www.djangoproject.com/>
   [Mux]: <http://www.gorillatoolkit.org/pkg/mux>
   [GORM]: <https://gorm.io/> 
   [Azure for go]: <https://docs.microsoft.com/en-us/azure/developer/go/>
   [Swagger]: <https://swagger.io/>